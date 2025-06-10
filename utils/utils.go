package utils

import (
	"errors"
	"fmt"
	"strconv"
)

func CheckArgs(args []string) (int, error) {
	if len(args) == 0 {
		return -1, errors.New("count args must be one")
	} else if argInt, err := strconv.Atoi(args[0]); err != nil {
		fmt.Println("Arg must be a number")
		return -1, errors.New("count args must be one")
	} else if argInt < 0 {
		return -1, errors.New("count args must be greater than zero")
	} else {
		return argInt, nil
	}
}
