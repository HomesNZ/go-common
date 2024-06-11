package logger

import (
	"sync"
)

type Buffer []byte

// Buffer pool to reuse byte buffers.
var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return (*Buffer)(&b)
	},
}
