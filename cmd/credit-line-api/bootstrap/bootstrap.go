package bootstrap

import (
	"fmt"

	"credit-line/internal/calculator"
	"credit-line/internal/controller"
	"credit-line/internal/service"
	"credit-line/pkg/env"
)

// Run retrieves the environment, init the database, builds the server router and starts the server
func Run() error {
	conf := env.LoadEnvironment()

	creditLimitCalculator := calculator.NewCreditLine(conf.Ratio)
	creditLimitService := service.NewCreditLine(creditLimitCalculator)
	creditLimitRouter := controller.NewCreditLineHandler(creditLimitService)

	router := newEchoRouter(creditLimitRouter)

	srv := newServer(router, conf.Server)
	err := srv.up()
	if err != nil {
		return fmt.Errorf("failed to init server, %v", err)
	}

	return nil
}
