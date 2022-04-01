package env

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Middlewares struct with middlewares values
type Middlewares struct {
	ApprovedRateLimitTime    uint   `envconfig:"APPROVED_RATE_LIMIT_TIME" default:"120"`
	ApprovedRateLimitRequest int64  `envconfig:"APPROVED_RATE_LIMIT_REQUEST" default:"2"`
	DeclineRateLimitTime     uint   `envconfig:"DECLINE_RATE_LIMIT_TIME" default:"30"`
	DeclineRateLimitRequest  int64  `envconfig:"DECLINE_RATE_LIMIT_REQUEST" default:"1"`
	DeclineRetriesAllowed    uint   `envconfig:"DECLINE_RETRIES_ALLOWED" default:"3"`
	DeclineRetriesMessage    string `envconfig:"DECLINE_RETRIES_MESSAGE" default:"A sales agent will contact you"`
}

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
	Server      *Server
	Ratio       *Ratios
	Middlewares *Middlewares
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

// RetrieveEnvVariables retrieve the env variables
func RetrieveEnvVariables() *Environment {
	conf := new(Environment)
	envconfig.Process("", conf)
	return conf
}
