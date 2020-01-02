# gostate

Simple stateful goroutine

## Purpose

Key value state using idea from [Stateful Goroutines](https://gobyexample.com/stateful-goroutines).

## CLI

```sh

go get github.com/reposandermets/gostate

go test -cover -v -race -count=1

```

## Usage example

```go

	GS.Start()

	type person struct {
		name string
		age  int
	}

	// by reference, danger
	person1 := person{name: "foo", age: 99}
	GS.Write("personByRef", &person1)
	person1 = person{} // set empty
  resRead, _ := GS.Read("personByRef")

	fmt.Println(resRead.(*person).name) // outputs: ""

	// by value, safe
	person2 := person{name: "baz", age: 33}
	GS.Write("personByVal", person2)
	person2 = person{} // set empty
	resRead, _ = GS.Read("personByVal")
	fmt.Println(resRead.(person).name) // outputs: baz

	person3 := person{name: "One", age: 53}
	person4 := person{name: "Two", age: 51}
	persons := make([]person, 0)
	persons = append(persons, person3)
	persons = append(persons, person4)

	GS.Write("personsByVal", persons)
	resRead, _ = GS.Read("personsByVal")
	fmt.Println(resRead.([]person)[0].name) // outputs: One

	m := make(map[string][]person)
	m["data"] = persons
	GS.Write("personsMap", m)
	resRead, _ = GS.Read("personsMap")
	fmt.Println(resRead.(map[string][]person)["data"][0].name) // outputs: One

	GS.Stop() // graceful shutdown

```
