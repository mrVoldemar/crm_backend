package main

import (
	"context"
	"log"

	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/app"
)

func main() {

	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}

}
