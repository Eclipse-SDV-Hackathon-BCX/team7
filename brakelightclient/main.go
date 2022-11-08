package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/transmitt0r/brakelightclient/kuksa/val/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	addr = flag.String("addr", "localhost:55555", "http service address")
)

func setBrakeLight(ctx context.Context, c pb.VALClient, val bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := c.Set(ctx, &pb.SetRequest{
		Updates: []*pb.EntryUpdate{
			{
				Entry: &pb.DataEntry{
					Path: "Vehicle.Body.Lights.IsBrakeOn",
					Value: &pb.Datapoint{
						Timestamp: timestamppb.Now(),
						Value:     &pb.Datapoint_Bool{Bool: val},
					},
				},
				Fields: []pb.Field{
					pb.Field_FIELD_VALUE,
				},
			},
		},
	})
	return err
}

func getBrakeLight(ctx context.Context, c pb.VALClient) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.Get(ctx, &pb.GetRequest{
		Entries: []*pb.EntryRequest{
			{
				Path: "Vehicle.Body.Lights.IsBrakeOn",
				View: pb.View_VIEW_TARGET_VALUE,
			},
		},
	})
	return resp, err
}

func brakeLightClient(ctx context.Context, c pb.VALClient) {
	for {
		resp, err := getBrakeLight(ctx, c)
		if err != nil {
			log.Fatal(err)
		}
		if len(resp.Entries) > 0 {
			val := resp.Entries[0].Value.GetBool()
			setBrakeLight(ctx, c, val)
		}
		time.Sleep(1 * time.Second)
	}

}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewVALClient(conn)

	fmt.Println("Connected to Databroker!")

	ctx := context.TODO()

	brakeLightClient(ctx, c)
}
