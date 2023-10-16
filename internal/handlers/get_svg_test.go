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
		pathEncoded  string
		queryEncoded string
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
		{"happy path with path param", state{pathEncoded: "qlDQtVOo5AIEAAD__w=="}, want{200, "<?xml version=\"1.0\" encoding=\"utf-8\"?>"}},
		{"happy path with query param", state{queryEncoded: "qlDQtVOo5AIEAAD__w=="}, want{200, "<?xml version=\"1.0\" encoding=\"utf-8\"?>"}},
		{"fail bad request with path param", state{pathEncoded: "qlDQtVOo5AIEAAD__w==&"}, want{400, "Invalid Base64 data."}},
		{"fail bad request with query param", state{queryEncoded: "qlDQtVOo5AIEAAD__w==&"}, want{400, "Invalid Base64 data."}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Controller{}

			req := httptest.NewRequest(http.MethodGet, "/svg", nil)
			q := req.URL.Query()
			if tc.state.pathEncoded != "" {
				q.Set(":encodedD2", tc.state.pathEncoded)
			}
			if tc.state.queryEncoded != "" {
				q.Set("script", tc.state.queryEncoded)
			}
			req.URL.RawQuery = q.Encode()

			resp := httptest.NewRecorder()

			c.GetD2SVGHandler(resp, req)
			assert.Equal(t, tc.want.code, resp.Code)
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
