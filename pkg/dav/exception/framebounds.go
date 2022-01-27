package exception

import "fmt"

type FrameBounds struct{}

func (err FrameBounds) Error() string {
	return "frame out of bounds"
}

var ErrorStopIterator error = fmt.Errorf("stop iter")
