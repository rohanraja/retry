package retry

import (
	"errors"
	"github.com/fatih/color"
	"time"
)

var Delay time.Duration
var NumRetries int
var PanicEnabled bool

func init() {
	Delay = 2 * time.Second
	NumRetries = 15
	PanicEnabled = false
}

// Do accepts a function argument that returns an error. It will keep executing this
// function NumRetries times until no error is returned.
// Optionally pass another function as an argument that it will execute before retrying
// in case an error is returned from the first argument fuction.
func Do(args ...interface{}) error {

	if len(args) == 0 {
		panic("Wrong number of arguments")
	}

	task, ok := args[0].(func() error)

	beforeRetry := func() {}
	if len(args) > 1 {
		beforeRetry, ok = args[1].(func())
	}

	if ok == false {
		panic("Wrong Type of Arguments given to retry.Do")
	}

	retries := NumRetries

	var atleastOneError error
	atleastOneError = nil

	err := errors.New("Non-Nil error")
	for retries > 0 && err != nil {

		err = task()

		if err != nil {
			atleastOneError = err
			color.Magenta("\nError: %v\nRetrying #%d after %v", err, (NumRetries - retries + 1), Delay)
			time.Sleep(Delay)
			beforeRetry()
		}

		retries = retries - 1
	}

	if err != nil {

		color.Red("\nError even after %d retries:\n%v", NumRetries, err)
		if PanicEnabled == true {
			panic(err)
		}
		return err
	}

	if atleastOneError != nil {
		color.Cyan("\nRecovered from error: %v in %d tries\n", atleastOneError, (NumRetries - retries - 1))
	}

	return err

}
