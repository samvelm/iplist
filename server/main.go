package main

import (
	"context"
	"io"
	"log"
	"net"

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

// ResolveUser implementation of the client server gRPC request
func (s *IplistServer) ResolveUser(ctx context.Context, in *pb.SourceIp) (*pb.UserGroups, error) {
	log.Printf("Received: %v", in.GetIpAddress())
	var groups = []string{"Engineering", "Finance"}

	userGroups := &pb.UserGroups{
		IpAddress: in.GetIpAddress(),
		UserId:    "john.doe",
		Groups:    groups,
	}

	return userGroups, nil
}

// ResolveIpAddress implementation of the bidirectional streaming
func (s *IplistServer) ResolveIpAddress(srv pb.IplistService_ResolveIpAddressServer) error {
	log.Println("start new server")
	ctx := srv.Context()

	// endless loop, will be be over once the context is done
	for {
		// exit if the context is done
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := srv.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		// just print what has been received from the stream
		log.Printf("received IP address from the client -> %s", req.GetIpAddress())
		// Generate the response with the newly generated IP address
		var groups = []string{"Engineering", "Finance", "Ops"}

		// For testing only.
		// Generate the new random IP address and put it back in the response
		// newIpAddress := u.GetRandomIpAddress()
		// log.Printf("just print the newly generated IP address -> %s\n\n", newIpAddress)
		userGroups := &pb.UserGroups{
			IpAddress: req.GetIpAddress(), // newIpAddress,
			UserId:    "john.bidirectional.doe",
			Groups:    groups,
		}
		if err := srv.Send(userGroups); err != nil {
			log.Printf("send error %v", err)
		}
	}
}

// Run the server listenning the specific port
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
