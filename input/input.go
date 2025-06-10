package input

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"worker-poll/models"
)

type InputManager struct {
	channels models.Channels
	commands chan models.Commands
	wg       *sync.WaitGroup
}

func NewInputManager(channels models.Channels, wg *sync.WaitGroup, commands chan models.Commands) *InputManager {
	return &InputManager{
		channels: channels,
		commands: commands,
		wg:       wg,
	}
}

func (im *InputManager) Start(commands []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Список команд:\n" +
		"list - список воркеров\n" +
		"addworker - добавить воркера\n" +
		"removeworker id - удалить воркера по его id")

	inputChan := make(chan string)
	go func() {
		line, err := reader.ReadString('\n')
		if err != nil {
			close(inputChan)
			return
		}
		inputChan <- line
	}()

	for {
		select {
		case <-im.channels.Interrupt.Done():
			im.wg.Done()
			fmt.Println("InputManager received interrupt")
			return
		case line, ok := <-inputChan:
			if !ok {
				return
			}
			args := strings.Fields(line)
			if len(args) == 0 {
				continue
			}
			if cmd := args[0]; slices.Contains(commands, cmd) {
				im.commands <- models.Commands{
					Cmd:  cmd,
					Args: args[1:],
				}
			} else {
				im.channels.InputData <- line
			}
		}
	}
}
