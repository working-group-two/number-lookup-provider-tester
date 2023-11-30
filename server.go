package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
	common "github.com/working-group-two/wgtwoapis/wgtwo/common/v0"
	pb "github.com/working-group-two/wgtwoapis/wgtwo/lookup/v0"
	"google.golang.org/grpc"
)

type printOptions struct {
	Requests  bool
	Responses bool
	Progress  bool
}

type grpcServer struct {
	pb.UnimplementedNumberLookupServiceServer
	rps          int64
	phoneNumbers []string
	print        *printOptions
}

func (s *grpcServer) NumberLookup(stream pb.NumberLookupService_NumberLookupServer) error {
	ctx := stream.Context()

	// Read messages from the stream:
	responses := make(chan *pb.NumberLookupResponse)
	defer close(responses)
	go func() {
		for {
			response, err := stream.Recv()
			if err == nil {
				responses <- response
			}
		}
	}()

	var (
		inFlight int
		rps      int
	)

	// Start pumping messages to the stream:
	durationPerRequest := time.Second / time.Duration(s.rps)

	sendTicker := time.NewTicker(durationPerRequest)
	defer sendTicker.Stop()

	rpsTicker := time.NewTicker(time.Second)
	defer rpsTicker.Stop()

	for {
		select {
		// stream closed, return
		case <-ctx.Done():
			return ctx.Err()

		// print responses as they arrive
		case response := <-responses:
			inFlight--
			s.printResponse(response)

		// log the current requests/s every second
		case <-rpsTicker.C:
			if s.print.Progress {
				slog.Info("Progress", "rps", rps)
			}
			rps = 0

		// send messages in regular intervals
		case <-sendTicker.C:
			req := &pb.NumberLookupRequest{
				Id: uuid.NewString(),
				Number: &common.PhoneNumber{
					E164: pickRandom(s.phoneNumbers),
				},
			}
			err := stream.Send(req)
			if err != nil {
				slog.Error("Failed to send request", slog.String("error", err.Error()))
				continue
			}
			rps++
			inFlight++
			if s.print.Requests {
				slog.Info(" Request",
					"id", req.Id,
					"number", req.Number.E164,
					"inFlight", inFlight,
				)
			}
		}
	}
}

func (s *grpcServer) printResponse(response *pb.NumberLookupResponse) {
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

func pickRandom(list []string) string {
	return list[rand.Intn(len(list))]
}

func startServer(ctx context.Context, listener string, rps int64, phoneNumbers []string, printOptions *printOptions) {
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
	go func() {
		<-ctx.Done()
		lis.Close()
		s.Stop()
	}()
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve", slog.String("error", err.Error()))
	}
}
