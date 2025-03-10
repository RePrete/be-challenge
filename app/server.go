package main

import (
    "context"

    "github.com/getsynq/entity-status-api/app/internal/run"
    "github.com/getsynq/entity-status-api/protos"
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
    err := e.service.InsertRun(ctx, &run.RunModel{})
    if err != nil {
        return nil, err
    }
    return &protos.InsertRunResponse{}, nil
}

func (e EntityStatusApi) BatchGetEntityStatus(ctx context.Context, request *protos.BatchGetEntityStatusRequest) (*protos.BatchGetEntityStatusResponse, error) {
    //TODO implement me
    panic("implement me")
}

func (e EntityStatusApi) GetEntityStatusSummary(ctx context.Context, request *protos.GetEntityStatusSummaryRequest) (*protos.GetEntityStatusSummaryResponse, error) {
    //TODO implement me
    panic("implement me")
}
