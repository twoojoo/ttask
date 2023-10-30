package main

import "errors"


func main () {
	first := errors.New("first-error")
	second := errors.New("second-error")

	err := errors.Join(first, second)

	panic(err)
}