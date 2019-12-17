package pool

import (
	"bytes"
	"sync"
)

type BytesBuffer struct {
	pool sync.Pool
}

func NewBytesBuffer() *BytesBuffer {
	return &BytesBuffer{pool: sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}}
}

func (p *BytesBuffer) Get() *bytes.Buffer {
	b := p.pool.Get().(*bytes.Buffer)
	b.Reset()

	return b
}

func (p *BytesBuffer) Put(b *bytes.Buffer) {
	p.pool.Put(b)
}
