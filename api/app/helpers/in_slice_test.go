package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInSliceWithString(t *testing.T) {
	data := []string{"ONE", "TWO", "FOUR"}

	testCases := []struct {
		toFind        string
		shouldBeFound bool
	}{
		{toFind: "ONE", shouldBeFound: true},
		{toFind: "TWO", shouldBeFound: true},
		{toFind: "THREE", shouldBeFound: false},
		{toFind: "FOUR", shouldBeFound: true},
	}

	for i, c := range testCases {
		found, err := InSlice(c.toFind, data)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, c.shouldBeFound, found, "Failed test %d", i)
	}
}
