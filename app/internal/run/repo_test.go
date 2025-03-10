package run

import (
    "context"
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "testing"
)

func TestEntityPostgresRepository_Upsert(t *testing.T) {
    ctx := context.Background()
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        "127.0.0.1", "5432", "postgres", "postgres", "entitystatus")
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    sut := EntityStatusService{
        runRepository:    RunPostgresRepository{db: db},
        entityRepository: AggregatePostgresRepository{db: db},
    }

    sut.ProjectEntityStatus(ctx, "m1")
    //sut := AggregatePostgresRepository{
    //    db: db,
    //}
    //migrator := db.Migrator()
    //
    //// Check if table exists
    //if !migrator.HasTable(&EntityStatusProjection{}) {
    //    // Create table
    //    err := migrator.CreateTable(&EntityStatusProjection{})
    //    if err != nil {
    //        panic("failed to create table")
    //    }
    //}
    //if !migrator.HasTable(&RunRecord{}) {
    //    err := migrator.CreateTable(&RunRecord{})
    //    if err != nil {
    //        panic("failed to create table")
    //    }
    //}
}
