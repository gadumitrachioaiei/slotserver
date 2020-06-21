package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	bet "github.com/gadumitrachioaiei/slotserver/proto/gen/go/bet/v1"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	con, err := grpc.Dial(":8080", grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&basicAuthCreds{"user1", "pass1"}))
	if err != nil {
		log.Fatalf("cannot connect to: %s, because: %v", ":8080", err)
	}
	b := bet.CreateBetRequest{
		Uid:   "xyz",
		Chips: 1000,
		Bet:   100,
	}
	c := bet.NewSlotMachineServiceClient(con)
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("key1", "val1"))
	result, err := c.CreateBet(ctx, &b)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			log.Println("details are:", s.Details())
			for _, d := range s.Details() {
				if br, ok := d.(*errdetails.BadRequest); ok {
					log.Fatalln("br", br)
				} else {
					log.Fatalf("type is: %T", d)
				}
			}
		}
		log.Fatalf("placing bet: %v", err)
	}
	fmt.Println(result)
}

// basicAuthCreds is an implementation of credentials.PerRPCCredentials
// that transforms the username and password into a base64 encoded value similar
// to HTTP Basic xxx
type basicAuthCreds struct {
	username, password string
}

// GetRequestMetadata sets the value for "authorization" key
func (b *basicAuthCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Basic " + basicAuth(b.username, b.password),
	}, nil
}

// RequireTransportSecurity should be true as even though the credentials are base64, we want to have it encrypted over the wire.
func (b *basicAuthCreds) RequireTransportSecurity() bool {
	return false
}

//helper function
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
