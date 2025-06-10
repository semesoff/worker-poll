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
		"add <n> - добавить воркера (n - кол-во)\n" +
		"remove <id> - удалить воркера (id - идентификатор)\n" +
		"count - кол-во воркеров\n" +
		"clear - удаление всех воркеров\n" +
		"exit - выход из программы\n" +
		": любое слово не из команд - входные данные (строки)")

	inputChan := make(chan string)
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				close(inputChan)
				return
			}
			inputChan <- line
		}
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
			if cmd := strings.ToLower(args[0]); slices.Contains(commands, cmd) {
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
