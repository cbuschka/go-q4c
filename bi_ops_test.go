package q4c

import (
	"testing"

	"github.com/cbuschka/go-q4c/types"
	"github.com/stretchr/testify/require"
)

func TestBiOps_JoinToSlice(t *testing.T) {
	type Person struct {
		Id   int
		Name string
	}

	type Address struct {
		PersonId int
		Street   string
	}

	idOfPerson := types.KeyFunc[Person, interface{}](func(p Person) interface{} {
		return p.Id
	})
	personIdOfAddress := types.KeyFunc[Address, interface{}](func(a Address) interface{} {
		return a.PersonId
	})

	harry := Person{Id: 1, Name: "Harry"}
	ron := Person{Id: 2, Name: "Ron"}
	hermione := Person{Id: 3, Name: "Hermione"}
	sirius := Person{Id: 4, Name: "Sirius"}
	allPersons := []Person{harry, ron, hermione, sirius}
	noPersons := []Person{}

	dormitory := Address{PersonId: harry.Id, Street: "Gryffindor Dormitory"}
	privetDrive := Address{PersonId: harry.Id, Street: "4 Privet Drive"}
	burrow := Address{PersonId: ron.Id, Street: "The Burrow"}
	shack := Address{PersonId: -1, Street: "The Shrieking Shack"}
	allAddresses := []Address{dormitory, privetDrive, burrow, shack}
	noAddresses := []Address{}

	noPairs := []types.Pair[Person, Address]{}

	tests := []struct {
		name      string
		persons   []Person
		addresses []Address
		expected  []types.Pair[Person, Address]
	}{
		{
			name:      "inner join, empy left",
			persons:   noPersons,
			addresses: allAddresses,
			expected:  noPairs,
		},
		{
			name:      "inner join, empy right",
			persons:   allPersons,
			addresses: noAddresses,
			expected:  noPairs,
		},
		{
			name:      "inner join, single, no match",
			persons:   []Person{harry},
			addresses: []Address{shack},
			expected:  noPairs,
		},
		{
			name:      "inner join, single, one to one",
			persons:   []Person{harry},
			addresses: []Address{dormitory},
			expected:  []types.Pair[Person, Address]{{Element1: harry, Element2: dormitory}},
		},
		{
			name:      "inner join, two, one to one",
			persons:   []Person{harry, ron},
			addresses: []Address{burrow, dormitory},
			expected: []types.Pair[Person, Address]{{Element1: harry, Element2: dormitory},
				{Element1: ron, Element2: burrow}},
		},
		{
			name:      "inner join, one, one to may",
			persons:   []Person{harry},
			addresses: []Address{privetDrive, dormitory},
			expected: []types.Pair[Person, Address]{{Element1: harry, Element2: privetDrive},
				{Element1: harry, Element2: dormitory}},
		},
		{
			name:      "inner join, all",
			persons:   allPersons,
			addresses: allAddresses,
			expected: []types.Pair[Person, Address]{
				{Element1: harry, Element2: dormitory},
				{Element1: harry, Element2: privetDrive},
				{Element1: ron, Element2: burrow},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBiSet[Person, Address]().SelectFrom(tt.persons).
				Join(tt.addresses).On(idOfPerson, personIdOfAddress).ToSlice()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}

}
