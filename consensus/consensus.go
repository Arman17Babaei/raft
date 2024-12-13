package consensus

import (
	"encoding/json"
	"fmt"
	"sync"

	pb "github.com/Arman17Babaei/raft/grpc/proto"
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
	requestCh   chan raft.Request
	commandCh   chan string
}

type Command struct {
	Value     int    `json:"value"`
	Operation string `json:"operation"`
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
	cm.requestCh = make(chan raft.Request, 20)
	cm.commandCh = make(chan string)
	pleaseCh := make(chan *pb.RequestPleaseRequest, 20)
	node := raft.NewNode(cm.myPort, cm.otherPorts, cm.requestCh, cm.commandCh, pleaseCh)
	cm.raftCluster = &RaftCluster{
		Node:        node,
		RaftService: grpc.NewRaftService(node, cm.myPort, pleaseCh),
	}
	cm.raftCluster.Node.StartElection()
	go cm.listenForCommands()
	cm.status = "Raft Started"
}

func (cm *ConsensusManager) Get() int {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.value
}

func (cm *ConsensusManager) Add(value int) error {
	return cm.submitCommand(Command{Operation: "add", Value: value})
}

func (cm *ConsensusManager) Sub(value int) error {
	return cm.submitCommand(Command{Operation: "sub", Value: value})
}

func (cm *ConsensusManager) Mul(value int) error {
	return cm.submitCommand(Command{Operation: "mul", Value: value})
}

func (cm *ConsensusManager) GetStatus() string {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.status
}

func (cm *ConsensusManager) submitCommand(command Command) error {
	commandS, err := json.Marshal(command)
	if err != nil {
		fmt.Printf("Failed to marshal command: %v\n", err)
		return err
	}
	cb := make(chan bool)
	cm.requestCh <- raft.Request{Command: string(commandS), Callback: cb}
	<-cb
	return nil
}

func (cm *ConsensusManager) listenForCommands() {
	for command := range cm.commandCh {
		fmt.Printf("Received command: %s\n", command)
		cm.applyCommand(command)
		fmt.Printf("Current value: %d\n", cm.value)
	}
}

func (cm *ConsensusManager) applyCommand(command string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	c := Command{}
	err := json.Unmarshal([]byte(command), &c)
	if err != nil {
		fmt.Printf("Failed to unmarshal command: %v\n", err)
		return
	}

	fmt.Printf("Applying command: %s %d\n", c.Operation, c.Value)

	switch c.Operation {
	case "add":
		cm.value += c.Value
	case "sub":
		cm.value -= c.Value
	case "mul":
		cm.value *= c.Value
	}

	cm.status = fmt.Sprintf("Applied command: %s %d", c.Operation, c.Value)
	fmt.Println("Command committed.")
}
