package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStrToDate(t *testing.T) {
	tests := []struct {
		input string
		ans   string
	}{
		{
			input: "1990-12-30",
			ans:   "1990-12-30",
		}, {
			input: "1994-1-3",
			ans:   "1994-01-03",
		}, {
			input: "2001-3-29",
			ans:   "2001-03-29",
		}, {
			input: "2004-12-1",
			ans:   "2004-12-01",
		},
	}

	for _, test := range tests {
		res := ConvertStrToDate(test.input)

		date := res.Format("2006-01-02")
		assert.Equal(t, test.ans, date)
		t.Log(res)
	}
}
