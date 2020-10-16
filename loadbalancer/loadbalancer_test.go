// Package loadbalancer 负载均衡测试
package loadbalancer

import (
	"sync"
	"testing"

	"github.com/sssvip/goutil/strutil"
	"github.com/stretchr/testify/assert"
)

func TestRobinLoadBalancer_Select(t *testing.T) {
	b := NewRobinLoadBalancer()
	var nodes []Node
	prefix := "server"
	maxNodes := 100
	maxCount := 0
	dstCounter := make(map[string]int)
	for i := range make([]bool, maxNodes) {
		maxCount += maxNodes
		nodeKey := strutil.Format("%v%v", prefix, i+1)
		dstCounter[nodeKey] = maxNodes
		nodes = append(nodes, Node{NodeKey: nodeKey, Weight: maxNodes})
		maxNodes--
	}
	b.InitNodes(nodes)
	counter := make(map[string]int)
	counterMutux := sync.Mutex{}
	wg := sync.WaitGroup{}
	for range make([]bool, maxCount) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			nodeKey := b.Select()
			counterMutux.Lock()
			if v, ok := counter[nodeKey]; ok {
				counter[nodeKey] = v + 1
			} else {
				counter[nodeKey] = 1
			}
			counterMutux.Unlock()
		}()
	}
	wg.Wait()
	assert.Equal(t, counter, dstCounter)
}
