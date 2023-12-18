package authz

import (
	"context"
	"fmt"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

const Schema = `
definition operator/user {}

definition operator/application {
	relation reader: operator/user 
	relation writer: operator/user

	permission read = reader + writer
	permission write = writer
}
`

type Client interface {
	CheckPermission(ctx context.Context, param CheckPermission) (bool, error)
	SaveRelationship(ctx context.Context, relationships Relationship) (string, error)
	DeleteRelationship(ctx context.Context, relationships Relationship) (string, error)
	ApplySchema(schema string) error
}

type client struct {
	authzedClient *authzed.Client
}

type CheckPermission struct {
	Resource   Resource
	Permission string
	Subject    Subject
}

type Relationship struct {
	Resource Resource
	Relation string
	Subject  Subject
}

type Resource struct {
	Type string
	Id   string
}

type Subject struct {
	Type string
	Id   string
}

func (a *client) CheckPermission(ctx context.Context, param CheckPermission) (bool, error) {

	contextMap := map[string]interface{}{
		"day_of_week": "monday",
	}

	contextStruct, err := structpb.NewStruct(contextMap)
	if err != nil {
		println(err.Error())
		return false, err
	}

	resource := &pb.ObjectReference{
		ObjectType: param.Resource.Type,
		ObjectId:   param.Resource.Id,
	}

	subject := &pb.SubjectReference{Object: &pb.ObjectReference{
		ObjectType: param.Subject.Type,
		ObjectId:   param.Subject.Id,
	}}

	request := &pb.CheckPermissionRequest{
		Resource:   resource,
		Permission: param.Permission,
		Subject:    subject,
		Context:    contextStruct,
	}

	println(request.Permission, param.Permission)

	fmt.Printf("%+v\n", request)

	resp, err := a.authzedClient.CheckPermission(ctx, request)
	if err != nil {
		println(err.Error())
		return false, err
	}
	response := resp.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION
	return response, nil
}

func (a *client) ApplySchema(schema string) error {
	request := &pb.WriteSchemaRequest{Schema: schema}
	_, err := a.authzedClient.WriteSchema(context.Background(), request)
	if err != nil {
		return err
	}
	return nil
}

func (a *client) writeRelationship(ctx context.Context, relationship Relationship, operation pb.RelationshipUpdate_Operation) (string, error) {
	request := &pb.WriteRelationshipsRequest{Updates: []*pb.RelationshipUpdate{
		{
			Operation: operation,
			Relationship: &pb.Relationship{
				Resource: &pb.ObjectReference{
					ObjectType: relationship.Resource.Type,
					ObjectId:   relationship.Resource.Id,
				},
				Relation: relationship.Relation,
				Subject: &pb.SubjectReference{
					Object: &pb.ObjectReference{
						ObjectType: relationship.Subject.Type,
						ObjectId:   relationship.Subject.Id,
					},
				},
			},
		},
	}}

	resp, err := a.authzedClient.WriteRelationships(context.Background(), request)
	if err != nil {
		return "", err
	}
	return resp.WrittenAt.Token, err
}

func (a *client) SaveRelationship(ctx context.Context, relationship Relationship) (string, error) {
	return a.writeRelationship(ctx, relationship, pb.RelationshipUpdate_OPERATION_CREATE)
}

func (a *client) DeleteRelationship(ctx context.Context, relationship Relationship) (string, error) {
	return a.writeRelationship(ctx, relationship, pb.RelationshipUpdate_OPERATION_DELETE)
}

func NewAuthZClient(host string, token string) (Client, error) {
	c, err := authzed.NewClient(
		host,
		grpcutil.WithInsecureBearerToken(token),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return &client{}, err
	}

	return &client{
		authzedClient: c,
	}, nil
}
