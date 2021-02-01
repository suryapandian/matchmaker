package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/suryapandian/matchmaker/games"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	var testCases = []struct {
		desc               string
		expectedStatusCode int
	}{
		{
			"sanity",
			http.StatusOK,
		},
	}

	var dummyMatchMaker games.Matchmaker

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			a := assert.New(t)
			r := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()
			GetRouter(&dummyMatchMaker).ServeHTTP(w, r)
			response := w.Result()
			a.Equal(testCase.expectedStatusCode, response.StatusCode)
		})
	}

}
