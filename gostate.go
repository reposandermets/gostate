package gostate

import (
	"log"
	"time"
)

// ReadState structure type for reading the state
type ReadState struct {
	Key            string
	ChReadResponse chan interface{}
}

// WriteState structure type for reading state
type WriteState struct {
	Key             string
	Data            interface{}
	ChWriteResponse chan interface{}
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
	// CyclePeriod frequency of listen loop
	CyclePeriod time.Duration
)

func init() {
	GS = &goState{}
	Log = false
	CyclePeriod = time.Millisecond
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
				rState.ChReadResponse <- state[rState.Key]
				if Log {
					log.Println("StateHolder READ")
				}
			case wState := <-chWriteState:
				state[wState.Key] = wState.Data
				wState.ChWriteResponse <- state[wState.Key]
				if Log {
					log.Println("StateHolder Write")
				}
			default:
				time.Sleep(CyclePeriod)
			}
		}
	}()
}
