package models

import (
	"context"
)

type Channels struct {
	InputData chan string
	Interrupt context.Context
}

type Commands struct {
	Cmd  string
	Args []string
}
