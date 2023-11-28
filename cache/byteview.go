package cache

type ByteView struct {
	data []byte
}

func (b ByteView) Len() int {
	return len(b.data)
}

func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.data)
}

func (b ByteView) Value() string {
	return string(b.data)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
