package app

import (
	"context"

	"github.com/authz-spicedb/internal/authz"
	"github.com/authz-spicedb/internal/controller"
	"github.com/gin-gonic/gin"
)

type Application interface {
	Run(ctx context.Context) error
}

type application struct {
}

func (a *application) Run(ctx context.Context) error {
	println("Initiating AuthZ...")
	client, err := authz.NewAuthZClient("localhost:50051", "somerandomkeyhere")
	if err != nil {
		return err
	}

	err = client.ApplySchema(authz.Schema)
	if err != nil {
		return err
	}

	authzController := controller.NewAuthzController(client)

	r := gin.Default()
	v1 := r.Group("/authz")
	{
		v1.POST("/relationship", authzController.SaveRelationship)
		v1.DELETE("/relationship", authzController.DeleteRelationship)
		v1.POST("/check-permission", authzController.CheckPermission)
	}

	return r.Run(":7070")
}

func NewApplication() Application {
	return &application{}
}
