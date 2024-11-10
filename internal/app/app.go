package app

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/config"
	"github.com/valek177/auth/internal/interceptor"
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
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
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
	wg.Add(4)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		log.Printf("Started user saver consumer")

		consumer, err := a.serviceProvider.UserSaverConsumer(ctx)
		if err != nil {
			log.Printf("failed to create consumer: %s", err.Error())
		}
		err = consumer.RunConsumer(ctx)
		if err != nil {
			log.Printf("failed to run consumer: %s", err.Error())
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
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
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor((interceptor.ValidateInterceptor)),
	)

	reflection.Register(a.grpcServer)

	authImpl, err := a.serviceProvider.AuthImpl(ctx)
	if err != nil {
		return err
	}

	user_v1.RegisterUserV1Server(a.grpcServer, authImpl)

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

func (a *App) runGRPCServer() error {
	grpcConfig, err := a.serviceProvider.GRPCConfig()
	if err != nil {
		return err
	}
	log.Printf("GRPC server is running on %s", grpcConfig.Address())

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
	log.Printf("HTTP server is running on %s", httpConfig.Address())

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

	log.Printf("Swagger server is running on %s", cfg.Address())

	err = a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			_ = file.Close()
		}()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
