package main

import (
	"context"
	"github.com/authz-spicedb/internal/app"
)

func main() {

	application := app.NewApplication()
	application.Run(context.Background())

}
