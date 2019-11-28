package network

import (
	"encoding/binary"
	"io"
	"net"
)

// ReadMessage read data from connection expecting first data size (big endian)
func ReadMessage(conn *net.Conn) ([]byte, error) {
	dataSizeBuffer := make([]byte, 4)
	(*conn).Read(dataSizeBuffer)
	maxBytes := binary.BigEndian.Uint32(dataSizeBuffer)

	bufferSize, reqInitialBufferSize := 512, 1024
	buffer, req := make([]byte, bufferSize, bufferSize), make([]byte, reqInitialBufferSize, reqInitialBufferSize)
	lowIndex := 0
	for {
		lenRead, err := (*conn).Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return []byte{}, err
		}

		// final buffer needs more space
		if lenRead > len(req)-lowIndex+1 {
			req = extendBuffer(req, lenRead)
		}

		copy(req[lowIndex:], buffer)
		lowIndex += lenRead

		if lowIndex >= int(maxBytes) {
			break
		}
	}

	return req[:lowIndex], nil
}
