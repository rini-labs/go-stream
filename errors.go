package stream

type Error int

const (
	ErrDone Error = iota
	ErrIteratorNotStarted
)

var (
	errorMessages = []string{
		"done",
		"iterator not started",
	}
)

func (e Error) Error() string {
	return errorMessages[e]
}
