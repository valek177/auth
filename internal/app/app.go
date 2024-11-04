package app

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/config"
	"github.com/valek177/platform-common/pkg/closer"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// App contains application object
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
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
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

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

	// go func() {
	// 	defer wg.Done()

	// 	err := a.runSwaggerServer()
	// 	if err != nil {
	// 		log.Fatalf("failed to run Swagger server: %v", err)
	// 	}
	// }()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
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
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

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
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	httpCfg, err := a.serviceProvider.HTTPConfig()
	if err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    httpCfg.Address(),
		Handler: corsMiddleware.Handler(mux),
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
