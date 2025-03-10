package run

import (
    "context"
    "errors"
    "fmt"
    "gorm.io/gorm/clause"
    "time"

    "gorm.io/gorm"

    "google.golang.org/protobuf/types/known/timestamppb"
)

type RunRecord struct {
    gorm.Model
    Path       string    `gorm:primary_key`
    At         time.Time `gorm:primary_key`
    ProcessId  string    `gorm:primary_key`
    Status     int
    IsIndirect bool
    IsDeletion bool
}

type EntityStatusProjection struct {
    gorm.Model
    ID     string
    Status int
    At     time.Time
}

func timestampToRecord(ts *timestamppb.Timestamp) string {
    if ts == nil {
        return "-"
    }
    return ts.AsTime().Format(time.RFC3339)
}

func newRunKey(entity string, at *timestamppb.Timestamp) string {
    return fmt.Sprintf("%s:%s", entity, timestampToRecord(at))
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

func (r RunPostgresRepository) GetRunsWorstStatusGroupedByPath(ctx context.Context, id string, since *time.Time, includeDeleted bool) (*RunRecord, error) {
    var result RunRecord

    err := r.db.Raw(`
        SELECT t1.*
        FROM run_records t1
        JOIN (SELECT process_id, path, MAX(at) AS latest_at
            FROM run_records
            WHERE path = ?
            GROUP BY process_id, path) t2
        ON t1.process_id = t2.process_id
        AND t1.path = t2.path
        AND t1.at = t2.latest_at;
    `, id).Scan(&result).Error

    if err != nil {
        return nil, err
    }

    return &result, err
}

func NewEntityPostgresRepository(db *gorm.DB) *EntityPostgresRepository {
    return &EntityPostgresRepository{db: db}
}

type EntityPostgresRepository struct {
    db *gorm.DB
}

func (e EntityPostgresRepository) Upsert(ctx context.Context, entity *EntityStatusProjection) error {
    return e.db.Clauses(clause.OnConflict{
        Columns: []clause.Column{
            {Name: "id"},
        },
        UpdateAll: true,
    }).Create(&entity).Error
}

func (e EntityPostgresRepository) GetEntityStatus(ctx context.Context, path string) (*EntityStatusProjection, error) {
    var entity EntityStatusProjection
    err := e.db.Where("id = ?", path).First(&entity).Error
    switch {
    case errors.Is(err, gorm.ErrRecordNotFound):
        return nil, nil
    case err != nil:
        return nil, err
    }
    return &entity, nil
}

func (e EntityPostgresRepository) GetEntitiesStatus(ctx context.Context, paths []string) ([]*EntityStatusProjection, error) {
    var result []*EntityStatusProjection
    err := e.db.Where("id IN (?)", paths).Find(&result).Error
    return result, err
}
