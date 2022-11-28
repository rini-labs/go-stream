package stream

type ErrorCodes int

const (
	Done ErrorCodes = iota
)

var errorCodesMessages []string

func (e ErrorCodes) Error() string {
	return errorCodesMessages[e]
}

func init() {
	errorCodesMessages = []string{
		"Stream finished",
	}
}
