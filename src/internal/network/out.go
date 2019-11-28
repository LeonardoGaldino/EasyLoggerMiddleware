package network

import (
	"encoding/binary"
	"net"
)

// WriteMessage writes data buffer to the connection specifying the data size first (big endian)
func WriteMessage(conn *net.Conn, data []byte) error {
	index := 0
	dataSizeBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(dataSizeBuffer, uint32(len(data)))
	(*conn).Write(dataSizeBuffer)

	for {
		wrote, err := (*conn).Write(data[index:])
		if err != nil {
			return err
		}

		index += wrote
		if index == len(data) {
			break
		}
	}
	return nil
}
