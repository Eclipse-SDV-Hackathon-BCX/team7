package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/transmitt0r/currentlocationclient/kuksa/val/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	addr = flag.String("addr", "localhost:55555", "http service address")
)

func setCurrentLatitude(ctx context.Context, c pb.VALClient, val float64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := c.Set(ctx, &pb.SetRequest{
		Updates: []*pb.EntryUpdate{
			{
				Entry: &pb.DataEntry{
					Path: "Vehicle.CurrentLocation.Latitude",
					Value: &pb.Datapoint{
						Timestamp: timestamppb.Now(),
						Value:     &pb.Datapoint_Double{Double: val},
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

func getCurrentLongitude(ctx context.Context, c pb.VALClient) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.Get(ctx, &pb.GetRequest{
		Entries: []*pb.EntryRequest{
			{
				Path: "Vehicle.CurrentLocation.Latitude",
				View: pb.View_VIEW_TARGET_VALUE,
			},
		},
	})
	return resp, err
}

func setCurrentLongitude(ctx context.Context, c pb.VALClient, val float64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := c.Set(ctx, &pb.SetRequest{
		Updates: []*pb.EntryUpdate{
			{
				Entry: &pb.DataEntry{
					Path: "Vehicle.CurrentLocation.Latitude",
					Value: &pb.Datapoint{
						Timestamp: timestamppb.Now(),
						Value:     &pb.Datapoint_Double{Double: val},
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

func getCurrentLatitude(ctx context.Context, c pb.VALClient) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.Get(ctx, &pb.GetRequest{
		Entries: []*pb.EntryRequest{
			{
				Path: "Vehicle.CurrentLocation.Latitude",
				View: pb.View_VIEW_TARGET_VALUE,
			},
		},
	})
	return resp, err
}

func currentLatitudeClient(ctx context.Context, c pb.VALClient, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		resp, err := getCurrentLatitude(ctx, c)
		if err != nil {
			log.Fatal(err)
		}
		if len(resp.Entries) > 0 {
			val := resp.Entries[0].Value.GetDouble()
			setCurrentLatitude(ctx, c, val)
		}
		time.Sleep(1 * time.Second)
	}
}

func currentLongitudeClient(ctx context.Context, c pb.VALClient, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		resp, err := getCurrentLongitude(ctx, c)
		if err != nil {
			log.Fatal(err)
		}
		if len(resp.Entries) > 0 {
			val := resp.Entries[0].Value.GetDouble()
			setCurrentLongitude(ctx, c, val)
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

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go currentLatitudeClient(ctx, c, wg)
	go currentLongitudeClient(ctx, c, wg)

	wg.Wait()
}
