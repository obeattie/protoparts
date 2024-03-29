package protoparts

import (
	"bytes"
	"sync"
)

var bufferPool = bufferPoolT{
	p: &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		}}}

type bufferPoolT struct {
	p *sync.Pool
}

func (b bufferPoolT) Get() *bytes.Buffer {
	buf := b.p.Get().(*bytes.Buffer)
	return buf
}

func (b bufferPoolT) Put(buf *bytes.Buffer) {
	buf.Reset()
	b.p.Put(buf)
}
