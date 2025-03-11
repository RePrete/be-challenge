package run

import (
    "context"
    "fmt"
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

// StatusCount is a DTO returned from repo layer, not a db record itself
type StatusCount struct {
    Status int
    Count  int
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

func (r AggregatePostgresRepository) GetEntityStatusSummary(ctx context.Context, paths []string) ([]*StatusCount, error) {
    result := []*StatusCount{}

    sql := `
        SELECT t1.status, COUNT(t1.status) as count
        FROM (
            SELECT path, MAX(status) as status
            FROM aggregate_entity_process_statuses %s
            GROUP BY path
        ) as t1
        GROUP BY t1.status
    `
    args := []interface{}{}
    whereCondition := ""
    if len(paths) > 0 {
        whereCondition = " AND t1.path IN (?)"
        args = append(args, paths)
    }
    sql = fmt.Sprintf(sql, whereCondition)

    err := r.db.Raw(sql, args...).Scan(&result).Error

    if err != nil {
        return result, err
    }
    return result, nil
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
