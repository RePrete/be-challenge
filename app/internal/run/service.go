package run

import (
    "context"
)

type RunModel struct {
}

type RunRepository interface {
    NewRun(ctx context.Context, record RunRecord) error
}

func NewEntityStatusService(repository RunRepository) *EntityStatusService {
    return &EntityStatusService{
        repo: repository,
    }
}

type EntityStatusService struct {
    repo RunRepository
}

func (e *EntityStatusService) InsertRun(ctx context.Context, run *RunModel) error {
    return e.repo.NewRun(ctx, RunRecord{
        ProcessId:     "sada",
        DirectPath:    "asd",
        IndirectPaths: nil,
        Status:        "",
        At:            nil,
        IsDeletion:    false,
    })
}
