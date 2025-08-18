package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fhir-api/controllers"
	"fhir-api/middleware"
	"fhir-api/services"
	"fhir-api/utils"
)

type App struct {
	serverPort string
	mongoURI   string
	dbName     string
	jwtSecret  string
	jwtClient  string
	router     *gin.Engine
	logger     *logrus.Logger
	mongo      *mongo.Client
}

func RunApp() *App {
	logger, err := setupLogger()
	if err != nil {
		log.Fatalf("failed to setup logger: %v", err)
	}

	cfg := loadEnvConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(cfg.mongoURI).
		SetAuth(options.Credential{
			Username: cfg.dbUser,
			Password: cfg.dbPwd,
		}))

	if err != nil {
		logger.Fatalf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatalf("failed to ping MongoDB: %v", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(utils.GinLogger(logger))

	return &App{
		serverPort: cfg.serverPort,
		mongoURI:   cfg.mongoURI,
		dbName:     cfg.dbName,
		jwtSecret:  cfg.jwtSecret,
		jwtClient:  cfg.jwtClient,
		router:     router,
		logger:     logger,
		mongo:      client,
	}
}

func setupLogger() (*logrus.Logger, error) {
	logLevel := os.Getenv("LOG_LEVEL")
	log.Printf(logLevel)

	logger := logrus.New()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)

	if os.Getenv("LOG_FORMAT") == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	if logPath := os.Getenv("LOG_PATH"); logPath != "" {
		rotationTime, err := time.ParseDuration(os.Getenv("LOG_ROTATION_TIME"))
		if err != nil {
			log.Printf("Nível de log inválido '%s', usando 'info' como padrão", logLevel)
			level = logrus.InfoLevel
		}
		maxAge, err := time.ParseDuration(os.Getenv("LOG_MAX_AGE"))
		if err != nil {
			return nil, fmt.Errorf("LOG_MAX_AGE inválido: %v", err)
		}

		log.Printf(maxAge.String())
		log.Printf(rotationTime.String())

		writer, err := rotatelogs.New(
			logPath+"/fhir-api.%Y-%m-%d.log",
			rotatelogs.WithLinkName(logPath+"/fhir-api.log"),
			rotatelogs.WithRotationTime(rotationTime),
			rotatelogs.WithMaxAge(maxAge),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file: %v", err)
		}
		logger.SetOutput(io.MultiWriter(os.Stdout, writer))
	} else {
		logger.SetOutput(os.Stdout)
	}

	return logger, nil
}

type envConfig struct {
	serverPort string
	mongoURI   string
	dbName     string
	dbUser     string
	dbPwd      string
	jwtSecret  string
	jwtClient  string
}

func loadEnvConfig() envConfig {
	return envConfig{
		serverPort: os.Getenv("SERVER_PORT"),
		mongoURI:   os.Getenv("DB_URI"),
		dbName:     os.Getenv("DB_NAME"),
		dbUser:     os.Getenv("DB_USER"),
		dbPwd:      os.Getenv("DB_PWD"),
		jwtSecret:  os.Getenv("JWT_SECRET_KEY"),
		jwtClient:  os.Getenv("JWT_CLIENT_CODE"),
	}
}

// @title Go API
// @version 1.0
// @description API para recursos FHIR
// @host api.local.<client>:8082
// @BasePath /api/v1
func (a *App) Run() error {
	db := a.mongo.Database(a.dbName)

	authService := services.NewAuthService(
		a.jwtSecret,
		a.jwtClient,
		24*time.Hour,
	)

	encounterService := services.NewEncounterService(db, a.logger)
	encounterController := controllers.NewEncounterController(encounterService)

	patientService := services.NewPatientService(db, a.logger)
	patientController := controllers.NewPatientController(patientService)

	practitionerservice := services.NewPractitionerService(db, a.logger)
	practitionerController := controllers.NewPractitionerController(practitionerservice)

	router := a.router
	api := router.Group("/api/v1")
	{

		// Rota do Swagger
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Prometheus
		// api.Use(middleware.PrometheusMiddleware())
		api.GET("/metrics", gin.WrapH(promhttp.Handler()))

		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		api.POST("/auth/token", func(c *gin.Context) {
			token, err := authService.GenerateToken()
			if err != nil {
				a.logger.WithError(err).Error("failed to generate token")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": token})
		})

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(a.jwtSecret, a.jwtClient))
		{

			protected.GET("/patients/:id", patientController.GetPatient)
			protected.GET("/practitioners/:id", practitionerController.GetPractitioner)
			protected.GET("/encounters/:id", encounterController.GetEncounter)
			protected.POST("/encounters/:id/review-request", encounterController.UpdateEncounterStatus)
		}
	}

	// Configurar e iniciar servidor (mesmo conteúdo anterior)
	srv := &http.Server{
		Addr:    ":" + a.serverPort,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Infof("Server is running on port %s", a.serverPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	a.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.mongo.Disconnect(ctx); err != nil {
		a.logger.Errorf("Failed to disconnect from MongoDB: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.logger.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	a.logger.Info("Server exited properly")
	return nil
}
