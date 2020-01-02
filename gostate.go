package gostate

import (
	"errors"
	"log"
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

type goState struct {
}

// IGoState interface for the gostate
type IGoState interface {
	// State is core goroutine holding the state structure.
	State(chWriteState <-chan WriteState, chReadState <-chan ReadState, chQuitState chan bool)
	// Start state if state  isn't running
	Start() error
	// Stop state if state is running
	Stop() error
	// Restart state if state is running
	Restart() error
	// Write value to state if state is running
	Write(key string, data interface{}) (responseData interface{}, Error error)
	// Read value from state if state is running
	Read(key string) (responseData interface{}, Error error)
}

var (
	// GS interface holder
	GS IGoState
	// Log ops, defaults false
	Log bool
	// stateRunning status of state
	stateRunning bool
	// chReadState reading state
	chReadState chan ReadState
	// chWriteState writing state
	chWriteState chan WriteState
	// chQuitState quitting state
	chQuitState chan bool
)

func init() {
	GS = &goState{}
	Log = false
	stateRunning = false
	chReadState = make(chan ReadState)
	chWriteState = make(chan WriteState)
	chQuitState = make(chan bool)
}

// State the core goroutine holding the state structure
func (gs *goState) State(chWrite <-chan WriteState, chRead <-chan ReadState, chQuit chan bool) {

	go func() {
		var state = make(map[string]interface{})
		for {
			select {
			case <-chQuit:
				if Log {
					log.Println("StateHolder STOP")
				}
				return
			case rState := <-chRead:
				rState.ChReadResponse <- state[rState.Key]
				if Log {
					log.Println("StateHolder READ")
				}
			case wState := <-chWrite:
				state[wState.Key] = wState.Data
				wState.ChWriteResponse <- state[wState.Key]
				if Log {
					log.Println("StateHolder Write")
				}
			}
		}
	}()
}

// Start state if state  isn't running
func (gs *goState) Start() error {
	if stateRunning {
		if Log {
			log.Println("State already running, nothing to start")
		}
		return errors.New("State already running, nothing to start")
	}

	GS.State(chWriteState, chReadState, chQuitState)
	stateRunning = true
	return nil
}

// Stop state if state is running
func (gs *goState) Stop() error {
	if !stateRunning {
		if Log {
			log.Println("State not running, nothing to stop")
		}
		return errors.New("State not running, nothing to stop")
	}
	chQuitState <- true
	stateRunning = false
	return nil
}

// Restart state if state is running
func (gs *goState) Restart() error {
	if !stateRunning {
		if Log {
			log.Println("State not running, nothing to restart")
		}
		return errors.New("State not running, nothing to restart")
	}

	chQuitState <- true
	stateRunning = false
	GS.State(chWriteState, chReadState, chQuitState)
	stateRunning = true
	return nil
}

// Write value to state if state is running
func (gs *goState) Write(key string, data interface{}) (responseData interface{}, Error error) {
	if !stateRunning {
		if Log {
			log.Println("State not running, no write")
		}
		return responseData, errors.New("State not running, no write")
	}

	w := WriteState{
		Key:             key,
		Data:            data,
		ChWriteResponse: make(chan interface{}),
	}
	chWriteState <- w
	select {
	case responseData = <-w.ChWriteResponse:
		return responseData, nil
	}
}

// Read value from state if state is running
func (gs *goState) Read(key string) (responseData interface{}, Error error) {
	if !stateRunning {
		if Log {
			log.Println("State not running, no read")
		}
		return responseData, errors.New("State not running, no read")
	}
	r := ReadState{
		Key:            key,
		ChReadResponse: make(chan interface{}),
	}
	chReadState <- r
	select {
	case responseData = <-r.ChReadResponse:
		return responseData, nil
	}
}
