package gostate

import (
	"log"
	"os"
	"testing"
	"time"
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

func TestWriteReadStateBuffered(t *testing.T) {
	log.Println("Test write and read with buffered channels")

	Log = true

	chReadState := make(chan ReadState, 100)
	chWriteState := make(chan WriteState, 100)
	chResetState := make(chan bool)
	chQuitState := make(chan bool)

	GS.State(chWriteState, chReadState, chResetState, chQuitState)

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

	waitFlag := true
	for waitFlag {
		select {
		case res := <-r.ChReadResponse:
			waitFlag = false
			if res != "Bar" {
				t.Errorf("wanted 'Bar' got:\n%v", res)
			}

		default:
			time.Sleep(time.Millisecond)
		}
	}

	chQuitState <- true
}

func TestWriteReadStateBufferedOrder(t *testing.T) {
	log.Println("Test the order of writes")

	Log = true

	chReadState := make(chan ReadState, 100)
	chWriteState := make(chan WriteState, 100)
	chResetState := make(chan bool)
	chQuitState := make(chan bool)

	GS.State(chWriteState, chReadState, chResetState, chQuitState)

	w1 := WriteState{
		Key:             "Foo",
		Data:            "Bar1",
		ChWriteResponse: make(chan interface{}),
	}
	w2 := WriteState{
		Key:             "Foo",
		Data:            "Bar2",
		ChWriteResponse: make(chan interface{}),
	}
	w3 := WriteState{
		Key:             "Foo",
		Data:            "Bar3",
		ChWriteResponse: make(chan interface{}),
	}

	r := ReadState{
		Key:            "Foo",
		ChReadResponse: make(chan interface{}, 1),
	}
	chWriteState <- w1
	select {
	case res := <-w1.ChWriteResponse:
		if res != "Bar1" {
			t.Errorf("wanted 'Bar1' got:\n%v", res)
		}
	}
	chWriteState <- w2
	select {
	case res := <-w2.ChWriteResponse:
		if res != "Bar2" {
			t.Errorf("wanted 'Bar2' got:\n%v", res)
		}
	}
	chWriteState <- w3
	select {
	case res := <-w3.ChWriteResponse:
		if res != "Bar3" {
			t.Errorf("wanted 'Bar3' got:\n%v", res)
		}
	}
	chReadState <- r

	waitFlag := true

	for waitFlag {
		select {
		case res := <-r.ChReadResponse:
			waitFlag = false

			if res != "Bar3" {
				t.Errorf("wanted 'Bar3' got:\n%v", res)
			}

		default:
			time.Sleep(time.Millisecond)
		}
	}

	chQuitState <- true
}

func TestResetState(t *testing.T) {
	log.Println("Test the order of writes")

	Log = true

	chReadState := make(chan ReadState, 100)
	chWriteState := make(chan WriteState, 100)
	chResetState := make(chan bool)
	chQuitState := make(chan bool)

	GS.State(chWriteState, chReadState, chResetState, chQuitState)

	w1 := WriteState{
		Key:             "Foo",
		Data:            "Bar1",
		ChWriteResponse: make(chan interface{}),
	}

	r := ReadState{
		Key:            "Foo",
		ChReadResponse: make(chan interface{}, 1),
	}
	chWriteState <- w1
	select {
	case res := <-w1.ChWriteResponse:
		if res != "Bar1" {
			t.Errorf("wanted 'Bar1' got:\n%v", res)
		}
	}

	chResetState <- true

	chReadState <- r

	waitFlag := true

	for waitFlag {
		select {
		case res := <-r.ChReadResponse:
			waitFlag = false

			if res != nil {
				t.Errorf("wanted 'nil' got:\n%v", res)
			}

		default:
			time.Sleep(time.Millisecond)
		}
	}

	chQuitState <- true
}
