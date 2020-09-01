package functions

import (
	"io"
	"net/http"
)

// pass CMD output to HTTP
func HttpExecPipe(res http.ResponseWriter, pipeReader *io.PipeReader) {
	BUFLEN := 1024 // for
	buffer := make([]byte, BUFLEN)
	defer Recover("Http Flush Panic")
	for {
		n, err := pipeReader.Read(buffer)
		if err != nil {
			pipeReader.Close()
			break
		}

		data := buffer[0:n]
		res.Write(data)
		f, ok := res.(http.Flusher)
		if ok {
			f.Flush()
		}

		//reset buffer
		for i := 0; i < n; i++ {
			buffer[i] = 0
		}
	}
}
