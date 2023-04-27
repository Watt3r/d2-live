package urlenc

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errFail = errors.New("unit test failure")

func TestController_GetD2SVGHandler(t *testing.T) {
	type state struct {
		encoded string
	}
	type want struct {
		output string
		err    error
	}
	testCases := []struct {
		name string
		state
		want
	}{
		{"happy path", state{encoded: "qlDQtVOo5AIEAAD__w=="}, want{"x -> y\n", nil}},
		{"fail bad request", state{encoded: "qlDQtVOo5AIEAAD__w==&"}, want{"", base64.CorruptInputError(20)}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			decoded, err := Decode(tc.state.encoded)
			assert.Equal(t, decoded, tc.want.output)
			assert.Equal(t, err, tc.want.err)
		})
	}
}
