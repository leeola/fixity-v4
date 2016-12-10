package database

import (
	"errors"

	"github.com/leeola/kala/index"
)

var (
	// ErrNoRecord is returned if the database could not find the requested record
	ErrNoRecord = errors.New("database could not find the requested record")
)

// Database is a db interface for Node, Index, and Peer data and metadata.
//
// Note that the Database interface is an interface for a single Node. If
// the actual database is shared by multiple Nodes, the initialization of that
// database connection must be provided the information needed to discern this
// Nodes data from the rest.
//
// For example `Database.NodeId()` returns this nodes Id without any params.
// Mysql could be configured to prefix all table requests with a node specifc name
// as specified by the nodes config and etc.
type Database interface {
	// SetNodeId sets a nodes Id in the db for the running Node.
	SetNodeId(nodeId string) error

	// GetNodeId gets the NodeId for the running Node.
	GetNodeId() (string, error)

	// SetPeerPinLastEntry sets the lastEntry for the given pin query from the peer.
	//
	// This is used to ensure a Pinning can keep track of the last successful pin
	// entry it received from a peer.
	SetPeerPinLastEntry(peerId string, pin index.PinQuery, lastEntry int) error

	// GetPeerPinLastEntry gets the lastEntry for the given pin query from the peer.
	GetPeerPinLastEntry(peerId string, pin index.PinQuery) (int, error)
}
