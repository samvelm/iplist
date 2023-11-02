package main

import (
	"context"
	"log"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	pb "illumio.com/iplist/proto"
)

const (
	// Port for gRPC server to listen to
	PORT = ":50051"
)

type IplistServer struct {
	pb.UnimplementedIplistServiceServer
}

func (s *IplistServer) CreateIplist(ctx context.Context, in *pb.NewIplist) (*pb.Iplist, error) {
	log.Printf("Received: %v", in.GetName())
	iplist := &pb.Iplist{
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Done:        false,
		Id:          uuid.New().String(),
	}

	return iplist, nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterIplistServiceServer(s, &IplistServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
