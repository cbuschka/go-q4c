package testsuite

import (
	"github.com/cbuschka/go-q4c"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestUniOps(t *testing.T) {
	type Person struct {
		Id   int
		Name string
	}

	jane := Person{Id: 1, Name: "Jane"}
	tarzan := Person{Id: 2, Name: "Tarzan"}
	allPersons := []Person{
		jane,
		tarzan,
	}

	noFilter := func(Person) bool {
		return true
	}

	tests := []struct {
		name     string
		persons  []Person
		filter   func(Person) bool
		expected []Person
	}{
		{
			name:     "all, no filter",
			persons:  allPersons,
			filter:   noFilter,
			expected: allPersons,
		},
		{
			name:    "all, filtered",
			persons: allPersons,
			filter: func(p Person) bool {
				return p == jane
			},
			expected: []Person{jane},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := q4c.SelectFrom(allPersons).Where(tt.filter).ToSlice()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)

			collected := make([]Person, 0)
			for person, err := range q4c.SelectFrom(allPersons).Where(tt.filter).Stream() {
				require.NoError(t, err)
				collected = append(collected, person)
			}
			require.Equal(t, tt.expected, collected)

			first, found, err := q4c.SelectFrom(allPersons).Where(tt.filter).First()
			require.NoError(t, err)
			require.Equal(t, len(tt.expected) > 0, found)
			if found {
				require.Equal(t, tt.expected[0], first)
			}
		})
	}
}
