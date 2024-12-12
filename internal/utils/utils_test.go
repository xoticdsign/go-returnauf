package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Unit тест для функции RandInt
func TestUnitRandInt(t *testing.T) {
	cases := []struct {
		name  string
		input int
	}{
		{
			name:  "general case",
			input: 40,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			mockSupport := new(Support)

			var gotInts []int
			var gotStrs []string

			for x := 0; x < 5; x++ {
				gotInt, gotStr := mockSupport.RandInt(cs.input)

				gotInts = append(gotInts, gotInt)
				gotStrs = append(gotStrs, gotStr)

				if x < 0 {
					assert.NotEqual(t, gotInt, gotInts[x-1])
					assert.NotEqual(t, gotStr, gotStrs[x-1])
				}
			}
		})
	}
}
