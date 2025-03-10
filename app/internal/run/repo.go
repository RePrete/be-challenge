package run

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RunRecord struct {
	ProcessId     string                 `json:"process_id,omitempty"`
	DirectPath    string                 `json:"direct_path,omitempty"`
	IndirectPaths []string               `json:"indirect_paths,omitempty"`
	Status        string                 `json:"status,omitempty"`
	At            *timestamppb.Timestamp `json:"at,omitempty"`
	IsDeletion    bool                   `json:"is_deletion,omitempty"`
}

func NewRunDynamoDB(client *dynamodb.Client) *RunDynamoDB {
	return &RunDynamoDB{
		client: client,
	}
}

type RunDynamoDB struct {
	client *dynamodb.Client
}

func (r RunDynamoDB) NewRun(ctx context.Context, record RunRecord) error {
	// Define the table name
	tableName := "RunRecords"

	// Create the item to insert
	item := map[string]types.AttributeValue{
		"process_id":     &types.AttributeValueMemberS{Value: record.ProcessId},
		"direct_path":    &types.AttributeValueMemberS{Value: record.DirectPath},
		"indirect_paths": &types.AttributeValueMemberSS{Value: record.IndirectPaths},
		"status":         &types.AttributeValueMemberS{Value: record.Status},
		"at":             &types.AttributeValueMemberS{Value: record.At.AsTime().Format(time.RFC3339)},
		"is_deletion":    &types.AttributeValueMemberBOOL{Value: record.IsDeletion},
	}

	// Create the PutItem input
	input := &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	}

	// Put the item into the table
	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
