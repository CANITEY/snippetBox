package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"caniteySnippetBox/internal/assert"
)

// func TestPing(t *testing.T) {
// 	t.Run("Ping", func(t *testing.T) {
// 		rr := httptest.NewRecorder()
// 		r, err := http.NewRequest(http.MethodGet, "/", nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		ping(rr, r)
// 		rs := rr.Result()

// 		assert.Equal(t, rs.StatusCode, http.StatusOK)
// 		defer rs.Body.Close()

// 		body, err := io.ReadAll(rs.Body)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		body = bytes.TrimSpace(body)
// 		assert.Equal(t, string(body), "OK")
// 	})
// }

func TestPingNew(t *testing.T) {
	t.Run("Handler Test", func(t *testing.T) {
		app := &application{
			errorLog: log.New(io.Discard, "", 0),
			infoLog: log.New(io.Discard, "", 0),
		}

		ts := httptest.NewTLSServer(app.routes())
		defer ts.Close()

		rs, err := ts.Client().Get(ts.URL+"/ping")
		if err != nil {
			t.Fatal(err)
		}


		assert.Equal(t, rs.StatusCode, http.StatusOK)
	})
}
