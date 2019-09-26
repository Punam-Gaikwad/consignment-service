package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Punam-Gaikwad/microservices/consignment-service/proto/consignment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Repository -
type Repository interface {
	Create(*pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
}

//MongoRepository Implementation
type MongoRepository struct {
	collection *mongo.Collection
}

//Create -
func (repository *MongoRepository) Create(consignment *pb.Consignment) error {
	_, err := repository.collection.InsertOne(context.Background(), consignment)
	return err
}

//GetAll -
func (repository *MongoRepository) GetAll() ([]*pb.Consignment, error) {
	fmt.Println("entered to GetAll method")

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(10)

	var consignments []*pb.Consignment

	cur, err := repository.collection.Find(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.Background()) {
		var consignment *pb.Consignment

		if err := cur.Decode(&consignment); err != nil {
			return nil, err
		}
		consignments = append(consignments, consignment)
	}
	cur.Close(context.TODO())
	return consignments, err
}
