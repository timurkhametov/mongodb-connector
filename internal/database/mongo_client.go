package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongo-db-connector/internal/infrastructure"
	"reflect"
	"strings"
)

func BuildClient(ctx context.Context, config *infrastructure.DBConfig) (*mongo.Client, error) {
	if strings.TrimSpace(config.Username) == "" {
		return nil, fmt.Errorf("username is required")
	}

	if strings.TrimSpace(config.Password) == "" {
		return nil, fmt.Errorf("password is required")
	}

	if strings.TrimSpace(config.Host) == "" {
		return nil, fmt.Errorf("host is required")
	}

	if config.Port <= 0 {
		return nil, fmt.Errorf("port is required")
	}

	// The code below is used to leverage json tags instead of bson structures.
	// https://stackoverflow.com/questions/33643442/can-i-use-json-tags-as-bson-tags-in-mgo
	structcodec, _ := bsoncodec.NewStructCodec(bsoncodec.JSONFallbackStructTagParser)
	rb := bson.NewRegistryBuilder()
	rb.RegisterDefaultEncoder(reflect.Struct, structcodec)
	rb.RegisterDefaultDecoder(reflect.Struct, structcodec)

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d", config.Username, config.Password, config.Host, config.Port)

	opts := options.Client().
		SetRegistry(rb.Build()).
		ApplyURI(connectionString)

	return mongo.Connect(ctx, opts)
}
