package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "illumio.com/iplist/proto"
	u "illumio.com/iplist/utils"
)

const (
	ADDRESS           = "localhost:50051"
	RANDOM_ITERATIONS = 2
)

type IplistTask struct {
	Name        string
	Description string
	Done        bool
}

func main() {
	// dail server
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}
	// close connection when done
	defer conn.Close()

	// create stream
	client := pb.NewIplistServiceClient(conn)
	stream, err := client.ResolveIpAddress(context.Background())
	if err != nil {
		log.Fatalf("opening stream error %v", err)
	}

	ctx := stream.Context()
	done := make(chan bool)

	// first goroutine creates 10 random IP addresses and sends to stream
	// closes it after 10 iterations
	go func() {
		for i := 1; i <= RANDOM_ITERATIONS; i++ {
			// generates random IP address and sends it to stream
			sourceIp := pb.SourceIp{IpAddress: u.GetRandomIpAddress()}
			if err := stream.Send(&sourceIp); err != nil {
				log.Fatalf("can not send %v", err)
			}
			// print the transferred IP address and sleep for 200 milliseconds
			log.Printf("The generated source IP sent -> %s", sourceIp.IpAddress)
			time.Sleep(time.Millisecond * 200)
		}
		if err := stream.CloseSend(); err != nil {
			log.Println(err)
		}
	}()

	// the second goroutine receives data from the stream and prints it.
	// if the stream is finished closes done channel
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				// the stream is finished. close the done channel
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			log.Printf(`User groups based on requested source IP
				userId : %s
				ipAddress : %s
				groups : %v,
			`, resp.GetUserId(), resp.GetIpAddress(), resp.GetGroups())
		}
	}()

	// third goroutine closes done channel if context is done
	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Println(err)
		}
	}()
	<-done
	log.Printf("processed all calls")

}
