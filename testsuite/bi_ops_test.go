package testsuite

import (
	"testing"

	"github.com/cbuschka/go-q4c"
	"github.com/stretchr/testify/require"
)

func TestBiOps(t *testing.T) {
	type Person struct {
		Id   int
		Name string
	}

	type Address struct {
		PersonId int
		Street   string
	}

	idOfPerson := q4c.KeyFunc[Person, interface{}](func(p Person) interface{} {
		return p.Id
	})
	personIdOfAddress := q4c.KeyFunc[Address, interface{}](func(a Address) interface{} {
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

	noPairs := []q4c.Pair[Person, Address]{}

	noFilter := func(Person, Address) bool {
		return true
	}

	tests := []struct {
		name      string
		persons   []Person
		addresses []Address
		filter    func(Person, Address) bool
		expected  []q4c.Pair[Person, Address]
	}{
		{
			name:      "inner join, empy left",
			persons:   noPersons,
			addresses: allAddresses,
			filter:    noFilter,
			expected:  noPairs,
		},
		{
			name:      "inner join, empy right",
			persons:   allPersons,
			addresses: noAddresses,
			filter:    noFilter,
			expected:  noPairs,
		},
		{
			name:      "inner join, single, no match",
			persons:   []Person{harry},
			addresses: []Address{shack},
			filter:    noFilter,
			expected:  noPairs,
		},
		{
			name:      "inner join, single, one to one",
			persons:   []Person{harry},
			addresses: []Address{dormitory},
			filter:    noFilter,
			expected:  []q4c.Pair[Person, Address]{{Element1: harry, Element2: dormitory}},
		},
		{
			name:      "inner join, two, one to one",
			persons:   []Person{harry, ron},
			addresses: []Address{burrow, dormitory},
			filter:    noFilter,
			expected: []q4c.Pair[Person, Address]{{Element1: harry, Element2: dormitory},
				{Element1: ron, Element2: burrow}},
		},
		{
			name:      "inner join, one, one to may",
			persons:   []Person{harry},
			addresses: []Address{privetDrive, dormitory},
			filter:    noFilter,
			expected: []q4c.Pair[Person, Address]{{Element1: harry, Element2: privetDrive},
				{Element1: harry, Element2: dormitory}},
		},
		{
			name:      "inner join, all",
			persons:   allPersons,
			addresses: allAddresses,
			filter:    noFilter,
			expected: []q4c.Pair[Person, Address]{
				{Element1: harry, Element2: dormitory},
				{Element1: harry, Element2: privetDrive},
				{Element1: ron, Element2: burrow},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := q4c.NewBiSet[Person, Address]().SelectFrom(tt.persons).
				Join(tt.addresses).On(idOfPerson, personIdOfAddress).Where(tt.filter).ToSlice()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)

			collected := make([]q4c.Pair[Person, Address], 0)
			for pair, err := range q4c.NewBiSet[Person, Address]().SelectFrom(tt.persons).
				Join(tt.addresses).On(idOfPerson, personIdOfAddress).Where(tt.filter).Stream() {
				require.NoError(t, err)
				collected = append(collected, pair)
			}
			require.Equal(t, tt.expected, collected)

			first, found, err := q4c.NewBiSet[Person, Address]().SelectFrom(tt.persons).
				Join(tt.addresses).On(idOfPerson, personIdOfAddress).Where(tt.filter).First()
			require.NoError(t, err)
			require.Equal(t, len(tt.expected) > 0, found)
			if found {
				require.Equal(t, tt.expected[0], first)
			}

		})
	}

}

func TestLeftOuterJoin(t *testing.T) {
	type Person struct {
		Id   int
		Name string
	}

	type Address struct {
		PersonId int
		Street   string
	}

	idOfPerson := q4c.KeyFunc[Person, interface{}](func(p Person) interface{} {
		return p.Id
	})
	personIdOfAddress := q4c.KeyFunc[Address, interface{}](func(a Address) interface{} {
		return a.PersonId
	})

	harry := Person{Id: 1, Name: "Harry"}
	ron := Person{Id: 2, Name: "Ron"}
	hermione := Person{Id: 3, Name: "Hermione"}
	allPersons := []Person{harry, ron, hermione}
	noPersons := []Person{}

	dormitory := Address{PersonId: harry.Id, Street: "Gryffindor Dormitory"}
	privetDrive := Address{PersonId: harry.Id, Street: "4 Privet Drive"}
	burrow := Address{PersonId: ron.Id, Street: "The Burrow"}
	shack := Address{PersonId: -1, Street: "The Shrieking Shack"}
	allAddresses := []Address{dormitory, privetDrive, burrow, shack}
	noAddresses := []Address{}

	noPairs := []q4c.Pair[Person, Address]{}
	emptyAddress := Address{}

	noFilter := func(Person, Address) bool {
		return true
	}

	tests := []struct {
		name      string
		persons   []Person
		addresses []Address
		filter    func(Person, Address) bool
		expected  []q4c.Pair[Person, Address]
	}{
		{
			name:      "left join, empty left",
			persons:   noPersons,
			addresses: allAddresses,
			filter:    noFilter,
			expected:  noPairs,
		},
		{
			name:      "left join, empty right keeps left with zero value right",
			persons:   []Person{harry, ron},
			addresses: noAddresses,
			filter:    noFilter,
			expected: []q4c.Pair[Person, Address]{
				{Element1: harry, Element2: emptyAddress},
				{Element1: ron, Element2: emptyAddress},
			},
		},
		{
			name:      "left join, single, no match keeps left",
			persons:   []Person{hermione},
			addresses: allAddresses,
			filter:    noFilter,
			expected:  []q4c.Pair[Person, Address]{{Element1: hermione, Element2: emptyAddress}},
		},
		{
			name:      "left join, one to many emits each right match",
			persons:   []Person{harry},
			addresses: []Address{privetDrive, dormitory},
			filter:    noFilter,
			expected: []q4c.Pair[Person, Address]{
				{Element1: harry, Element2: privetDrive},
				{Element1: harry, Element2: dormitory},
			},
		},
		{
			name:      "left join, all keeps unmatched left rows",
			persons:   allPersons,
			addresses: allAddresses,
			filter:    noFilter,
			expected: []q4c.Pair[Person, Address]{
				{Element1: harry, Element2: dormitory},
				{Element1: harry, Element2: privetDrive},
				{Element1: ron, Element2: burrow},
				{Element1: hermione, Element2: emptyAddress},
			},
		},
		{
			name:      "left join, filtered after join",
			persons:   allPersons,
			addresses: allAddresses,
			filter: func(person Person, address Address) bool {
				return person == hermione || address.Street == "The Burrow"
			},
			expected: []q4c.Pair[Person, Address]{
				{Element1: ron, Element2: burrow},
				{Element1: hermione, Element2: emptyAddress},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := q4c.NewBiSet[Person, Address]().SelectFrom(tt.persons).
				LeftOuterJoin(tt.addresses).On(idOfPerson, personIdOfAddress).Where(tt.filter).ToSlice()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)

			collected := make([]q4c.Pair[Person, Address], 0)
			for pair, err := range q4c.NewBiSet[Person, Address]().SelectFrom(tt.persons).
				LeftOuterJoin(tt.addresses).On(idOfPerson, personIdOfAddress).Where(tt.filter).Stream() {
				require.NoError(t, err)
				collected = append(collected, pair)
			}
			require.Equal(t, tt.expected, collected)

			first, found, err := q4c.NewBiSet[Person, Address]().SelectFrom(tt.persons).
				LeftOuterJoin(tt.addresses).On(idOfPerson, personIdOfAddress).Where(tt.filter).First()
			require.NoError(t, err)
			require.Equal(t, len(tt.expected) > 0, found)
			if found {
				require.Equal(t, tt.expected[0], first)
			}
		})
	}
}
