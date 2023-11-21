package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/working-group-two/number-lookup-provider-tester/internal/counters"
	common "github.com/working-group-two/wgtwoapis/wgtwo/common/v0"
	pb "github.com/working-group-two/wgtwoapis/wgtwo/lookup/v0"
	"google.golang.org/grpc"
	"log/slog"
	"math/rand"
	"net"
	"time"
)

var nanosPerSecond = int64(time.Second / time.Nanosecond)

type PrintOptions struct {
	Requests  bool
	Responses bool
	Progress  bool
}

type grpcServer struct {
	pb.UnimplementedNumberLookupServiceServer
	rps          uint32
	phoneNumbers []string
	print        *PrintOptions
}

func (s *grpcServer) NumberLookup(stream pb.NumberLookupService_NumberLookupServer) error {
	ctx := stream.Context()
	durationPerRequest := time.Duration(nanosPerSecond / int64(s.rps))
	ticker := time.NewTicker(durationPerRequest)

	inFlight := counters.NewInFlightCounter()
	rps := counters.NewRpsCounter()

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				go func() {
					req := &pb.NumberLookupRequest{
						Id: uuid.NewString(),
						Number: &common.PhoneNumber{
							E164: PickRandom(s.phoneNumbers),
						},
					}
					err := stream.Send(req)
					rps.Increase()
					inFlight.Increase()
					if s.print.Progress && rps.GetCounter()%100 == 0 {
						slog.Info("Progress",
							"rps", rps.GetAndReset(),
						)
					}
					if s.print.Requests {
						slog.Info(" Request",
							"id", req.Id,
							"number", req.Number.E164,
							"inFlight", inFlight.Get(),
						)
					}
					if err != nil {
						slog.Error("Failed to send request", slog.String("error", err.Error()))
					}
				}()
			}
		}
	}()

	for {
		response, err := stream.Recv()
		inFlight.Decrease()
		if err != nil {
			return err
		}

		replyTo := response.NumberLookupRequest
		switch r := response.Reply.(type) {
		case *pb.NumberLookupResponse_Result:
			if s.print.Responses {
				slog.Info("Response",
					"id", replyTo.Id,
					"number", replyTo.Number.E164,
					"result", r.Result.Name,
				)
			}
		case *pb.NumberLookupResponse_Error:
			if s.print.Responses {
				slog.Warn(
					"Response",
					"id", replyTo.Id,
					"number", replyTo.Number.E164,
					"error", r.Error.Message,
				)
			}
		default:
			slog.Warn(
				"Response",
				"id", replyTo.Id,
				"number", replyTo.Number.E164,
				"error", "unknown response type",
				"type", fmt.Sprintf("%T", r),
			)
		}
	}
}

func PickRandom(list []string) string {
	return list[rand.Intn(len(list))]
}

func Start(listener string, rps uint32, phoneNumbers []string, printOptions *PrintOptions) {
	lis, err := net.Listen("tcp", listener)
	if err != nil {
		slog.Error("Failed to listen", slog.String("error", err.Error()))
	}
	s := grpc.NewServer()
	pb.RegisterNumberLookupServiceServer(s, &grpcServer{
		rps:          rps,
		phoneNumbers: phoneNumbers,
		print:        printOptions,
	})
	slog.Info("Starting server",
		"listener", lis.Addr(),
		"rps", rps,
		"#phoneNumbers", len(phoneNumbers),
	)
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve", slog.String("error", err.Error()))
	}
}
