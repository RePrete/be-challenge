package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/getsynq/entity-status-api/app/internal/run"

    "google.golang.org/grpc"

    "github.com/getsynq/entity-status-api/protos"
)

func main() {
    ctx := context.Background()
    port := os.Getenv("PORT")
    if len(port) == 0 {
        port = "8080"
    }

    dynamo := CreateDynamoDBClient(ctx)

    lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    var opts []grpc.ServerOption
    grpcServer := grpc.NewServer(opts...)
    protos.RegisterEntityStatusServiceServer(
        grpcServer,
        NewEntityStatusApi(
            run.NewEntityStatusService(
                run.NewRunDynamoDB(dynamo),
            ),
        ),
    )
    grpcServer.Serve(lis)
}

func CreateDynamoDBClient(ctx context.Context) *dynamodb.Client {
    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion("us-east-1"),
        config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
            return aws.Endpoint{
                URL: "localhost:8000",
            }, nil
        })),
    )

    if err != nil {
        panic(err)
    }

    return dynamodb.NewFromConfig(cfg)
}
