package main

import (
	"fmt"

	"github.com/marianolg/enums"
)

const (
	pending int = iota
	processing
	done
	cancelled
)

func main() {
	statuses := enums.New(
		pending,
		processing,
		done,
		cancelled,
	)

	if statuses.IsValid(1) {
		fmt.Println(1, "It's safe to accept this status input!")
	}

	if !statuses.IsValid(5) {
		fmt.Println(5, "Not a valid status - this should fail!")
	}

	if statuses.IsAnyValid(1, 5) {
		fmt.Println([]int{1, 5}, "Well at least one of them it's ok")
	}

	if !statuses.IsAnyValid(4, 5) {
		fmt.Println([]int{4, 5}, "Not anymore...")
	}

	if !statuses.AreAllValid(1, 5) {
		fmt.Println([]int{1, 5}, "Yikes! I was hoping all of them to be statuses")
	}

	if statuses.AreAllValid(1, 2) {
		fmt.Println([]int{1, 2}, "Now we are talking!")
	}

	var isValidStatus func(int) bool
	statuses.SetTypedIsValid(&isValidStatus)

	if isValidStatus(1) {
		fmt.Println(1, "Can't send wrong type no more")
	}

	statusesConvert := enums.NewConvert(
		pending,
		processing,
		done,
		cancelled,
	)

	if statusesConvert.IsValid(1.0) {
		fmt.Println(1.0, "Got float from client - still can tell that's a valid status")
	}
}
