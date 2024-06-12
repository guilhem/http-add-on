package middleware

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("loggingResponseWriter", func() {
	Context("New", func() {
		It("returns new object with expected field values set", func() {
			var (
				w = httptest.NewRecorder()
			)

			lrw := newLoggingResponseWriter(w)
			Expect(lrw).NotTo(BeNil())
			Expect(lrw.downstreamResponseWriter).To(Equal(w))
			Expect(lrw.bytesWritten).To(Equal(0))
			Expect(lrw.statusCode).To(Equal(0))
		})
	})

	Context("BytesWritten", func() {
		It("returns the expected value", func() {
			const (
				bw = 128
			)

			lrw := &loggingResponseWriter{
				bytesWritten: bw,
			}

			ret := lrw.BytesWritten()
			Expect(ret).To(Equal(bw))
		})
	})

	Context("StatusCode", func() {
		It("returns the expected value", func() {
			const (
				sc = http.StatusTeapot
			)

			lrw := &loggingResponseWriter{
				statusCode: sc,
			}

			ret := lrw.StatusCode()
			Expect(ret).To(Equal(sc))
		})
	})

	Context("Header", func() {
		It("returns downstream method call", func() {
			var (
				w = httptest.NewRecorder()
			)

			lrw := &loggingResponseWriter{
				downstreamResponseWriter: w,
			}

			h := w.Header()
			h.Set("Content-Type", "application/json")

			ret := lrw.Header()
			Expect(ret).To(Equal(h))
		})
	})

	Context("Write", func() {
		It("invokes downstream method, increases bytesWritten accordingly, and returns expected values", func() {
			const (
				body      = "KEDA"
				bodyLen   = len(body)
				initialBW = 60
			)

			var (
				w = httptest.NewRecorder()
			)

			lrw := &loggingResponseWriter{
				bytesWritten:             initialBW,
				downstreamResponseWriter: w,
			}

			n, err := lrw.Write([]byte(body))
			Expect(err).To(BeNil())
			Expect(n).To(Equal(bodyLen))

			Expect(lrw.bytesWritten).To(Equal(initialBW + bodyLen))

			Expect(w.Body.String()).To(Equal(body))
		})
	})

	Context("WriteHeader", func() {
		It("invokes downstream method and records the value", func() {
			const (
				sc = http.StatusTeapot
			)

			var (
				w = httptest.NewRecorder()
			)

			lrw := &loggingResponseWriter{
				statusCode:               http.StatusOK,
				downstreamResponseWriter: w,
			}
			lrw.WriteHeader(sc)

			Expect(lrw.statusCode).To(Equal(sc))

			Expect(w.Code).To(Equal(sc))
		})
	})

	Context("Hijack", func() {
		It("invokes downstream method when it implements http.Hijacker", func() {
			var (
				hj  = &hijackTester{}
				lrw = &loggingResponseWriter{
					downstreamResponseWriter: hj,
				}
				_, _, err = lrw.Hijack()
			)

			Expect(err).To(BeNil())
			Expect(hj.called).To(BeTrue())
		})
	})

	Context("Hijack", func() {
		It("returns error when downstreamResponseWriter does not implement http.Hijacker", func() {
			var (
				w   = httptest.NewRecorder()
				lrw = &loggingResponseWriter{
					downstreamResponseWriter: w,
				}
				_, _, err = lrw.Hijack()
			)

			Expect(err).NotTo(BeNil())
		})
	})
})

type hijackTester struct {
	called bool
}

func (hj *hijackTester) Header() http.Header {
	return http.Header{}
}

func (hj *hijackTester) Write([]byte) (int, error) {
	return 0, nil
}

func (hj *hijackTester) WriteHeader(int) {
}

func (hj *hijackTester) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj.called = true
	return nil, nil, nil
}
