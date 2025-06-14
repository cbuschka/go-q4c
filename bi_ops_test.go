package q4c

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

	persons := []Person{
		{Id: 1, Name: "Jane"},
		{Id: 2, Name: "Tarzan"},
	}

	addresses := []Address{
		{PersonId: 1, Street: "Jane's home"},
		{PersonId: 2, Street: "Tarzan's home"},
	}

	idOfPerson := func(p Person) int {
		return p.Id
	}
	personIdOfAddress := func(a Address) int {
		return a.PersonId
	}
	result := NewBiSet[Person, Address]().SelectFrom(persons).
		Join(addresses).On(idOfPerson, personIdOfAddress).ToSlice()
	assert.Equal(t, result, persons)
}
