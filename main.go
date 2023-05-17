package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/konjure-exampels/go-aws/internal/user"
)

const (
	tableNameEnv  = "DYNAMODB_TABLE"
	streamNameEnv = "KINESIS_STREAM"
)

func main() {
	ctx := context.Background()

	tableName := os.Getenv(tableNameEnv)
	streamName := os.Getenv(streamNameEnv)
	port := os.Getenv("PORT")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	repository := user.NewRepository(cfg, tableName)
	stream := user.NewEventStream(cfg, streamName)

	handler := user.NewHandler(repository, stream)

	err = http.ListenAndServe(port, handler)
	if err != nil {
		log.Fatal(err)
	}
}
