package logger

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	megabyte = 1 << 20
)

// Logger is a middleware that logs the each request.
// endpoint is a list of endpoint that will not be logged
func Logger(log *logrus.Entry, endpoints ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode := 0
		start := time.Now()
		path := c.Request.URL.Path
		host := c.Request.Host
		method := c.Request.Method
		contentLength := c.Request.ContentLength

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// Process request
		c.Next()

		if shouldLogEndpoint(path, endpoints) {
			statusCode = c.Writer.Status()
			errMsg := ""
			//extract message error from the response body
			if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError || statusCode >= http.StatusInternalServerError {
				//decompress errMsg
				var buf2 bytes.Buffer
				gr, _ := gzip.NewReader(bytes.NewBuffer(w.body.Bytes()))
				defer gr.Close()
				out, _ := ioutil.ReadAll(gr)
				buf2.Write(out)
				errMsg = fmt.Sprintf("%v", buf2.String())
			}

			end := time.Now()
			latency := end.Sub(start)
			timeStamp := time.Now()

			msg := fmt.Sprintf("%v | %3d | %13v | %s |%-7s %#v | %.2f MB\n%s",
				timeStamp.Format(time.RFC1123),
				statusCode,
				latency,
				host,
				method,
				path,
				float64(contentLength)/megabyte,
				errMsg)

			switch {
			case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
				{
					log.Warn(msg)
				}
			case statusCode >= http.StatusInternalServerError:
				{
					log.Error(msg)
				}
			default:
				log.Info(msg)
			}
		}
	}
}

func shouldLogEndpoint(str string, arr []string) bool {
	for _, elm := range arr {
		if strings.EqualFold(elm, str) {
			return false
		}
	}
	return true
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
