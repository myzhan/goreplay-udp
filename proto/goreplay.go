package proto

import (
	"bytes"
	"strconv"
)

const (
	RequestPayload          = '1'
	ResponsePayload         = '2'
	ReplayedResponsePayload = '3'
)

var PayloadSeparator = "\n🐵🙈🙉\n"

func PayloadHeader(payloadType byte, uuid []byte, timing int64) (header []byte) {
	var sTime string

	sTime = strconv.FormatInt(timing, 10)

	//Example:
	// 3 f45590522cd1838b4a0d5c5aab80b77929dea3b3 1231\n
	// `+ 1` indicates space characters or end of line
	headerLen := 1 + 1 + len(uuid) + 1 + len(sTime) + 1

	header = make([]byte, headerLen)
	header[0] = payloadType
	header[1] = ' '
	header[2+len(uuid)] = ' '
	header[len(header)-1] = '\n'

	copy(header[2:], uuid)
	copy(header[3+len(uuid):], sTime)

	return header
}

func PayloadBody(payload []byte) []byte {
	headerSize := bytes.IndexByte(payload, '\n')
	return payload[headerSize+1:]
}

func PayloadMeta(payload []byte) [][]byte {
	headerSize := bytes.IndexByte(payload, '\n')
	if headerSize < 0 {
		headerSize = 0
	}
	return bytes.Split(payload[:headerSize], []byte{' '})
}

func IsRequestPayload(payload []byte) bool {
	return payload[0] == RequestPayload
}
