package problems

import "fmt"

type Problem struct {
	id      string // unique identifier for the problem
	code    string // the code of the problem
	message string // the message to display to the user
}

func New(serviceCode, problemCode, message string) Problem {
	return Problem{
		id:      fmt.Sprintf("%s.%s", serviceCode, problemCode),
		code:    problemCode,
		message: message,
	}
}

func (p Problem) Error() string {
	return p.String()
}

func (p Problem) String() string {
	return p.message
}
