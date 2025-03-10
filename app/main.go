package main

import (
    "fmt"
    "log"
    "net"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "github.com/RePrete/entity-status-api/app/internal/run"
    _ "github.com/lib/pq"

    "google.golang.org/grpc"

    "github.com/RePrete/entity-status-api/protos"
)

func main() {
    //ctx := context.Background()
    port := os.Getenv("PORT")
    if len(port) == 0 {
        port = "8080"
    }

    db := CreatePostgresConnection()
    db.AutoMigrate(&run.RunRecord{})
    db.AutoMigrate(&run.AggregatePostgresRepository{})

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
                run.NewRunPostgresRepository(db),
                run.NewAggregatePostgresRepository(db),
            ),
        ),
    )
    grpcServer.Serve(lis)
}

func CreatePostgresConnection() *gorm.DB {
    //host := os.Getenv("POSTGRES_HOST")
    //port := os.Getenv("POSTGRES_PORT")
    //user := os.Getenv("POSTGRES_USER")
    //password := os.Getenv("POSTGRES_PASSWORD")
    //dbname := os.Getenv("POSTGRES_DB")
    //
    //dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    dsn := "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=entitystatus sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }

    migrator := db.Migrator()

    // Check if table exists
    if !migrator.HasTable(&run.RunRecord{}) {
        err := migrator.CreateTable(&run.RunRecord{})
        if err != nil {
            panic("failed to create table")
        }
    }
    if !migrator.HasTable(&run.AggregateEntityProcessStatus{}) {
        err := migrator.CreateTable(&run.AggregateEntityProcessStatus{})
        if err != nil {
            panic("failed to create table")
        }
    }

    return db
}
