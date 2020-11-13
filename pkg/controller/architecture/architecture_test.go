package architecture

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
)

func TestAdd(t *testing.T) {
	var wg sync.WaitGroup // using wait group to know if task was completed
	var commonMap = New()
	var expectedNumberOfNodes = 100

	// Adding random nodes to the concurrent map
	for i := 0; i < expectedNumberOfNodes; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			node, err := generateNode("node", "amd64", "linux")
			if err != nil {
				t.Errorf("unxepected error when generating the node info: %s", err)
			}
			err = commonMap.AddNode(node)
			if err != nil {
				t.Errorf("unxepected error when appending the node info: %s", err)
			}
		}()
	}
	wg.Wait()

	// checking if all operations occurred without any issues on the locks
	if commonMap.Length() != expectedNumberOfNodes {
		t.Errorf("unxepected number of node information was found: exp=%d act=%d", expectedNumberOfNodes, commonMap.Length())
	}

}

func TestRead(t *testing.T) {
	var wg sync.WaitGroup // using wait group to know if task was completed
	var commonMap = New()
	var expectedNumberOfNodes = 1000
	var nodes []*corev1.Node

	for i := 0; i < expectedNumberOfNodes; i++ {
		node, err := generateNode("node", "amd64", "linux")
		if err != nil {
			t.Errorf("unxepected error when generating the node info: %s", err)
		}
		err = commonMap.AddNode(node)
		if err != nil {
			t.Errorf("unxepected error when appending the node info: %s", err)
		}
		nodes = append(nodes, node)
	}

	for _, node := range nodes {
		wg.Add(1)
		go func() {
			arch, err := commonMap.GetNodeArchitecture(node.Name)
			if err != nil {
				t.Errorf("unxepected error when reading node information: %s", err)
			}
			if expectedArch, ok := node.Labels["kubernetes.io/arch"]; ok {
				if arch.Architecture != expectedArch {
					t.Errorf("unxepected node architecture was found: exp=%q act=%q", expectedArch, arch.Architecture)
				}
			}
			if expectedOS, ok := node.Labels["kubernetes.io/os"]; ok {
				if arch.OS != expectedOS {
					t.Errorf("unxepected node architecture was found: exp=%q act=%q", expectedOS, arch.OS)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if commonMap.Length() != expectedNumberOfNodes {
		t.Errorf("unxepected number of node information was found: exp=%d act=%d", expectedNumberOfNodes, commonMap.Length())
	}

}

func TestDelete(t *testing.T) {
	var wg sync.WaitGroup // using wait group to know if task was completed
	var commonMap = New()
	var expectedNumberOfNodes = 1000
	var nodes []*corev1.Node

	for i := 0; i < expectedNumberOfNodes; i++ {
		node, err := generateNode("node", "amd64", "linux")
		if err != nil {
			t.Errorf("unxepected error when generating the node info: %s", err)
		}
		err = commonMap.AddNode(node)
		if err != nil {
			t.Errorf("unxepected error when appending the node info: %s", err)
		}
		nodes = append(nodes, node)
	}

	for _, node := range nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := commonMap.DeleteNode(node.Name)
			if err != nil {
				t.Errorf("unxepected error when deleting the node info: %s", err)
			}
		}()
		wg.Wait()
	}

	// checking if all operations occurred without any issues on the locks
	if commonMap.Length() != 0 {
		t.Errorf("unxepected number of node information was found: exp=%d act=%d", 0, commonMap.Length())
	}

}

func generateNode(name, arch, os string) (*corev1.Node, error) {
	suffix, err := uuid.NewRandom()
	node := &corev1.Node{}
	node.Name = fmt.Sprintf("%s%d", name, suffix.ID())
	node.Labels = map[string]string{
		"kubernetes.io/arch": arch,
		"kubernetes.io/os":   os,
	}
	return node, err
}
