package main

import (
    "context"
    "flag"
    "google.golang.org/protobuf/types/known/timestamppb"
    "os"
    "slices"
    "testing"
    "time"

    "github.com/RePrete/entity-status-api/protos"
    "github.com/stretchr/testify/suite"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func TestEntityStatusApi(t *testing.T) {
    suite.Run(t, new(TestEntityStatusApiTestSuite))
}

type TestEntityStatusApiTestSuite struct {
    suite.Suite

    conn   *grpc.ClientConn
    client protos.EntityStatusServiceClient
}

func (s *TestEntityStatusApiTestSuite) SetupSuite() {
    serverAddr := flag.String("server_addr", "localhost:8080", "The server address in the format of host:port")
    flag.Parse()

    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    }

    var err error
    s.conn, err = grpc.NewClient(*serverAddr, opts...)
    if err != nil {
        panic(err)
    }

    s.client = protos.NewEntityStatusServiceClient(s.conn)
}

func (s *TestEntityStatusApiTestSuite) TearDownSuite() {
    if s.conn != nil {
        s.conn.Close()
    }
}

type Runner struct {
    processOneId   string
    processTwoId   string
    processThreeId string

    monitorOneId string
    monitorTwoId string

    tableId string

    t0 time.Time
}

func NewRunner() *Runner {
    return &Runner{
        processOneId:   "p1",
        processTwoId:   "p2",
        processThreeId: "p3",

        //monitorOneId: uuid.NewString(),
        //monitorTwoId: uuid.NewString(),
        //
        //tableId: uuid.NewString(),
        monitorOneId: "m1",
        monitorTwoId: "m2",

        tableId: "t1",

        t0: time.Now().Add(-time.Hour * 48),
    }
}

// There is a database table that is monitored by two processes.
// There are two processes monitoring a database table.
// Process 1 runs every 6 hours.
// Process 2 runs every 12 hours.
// Another process runs every 24 hours to check database table connectivity.
// 00h : process 1 - OK, process 2 - OK, process 3 - OK
// 06h : process 1 - FATAL
// 12h : process 1 - FATAL, process 2 - OK
// 18h : process 1 - ERROR
// 24h : process 1 - ERROR, process 2 - OK, process 3 - WARN
// 25h : process 1 - deleted
// 36h : process 2 - OK
// 48h : process 2 - ERROR, process 3 - OK
func (s *TestEntityStatusApiTestSuite) insertRuns(runner *Runner, hours []int) {
    ctx := context.Background()

    if slices.Contains(hours, 0) {
        // 00h : process 1 - OK, process 2 - OK, process 3 - OK
        at := runner.t0
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processOneId,
                    DirectPath:    runner.monitorOneId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_OK,
                    At:            timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processTwoId,
                    DirectPath:    runner.monitorTwoId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_OK,
                    At:            timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:  runner.processThreeId,
                    DirectPath: runner.tableId,
                    Status:     protos.Status_STATUS_OK,
                    At:         timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
    }

    if slices.Contains(hours, 6) {
        // 06h : process 1 - FATAL
        at := runner.t0.Add(time.Hour * 6)
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processOneId,
                    DirectPath:    runner.monitorOneId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_FATAL,
                    At:            timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
    }

    if slices.Contains(hours, 12) {
        // 12h : process 1 - FATAL, process 2 - OK
        at := runner.t0.Add(time.Hour * 12)
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processOneId,
                    DirectPath:    runner.monitorOneId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_FATAL,
                    At:            timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processTwoId,
                    DirectPath:    runner.monitorTwoId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_OK,
                    At:            timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
    }

    if slices.Contains(hours, 18) {
        // 18h : process 1 - ERROR
        at := runner.t0.Add(time.Hour * 18)
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processOneId,
                    DirectPath:    runner.monitorOneId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_ERROR,
                    At:            timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
    }

    if slices.Contains(hours, 24) {
        // 24h : process 1 - ERROR, process 2 - OK
        at := runner.t0.Add(time.Hour * 24)
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processOneId,
                    DirectPath:    runner.monitorOneId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_ERROR,
                    At:            timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processTwoId,
                    DirectPath:    runner.monitorTwoId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_OK,
                    At:            timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:  runner.processThreeId,
                    DirectPath: runner.tableId,
                    Status:     protos.Status_STATUS_WARN,
                    At:         timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
    }

    if slices.Contains(hours, 25) {
        // 25h : process 1 - deleted
        at := runner.t0.Add(time.Hour * 25)
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processOneId,
                    DirectPath:    runner.monitorOneId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_UNSPECIFIED,
                    At:            timestamppb.New(at),
                    IsDeletion:    true,
                },
            })
            s.Require().NoError(err)
        }
    }

    if slices.Contains(hours, 36) {
        // 36h : process 2 - OK
        at := runner.t0.Add(time.Hour * 36)
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processTwoId,
                    DirectPath:    runner.monitorTwoId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_OK,
                    At:            timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
    }

    if slices.Contains(hours, 48) {
        // 48h : process 2 - ERROR
        at := runner.t0.Add(time.Hour * 48)
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:     runner.processTwoId,
                    DirectPath:    runner.monitorTwoId,
                    IndirectPaths: []string{runner.tableId},
                    Status:        protos.Status_STATUS_ERROR,
                    At:            timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
        {
            _, err := s.client.InsertRun(ctx, &protos.InsertRunRequest{
                Run: &protos.Run{
                    ProcessId:  runner.processThreeId,
                    DirectPath: runner.tableId,
                    Status:     protos.Status_STATUS_OK,
                    At:         timestamppb.New(at),
                },
            })
            s.Require().NoError(err)
        }
    }
}

func (s *TestEntityStatusApiTestSuite) TestDirectStatus() {
    runner := NewRunner()

    // insert statuses for first 12 hours
    {
        s.insertRuns(runner, []int{0, 6, 12})

        resp, err := s.client.BatchGetEntityStatus(context.Background(), &protos.BatchGetEntityStatusRequest{
            Paths: []string{runner.monitorOneId, runner.monitorTwoId},
        })
        s.Require().NoError(err)

        s.Require().NotNil(resp.Statuses[runner.monitorOneId])
        s.Require().Equal(protos.Status_STATUS_FATAL, resp.Statuses[runner.monitorOneId].Status)
        s.Require().Equal(false, resp.Statuses[runner.monitorOneId].IsDeleted)

        s.Require().NotNil(resp.Statuses[runner.monitorTwoId])
        s.Require().Equal(protos.Status_STATUS_OK, resp.Statuses[runner.monitorTwoId].Status)
        s.Require().Equal(false, resp.Statuses[runner.monitorTwoId].IsDeleted)
        os.Exit(0)
    }

    // insert statuses till 24h
    {
        s.insertRuns(runner, []int{18, 24})
        resp, err := s.client.BatchGetEntityStatus(context.Background(), &protos.BatchGetEntityStatusRequest{
            Paths: []string{runner.monitorOneId, runner.monitorTwoId},
        })
        s.Require().NoError(err)

        s.Require().NotNil(resp.Statuses[runner.monitorOneId])
        s.Require().Equal(protos.Status_STATUS_ERROR, resp.Statuses[runner.monitorOneId].Status)
        s.Require().Equal(false, resp.Statuses[runner.monitorOneId].IsDeleted)

        s.Require().NotNil(resp.Statuses[runner.monitorTwoId])
        s.Require().Equal(protos.Status_STATUS_OK, resp.Statuses[runner.monitorTwoId].Status)
        s.Require().Equal(false, resp.Statuses[runner.monitorTwoId].IsDeleted)
    }

    // insert statuses till 48h
    {
        s.insertRuns(runner, []int{25, 36, 48})
        resp, err := s.client.BatchGetEntityStatus(context.Background(), &protos.BatchGetEntityStatusRequest{
            Paths: []string{runner.monitorOneId, runner.monitorTwoId},
        })
        s.Require().NoError(err)

        s.Require().NotNil(resp.Statuses[runner.monitorOneId])
        s.Require().Equal(protos.Status_STATUS_UNSPECIFIED, resp.Statuses[runner.monitorOneId].Status)
        s.Require().Equal(true, resp.Statuses[runner.monitorOneId].IsDeleted)

        s.Require().NotNil(resp.Statuses[runner.monitorTwoId])
        s.Require().Equal(protos.Status_STATUS_ERROR, resp.Statuses[runner.monitorTwoId].Status)
        s.Require().Equal(false, resp.Statuses[runner.monitorTwoId].IsDeleted)
    }
}

func (s *TestEntityStatusApiTestSuite) TestIndirectStatus() {
    runner := NewRunner()

    // insert statuses for first 12 hours
    {
        s.insertRuns(runner, []int{0, 6, 12})

        resp, err := s.client.BatchGetEntityStatus(context.Background(), &protos.BatchGetEntityStatusRequest{
            Paths: []string{runner.tableId},
        })
        s.Require().NoError(err)

        s.Require().NotNil(resp.Statuses[runner.tableId])
        s.Require().Equal(protos.Status_STATUS_OK, resp.Statuses[runner.tableId].Status)
        s.Require().Equal(false, resp.Statuses[runner.tableId].IsDeleted)
    }

    // insert statuses till 24h
    {
        s.insertRuns(runner, []int{18, 24})
        resp, err := s.client.BatchGetEntityStatus(context.Background(), &protos.BatchGetEntityStatusRequest{
            Paths: []string{runner.tableId},
        })
        s.Require().NoError(err)

        s.Require().NotNil(resp.Statuses[runner.tableId])
        s.Require().Equal(protos.Status_STATUS_ERROR, resp.Statuses[runner.tableId].Status)
        s.Require().Equal(false, resp.Statuses[runner.tableId].IsDeleted)
    }

    // insert statuses till 48h
    {
        s.insertRuns(runner, []int{25, 36, 48})
        resp, err := s.client.BatchGetEntityStatus(context.Background(), &protos.BatchGetEntityStatusRequest{
            Paths: []string{runner.tableId},
        })
        s.Require().NoError(err)

        s.Require().NotNil(resp.Statuses[runner.tableId])
        s.Require().Equal(protos.Status_STATUS_ERROR, resp.Statuses[runner.tableId].Status)
        s.Require().Equal(false, resp.Statuses[runner.tableId].IsDeleted)
    }
}

func (s *TestEntityStatusApiTestSuite) TestSummary() {
    runner := NewRunner()

    // insert statuses at 00h
    {
        s.insertRuns(runner, []int{0})

        resp, err := s.client.GetEntityStatusSummary(context.Background(), &protos.GetEntityStatusSummaryRequest{})
        s.Require().NoError(err)

        s.Require().Len(resp.Summary, 1)
        s.Require().Equal(protos.Status_STATUS_OK, resp.Summary[0].Status)
        s.Require().Equal(int32(3), resp.Summary[0].Count)
    }

    // insert statuses till 12h
    {
        s.insertRuns(runner, []int{6, 12})

        resp, err := s.client.GetEntityStatusSummary(context.Background(), &protos.GetEntityStatusSummaryRequest{})
        s.Require().NoError(err)

        s.Require().Len(resp.Summary, 2)
        summaries := map[protos.Status]*protos.SummaryItem{}
        for _, summary := range resp.Summary {
            summaries[summary.Status] = summary
        }
        s.Require().NotNil(summaries[protos.Status_STATUS_FATAL])
        s.Require().Equal(int32(2), summaries[protos.Status_STATUS_FATAL].Count)
        s.Require().NotNil(summaries[protos.Status_STATUS_OK])
        s.Require().Equal(int32(1), summaries[protos.Status_STATUS_OK].Count)
    }

    // insert statuses till 24h
    {
        s.insertRuns(runner, []int{18, 24})

        resp, err := s.client.GetEntityStatusSummary(context.Background(), &protos.GetEntityStatusSummaryRequest{})
        s.Require().NoError(err)

        s.Require().Len(resp.Summary, 2)
        summaries := map[protos.Status]*protos.SummaryItem{}
        for _, summary := range resp.Summary {
            summaries[summary.Status] = summary
        }
        s.Require().NotNil(summaries[protos.Status_STATUS_ERROR])
        s.Require().Equal(int32(2), summaries[protos.Status_STATUS_ERROR].Count)
        s.Require().NotNil(summaries[protos.Status_STATUS_OK])
        s.Require().Equal(int32(1), summaries[protos.Status_STATUS_OK].Count)
    }

    // insert statuses till 48h
    {
        s.insertRuns(runner, []int{25, 36, 48})

        resp, err := s.client.GetEntityStatusSummary(context.Background(), &protos.GetEntityStatusSummaryRequest{})
        s.Require().NoError(err)

        s.Require().Len(resp.Summary, 1)
        s.Require().Equal(protos.Status_STATUS_ERROR, resp.Summary[0].Status)
        s.Require().Equal(int32(2), resp.Summary[0].Count)
    }
}
