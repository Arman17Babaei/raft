package raft

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Arman17Babaei/raft/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func SendRPCToPeer(peerID int, method string, request interface{}) (interface{}, error) {
	c, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", peerID), 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return nil, err
	}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300 * time.Millisecond)
	defer cancel()

	switch method {
	case "AppendEntries":
		return pb.NewRaftClient(c).AppendEntries(ctx, request.(*pb.AppendEntriesRequest))
	case "RequestVote":
		return pb.NewRaftClient(c).RequestVote(ctx, request.(*pb.RequestVoteRequest))
	default:
		return nil, fmt.Errorf("unknown method: %s", method)
	}
}