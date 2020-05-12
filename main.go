package main

import (
	"context"
	pb "github.com/seftomsk/shippy-service-consignment/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

const port = ":50051"

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

// Create a new consignment
func (r *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	r.mu.Lock()
	updated := append(r.consignments, consignment)
	r.consignments = updated
	r.mu.Unlock()
	return consignment, nil
}

// GetAll consignments
func (r *Repository) GetAll() []*pb.Consignment {
	return r.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo repository
}

func (s *service) CreateConsignment(_ context.Context, in *pb.Consignment) (*pb.Response, error) {
	// Save our consignment
	consignment, err := s.repo.Create(in)
	if err != nil {
		return nil, err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func (s *service) GetConsignments(ctx context.Context, in *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()
	return &pb.Response{Consignments: consignments}, nil
}

func main() {
	repo := &Repository{}

	// Set-up our gRPC server.
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterShippingServiceServer(server, &service{repo})

	// Register reflection service on gRPC server.
	reflection.Register(server)

	log.Println("Running on port: ", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
