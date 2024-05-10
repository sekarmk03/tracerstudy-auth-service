package config

import (
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME,default=tracer-study-grpc"`
	Port        Port
	MySQL       MySQL
	// Redis       Redis
	// JWTConfig   JWTConfig
	JWT         JWTConfig
	SIAK_API    SIAK_API
	ClientURL ClientURL
}

type Port struct {
	GRPC string `env:"PORT_GRPC,default=8081"`
	REST string `env:"PORT_REST,default=8080"`
}

type MySQL struct {
	Host     string `env:"MYSQL_HOST,default=localhost"`
	Port     string `env:"MYSQL_PORT,default=3306"`
	User     string `env:"MYSQL_USER,default=root"`
	Password string `env:"MYSQL_PASSWORD,default=skrmk372"`
	Name     string `env:"MYSQL_NAME,default=new_tracer"`
}

type ClientURL struct {
	Pkts string `env:"CLIENT_URL_PKTS"`
	MhsBiodata string `env:"CLIENT_URL_MHSBIODATA"`
}

// type Redis struct {
// 	Address  string `env:"REDIS_ADDRESS,required"`
// 	Password string `env:"REDIS_PASSWORD"`
// }

// type JWTConfig struct {
// 	SecretKey string `env:"JWT_SECRET_KEY"`
// }

type SIAK_API struct {
	URL string `env:"SIAK_API_URL"`
	KEY string `env:"SIAK_API_KEY"`
}

type JWTConfig struct {
	JwtSecretKey  string        `env:"JWT_SECRET_KEY"`
	TokenDuration time.Duration `env:"JWT_DURATION,default=30m"`
}

func NewConfig(env string) (*Config, error) {
	_ = godotenv.Load(env)

	var config Config
	if err := envdecode.Decode(&config); err != nil {
		return nil, errors.Wrap(err, "ERROR: [NewConfig] Error while decoding env")
	}

	return &config, nil
}
