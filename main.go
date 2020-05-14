package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	pb "github.com/seftomsk/shippy-service-consignment/proto/consignment"
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

//func (s *service) CreateConsignment(_ context.Context, in *pb.Consignment) (*pb.Response, error) {
//	// Save our consignment
//	consignment, err := s.repo.Create(in)
//	if err != nil {
//		return nil, err
//	}
//
//	// Return matching the `Response` message we created in our
//	// protobuf definition.
//
//	return &pb.Response{Created: true, Consignment: consignment}, nil
//}
//
//func (s *service) GetConsignments(_ context.Context, _ *pb.GetRequest) (*pb.Response, error) {
//	consignments := s.repo.GetAll()
//	return &pb.Response{Consignments: consignments}, nil
//}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	// Create a new service. Optionally include some options here.
	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("shippy.service.consignment"),
	)

	// Init will parse the command line flags.
	srv.Init()

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}


	//// Set-up our gRPC server.
	//listener, err := net.Listen("tcp", port)
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//server := grpc.NewServer()
	//
	//// Register our service with the gRPC server, this will tie our
	//// implementation into the auto-generated interface code for our
	//// protobuf definition.
	//pb.RegisterShippingServiceServer(server, &service{repo})
	//
	//// Register reflection service on gRPC server.
	//reflection.Register(server)
	//
	//log.Println("Running on port: ", port)
	//if err := server.Serve(listener); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
}
