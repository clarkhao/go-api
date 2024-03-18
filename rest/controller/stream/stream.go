package stream

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func GetChunkHandler(c *gin.Context) {
	// Set the headers for streaming response
	c.Header("Content-Type", "text/plain")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Cache-Control", "no-cache")
	c.Header("X-Content-Type-Options", "nosniff")

	// Get the ResponseWriter
	w := c.Writer

	// Implement streaming logic
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 500)
		data := []byte("Streaming data chunk " + string(rune(i+65)) + "\n")
		w.Write(data)

		// Flush the buffer
		f, ok := w.(http.Flusher)
		if ok {
			f.Flush()
		}
	}
	//w.Header() can be used to add more response headers
	// Signal the end of the stream
	w.WriteHeader(http.StatusOK)
}

func GetStreamHandler(c *gin.Context) {
	filename := c.Param("filename")
	file, err := os.Open("videos/" + filename)
	if err != nil {
		c.String(http.StatusNotFound, "Video not found.")
		return
	}
	defer file.Close()

	c.Header("Content-Type", "video/mp4")
	buffer := make([]byte, 64*1024) // 64KB buffer size
	io.CopyBuffer(c.Writer, file, buffer)
}
