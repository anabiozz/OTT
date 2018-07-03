package raft

import (
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/graft"
)

func Raft() {
	ci := graft.ClusterInfo{Name: "health_manager", Size: 3}
	rpc, err := graft.NewNatsRpc(&nats.DefaultOptions)
	errChan := make(chan error)
	stateChangeChan := make(chan StateChange)
	handler := graft.NewChanHandler(stateChangeChan, errChan)

	node, err := graft.New(ci, handler, rpc, "/tmp/graft.log")

	// ...

	if node.State() == graft.LEADER {
		// Process as a LEADER
	}

	select {

	case sc := <-stateChangeChan:
		// Process a state change
	case err := <-errChan:
		// Process an error, log etc.
	}
}
