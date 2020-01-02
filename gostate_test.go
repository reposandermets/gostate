package gostate

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("start running tests")
	exitVal := m.Run()
	log.Println("tests done")

	os.Exit(exitVal)
}

func TestStateDirect(t *testing.T) {
	log.Println("Test write and read with unbuffered channels")

	Log = true

	chReadState := make(chan ReadState)
	chWriteState := make(chan WriteState)
	chQuitState := make(chan bool)

	GS.State(chWriteState, chReadState, chQuitState)

	w := WriteState{
		Key:             "Foo",
		Data:            "Bar",
		ChWriteResponse: make(chan interface{}),
	}
	chWriteState <- w
	select {
	case res := <-w.ChWriteResponse:
		if res != "Bar" {
			t.Errorf("wanted 'Bar' got:\n%v", res)
		}
	}
	r := ReadState{
		Key:            "Foo",
		ChReadResponse: make(chan interface{}),
	}
	chReadState <- r
	select {
	case res := <-r.ChReadResponse:
		if res != "Bar" {
			t.Errorf("wanted 'Bar' got:\n%v", res)
		}
	}
	chQuitState <- true
}

func TestErrorMessages(t *testing.T) {
	GS.Start()
	err := GS.Start()
	if err.Error() != "State already running, nothing to start" {
		t.Errorf("wanted 'State already running, nothing to start' got:\n%v", err.Error())
	}

	resWrite, _ := GS.Write("foo", "bar")
	if resWrite != "bar" {
		t.Errorf("wanted 'bar' got:\n%v", resWrite)
	}

	resRead, _ := GS.Read("foo")
	if resRead != "bar" {
		t.Errorf("wanted 'bar' got:\n%v", resRead)
	}

	GS.Stop()
	err = GS.Stop()
	if err.Error() != "State not running, nothing to stop" {
		t.Errorf("wanted 'State not running, nothing to stop' got:\n%v", err.Error())
	}

	_, err = GS.Write("foo", "bar")
	if err.Error() != "State not running, no write" {
		t.Errorf("wanted 'State not running, no write' got:\n%v", err.Error())
	}

	_, err = GS.Read("foo")
	if err.Error() != "State not running, no read" {
		t.Errorf("wanted 'State not running, no read' got:\n%v", err.Error())
	}

	err = GS.Restart()
	if err.Error() != "State not running, nothing to restart" {
		t.Errorf("wanted 'State not running, nothing to restart' got:\n%v", err.Error())
	}

	GS.Start()
	resWrite, _ = GS.Write("foo", "bar")
	if resWrite != "bar" {
		t.Errorf("wanted 'bar' got:\n%v", resWrite)
	}

	err = GS.Restart()
	if err != nil {
		t.Errorf("wanted err='nil' got:\n%v", err.Error())
	}
	resRead, _ = GS.Read("foo")
	if resRead != nil {
		t.Errorf("wanted resRead=nil got:\n%v", resRead)
	}
	GS.Stop()
}

func TestUsage(t *testing.T) {
	GS.Start()

	type person struct {
		name string
		age  int
	}

	// by reference, danger
	person1 := person{name: "foo", age: 99}
	GS.Write("personByRef", &person1)
	person1 = person{}
	resRead, _ := GS.Read("personByRef")
	fmt.Println(resRead.(*person).name) // outputs: ""
	if resRead.(*person).name != "" {
		t.Errorf("wanted '' got:\n%v", resRead)
	}

	// by value, safe
	person2 := person{name: "baz", age: 33}
	GS.Write("personByVal", person2)
	person2 = person{}
	resRead, _ = GS.Read("personByVal")
	fmt.Println(resRead.(person).name) // outputs: baz
	if resRead.(person).name != "baz" {
		t.Errorf("wanted baz got:\n%v", resRead)
	}

	person3 := person{name: "One", age: 53}
	person4 := person{name: "Two", age: 51}
	persons := make([]person, 0)
	persons = append(persons, person3)
	persons = append(persons, person4)

	GS.Write("personsByVal", persons)
	resRead, _ = GS.Read("personsByVal")
	fmt.Println(resRead.([]person)[0].name) // outputs: One
	if resRead.([]person)[0].name != "One" {
		t.Errorf("wanted One got:\n%v", resRead.([]person)[0].name)
	}

	m := make(map[string][]person)
	m["data"] = persons
	GS.Write("personsMap", m)
	resRead, _ = GS.Read("personsMap")
	fmt.Println(resRead.(map[string][]person)["data"][0].name) // outputs: One
	if resRead.(map[string][]person)["data"][0].name != "One" {
		t.Errorf("wanted One got:\n%v", resRead.(map[string][]person)["data"][0].name)
	}

	GS.Stop() // graceful shutdown
}
