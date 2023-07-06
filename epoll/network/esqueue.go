package network

import "sync"

type esqueue struct {
	sync.Mutex
	nodes []interface{}
}

func (q *esqueue) Add(node interface{}) (one bool) {
	q.Lock()
	q.nodes = append(q.nodes, node)
	n := len(q.nodes)
	q.Unlock()
	return n == 1
}

func (q *esqueue) ForEach(iter func(note interface{}) error) error {
	q.Lock()
	if len(q.nodes) == 0 {
		q.Unlock()
		return nil
	}
	nodes := q.nodes
	q.nodes = nil
	q.Unlock()

	for _, node := range nodes {
		if err := iter(node); err != nil {
			return err
		}
	}

	q.Lock()
	if q.nodes == nil {
		for i := range nodes {
			nodes[i] = nil
		}
		q.nodes = nodes[:0]
	}
	q.Unlock()
	return nil
}
