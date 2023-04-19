package app

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/tabalt/gracehttp"

	"github.com/anonymousMoonPrince/file-service/internal/app/client/database/postgres"
	"github.com/anonymousMoonPrince/file-service/internal/app/client/storage/minio"
	"github.com/anonymousMoonPrince/file-service/internal/app/config"
	"github.com/anonymousMoonPrince/file-service/internal/app/controller"
	"github.com/anonymousMoonPrince/file-service/internal/app/repository"
	"github.com/anonymousMoonPrince/file-service/internal/app/service"
)

type App struct {
	port string

	// controllers
	fileController *controller.FileController
}

func NewApp() *App {
	cfg := config.Get()

	databaseClient := postgres.NewClient(cfg.PostgresConfig)

	fileRepository := repository.NewFileRepository(databaseClient)
	sizeByURL, err := fileRepository.GetSizeByURL(context.Background())
	if err != nil {
		logrus.WithError(err).Fatal("get bucket size failed")
	}

	storageClient := minio.NewClient(cfg.MinioConfigs, sizeByURL)

	config.AddConfigHook(func(cfg config.Config) {
		sizeByURL, err := fileRepository.GetSizeByURL(context.Background())
		if err != nil {
			return
		}

		storageClient.UpdateClient(cfg.MinioConfigs, sizeByURL)
	})

	fileService := service.NewFileService(storageClient, fileRepository, cfg.BusinessConfig.Bucket)
	fileController := controller.NewFileController(fileService)

	return &App{
		port:           cfg.ServerConfig.Port,
		fileController: fileController,
	}
}

func (a *App) Run() {
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)
	r.Put("/", a.fileController.Upload)
	r.Get("/{uuid}", a.fileController.Download)

	logrus.Infof("Listening port [:%s]", a.port)
	logrus.Infof("Server closed: %v", gracehttp.NewServer(":"+a.port, r, 0, 0).ListenAndServe())
}
