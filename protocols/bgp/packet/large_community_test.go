package packet

import (
	"errors"
	"testing"

	"fmt"
	"math"

	"strconv"

	"strings"

	"github.com/stretchr/testify/assert"
)

func TestParseLargeCommunityString(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected LargeCommunity
		err      error
	}{
		{
			name: "normal large community",
			in:   "(1,2,3)",
			expected: LargeCommunity{
				GlobalAdministrator: 1,
				DataPart1:           2,
				DataPart2:           3,
			},
			err: nil,
		},
		{
			name:     "too short community",
			in:       "(1,2)",
			expected: LargeCommunity{},
			err:      errors.New("can not parse large community 1,2"),
		},
		{
			name: "missing parentheses large community",
			in:   "1,2,3",
			expected: LargeCommunity{
				GlobalAdministrator: 1,
				DataPart1:           2,
				DataPart2:           3,
			},
			err: nil,
		},
		{
			name:     "malformed large community",
			in:       "[1,2,3]",
			expected: LargeCommunity{},
			err:      &strconv.NumError{Func: "ParseUint", Num: "[1", Err: strconv.ErrSyntax},
		},
		{
			name:     "missing digit",
			in:       "(,2,3)",
			expected: LargeCommunity{},
			err:      &strconv.NumError{Func: "ParseUint", Num: "", Err: strconv.ErrSyntax},
		},
		{
			name:     "too big global administrator",
			in:       fmt.Sprintf("(%d,1,2)", math.MaxInt64),
			expected: LargeCommunity{},
			err:      &strconv.NumError{Func: "ParseUint", Num: fmt.Sprintf("%d", math.MaxInt64), Err: strconv.ErrRange},
		},
		{
			name:     "too big data part 1",
			in:       fmt.Sprintf("(1,%d,2)", math.MaxInt64),
			expected: LargeCommunity{1, 0, 0},
			err:      &strconv.NumError{Func: "ParseUint", Num: fmt.Sprintf("%d", math.MaxInt64), Err: strconv.ErrRange},
		},
		{
			name:     "too big data part 2",
			in:       fmt.Sprintf("(1,2,%d)", math.MaxInt64),
			expected: LargeCommunity{1, 2, 0},
			err:      &strconv.NumError{Func: "ParseUint", Num: fmt.Sprintf("%d", math.MaxInt64), Err: strconv.ErrRange},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			com, err := ParseLargeCommunityString(test.in)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, test.expected, com)
		})
	}
}

func BenchmarkParseLargeCommunityString(b *testing.B) {
	for _, i := range []int{1, 2, 4, 8, 16, 32, 64} {
		str := getNNumbers(i)
		input := strings.Join([]string{str, str, str}, ",")
		b.Run(fmt.Sprintf("CommunitySize-%d-numbers", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				ParseLargeCommunityString(input)
			}
		})
	}
}

func getNNumbers(n int) (ret string) {
	var numbers string
	for i := 0; i < n; i++ {
		numbers += strconv.Itoa(i % 10)
	}
	return numbers
}
