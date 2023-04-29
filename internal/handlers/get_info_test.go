package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestController_GetInfoHandler(t *testing.T) {
	type state struct {
	}
	type want struct {
		code     int
		jsonResp string
	}
	testCases := []struct {
		name string
		state
		want
	}{
		{"happy path", state{}, want{200, "dev"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Controller{
				Version: "dev",
			}

			req := httptest.NewRequest(http.MethodGet, "/info", nil)
			resp := httptest.NewRecorder()

			c.GetInfoHandler(resp, req)
			assert.Equal(t, resp.Code, tc.want.code)
			if !strings.Contains(resp.Body.String(), tc.want.jsonResp) {
				t.Errorf(
					`response body "%s" does not contain "%s"`,
					resp.Body.String(),
					tc.want.jsonResp,
				)
			}

		})
	}
}
