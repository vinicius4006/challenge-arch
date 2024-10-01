package configs

import (
	"errors"

	"github.com/joeshaw/envdecode"
)

type conf struct {
	DBDriver          string `env:"DB_DRIVER,required"`
	DBHost            string `env:"DB_HOST,required"`
	DBPort            string `env:"DB_PORT,required"`
	DBUser            string `env:"DB_USER,required"`
	DBPassword        string `env:"DB_PASSWORD,required"`
	DBName            string `env:"DB_NAME,required"`
	WebServerPort     string `env:"HTTP_PORT,required"`
	GRPCServerPort    string `env:"GRPC_PORT,required"`
	GraphQLServerPort string `env:"GRAPHQL_PORT,required"`
}

func LoadConfig() (*conf, error) {
	var cfg conf
	if err := envdecode.Decode(&cfg); err != nil {
		return nil, errors.New("error to decode config: " + err.Error())
	}
	return &cfg, nil
}
