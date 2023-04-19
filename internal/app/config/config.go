package config

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"

	"github.com/anonymousMoonPrince/file-service/internal/app/client/database/postgres"
	"github.com/anonymousMoonPrince/file-service/internal/app/client/storage/minio"
)

type Config struct {
	MinioConfigs   []minio.Config  `mapstructure:"minio"`
	PostgresConfig postgres.Config `mapstructure:"postgres"`
	ServerConfig   struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	BusinessConfig struct {
		ChunkCount int    `mapstructure:"chunk_count"`
		Bucket     string `mapstructure:"bucket"`
	} `mapstructure:"business"`
}

var config Config

func init() {
	flag.Parse()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Fatal("read config failed")
	}

	if err := viper.Unmarshal(&config); err != nil {
		logrus.WithError(err).Fatal("unmarshal config failed")
	}

	viper.WatchConfig()
}

func AddConfigHook(hook func(cfg Config)) {
	viper.OnConfigChange(func(in fsnotify.Event) {
		hook(Get())
	})
}

func Get() Config {
	var localConfig Config
	if err := viper.Unmarshal(&localConfig); err != nil {
		logrus.WithError(err).Error("unmarshal config failed")
		return config
	}
	config = localConfig
	return localConfig
}
