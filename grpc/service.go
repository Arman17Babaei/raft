package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/Arman17Babaei/raft/grpc/proto"
	"github.com/Arman17Babaei/raft/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RaftService struct {
	mu sync.Mutex
	node *raft.Node
	pleaseCh chan *pb.RequestPleaseRequest

	pb.UnimplementedRaftServer
}

func NewRaftService(node *raft.Node, port int, pleaseCh chan *pb.RequestPleaseRequest) *RaftService {
	raftService := &RaftService{node: node, pleaseCh: pleaseCh}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpc.WithTransportCredentials(insecure.NewCredentials())
	grpcServer := grpc.NewServer()
	pb.RegisterRaftServer(grpcServer, raftService)

	go func() {
		fmt.Printf("Starting gRPC server on port %d...", port)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return raftService
}
func (rs *RaftService) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	fmt.Printf("AppendEntries received: %+v\n", req)

	return rs.node.AppendEntriesHandler(req)
}

func (rs *RaftService) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	
	fmt.Printf("RequestVote received: %+v\n", req)

	return rs.node.RequestVoteHandler(req)
}

func (rs *RaftService) PleaseDoThis(ctx context.Context, req *pb.RequestPleaseRequest) (*pb.RequestPleaseResponse, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	
	fmt.Printf("PleaseDoThis received: %+v\n", req)

	rs.pleaseCh <- req

	return &pb.RequestPleaseResponse{Success: true}, nil
}