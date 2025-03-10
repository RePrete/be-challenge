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
}

type AggregateRepository interface {
    GetCurrentStatus(ctx context.Context, paths []string) ([]*AggregateEntityProcessStatus, error)
    Upsert(ctx context.Context, processId, path string, status int, at time.Time) error
}

func NewEntityStatusService(run RunRepository, entity AggregateRepository) *EntityStatusService {
    return &EntityStatusService{
        runRepository:    run,
        entityRepository: entity,
    }
}

type EntityStatusService struct {
    runRepository    RunRepository
    entityRepository AggregateRepository
}

func (e *EntityStatusService) BatchGetEntityStatus(ctx context.Context, paths []string) (map[string]StatusModel, error) {
    entities, err := e.entityRepository.GetCurrentStatus(ctx, paths)
    if err != nil {
        return nil, err
    }

    result := map[string]StatusModel{}
    for _, record := range entities {
        result[record.Path] = StatusModel{
            Path:       record.Path,
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
        e.AggregateEntityProcessStatus(ctx, run.ProcessId, path, run.Status, run.At)
    }
    return nil
}

func (e *EntityStatusService) AggregateEntityProcessStatus(ctx context.Context, process, path string, status int, at time.Time) error {
    return e.entityRepository.Upsert(ctx, process, path, status, at)
}
