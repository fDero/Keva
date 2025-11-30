package misc

import (
	"fmt"
	"os"
)

func FailOnError[T any](data T, err error) T {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return data
}

func FirstOfManyErrorsOrNone(elems []error) error {
	for _, elem := range elems {
		if elem != nil {
			return elem
		}
	}
	return nil
}

func DiscardReturn[T any](value T, err error) error {
	return err
}

func DiscardError[T any](value T, err error) T {
	return value
}
