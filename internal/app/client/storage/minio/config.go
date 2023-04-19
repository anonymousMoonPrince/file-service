package minio

type Config struct {
	URL       string `mapstructure:"url"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Secure    bool   `mapstructure:"secure"`
	Token     string `mapstructure:"token"`
	Region    string `mapstructure:"region"`
}
