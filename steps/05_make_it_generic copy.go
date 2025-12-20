package steps

import "fmt"

type handlerFunc func() error

func process(handler handlerFunc) error {
	err := handler()
	if err != nil {
		//which processor errored
		return fmt.Errorf("error processing: %w", err)
	}
	return nil
}
