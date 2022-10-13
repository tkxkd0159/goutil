package ds

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterList(t *testing.T) {
	l := NewIterList[int]()
	require.Equal(t, &IterList[int]{}, l)

	tcs := []struct {
		elem int
		idx  int
	}{
		{2, 1},
		{1, 0},
		{3, 0},
	}

	l.Append(1)
	l.Append(2)
	l.Append(3)
	require.Equal(t, len(tcs), l.Size())

	i := 0
	for e := range l.Iter() {
		require.Equal(t, l.Get(i), e)
		i++
	}

	for _, tc := range tcs {
		require.Equal(t, tc.elem, l.Get(tc.idx))
		l.Delete(tc.idx)
	}

	require.Equal(t, &IterList[int]{size: 0, items: []int{}}, l)
}
