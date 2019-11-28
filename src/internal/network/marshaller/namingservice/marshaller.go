package namingservice

import (
	"fmt"
	"strings"
)

// MarshallRequest encodes a NamingService RequestMessage to bytes
func MarshallRequest(msg *RequestMessage) []byte {
	return []byte(fmt.Sprintf("%s:%s", msg.Op.String(), msg.Data))
}

// UnmarshallRequest decodes bytes into a NamingService RequestMessage
func UnmarshallRequest(data []byte) *RequestMessage {
	msg := strings.Split(string(data), ":")
	operation := StringToOperation(msg[0])
	return &RequestMessage{
		Op:   operation,
		Data: msg[1],
	}
}

// MarshallResponse encodes a NamingService ResponseMessage to bytes
func MarshallResponse(msg *ResponseMessage) []byte {
	return []byte(fmt.Sprintf("%s:%s", msg.Res.String(), msg.Data))
}

// UnmarshallResponse decodes bytes into a NamingService ResponseMessage
func UnmarshallResponse(data []byte) *ResponseMessage {
	msg := strings.Split(string(data), ":")
	result := StringToResult(msg[0])
	return &ResponseMessage{
		Res:  result,
		Data: fmt.Sprintf("%s:%s", msg[1], msg[2]),
	}
}
