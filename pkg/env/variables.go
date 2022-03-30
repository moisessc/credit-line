package env

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Server struct with server values
type Server struct {
	Port            uint16 `envconfig:"SERVER_PORT" default:"3000"`
	ShutdownTimeOut uint16 `envconfig:"SERVER_SHUTDOWN_TIMEOUT" default:"10"`
}

// Ratios struct with ratios values
type Ratios struct {
	CashBalance    float64 `envconfig:"CASH_BALANCE_RATIO" default:"3"`
	MonthlyRevenue float64 `envconfig:"MONTHLY_REVENUE_RATIO" default:"5"`
}

// Environment struct with the environment values
type Environment struct {
	Server *Server
	Ratio  *Ratios
}

// LoadEnvironment loads a .env file and set the environment variables
func LoadEnvironment() *Environment {
	err := godotenv.Load()
	if err != nil {
		log.Println("config file not found")
	}

	conf := new(Environment)
	envconfig.Process("", conf)
	return conf
}
