package api

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/swaggo/swag"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/time/rate"

	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"github.com/LensPlatform/Lens-users-svc/pkg/amqp"
	_ "github.com/LensPlatform/Lens-users-svc/pkg/api/docs"
	"github.com/LensPlatform/Lens-users-svc/pkg/config"
	"github.com/LensPlatform/Lens-users-svc/pkg/fscache"
	"github.com/LensPlatform/Lens-users-svc/pkg/middleware"
	"github.com/LensPlatform/Lens-users-svc/pkg/tables"

	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"

	"github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Lens Platform API
// @version 2.0
// @description Go microservice for Kubernetes.

// @contact.name Source Code
// @contact.url https://github.com/LensPlatform/Lens-users-svc

// @license.name MIT License
// @license.url https://github.com/stefanprodan/podinfo/blob/master/LICENSE

// @host localhost:9898
// @BasePath /
// @schemes http https

var (
	healthy int32
	ready   int32
	watcher *fscache.Watcher
)

type FluxConfig struct {
	GitUrl    string `mapstructure:"git-url"`
	GitBranch string `mapstructure:"git-branch"`
}

type Config struct {
	HttpClientTimeout         time.Duration `mapstructure:"http-client-timeout"`
	HttpServerTimeout         time.Duration `mapstructure:"http-server-timeout"`
	HttpServerShutdownTimeout time.Duration `mapstructure:"http-server-shutdown-timeout"`
	BackendURL                []string      `mapstructure:"backend-url"`
	UILogo                    string        `mapstructure:"ui-logo"`
	UIMessage                 string        `mapstructure:"ui-message"`
	UIColor                   string        `mapstructure:"ui-color"`
	UIPath                    string        `mapstructure:"ui-path"`
	DataPath                  string        `mapstructure:"data-path"`
	ConfigPath                string        `mapstructure:"config-path"`
	Port                      string        `mapstructure:"port"`
	PortMetrics               int           `mapstructure:"port-metrics"`
	Hostname                  string        `mapstructure:"hostname"`
	H2C                       bool          `mapstructure:"h2c"`
	RandomDelay               bool          `mapstructure:"random-delay"`
	RandomError               bool          `mapstructure:"random-error"`
	JWTSecret                 string        `mapstructure:"jwt-secret"`
}

type Server struct {
	router *mux.Router
	logger *zap.Logger
	config *Config
	zipkinEndpoint string
}

func NewServer(config *Config, logger *zap.Logger, endpoint string) (*Server, error) {
	srv := &Server{
		router: mux.NewRouter(),
		logger: logger,
		config: config,
		zipkinEndpoint: endpoint,
	}

	logger.Info("Tracer", zap.String("type of tracer", "zipkin"), zap.String("URL", endpoint))

	return srv, nil
}

func (s *Server) registerHandlers() {
	s.router.Handle("/metrics", promhttp.Handler())
	s.router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	s.router.HandleFunc("/", s.indexHandler).HeadersRegexp("User-Agent", "^Mozilla.*").Methods("GET")
	s.router.HandleFunc("/", s.infoHandler).Methods("GET")
	s.router.HandleFunc("/version", s.versionHandler).Methods("GET")
	s.router.HandleFunc("/echo", s.echoHandler).Methods("POST")
	s.router.HandleFunc("/env", s.envHandler).Methods("GET", "POST")
	s.router.HandleFunc("/headers", s.echoHeadersHandler).Methods("GET", "POST")
	s.router.HandleFunc("/delay/{wait:[0-9]+}", s.delayHandler).Methods("GET").Name("delay")
	s.router.HandleFunc("/healthz", s.healthzHandler).Methods("GET")
	s.router.HandleFunc("/readyz", s.readyzHandler).Methods("GET")
	s.router.HandleFunc("/readyz/enable", s.enableReadyHandler).Methods("POST")
	s.router.HandleFunc("/readyz/disable", s.disableReadyHandler).Methods("POST")
	s.router.HandleFunc("/panic", s.panicHandler).Methods("GET")
	s.router.HandleFunc("/status/{code:[0-9]+}", s.statusHandler).Methods("GET", "POST", "PUT").Name("status")
	s.router.HandleFunc("/store", s.storeWriteHandler).Methods("POST")
	s.router.HandleFunc("/store/{hash}", s.storeReadHandler).Methods("GET").Name("store")
	s.router.HandleFunc("/configs", s.configReadHandler).Methods("GET")
	s.router.HandleFunc("/token", s.tokenGenerateHandler).Methods("POST")
	s.router.HandleFunc("/token/validate", s.tokenValidateHandler).Methods("GET")
	s.router.HandleFunc("/api/info", s.infoHandler).Methods("GET")
	s.router.HandleFunc("/api/echo", s.echoHandler).Methods("POST")
	s.router.HandleFunc("/ws/echo", s.echoWsHandler)
	s.router.HandleFunc("/chunked", s.chunkedHandler)
	s.router.HandleFunc("/chunked/{wait:[0-9]+}", s.chunkedHandler)
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	s.router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc()
		if err != nil {
			s.logger.Error("swagger error", zap.Error(err), zap.String("path", "/swagger.json"))
		}
		w.Write([]byte(doc))
	})
}

func (s *Server) registerMiddlewares() {
	promMw := middleware.NewPrometheusMiddleware()
	s.router.Use(promMw.Handler)

	httpLoggerMw := middleware.NewLoggingMiddleware(s.logger)
	s.router.Use(httpLoggerMw.Handler)

	rateLimitterMw := middleware.NewRateLimitMiddleware(rate.Every(time.Second), 5)
	s.router.Use(rateLimitterMw.Handler)

	instrumentationMw := middleware.NewInstrumentationMiddleware("API")
	s.router.Use(instrumentationMw.Handler)

	tracer, _ := s.NewZipkinTracer()
	zipkinMw := middleware.NewZipKinTracerMiddleware("API", tracer)
	s.router.Use(zipkinMw.Handler)

	panicMw := middleware.NewPanicRecovery(*s.logger)
	s.router.Use(panicMw.Handler)

	circuitBreakerMw := middleware.NewCircuitBreaker("users_microservice", 5, 0, time.Duration(60) * time.Second, nil, *s.logger)
	s.router.Use(circuitBreakerMw.Handler)

	s.router.Use(versionMiddleware)

	if s.config.RandomDelay {
		s.router.Use(randomDelayMiddleware)
	}
	if s.config.RandomError {
		s.router.Use(randomErrorMiddleware)
	}
}

func(s *Server) NewZipkinTracer() (*zipkin.Tracer, error) {
	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(s.zipkinEndpoint)

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: "Lens-users-svc", Port: 8080}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 1 (100%) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	return t, err
}

// InitDbConnection initializes a database connection and creates associated tables/migrates schemas
func (s *Server) ConnectToDatabase() (*gorm.DB, error) {
	config.LoadConfig()
	connString := config.Config.GetDatabaseConnectionString()
	db, err := gorm.Open("postgres", connString)
	if err != nil {
		s.logger.Error(err.Error())
		os.Exit(1)
	}

	s.logger.Info("successfully connected to database")
	db.SingularTable(true)
	db.LogMode(false)
	s.CreateTablesOrMigrateSchemas(db)
	return db, err
}

// CreateTablesOrMigrateSchemas creates a given set of tables based on a schema
// if it does not exist or migrates the table schemas to the latest version
func (s *Server)CreateTablesOrMigrateSchemas(db *gorm.DB) {
	var userTable tables.UserTable
	var teamsTable tables.TeamTable
	var groupTable tables.GroupTable
	userTable.MigrateSchemaOrCreateTable(db, s.logger)
	teamsTable.MigrateSchemaOrCreateTable(db, s.logger)
	groupTable.MigrateSchemaOrCreateTable(db, s.logger)
}

// InitQueues initializes a set of producer and consumer amqp queues to be used for things such as
// account registration emails amongst many others.
func (s *Server)ConnectToQueues() (amqp.Queue, amqp.Queue) {
	amqpConnString := "amqp://user:bitnami@stats/"
	producerQueueNames := []string{"lens_welcome_email", "lens_password_reset_email", "lens_email_reset_email"}
	consumerQueueNames := []string{"user_inactive"}
	amqpproducerconn, err := amqp.NewAmqpConnection(amqpConnString, producerQueueNames)
	if err != nil {
		s.logger.Error(err.Error())
	}
	amqpconsumerconn, err := amqp.NewAmqpConnection(amqpConnString, consumerQueueNames)
	if err != nil {
		s.logger.Error(err.Error())
	}
	return amqpproducerconn, amqpconsumerconn
}

func (s *Server) ListenAndServe(stopCh <-chan struct{}) {
	go s.startMetricsServer()

	s.registerHandlers()
	s.registerMiddlewares()
	_, err := s.ConnectToDatabase()
	if err != nil {
		s.logger.Error(err.Error())
	}
	s.ConnectToQueues()

	var handler http.Handler
	if s.config.H2C {
		handler = h2c.NewHandler(s.router, &http2.Server{})
	} else {
		handler = s.router
	}

	srv := &http.Server{
		Addr:         ":" + s.config.Port,
		WriteTimeout: s.config.HttpServerTimeout,
		ReadTimeout:  s.config.HttpServerTimeout,
		IdleTimeout:  2 * s.config.HttpServerTimeout,
		Handler:      handler,
	}

	//s.printRoutes()

	// load configs in memory and start watching for changes in the config dir
	if stat, err := os.Stat(s.config.ConfigPath); err == nil && stat.IsDir() {
		var err error
		watcher, err = fscache.NewWatch(s.config.ConfigPath)
		if err != nil {
			s.logger.Error("config watch error", zap.Error(err), zap.String("path", s.config.ConfigPath))
		} else {
			watcher.Watch()
		}
	}

	// run server in background
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal("HTTP server crashed", zap.Error(err))
		}
	}()

	// signal Kubernetes the server is ready to receive traffic
	atomic.StoreInt32(&healthy, 1)
	atomic.StoreInt32(&ready, 1)

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), s.config.HttpServerShutdownTimeout)
	defer cancel()

	// all calls to /healthz and /readyz will fail from now on
	atomic.StoreInt32(&healthy, 0)
	atomic.StoreInt32(&ready, 0)

	s.logger.Info("Shutting down HTTP server", zap.Duration("timeout", s.config.HttpServerShutdownTimeout))

	// wait for Kubernetes readiness probe to remove this instance from the load balancer
	// the readiness check interval must be lower than the timeout
	if viper.GetString("level") != "debug" {
		time.Sleep(3 * time.Second)
	}

	// attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Warn("HTTP server graceful shutdown failed", zap.Error(err))
	} else {
		s.logger.Info("HTTP server stopped")
	}
}

func (s *Server) startMetricsServer() {
	if s.config.PortMetrics > 0 {
		mux := http.DefaultServeMux
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%v", s.config.PortMetrics),
			Handler: mux,
		}

		srv.ListenAndServe()
	}
}

func (s *Server) printRoutes() {
	s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})
}

type ArrayResponse []string
type MapResponse map[string]string
