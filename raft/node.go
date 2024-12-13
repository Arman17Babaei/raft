package raft

import "sync"

type Role string

const (
    Follower  Role = "Follower"
    Candidate Role = "Candidate"
    Leader    Role = "Leader"
)

type Node struct {
    Mu          sync.Mutex
    Id          int          // Unique ID for the node
    Role        Role         // Current role: Follower, Candidate, or Leader
    CurrentTerm int          // Latest term seen
    VotedFor    *int         // Candidate ID voted for in the current term
    Log         []LogEntry   // Log entries
    CommitIndex int          // Index of the highest log entry known to be committed
    LastApplied int          // Index of the highest log entry applied to the state machine
    NextIndex   map[int]int  // For Leader: next log index to send to each follower
    MatchIndex  map[int]int  // For Leader: highest log entry index known to be replicated on each follower

    OthersIds []int // IDs of other nodes in the cluster
}

type LogEntry struct {
    Term    int         // Term when the entry was received by the leader
    Command interface{} // Command to apply to the state machine
}

func NewNode(id int, others []int) *Node {
    return &Node{
        Id:          id,
        Role:        Follower,
        CurrentTerm: 0,
        VotedFor:    nil,
        Log:         make([]LogEntry, 0),
        CommitIndex: 0,
        LastApplied: 0,
        NextIndex:   make(map[int]int),
        MatchIndex:  make(map[int]int),

        OthersIds: others,
    }
}
