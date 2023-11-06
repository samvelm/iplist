package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "illumio.com/iplist/proto"
	u "illumio.com/iplist/utils"
)

const (
	ADDRESS = "localhost:50051"
)

type IplistTask struct {
	Name        string
	Description string
	Done        bool
}

func main() {
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}

	defer conn.Close()

	client := pb.NewIplistServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := client.ResolveUser(ctx, &pb.SourceIp{IpAddress: u.GetRandomIpAddress()})

	if err != nil {
		log.Fatalf("could not resolve user: %v", err)
	}

	log.Printf(`
		userId : %s
		ipAddress : %s
		groups : %v,
	`, res.GetUserId(), res.GetIpAddress(), res.GetGroups())

}
