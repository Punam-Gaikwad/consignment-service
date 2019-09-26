build:
	protoc -I. --go_out=plugins=grpc:$(GOPATH)/src/github.com/Punam-Gaikwad/consignment-service \
	  proto/consignment/consignment.proto
