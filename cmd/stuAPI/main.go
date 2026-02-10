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
	"github.com/Sarthak-D97/go_stuAPI/internal/storage/sqlite"
	"github.com/Sarthak-D97/go_stuAPI/middlewares"
	"github.com/Sarthak-D97/go_stuAPI/repository"
	"github.com/Sarthak-D97/go_stuAPI/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	gindump "github.com/tpkeeper/gin-dump"
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func main() {
	setupLogOutput()
	cfg := config.MustLoad()
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("SQLite setup failed:", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	jwtService := service.NewJWTService()
	LoginService := service.NewLoginService()
	loginController := controller.NewLoginController(*LoginService, jwtService)

	studentService := service.NewStudentService(storage, rdb)
	studentController := controller.NewStudentController(studentService)

	videoRepository := repository.NewVideoRepository()
	defer videoRepository.CloseDB()

	videoService := service.NewVideoService(videoRepository)
	videoController := controller.New(videoService)
	router := gin.New()
	router.Use(gin.Recovery(), middlewares.Logger(), gindump.Dump())
	router.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{"token": token})
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	})
	api := router.Group("/api", middlewares.AuthorizeJWT())
	{
		students := api.Group("/students")
		{
			students.POST("/", studentController.Create)
			students.GET("/:id", studentController.GetByID)
			students.PUT("/:id", studentController.Update)
			students.GET("/", studentController.GetList)
			students.DELETE("/:id", studentController.Delete)

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
	srv := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
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
