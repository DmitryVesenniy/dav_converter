package exception

type FewBytes struct{}

func (err FewBytes) Error() string {
	return "unable to read enough bytes"
}
