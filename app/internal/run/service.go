package run

import (
    "context"
    "slices"
    "time"
)

type RunModel struct {
    ProcessId     string
    DirectPath    string
    IndirectPaths []string
    Status        int
    At            time.Time
    IsDeletion    bool
}

type StatusModel struct {
    Path       string `gorm:primary_key`
    Status     int
    IsDeletion bool
    At         time.Time
}

type RunRepository interface {
    CreateRuns(ctx context.Context, records []RunRecord) error
    GetRunsWorstStatusGroupedByPath(ctx context.Context, id string, since *time.Time, includeDeleted bool) (*RunRecord, error)
}

type EntityRepository interface {
    GetEntityStatus(ctx context.Context, path string) (*EntityStatusProjection, error)
    GetEntitiesStatus(ctx context.Context, paths []string) ([]*EntityStatusProjection, error)
    Upsert(ctx context.Context, entity *EntityStatusProjection) error
}

func NewEntityStatusService(run RunRepository, entity EntityRepository) *EntityStatusService {
    return &EntityStatusService{
        runRepository:    run,
        entityRepository: entity,
    }
}

type EntityStatusService struct {
    runRepository    RunRepository
    entityRepository EntityRepository
}

func (e *EntityStatusService) BatchGetEntityStatus(ctx context.Context, paths []string) (map[string]StatusModel, error) {
    entities, err := e.entityRepository.GetEntitiesStatus(ctx, paths)
    if err != nil {
        return nil, err
    }

    result := map[string]StatusModel{}
    for _, record := range entities {
        result[record.ID] = StatusModel{
            Path:       record.ID,
            Status:     record.Status,
            IsDeletion: false,
            At:         record.At,
        }
    }
    return result, nil
}

func (e *EntityStatusService) InsertRun(ctx context.Context, run *RunModel) error {
    records := []RunRecord{}
    records = append(records, RunRecord{
        ProcessId:  run.ProcessId,
        Path:       run.DirectPath,
        Status:     run.Status,
        At:         run.At,
        IsIndirect: false,
        IsDeletion: run.IsDeletion,
    })
    for _, indirectPath := range run.IndirectPaths {
        records = append(records, RunRecord{
            ProcessId:  run.ProcessId,
            Path:       indirectPath,
            Status:     run.Status,
            At:         run.At,
            IsIndirect: true,
            IsDeletion: run.IsDeletion,
        })
    }
    err := e.runRepository.CreateRuns(ctx, records)
    if err != nil {
        return err
    }

    affectedPaths := slices.Concat([]string{run.DirectPath}, run.IndirectPaths)
    for _, path := range affectedPaths {
        e.ProjectEntityStatus(ctx, path)
    }
    return nil
}

func (e *EntityStatusService) ProjectEntityStatus(ctx context.Context, path string) {
    entity, err := e.entityRepository.GetEntityStatus(nil, path)
    if err != nil {
        // something has to be done here
        return
    }

    var at *time.Time
    if entity != nil {
        at = &entity.At
    } else {
        entity = &EntityStatusProjection{
            ID:     path,
            Status: 0,
        }
    }

    runs, err := e.runRepository.GetRunsWorstStatusGroupedByPath(ctx, path, at, false)
    if err != nil || runs == nil {
        return
    }

    entity.Status = runs.Status
    entity.At = runs.At

    err = e.entityRepository.Upsert(ctx, entity)
    if err != nil {
        return
    }
}
