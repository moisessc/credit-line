package bootstrap

import (
	"fmt"

	"credit-line/pkg/env"
)

// Run retrieves the environment, init the database, builds the server router and starts the server
func Run() error {
	env.LoadEnvironment()

	router := newEchoRouter()

	srv := newServer(router)
	err := srv.up()
	if err != nil {
		return fmt.Errorf("failed to init server, %v", err)
	}

	return nil
}
