package grpc

import (
	"log"

	"google.golang.org/grpc"
)

type GRPCCLIENT struct {
	*grpc.ClientConn
}

func NewClient(host string) *GRPCCLIENT {
	cli, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
	}
	return &GRPCCLIENT{cli}
}
