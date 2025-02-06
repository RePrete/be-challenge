# Entity Status API

Hi there, we're excited to see you here. Thank you for taking the time to apply to SYNQ!

The objective of this challenge is to understand your philosophy and approach to solving problems. We're not looking for you to write a lot of code from scratch. We wanted to make this two-way evaluation, so we've extracted a real challenge from our codebase to give you flavour of type of problems we are solving in our day-to-day. Hopefully what follows is something you're excited to solve.

## **Problem Statement**

At **SYNQ** we gather the statuses of various entities from several sources. These statuses are emitted as `Runs` from multiple processes. The runs are put into pub/sub-queues and are consumed by asynchronous processors. They arrive in possibly out-of-order `at`.

```protobuf
enum Status {
	STATUS_OK = 1;
	STATUS_WARN = 2;
	STATUS_ERROR = 3;
	STATUS_FATAL = 4;
}

message Run {
	string process_id = 1;
	string direct_path = 2; // the paths of the entities directly affected by this run
	repeated string indirect_paths = 3; // the paths of the entities indirectly affected by this run
	Status status = 4;  // status of this run
	google.protobuf.Timestamp at = 5; // timestamp of when this run occured
	bool is_deletion = 6; // boolean flag to specify deletion
}
```

Each `Entity` object represents a unique entity in the customer ecosystem. This could be a `dbt Test` or the `database table` that it acts on, an `anomaly detection monitor`, or the `data segment` that it monitors. It is always uniquely identified by its `path`.

Processes (like `dbt tests` or `anomaly detection`) occur at irregular intervals and emit a `Run` object with their result. The `Run` object carries the paths of the entities that are affected by it. There are two ways a `Run` affects an entity: `directly` or `indirectly`.

For example, a test failure directly affects the test and indirectly affects the table that is being tested.

The current status of an entity is the worst of the latest available statuses (direct and indirect). Deleted entities do not contribute towards the status.

Implement the following API for entity status.

```protobuf
service EntityStatusService {
	// Insert the run and consume the status update.
	rpc InsertRun(InsertRunRequest)
		returns (InsertRunResponse) {}
	// Get the statuses of entities by their paths.
	rpc BatchGetEntityStatus(BatchGetEntityStatusRequest)
		returns (BatchGetEntityStatusResponse) {}
	// Get a summary of the status of entities.
	rpc GetEntityStatusSummary(GetEntityStatusSummaryRequest)
		returns (GetEntityStatusSummaryResponse) {}
}

message InsertRunRequest {
	Run run = 1;
}

message InsertRunResponse {}

message EntityStatus {
	string path = 1;
	bool is_deleted = 2;
	Status status = 3;
  google.protobuf.Timestamp last_at = 4;
}

message BatchGetEntityStatusRequest {
	repeated string paths = 1;
}

message BatchGetEntityStatusResponse {
	map<string, EntityStatus> statuses = 1;
}

message GetEntityStatusSummaryRequest {
	repeated string paths = 1; // empty list == all
}

message SummaryItem {
	Status status = 1;
	int32 count = 2;
}

message GetEntityStatusSummaryResponse {
	repeated SummaryItem summary = 1;
}
```

### **Bonus**

> 💡Attempt only if you have completed the basic problem.

Implement a historical lookup for a daily status summary.

```protobuf
service EntityStatusService {
	...
	rpc GetEntityStatusSummaryHistory(GetEntityStatusSummaryHistoryRequest)
	  returns (GetEntityStatusSummaryHistoryResponse) {}
}

message SummaryItem {
	Status status = 1;
	int32 count = 2;
}

message GetDailyEntityStatusSummaryRequest {
	repeated string paths = 1;  // empty list == all
  optional google.protobuf.Timestamp from = 2;  // missing == unbound
  optional google.protobuf.Timestamp to = 3;  // missing == unbound
}

message GetDailyEntityStatusSummaryResponse {
	message DailyItem {
		google.protobuf.Timestamp date = 1;
		repeated SummaryItem summary = 2;
	}
	repeated DailyItem daily = 1;
}
```

## **Instructions**

- Your code should be executable as a `gRPC` server.
- We prefer `Golang`, but you can code in any strongly typed language.
- You can choose a database of your choice. Please include instructions on how to set up the DB (preferably in a Docker container).
- Please ensure that you run the test suite provided.
- We’d appreciate a well-commented code, but that’s secondary to a functional and complete code.
- Feel free to use any tool you have available (web search, AI agents, friend of a friend, etc.).
- The task should not take more than 2-3 hours of coding. If you are done sooner, try out the [Bonus](https://www.notion.so/Entity-Status-API-18ff7d6176b3800ea671d901b149d367?pvs=21) section.
- Please make notes (within the code base or otherwise) which could help you with the follow-up discussion.

## Running the project

### Docker

```
$ docker build . -t entity-status-api:latest
$ docker run --rm --detach -p 8080:8080 entity-status-api:latest
$ cd app && go test -v
```
