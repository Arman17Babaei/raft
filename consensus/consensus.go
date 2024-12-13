package consensus

import (
	"fmt"
	"sync"

	"github.com/Arman17Babaei/raft/grpc"
	"github.com/Arman17Babaei/raft/raft"
)

type ConsensusManager struct {
	mu          sync.Mutex
	value       int
	status      string
	myPort      int
	otherPorts  []int
	raftStarted bool
	raftCluster *RaftCluster
}

type RaftCluster struct {
	Node        *raft.Node
	RaftService *grpc.RaftService
}

func NewConsensusManager(initValue int, myPort int, otherPorts []int) *ConsensusManager {
	return &ConsensusManager{
		value:      initValue,
		status:     "Initialized",
		myPort:     myPort,
		otherPorts: otherPorts,
	}
}

func (cm *ConsensusManager) Start() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.raftStarted {
		fmt.Println("Raft already started")
		return
	}
	fmt.Printf("Starting Raft on port %d with peers %v...\n", cm.myPort, cm.otherPorts)
	
	cm.raftStarted = true
	cm.raftCluster = &RaftCluster{
		Node:        raft.NewNode(cm.myPort, cm.otherPorts),
		RaftService: grpc.NewRaftService(nil, cm.myPort),
	}
	cm.status = "Raft Started"
}

func (cm *ConsensusManager) Get() int {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.value
}

func (cm *ConsensusManager) Add(value int) {
	cm.applyCommand("add", value)
}

func (cm *ConsensusManager) Sub(value int) {
	cm.applyCommand("sub", value)
}

func (cm *ConsensusManager) Mul(value int) {
	cm.applyCommand("mul", value)
}

func (cm *ConsensusManager) GetStatus() string {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.status
}

func (cm *ConsensusManager) applyCommand(command string, value int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	fmt.Printf("Applying command: %s %d\n", command, value)

	switch command {
	case "add":
		cm.value += value
	case "sub":
		cm.value -= value
	case "mul":
		cm.value *= value
	}

	cm.status = fmt.Sprintf("Applied command: %s %d", command, value)
	fmt.Println("Command committed.")
}
