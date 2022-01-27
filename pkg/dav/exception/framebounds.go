package exception

type FrameBounds struct{}

func (err FrameBounds) Error() string {
	return "frame out of bounds"
}
