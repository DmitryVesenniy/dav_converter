package worker

import (
	"runtime"
	"time"
)

type protocol struct {
	id   int
	err  error
	done bool
}

// Worker выполняет задачи в отдельных горутинах
func Worker(
	defs []func() error,
	countProcess int,
	stop <-chan struct{},
	done chan<- struct{},
	errorCh chan<- error,
) {
	countActive := 0
	pipeline := make(chan protocol)
	indexProcess := 1

LOOP:
	for {
		select {
		case <-stop:
			break LOOP

		case state := <-pipeline:
			if state.err != nil {
				errorCh <- state.err
			}
			countActive--

		default:
			if len(defs) > 0 {
				if countActive < countProcess {
					countDelta := min(countProcess-countActive, len(defs))
					for i := 0; i < countDelta; i++ {
						_func := defs[0]
						defs = defs[1:]
						go run(indexProcess, pipeline, _func)
						indexProcess++
						countActive++
					}
				}
			} else {
				if countActive == 0 {
					break LOOP
				}
			}

			runtime.Gosched()
			time.Sleep(100 * time.Millisecond)
		}
	}
	done <- struct{}{}
}

func run(id int, pipe chan<- protocol, fn func() error) {
	state := protocol{
		id: id,
	}
	err := fn()
	state.done = true
	if err != nil {
		state.err = err
	}
	pipe <- state
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
