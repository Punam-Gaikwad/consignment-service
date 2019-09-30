package main

import (
	"fmt"
	"log"
	"net"

	// Import the generated protobuf code
	pb "github.com/Punam-Gaikwad/consignment-service/proto/consignment"
	vesselProto "github.com/Punam-Gaikwad/vessel-service/proto/vessel"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	address     = "vessel-service:50052"
	port        = ":50051"
	defaultHost = "mongodb://datastore:27017/?ssl=false"
)

func main() {
	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err, port)
	}

	server := grpc.NewServer()
	fmt.Println("Connected to grpc server")

	//uri := os.Getenv("DB_HOST")
	//if uri == "" {
	uri := defaultHost
	//	}
	fmt.Println("uri: ", defaultHost)
	client, err := CreateClient(uri)
	if err != nil {
		log.Panic(err)
	}

	defer client.Disconnect(context.TODO())

	consignmentsCollection := client.Database("shippy").Collection("consignments")

	repository := &MongoRepository{consignmentsCollection}

	// Set up a connection to the vessel server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	fmt.Println("connection ", conn.GetState())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	vesselClient := vesselProto.NewVesselServiceClient(conn)

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterShippingServiceServer(server, &Handler{repository, vesselClient})

	// Register reflection service on gRPC server.
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
