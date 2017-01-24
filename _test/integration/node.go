package integration

import (
	"fmt"
	"path/filepath"

	"github.com/leeola/kala/multiload"
	"github.com/leeola/kala/node"
)

func DescribeNodeTest(n *NodeTest) string {
	// Disabled until i can decide how to print the complex peering information.
	//
	// var peersDesc string
	// for i, p := range n.Peers {
	// 	s := DescribeNodeTest(p)
	// 	lines := strings.Split(s, "\n")
	// 	// pad each line
	// 	for x, l := range lines {
	// 		lines[i] = "  " + l
	// 	}
	// 	peersDesc += strings.Join(lines, "\n")
	// }

	return fmt.Sprintf("config:%s", n.ConfigPath)
}

type NodeTest struct {
	ConfigPath string
	Node       *node.Node
}

// TODO(leeola): allow disabling tests via env vars
func NodeTests(tmpDir string) []*NodeTest {
	configRoot := "../fixtures/configs"

	return []*NodeTest{{
		ConfigPath: "basic.toml",
		Node:       mustLoadNode(filepath.Join(configRoot, "basic.toml"), tmpDir),
	}}
}

func mustLoadNode(c, t string) *node.Node {
	rc := multiload.RootConfig{
		RootPath: t,
	}

	n, err := multiload.LoadNodeWithDefault(c, rc, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	return n
}
