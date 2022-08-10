package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"localhost:8090"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8080"`
}

func NewConfig() (*Config, error) {
	conf := &Config{}

	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	regStringVar(&conf.RunAddress, "a", conf.RunAddress, "gophermart run address")
	regStringVar(&conf.DatabaseURI, "d", conf.DatabaseURI, "database connection")
	regStringVar(&conf.AccrualSystemAddress, "r", conf.AccrualSystemAddress, "broker address")
	flag.Parse()

	return conf, nil
}

func regStringVar(p *string, name string, value string, usage string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(p, name, value, usage)
	}
}
