package app

import (
	"context"
	"flag"
	"io"
	"net"
	"net/http"
	"os"
	"sync"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/natefinch/lumberjack"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/valek177/auth/grpc/pkg/access_v1"
	"github.com/valek177/auth/grpc/pkg/auth_v1"
	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/config"
	"github.com/valek177/auth/internal/interceptor"
	"github.com/valek177/auth/internal/logger"
	"github.com/valek177/auth/internal/metric"
	"github.com/valek177/auth/internal/tracing"
	_ "github.com/valek177/auth/statik" //nolint:revive
	"github.com/valek177/platform-common/pkg/closer"
)

var (
	configPath                string
	corsAllowedOriginsDefault = []string{"*"}
	corsAllowedMethodsDefault = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsAllowedHeadersDefault = []string{"Accept", "Content-Type", "Content-Length", "Authorization"}
)

const corsAllowCredentialsDefault = true

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// App contains application object
type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
}

// NewApp creates new App object
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run runs application
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(5)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logger.FatalWithMsg("failed to run GRPC server: ", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logger.FatalWithMsg("failed to run HTTP server: ", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			logger.FatalWithMsg("failed to run Swagger server: ", err)
		}
	}()

	go func() {
		defer wg.Done()

		logger.Debug("Started user saver consumer")

		consumer, err := a.serviceProvider.UserSaverConsumer(ctx)
		if err != nil {
			logger.ErrorWithMsg("failed to create consumer: ", err)
		}
		err = consumer.RunConsumer(ctx)
		if err != nil {
			logger.ErrorWithMsg("failed to run consumer: ", err)
		}
	}()

	go func() {
		defer wg.Done()

		if err := a.runPrometheusServer(); err != nil {
			logger.FatalWithMsg("failed to run prometheus HTTP server: ", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initJaegerTracing,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initMetrics,
		a.initPrometheusServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	grpcCfg, err := a.serviceProvider.GRPCConfig()
	if err != nil {
		return err
	}

	creds, err := credentials.NewServerTLSFromFile(grpcCfg.TLSCertFile(), grpcCfg.TLSKeyFile())
	if err != nil {
		return err
	}

	logger.Init(getCore(getAtomicLevel(grpcCfg.LogLevel())))

	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.LogInterceptor,
				interceptor.MetricsInterceptor,
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				interceptor.ValidateInterceptor,
			),
		),
	)

	reflection.Register(a.grpcServer)

	userImpl, err := a.serviceProvider.UserImpl(ctx)
	if err != nil {
		return err
	}

	authImpl, err := a.serviceProvider.AuthImpl(ctx)
	if err != nil {
		return err
	}

	accessImpl, err := a.serviceProvider.AccessImpl(ctx)
	if err != nil {
		return err
	}

	user_v1.RegisterUserV1Server(a.grpcServer, userImpl)
	auth_v1.RegisterAuthV1Server(a.grpcServer, authImpl)
	access_v1.RegisterAccessV1Server(a.grpcServer, accessImpl)

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	grpcCfg, err := a.serviceProvider.GRPCConfig()
	if err != nil {
		return err
	}

	err = user_v1.RegisterUserV1HandlerFromEndpoint(ctx, mux, grpcCfg.Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   corsAllowedOriginsDefault,
		AllowedMethods:   corsAllowedMethodsDefault,
		AllowedHeaders:   corsAllowedHeadersDefault,
		AllowCredentials: corsAllowCredentialsDefault,
	})

	httpCfg, err := a.serviceProvider.HTTPConfig()
	if err != nil {
		return err
	}

	a.httpServer = &http.Server{ //nolint:gosec
		Addr:    httpCfg.Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	cfg, err := a.serviceProvider.SwaggerConfig()
	if err != nil {
		return err
	}

	a.swaggerServer = &http.Server{ //nolint:gosec
		Addr:    cfg.Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	return metric.Init(ctx)
}

func (a *App) initPrometheusServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusCfg, err := a.serviceProvider.PrometheusConfig()
	if err != nil {
		return err
	}

	a.prometheusServer = &http.Server{
		Addr:    prometheusCfg.Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) initJaegerTracing(_ context.Context) error {
	cfg, err := a.serviceProvider.JaegerConfig()
	if err != nil {
		return err
	}

	return tracing.Init(cfg)
}

func (a *App) runGRPCServer() error {
	grpcConfig, err := a.serviceProvider.GRPCConfig()
	if err != nil {
		return err
	}
	logger.Info("GRPC server is running on " + grpcConfig.Address())

	list, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	httpConfig, err := a.serviceProvider.HTTPConfig()
	if err != nil {
		return err
	}
	logger.Info("HTTP server is running on " + httpConfig.Address())

	err = a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	cfg, err := a.serviceProvider.SwaggerConfig()
	if err != nil {
		return err
	}

	logger.Info("Swagger server is running on " + cfg.Address())

	err = a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		logger.Debug("Serving swagger file: " + path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug("Open swagger file: " + path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			_ = file.Close()
		}()

		logger.Debug("Read swagger file: " + path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug("Write swagger file: " + path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug("Served swagger file: " + path)
	}
}

func (a *App) runPrometheusServer() error {
	cfg, err := a.serviceProvider.PrometheusConfig()
	if err != nil {
		return err
	}
	logger.Info("Prometheus server is running on " + cfg.Address())

	err = a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel(logLevel string) zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		logger.FatalWithMsg("failed to set log level: ", err)
	}

	return zap.NewAtomicLevelAt(level)
}
