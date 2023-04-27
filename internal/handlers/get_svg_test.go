package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errFail = errors.New("unit test failure")

func TestController_GetD2SVGHandler(t *testing.T) {
	type state struct {
		encoded string
	}
	type want struct {
		code int
		err  string
	}
	testCases := []struct {
		name string
		state
		want
	}{
		{"happy path", state{encoded: "qlDQtVOo5AIEAAD__w=="}, want{200, "<?xml version=\"1.0\" encoding=\"utf-8\"?>"}},
		{"fail bad request", state{encoded: "qlDQtVOo5AIEAAD__w==&"}, want{400, "Invalid Base64 data."}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Controller{}

			req := httptest.NewRequest(http.MethodGet, "/svg", nil)
			q := req.URL.Query()
			q.Set(":encodedD2", tc.state.encoded)
			req.URL.RawQuery = q.Encode()
			resp := httptest.NewRecorder()

			c.GetD2SVGHandler(resp, req)
			assert.Equal(t, resp.Code, tc.want.code)
			if !strings.Contains(resp.Body.String(), tc.want.err) {
				t.Errorf(
					`response body "%s" does not contain "%s"`,
					resp.Body.String(),
					tc.want.err,
				)
			}

		})
	}
}
