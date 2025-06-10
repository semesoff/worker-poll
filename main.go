package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"worker-poll/input"
	"worker-poll/models"
	"worker-poll/utils"
	"worker-poll/worker"
)

type Manager struct {
	hashWorkers  map[int]chan struct{} // worker's id
	countWorkers int
	lastWorkerID int
	channels     models.Channels
	commands     chan models.Commands
	wg           *sync.WaitGroup
}

func NewManager(channels models.Channels, wg *sync.WaitGroup) *Manager {
	return &Manager{
		hashWorkers:  make(map[int]chan struct{}),
		countWorkers: 0,
		lastWorkerID: 0,
		channels:     channels,
		commands:     make(chan models.Commands, 1),
		wg:           wg,
	}
}

func (m *Manager) Start(commands []string, interrupt chan os.Signal) {
	for {
		select {
		case <-m.channels.Interrupt.Done():
			close(m.commands)
			m.wg.Done()
			fmt.Println("Manager received interrupt")
			return
		case v, ok := <-m.commands:
			if ok {
				switch v.Cmd {
				case commands[0]:
					m.ListWorkers()
				case commands[1]:
					if argInt, err := utils.CheckArgs(v.Args); err == nil {
						m.AddWorkers(argInt)
					} else {
						fmt.Println(err)
					}
				case commands[2]:
					if argInt, err := utils.CheckArgs(v.Args); err == nil {
						m.RemoveWorker(argInt)
					} else {
						fmt.Println(err)
					}
				case commands[3]:
					m.RemoveAllWorkers()
				case commands[4]:
					interrupt <- os.Interrupt
				}
			}
		}
	}
}

// AddWorkers - добавление воркера
func (m *Manager) AddWorkers(n int) {
	for i := 0; i < n; i++ {
		m.lastWorkerID++
		m.countWorkers++
		cancelChan := make(chan struct{}, 1)
		m.hashWorkers[m.lastWorkerID] = cancelChan
		m.wg.Add(1)
		go worker.StartWorker(&m.channels, cancelChan, m.wg, m.lastWorkerID)
		fmt.Printf("Added worker #%d\n", m.lastWorkerID)
	}
}

// RemoveWorker - удаление воркера по id
func (m *Manager) RemoveWorker(id int) {
	if cancelChan, ok := m.hashWorkers[id]; ok {
		delete(m.hashWorkers, id)
		cancelChan <- struct{}{}
		m.countWorkers--
	} else {
		fmt.Printf("Not found worker with id: %d\n", id)
	}
}

// RemoveAllWorkers - удаление всех воркеров
func (m *Manager) RemoveAllWorkers() {
	if m.countWorkers == 0 {
		fmt.Println("No workers to remove")
		return
	}
	for id := range m.hashWorkers {
		m.RemoveWorker(id)
	}
}

// ListWorkers - выводит список воркеров
func (m *Manager) ListWorkers() {
	if m.countWorkers == 0 {
		fmt.Println("No workers found")
		return
	}
	fmt.Println("List workers:")
	for i := range m.hashWorkers {
		fmt.Printf("worker #%d\n", i)
	}
}

func main() {
	const size = 10 // размер канала для строк (входных данных)
	channels := models.Channels{
		InputData: make(chan string, size), // канал строк (входные данные)
		Interrupt: context.Background(),    // канал завершения работы программы (Ctrl + C)
	}
	commands := []string{"list", "add", "remove", "clear", "exit"}
	wg := &sync.WaitGroup{}
	wg.Add(3) // Manager, InputManager, Interrupt Catcher

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	channels.Interrupt = ctx

	// Initialize and start Manager
	manager := NewManager(channels, wg)
	go manager.Start(commands, interrupt)

	// Initialize and start Input Manager
	inputManager := input.NewInputManager(channels, wg, manager.commands)
	go inputManager.Start(commands)

	// Start interrupt catcher
	go func() {
		for {
			<-interrupt
			cancel()
			wg.Done()
			return
		}
	}()

	wg.Wait()
	close(channels.InputData)
	fmt.Println("Program is shutting down")
}
