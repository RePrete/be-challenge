package main

import (
    "context"
    "fmt"
    "google.golang.org/protobuf/types/known/timestamppb"

    "github.com/RePrete/entity-status-api/app/internal/run"
    "github.com/RePrete/entity-status-api/protos"
)

func NewEntityStatusApi(service *run.EntityStatusService) *EntityStatusApi {
    return &EntityStatusApi{
        service: service,
    }
}

type EntityStatusApi struct {
    protos.UnimplementedEntityStatusServiceServer
    service *run.EntityStatusService
}

func (e EntityStatusApi) InsertRun(ctx context.Context, request *protos.InsertRunRequest) (*protos.InsertRunResponse, error) {
    if request == nil || request.Run == nil {
        return nil, fmt.Errorf("request is nil")
    }

    err := e.service.InsertRun(ctx, &run.RunModel{
        ProcessId:     request.Run.ProcessId,
        DirectPath:    request.Run.DirectPath,
        IndirectPaths: request.Run.IndirectPaths,
        Status:        int(request.Run.Status),
        At:            request.Run.At.AsTime(),
        IsDeletion:    request.Run.IsDeletion,
    })
    if err != nil {
        return nil, err
    }
    return &protos.InsertRunResponse{}, nil
}

func (e EntityStatusApi) BatchGetEntityStatus(ctx context.Context, request *protos.BatchGetEntityStatusRequest) (*protos.BatchGetEntityStatusResponse, error) {
    if request == nil {
        return nil, fmt.Errorf("request is nil")
    }

    entityStatuses, err := e.service.BatchGetEntityStatus(ctx, request.Paths)
    if err != nil {
        return &protos.BatchGetEntityStatusResponse{}, err
    }

    result := map[string]*protos.EntityStatus{}
    for entity, status := range entityStatuses {
        result[entity] = &protos.EntityStatus{
            Path:      status.Path,
            IsDeleted: status.IsDeletion,
            Status:    newPbStatusFromModel(status.Status),
            LastAt:    timestamppb.New(status.At),
        }
    }

    return &protos.BatchGetEntityStatusResponse{
        Statuses: result,
    }, nil
}

func newPbStatusFromModel(status int) protos.Status {
    switch status {
    case 1:
        return protos.Status_STATUS_OK
    case 2:
        return protos.Status_STATUS_WARN
    case 3:
        return protos.Status_STATUS_ERROR
    case 4:
        return protos.Status_STATUS_FATAL
    default:
        return protos.Status_STATUS_UNSPECIFIED
    }
}

func (e EntityStatusApi) GetEntityStatusSummary(ctx context.Context, request *protos.GetEntityStatusSummaryRequest) (*protos.GetEntityStatusSummaryResponse, error) {
    //TODO implement me
    panic("implement me")
}
