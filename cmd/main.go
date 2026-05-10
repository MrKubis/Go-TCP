package main

import (
	"fmt"
	"net"
	"time"
)

type WorkerPool struct {
	tasks chan net.Conn
}

func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{tasks: make(chan net.Conn, 100)}
	for i := 0; i < size; i++ {
		go pool.worker()
	}
	return pool
}
func (p *WorkerPool) worker() {
	for conn := range p.tasks {
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		fmt.Printf("Recieved: %s", buffer[:n])
		conn.Write([]byte("Lubie placki\n"))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	// Initialize Worker pool
	pool := NewWorkerPool(10)

	fmt.Println("Server listening on :8080")
	// Equivalent to while(true)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err)
			continue
		}
		pool.tasks <- conn
	}
}
