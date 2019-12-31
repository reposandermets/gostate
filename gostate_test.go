package gostate

import (
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

func TestQuitState(t *testing.T) {
	chReadState := make(chan ReadState, 100)
	chWriteState := make(chan WriteState, 100)
	chResetState := make(chan bool)
	chQuitState := make(chan bool)
	GS.State(chWriteState, chReadState, chResetState, chQuitState)
	chQuitState <- true
}

func TestWriteReadStateUnbuffered(t *testing.T) {
	log.Println("Test write and read with unbuffered channels")
	Log = true
	chReadState := make(chan ReadState)
	chWriteState := make(chan WriteState)
	chResetState := make(chan bool)
	chQuitState := make(chan bool)
	GS.State(chWriteState, chReadState, chResetState, chQuitState)

	w := WriteState{
		Key:  "Foo",
		Data: "Bar",
	}
	chWriteState <- w
	r := ReadState{
		Key:    "Foo",
		ChSend: make(chan interface{}),
	}
	chReadState <- r
	select {
	case res := <-r.ChSend:
		if res != "Bar" {
			t.Errorf("wanted 'Bar' got:\n%v", res)
		}
	}
	chQuitState <- true
}
