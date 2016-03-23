package controllers_test

import (
	"github.com/dugancathal/stuffs/controllers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"bytes"
)

var _ = Describe("EchoController", func() {
	It("returns empty if nothing has been set", func() {
		echoController := controllers.NewEchoController()

		req, _ := http.NewRequest("POST", "http://example.com/foo", bytes.NewBuffer([]byte(`bar`)))

		w := httptest.NewRecorder()

		echoController.HandleReq(w, req)
		Expect(w.Body.String()).To(Equal(``))
	})

	It("allows setting responses from the /set-response-for prefix", func() {
		echoController := controllers.NewEchoController()

		req, _ := http.NewRequest("PUT", "http://example.com/set-response-for/bar", bytes.NewBuffer([]byte(`foobar`)))
		sillyWriter := httptest.NewRecorder()

		echoController.HandleSetReq(sillyWriter, req)

		req2, _ := http.NewRequest("PUT", "http://example.com/bar", bytes.NewBuffer([]byte(``)))
		w := httptest.NewRecorder()

		echoController.HandleReq(w, req2)
		Expect(w.Body.String()).To(Equal(`foobar`))
	})

	It("does not pollute other responses", func() {
		echoController := controllers.NewEchoController()

		req, _ := http.NewRequest("PUT", "http://example.com/set-response-for/bar", bytes.NewBuffer([]byte(`foobar`)))
		barWriter := httptest.NewRecorder()

		echoController.HandleSetReq(barWriter, req)

		req2, _ := http.NewRequest("PUT", "http://example.com/set-response-for/foo", bytes.NewBuffer([]byte(`barfoo`)))
		fooWriter := httptest.NewRecorder()

		echoController.HandleSetReq(fooWriter, req2)

		req3, _ := http.NewRequest("PUT", "http://example.com/bar", bytes.NewBuffer([]byte(`whatever`)))
		w := httptest.NewRecorder()

		echoController.HandleReq(w, req3)
		Expect(w.Body.String()).To(Equal(`foobar`))
	})

	It("lets the user modify response", func() {
		echoController := controllers.NewEchoController()

		req, _ := http.NewRequest("PUT", "http://example.com/set-response-for/bar", bytes.NewBuffer([]byte(`foobar`)))
		barWriter := httptest.NewRecorder()

		echoController.HandleSetReq(barWriter, req)

		req2, _ := http.NewRequest("PUT", "http://example.com/set-response-for/bar", bytes.NewBuffer([]byte(`barfoo`)))
		fooWriter := httptest.NewRecorder()

		echoController.HandleSetReq(fooWriter, req2)

		req3, _ := http.NewRequest("PUT", "http://example.com/bar", bytes.NewBuffer([]byte(`whatever`)))
		w := httptest.NewRecorder()

		echoController.HandleReq(w, req3)
		Expect(w.Body.String()).To(Equal(`barfoo`))
	})

	It("handles long paths", func() {
		echoController := controllers.NewEchoController()

		req, _ := http.NewRequest("PUT", "http://example.com/set-response-for/foo/bar", bytes.NewBuffer([]byte(`foobar`)))
		barWriter := httptest.NewRecorder()

		echoController.HandleSetReq(barWriter, req)

		req2, _ := http.NewRequest("PUT", "http://example.com/set-response-for/bar", bytes.NewBuffer([]byte(`barfoo`)))
		fooWriter := httptest.NewRecorder()

		echoController.HandleSetReq(fooWriter, req2)

		req3, _ := http.NewRequest("PUT", "http://example.com/foo/bar", bytes.NewBuffer([]byte(`whatever`)))
		w := httptest.NewRecorder()

		echoController.HandleReq(w, req3)
		Expect(w.Body.String()).To(Equal(`foobar`))
	})

	It("cares about the HTTP Method", func() {
		echoController := controllers.NewEchoController()

		req, _ := http.NewRequest("PUT", "http://example.com/set-response-for/bar", bytes.NewBuffer([]byte(`foobar`)))
		sillyWriter := httptest.NewRecorder()

		echoController.HandleSetReq(sillyWriter, req)

		req2, _ := http.NewRequest("PATCH", "http://example.com/bar", bytes.NewBuffer([]byte(``)))
		w := httptest.NewRecorder()

		echoController.HandleReq(w, req2)
		Expect(w.Body.String()).To(Equal(``))
	})
})
