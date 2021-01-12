package GeeCache

//ByteView is read-only data,
//It will return copy of source data when you try to get access the data
type ByteView struct {
	bytes []byte
}

func (b ByteView) Len() int {
	return len(b.bytes)
}
func (b *ByteView) ByteSlice() []byte {
	return copyBytes(b.bytes)
}
func (b *ByteView) String() string {
	c := b.ByteSlice()
	return string(c)
}
func copyBytes(bytes []byte) []byte {
	c := make([]byte, len(bytes))
	copy(c, bytes)
	return c
}
