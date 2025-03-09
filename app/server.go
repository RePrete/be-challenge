package main

import (
    "context"
    "github.com/getsynq/entity-status-api/protos"
    "google.golang.org/grpc"
)

type EntityStatusApi struct {
    protos.UnimplementedEntityStatusServiceServer
}

func (e EntityStatusApi) InsertRun(ctx context.Context, in *protos.InsertRunRequest, opts ...grpc.CallOption) (*protos.InsertRunResponse, error) {
    //TODO implement me
    panic("implement me")
}

func (e EntityStatusApi) BatchGetEntityStatus(ctx context.Context, in *protos.BatchGetEntityStatusRequest, opts ...grpc.CallOption) (*protos.BatchGetEntityStatusResponse, error) {
    //TODO implement me
    panic("implement me")
}

func (e EntityStatusApi) GetEntityStatusSummary(ctx context.Context, in *protos.GetEntityStatusSummaryRequest, opts ...grpc.CallOption) (*protos.GetEntityStatusSummaryResponse, error) {
    //TODO implement me
    panic("implement me")
}
