package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn, redis *Redis) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			break
		}

		parts := strings.Split(input, " ")
		if len(parts) < 2 {
			conn.Write([]byte("Invalid command format.\n"))
			continue
		}

		command := parts[0]
		key := parts[1]
		value := interface{}(nil)
		if len(parts) > 2 {
			value = strings.Join(parts[2:], " ")
		}

		result := redis.HandleCommand(command, key, value)
		response, _ := json.Marshal(result)
		conn.Write(append(response, '\n'))
	}
}

func main() {
	redis := NewRedis()

	// Запускаем TCP-сервер на порту 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on :8080")

	// Основной цикл сервера для принятия подключений
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}

			fmt.Println("New connection from:", conn.RemoteAddr())
			go handleConnection(conn, redis)
		}
	}()

	// Основной цикл для чтения команд из консоли
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.ToLower(scanner.Text())

		if input == "exit" {
			break
		}

		parts := strings.Split(input, " ")
		if len(parts) < 2 {
			fmt.Println("Invalid command format.")
			continue
		}

		command := parts[0]
		key := parts[1]
		value := interface{}(nil)
		if len(parts) > 2 {
			value = strings.Join(parts[2:], " ")
		}

		result := redis.HandleCommand(command, key, value)
		fmt.Println(result)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
