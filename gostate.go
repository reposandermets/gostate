package gostate

import (
	"log"
	"time"
)

// ReadState structure type for reading the state
type ReadState struct {
	Key    string
	ChSend chan interface{}
}

// WriteState structure type for reading state
type WriteState struct {
	Key  string
	Data interface{}
}

type goState struct{}

// IGoState interface for the gostate
type IGoState interface {
	State(chWriteState <-chan WriteState, chReadState <-chan ReadState, chResetState chan bool, chQuitState chan bool)
}

var (
	// GS interface holder
	GS IGoState
	// Log ops, defaults false
	Log bool
)

func init() {
	GS = &goState{}
	Log = false
}

// State is core goroutine holding the state structure
func (gs *goState) State(chWriteState <-chan WriteState, chReadState <-chan ReadState, chResetState chan bool, chQuitState chan bool) {
	go func() {
		var state = make(map[string]interface{})
		for {
			select {
			case <-chQuitState:
				if Log {
					log.Println("StateHolder STOP")
				}
				return
			case <-chResetState:
				state = nil
				if Log {
					log.Println("StateHolder RESET")
				}
			case rState := <-chReadState:
				rState.ChSend <- state[rState.Key]
				if Log {
					log.Println("StateHolder READ")
				}
			case wState := <-chWriteState:
				state[wState.Key] = wState.Data
				if Log {
					log.Println("StateHolder Write")
				}
			default:
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
}
