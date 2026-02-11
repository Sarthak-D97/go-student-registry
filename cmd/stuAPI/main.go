package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sarthak-D97/go_stuAPI/controller"
	"github.com/Sarthak-D97/go_stuAPI/internal/config"
	db "github.com/Sarthak-D97/go_stuAPI/internal/platform/db"
	redisclient "github.com/Sarthak-D97/go_stuAPI/internal/platform/redis"
	"github.com/Sarthak-D97/go_stuAPI/middlewares"
	"github.com/Sarthak-D97/go_stuAPI/repository"
	studentRepoImpl "github.com/Sarthak-D97/go_stuAPI/repository"
	"github.com/Sarthak-D97/go_stuAPI/service"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	setupLogOutput()
	cfg := config.MustLoad()

	// 1. Initialize Postgres
	pgDB, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal("Postgres setup failed:", err)
	}

	// 2. Initialize Redis
	rdb := redisclient.NewClient()

	// 3. Initialize Services
	jwtService := service.NewJWTService()
	loginService := service.NewLoginService()
	loginController := controller.NewLoginController(*loginService, jwtService)

	// FIX 2: Create the Repository first, passing the DB connection to it
	// The service needs the Repository, NOT the raw pgDB connection
	studentRepo := studentRepoImpl.New(pgDB)

	// FIX 3: Pass the initialized repository to the service
	studentService := service.NewStudentService(studentRepo, rdb)
	studentController := controller.NewStudentController(studentService)

	videoRepository := repository.NewVideoRepository()
	defer videoRepository.CloseDB()

	videoService := service.NewVideoService(videoRepository)
	videoController := controller.New(videoService)

	// 4. Router Setup
	router := gin.New()
	router.Use(gin.Recovery(), middlewares.Logger(), gindump.Dump())

	// Public Routes
	router.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{"token": token})
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	})

	// Private Routes
	api := router.Group("/api", middlewares.AuthorizeJWT(jwtService))
	{
		students := api.Group("/students")
		{
			students.POST("/", studentController.Create)
			students.GET("/:id", studentController.GetByID)
			students.PUT("/:id", studentController.Update)
			students.GET("/", studentController.GetList)
			students.DELETE("/:id", studentController.Delete)
		}

		videos := api.Group("/videos")
		{
			videos.GET("/", videoController.FindAll)
			videos.POST("/", func(ctx *gin.Context) {
				_ = videoController.Save(ctx)
			})
			videos.PUT("/:id", func(ctx *gin.Context) {
				_ = videoController.Update(ctx)
			})
			videos.DELETE("/:id", func(ctx *gin.Context) {
				_ = videoController.Delete(ctx)
			})
		}
	}

	// 5. Server Startup
	srv := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	slog.Info("Server exiting")
}
