package bolt

import (
	"fmt"

	"github.com/leeola/kala/index"
)

func (b *Bolt) SetNodeId(nodeId string) error {
	return b.SetString(nodeBucketName, nodeIdKey, nodeId)
}

func (b *Bolt) GetNodeId() (string, error) {
	return b.GetString(nodeBucketName, nodeIdKey)
}

func (b *Bolt) GetIndexEntry(i int) (string, error) {
	return b.GetEntry(i)
}

func (b *Bolt) SetPeerPinLastEntry(peerId string, pin index.PinQuery, lastEntry int) error {
	key := []byte(fmt.Sprintf(peerId + pin.CommaString()))
	return b.SetInt(peersBucketName, key, lastEntry)
}

func (b *Bolt) GetPeerPinLastEntry(peerId string, pin index.PinQuery) (int, error) {
	key := []byte(fmt.Sprintf(peerId + pin.CommaString()))
	return b.GetInt(peersBucketName, key)
}

func (b *Bolt) SetNodeIndexVersion(version string) error {
	return b.SetString(nodeBucketName, nodeIndexVersionKey, version)
}

func (b *Bolt) GetNodeIndexVersion() (string, error) {
	return b.GetString(nodeBucketName, nodeIndexVersionKey)
}
