package WEST

import (
	"bytes"
	"net/http"
)

type WESTWriter struct {
	buffer      *bytes.Buffer
	status      int
	wroteStatus bool
}

func (writer WESTWriter) Header() http.Header {
	return http.Header{}
}

func (writer WESTWriter) Write(b []byte) (int, error) {
	return writer.buffer.Write(b)
}

func (writer WESTWriter) WriteHeader(status int) {
	if !writer.wroteStatus {
		writer.wroteStatus = true
		writer.status = status
	}
}

func (writer WESTWriter) Status() int {
	return writer.status
}

func MakeWestWriter() *WESTWriter {
	return &WESTWriter{
		buffer:      &bytes.Buffer{},
		status:      200,
		wroteStatus: false,
	}
}
