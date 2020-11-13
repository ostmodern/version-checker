package architecture

import (
	"fmt"
	"sync"

	corev1 "k8s.io/api/core/v1"
)

var _ NodeArchitectureMap = &defaultNodeMap{}

const (
	invalidFuncInput = "invalid function input"
)

// NodeMetadata metadata about a particular node
type NodeMetadata struct {
	OS           string
	Architecture string
}

type NodeArchitectureMap interface {
	GetNodeArchitecture(node string) (*NodeMetadata, error)
	AddNode(node *corev1.Node) error
	DeleteNode(node string) error
	Length() int
}

type defaultNodeMap struct {
	mu    sync.RWMutex
	nodes map[string]*NodeMetadata
}

func New() *defaultNodeMap {
	// might need to pass an initial map
	return &defaultNodeMap{
		nodes: make(map[string]*NodeMetadata),
	}
}

func (m *defaultNodeMap) GetNodeArchitecture(node string) (*NodeMetadata, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.nodes[node]; !ok {
		// no data about the node was found, return error
		return nil, fmt.Errorf("error fetching node's architecture data")
	}
	return m.nodes[node], nil
}

func (m *defaultNodeMap) AddNode(node *corev1.Node) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if node == nil {
		return fmt.Errorf("add node %s", invalidFuncInput)
	}

	arch, ok := node.Labels["kubernetes.io/arch"]
	if !ok {
		return fmt.Errorf("\"kubernetes.io/arch\" label not found on node: %s", node.Name)
	}

	os, ok := node.Labels["kubernetes.io/os"]
	if !ok {
		return fmt.Errorf("\"kubernetes.io/os\" label not found on node: %s", node.Name)
	}

	// change name to uid or selflink
	m.nodes[node.Name] = &NodeMetadata{
		OS:           os,
		Architecture: arch,
	}
	return nil
}

func (m *defaultNodeMap) DeleteNode(node string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if node == "" {
		return fmt.Errorf("delete node %s", invalidFuncInput)
	}
	// change name to uid or selflink
	delete(m.nodes, node)
	return nil
}

func (m *defaultNodeMap) Length() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.nodes)
}
