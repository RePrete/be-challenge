package run

import (
    "context"
    "gorm.io/gorm/clause"
    "time"

    "gorm.io/gorm"

    "google.golang.org/protobuf/types/known/timestamppb"
)

// RunRecord represent the status of a path, as reported from a process, at a specific timestamp
type RunRecord struct {
    Path       string    `gorm:primary_key`
    At         time.Time `gorm:primary_key`
    ProcessId  string    `gorm:primary_key`
    Status     int
    IsIndirect bool
    IsDeletion bool
}

func timestampToRecord(ts *timestamppb.Timestamp) string {
    if ts == nil {
        return "-"
    }
    return ts.AsTime().Format(time.RFC3339)
}

func NewRunPostgresRepository(db *gorm.DB) *RunPostgresRepository {
    return &RunPostgresRepository{db: db}
}

type RunPostgresRepository struct {
    db *gorm.DB
}

func (r RunPostgresRepository) CreateRuns(ctx context.Context, records []RunRecord) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        for _, record := range records {
            err := tx.Create(&record).Error
            if err != nil {
                return err
            }
        }
        return nil
    })
}

func NewAggregatePostgresRepository(db *gorm.DB) *AggregatePostgresRepository {
    return &AggregatePostgresRepository{db: db}
}

type AggregatePostgresRepository struct {
    db *gorm.DB
}

func (r AggregatePostgresRepository) GetCurrentStatus(ctx context.Context, paths []string) ([]*AggregateEntityProcessStatus, error) {
    result := []*AggregateEntityProcessStatus{}

    err := r.db.Select("path, MAX(status) as status").
        Table("aggregate_entity_process_statuses").
        Where("path IN (?)", paths).
        Group("path").
        Scan(&result).Error
    return result, err
}

// AggregateEntityProcessStatus represent the last path status reported by a process
type AggregateEntityProcessStatus struct {
    ProcessId string `gorm:"index:aggregate_id,unique"`
    Path      string `gorm:"index:aggregate_id,unique"`
    Status    int
    At        time.Time
}

func (r AggregatePostgresRepository) Upsert(ctx context.Context, aggregate AggregateEntityProcessStatus) error {
    return r.db.Clauses(clause.OnConflict{
        Columns: []clause.Column{
            {Name: "process_id"},
            {Name: "path"},
        },
        UpdateAll: true,
    }).Create(aggregate).Error
}
