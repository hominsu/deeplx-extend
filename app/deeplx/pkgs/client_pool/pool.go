package client_pool

import (
	"sync/atomic"

	"github.com/valyala/fasthttp"
)

type ClientPool struct {
	clients []*fasthttp.Client
	size    uint32
	index   uint32
}

func NewClientPool(clients ...*fasthttp.Client) (*ClientPool, func()) {
	cleanup := func() {
		for _, client := range clients {
			client.CloseIdleConnections()
		}
	}

	return &ClientPool{
		clients: clients,
		size:    uint32(len(clients)),
	}, cleanup
}

func (p *ClientPool) Get() *fasthttp.Client {
	index := atomic.AddUint32(&p.index, 1) % p.size
	return p.clients[index]
}

func (p *ClientPool) Add(client *fasthttp.Client) {
	atomic.AddUint32(&p.size, 1)
	p.clients = append(p.clients, client)
}

func (p *ClientPool) Len() int {
	return int(atomic.LoadUint32(&p.size))
}
