package main

import (
	"github.com/getsynq/entity-status-api/protos"
)

type EntityStatusApi struct {
	protos.UnimplementedEntityStatusServiceServer
}
