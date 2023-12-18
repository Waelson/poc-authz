package controller

import (
	"github.com/authz-spicedb/internal/authz"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthzController interface {
	SaveRelationship(c *gin.Context)
	DeleteRelationship(c *gin.Context)
	CheckPermission(c *gin.Context)
}

type authzController struct {
	authzCliente authz.Client
}

type RequestCheckPermission struct {
	Resource   RequestResource `json:"resource"`
	Permission string          `json:"permission"`
	Subject    RequestSubject  `json:"subject"`
}

type RequestRelationship struct {
	Resource RequestResource `json:"resource"`
	Relation string          `json:"relation"`
	Subject  RequestSubject  `json:"subject"`
}

type RequestResource struct {
	Namespace string `json:"namespace"`
	ID        string `json:"id"`
}

type RequestSubject struct {
	Namespace string `json:"namespace"`
	ID        string `json:"id"`
}

type ResponseRelationship struct {
	Token string `json:"token"`
}

type ResponseCheckPermission struct {
	Authorized bool `json:"authorized"`
}

func (a *authzController) DeleteRelationship(c *gin.Context) {
	input := RequestRelationship{}
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	param := authz.Relationship{
		Resource: authz.Resource{
			Type: input.Resource.Namespace,
			Id:   input.Resource.ID,
		},
		Relation: input.Relation,
		Subject: authz.Subject{
			Type: input.Subject.Namespace,
			Id:   input.Subject.ID,
		},
	}

	token, err := a.authzCliente.DeleteRelationship(c.Request.Context(), param)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusAccepted, &ResponseRelationship{
		Token: token,
	})
}

func (a *authzController) SaveRelationship(c *gin.Context) {
	input := RequestRelationship{}
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	param := authz.Relationship{
		Resource: authz.Resource{
			Type: input.Resource.Namespace,
			Id:   input.Resource.ID,
		},
		Relation: input.Relation,
		Subject: authz.Subject{
			Type: input.Subject.Namespace,
			Id:   input.Subject.ID,
		},
	}

	token, err := a.authzCliente.SaveRelationship(c.Request.Context(), param)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusAccepted, &ResponseRelationship{
		Token: token,
	})
}

func (a *authzController) CheckPermission(c *gin.Context) {
	input := RequestCheckPermission{}
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	param := authz.CheckPermission{
		Resource: authz.Resource{
			Type: input.Resource.Namespace,
			Id:   input.Resource.ID,
		},
		Permission: input.Permission,
		Subject: authz.Subject{
			Type: input.Subject.Namespace,
			Id:   input.Subject.ID,
		},
	}

	//fmt.Printf("%+v\n", param)

	authorized, err := a.authzCliente.CheckPermission(c.Request.Context(), param)
	if err != nil {
		println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusAccepted, &ResponseCheckPermission{
		Authorized: authorized,
	})
}

func NewAuthzController(client authz.Client) AuthzController {
	return &authzController{
		authzCliente: client,
	}
}
