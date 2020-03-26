package rpc

import (
	"net/rpc"
	"sync"
)

// rpc.Client 的连接池
type clientPool struct {
	clients []*rpc.Client

	sync.Mutex
}

// 连接池大小
func (p *clientPool) len() int {
	return len(p.clients)
}

// 读取连接池
func (p *clientPool) get() (client *rpc.Client) {
	p.Lock()
	if len(p.clients) > 0 {
		client = p.clients[0]
		p.clients = p.clients[1:]
	}
	p.Unlock()

	return
}

// 放回连接池
func (p *clientPool) put(client *rpc.Client) {
	p.Lock()
	p.clients = append(p.clients, client)
	p.Unlock()
}

// 关闭所有 rpc.Client 并清空连接池
func (p *clientPool) clean() (err error) {
	p.Lock()
	for _, c := range p.clients {
		err = c.Close()
	}
	p.clients = p.clients[:0]
	p.Unlock()

	return
}
