package network

func extendBuffer(buffer []byte, extra int) []byte {
	currentBufferSize := len(buffer)
	newBuffer := make([]byte, (2*currentBufferSize)+extra, (2*currentBufferSize)+extra)
	copy(newBuffer, buffer)
	return newBuffer
}
