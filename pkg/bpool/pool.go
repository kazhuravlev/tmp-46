package bpool

import (
	"bytes"
	"sync"
)

type Pool struct {
	pool *sync.Pool
}

func New() *Pool {
	return &Pool{
		pool: &sync.Pool{
			New: func() any { return new(bytes.Buffer) },
		},
	}
}

func (p *Pool) Get() *bytes.Buffer {
	item := p.pool.Get().(*bytes.Buffer)

	return item
}

func (p *Pool) Put(buf *bytes.Buffer) {
	buf.Reset()
	p.pool.Put(buf)
}
