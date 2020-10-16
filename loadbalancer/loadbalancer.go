// Package loadbalancer 负载均衡
package loadbalancer

import "sync"

// LoadBalancer 接口
type LoadBalancer interface {
	InitNodes(nodes []Node)
	Nodes() []Node
	Select() string
}

// Node 节点
type Node struct {
	NodeKey string
	Weight  int
}

// RobinLoadBalancer 结构
type RobinLoadBalancer struct {
	mu            sync.Mutex
	idx           int
	currentWeight int
	nodes         []Node
}

// Nodes 所有节点
func (l *RobinLoadBalancer) Nodes() []Node {
	return l.nodes
}

// InitNodes 初始化节点
func (l *RobinLoadBalancer) InitNodes(nodes []Node) {
	defer l.mu.Unlock()
	l.mu.Lock()
	l.nodes = nodes
}

// NewRobinLoadBalancer new
func NewRobinLoadBalancer() LoadBalancer {
	return &RobinLoadBalancer{}
}

// MaxWeight 最大权重
func (l *RobinLoadBalancer) MaxWeight() int {
	w := 0
	for _, node := range l.nodes {
		if node.Weight >= w {
			w = node.Weight
		}
	}
	return w
}

// Select 根据权重选择一个
func (l *RobinLoadBalancer) Select() string {
	defer l.mu.Unlock()
	l.mu.Lock()
	for {
		l.idx = (l.idx + 1) % len(l.nodes)
		if l.idx == 0 {
			l.currentWeight--
			if l.currentWeight <= 0 {
				l.currentWeight = l.MaxWeight()
				if l.currentWeight <= 0 {
					return ""
				}
			}
		}
		if l.nodes[l.idx].Weight >= l.currentWeight {
			return l.nodes[l.idx].NodeKey
		}
	}
}
