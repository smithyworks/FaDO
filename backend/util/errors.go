package util

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type ServerError struct {
	Lines []string
}

func (se *ServerError) Error() string {
	return fmt.Sprintf("Error: %v", strings.Join(se.Lines, "\n          "))
}

func (se *ServerError) Warning() string {
	return fmt.Sprintf("Warning: %v", strings.Join(se.Lines, "\n          "))
}

func PrintWarning(err error) error {
	if err == nil { return nil }

	_, file, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf("%v:%v", file, line)
	switch v := err.(type) {
	case *ServerError:
		v.Lines = append(v.Lines, msg)
		log.Println(v.Warning())
		return v
	default:
		se := &ServerError{Lines: []string{err.Error(), msg}}
		log.Println(se.Warning())
		return se
	}
}

func ProcessErr(err error, skips ...int) error {
	if err == nil { return nil }

	s := 1
	if skips != nil && len(skips) > 0 { s = skips[0] }
	_, file, line, _ := runtime.Caller(s)
	msg := fmt.Sprintf("%v:%v", file, line)
	switch v := err.(type) {
	case *ServerError:
		v.Lines = append(v.Lines, msg)
		return v
	default:
		return &ServerError{Lines: []string{err.Error(), msg}}
	}
}

func PrintErr(err error) error {
	if err = ProcessErr(err, 2); err != nil { log.Println(err.Error()) }
	return err
}