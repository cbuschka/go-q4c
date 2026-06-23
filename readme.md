# go-q4c - queries for collections in go

```go
type Person struct {
    Id   int
    Name string
}

persons := []Person{
    {Id: 1, Name: "Jane"},
    {Id: 2, Name: "Tarzan"},
}

selected = q4c.SelectFrom(persons).
	Where( (p Person) bool => p.Id == 1).
	ToSlice()
```

## License
Copyright (c) 2025 by the go-q4c maintainers.

[Apache License, Version 2.0](./license.txt)
