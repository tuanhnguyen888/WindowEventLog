package main

import (
	. "golang.org/x/sys/windows/svc/eventlog"
)

func main() {
	name := "Log test"
	l, err := Open(name)
	if err != nil {
		panic(err)
	}
	defer func(src string) {
		err := Remove(src)
		if err != nil {

		}
	}(name) // clear log between runs

	err = l.Info(800, "This is log test") // event id 200 == default event id 200 == "The code segment cannot be greater than or equal to 64K."
	if err != nil {
		panic(err)
	}
}
