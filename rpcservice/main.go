//go:generate protoc -I bet/ --go_out=bet/ --go_opt=paths=source_relative bet/bet.proto
//go:generate protoc -I bet/ --go-grpc_out=paths=source_relative:bet/ bet/bet.proto

// Package main implements a grpc server for our slot machine
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	bet "github.com/gadumitrachioaiei/slotserver/proto/gen/go/bet/v1"
	"github.com/gadumitrachioaiei/slotserver/slot"
)

func main() {
	log.SetFlags(log.Llongfile)
	grpcServer := grpc.NewServer(grpc.ConnectionTimeout(30 * time.Second))
	reflection.Register(grpcServer)
	bet.RegisterSlotMachineServiceServer(grpcServer, newSlotMachineServer())
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to open tcp port: %s %v", ":8080", err)
	}
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	defer grpcServer.Stop()
	// also start a REST proxy towards our grpc service
	log.Println("starting proxy")
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := bet.RegisterSlotMachineServiceHandlerFromEndpoint(context.Background(), mux, ":8080", opts); err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatal(err)
	}
}

type slotMachineServer struct {
	*bet.UnimplementedSlotMachineServiceServer
	m *slot.Machine
}

func newSlotMachineServer() slotMachineServer {
	return slotMachineServer{m: slot.NewMachine()}
}

// CreateBet places a bet and returns the outcome.
func (sms slotMachineServer) CreateBet(ctx context.Context, b *bet.CreateBetRequest) (*bet.CreateBetResponse, error) {
	// perform validation
	if b.Bet > b.Chips {
		s := status.New(codes.InvalidArgument, "placing the bet")
		s, _ = s.WithDetails(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{Field: "bet", Description: "it needs to be smaller than chips"},
			}})
		return nil, s.Err()
	}
	// just for learning about metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println(md)
	}
	r, err := sms.m.Bet(int(b.Chips), int(b.Bet))
	if err != nil {
		s := status.Newf(codes.InvalidArgument, "placing the bet: %v", err)
		return nil, s.Err()
	}
	// lots of unfortunate explicit conversion between what the business code returns and what we return from rpc
	// TODO how can this be avoided
	spins := make([]*bet.SpinDescription, len(r.Spins))
	for i, spin := range r.Spins {
		spins[i] = &bet.SpinDescription{
			Type: bet.SpinType(spin.TypeInt),
			Win:  int32(spin.Win),
		}
		spinLines := make([]*bet.SpinLine, len(spin.Stops))
		for i, stop := range spin.Stops {
			spinLines[i] = &bet.SpinLine{Value: make([]int32, len(stop))}
			for j := range stop {
				spinLines[i].Value[j] = int32(stop[j])
			}
		}
		spins[i].Lines = spinLines
		payLines := make([]*bet.PayLine, len(spin.PayLines))
		for i, payLine := range spin.PayLines {
			payLines[i] = &bet.PayLine{Value: make([]int32, len(payLine))}
			for j := range payLine {
				payLines[i].Value[j] = int32(payLine[j])
			}
		}
		spins[i].PayLines = payLines
	}
	return &bet.CreateBetResponse{
		Jwt:   b,
		Spins: spins,
		Win:   int32(r.Win),
		Chips: b.Chips,
		Bet:   b.Bet,
	}, nil
}
