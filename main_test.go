package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/leplasmo/ascetic"
)

// func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
// 	req, _ := http.NewRequest(method, path, nil)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)
// 	return w
// }

func TestTodoCRUD(t *testing.T) {
	th := newTodoHandler()
	t.Run("Retrieve non-existing ID", func(t *testing.T) {
		// w := performRequest(router, "GET", "/todos/-2")
		// assert.Equal(t, http.StatusBadRequest, w.Code)
		// assert.Equal(t, "{\"error\":\"Record not found!\"}", w.Body.String())
		wr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		th.get(wr, req)
		if wr.Code != http.StatusOK {
			t.Errorf("got HTTP status code %d, expected 200", wr.Code)
		}
	})
}
