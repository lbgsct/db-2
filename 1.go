package main

import (
	"fmt"
	"os"
	"strings"
	"bufio"
)

func main() {
	if len(os.Args) < 2 {
		interactiveMode()
	} else {
		fileName := os.Args[1]
		if fileName == "exit" {
			return
		}
		poolManager := NewPoolManager()
		schemaManager := NewSchemaManager()
		collectionManager := NewCollectionManager()
		commandProcessor := NewCommandProcessor(poolManager, schemaManager, collectionManager)

		if err := processCommandsFromFile(commandProcessor, fileName); err != nil {
			fmt.Println("Ошибки из файла:", err)
		}
	}
}

func interactiveMode() {
	poolManager := NewPoolManager()
	schemaManager := NewSchemaManager()
	collectionManager := NewCollectionManager()
	commandProcessor := NewCommandProcessor(poolManager, schemaManager, collectionManager)

	for {
		fmt.Print("Введите команду или имя файла ('exit' для выхода): ")
		var input string
		fmt.Scanln(&input)

		if input == "exit" {
			break
		}

		if _, err := os.Stat(input); err == nil {
			if err := processCommandsFromFile(commandProcessor, input); err != nil {
				fmt.Println("Ошибка выполнения из файла:", err)
			}
			continue
		}

		processCommand(commandProcessor, input)
	}
}

func processCommand(commandProcessor *CommandProcessor, input string) {
	parts := strings.Fields(input)
	if len(parts) < 1 {
		fmt.Println("Неправильный ввод. Пожалуйста, введите команду или имя файла.")
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "add_pool":
		if len(args) < 1 {
			fmt.Println("Неправильный ввод. Укажите название pool")
			return
		}
		commandProcessor.poolManager.AddPool(args[0])
	case "remove_pool":
		if len(args) < 1 {
			fmt.Println("Неправильный ввод. Укажите название pool")
			return
		}
		commandProcessor.poolManager.RemovePool(args[0])
	case "add_schema":
		if len(args) < 2 {
			fmt.Println("Неправильный ввод. Укажите название pool и название schema.")
			return
		}
		commandProcessor.schemaManager.AddSchema(args[0], args[1])
	case "remove_schema":
		if len(args) < 2 {
			fmt.Println("Неправильный ввод. Укажите название pool и название schema.")
			return
		}
		commandProcessor.schemaManager.RemoveSchema(args[0], args[1])
	case "add_collection":
		if len(args) < 3 {
			fmt.Println("Неправильный ввод. Укажите название pool, название schema и название collection.")
			return
		}
		commandProcessor.collectionManager.AddCollection(args[0], args[1], args[2])
	case "remove_collection":
		if len(args) < 3 {
			fmt.Println("Неправильный ввод. Укажите название pool, название schema и название collection.")
			return
		}
		commandProcessor.collectionManager.RemoveCollection(args[0], args[1], args[2])
	default:
		fmt.Println("Unknown command:", command)
	}
}

func processCommandsFromFile(commandProcessor *CommandProcessor, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		processCommand(commandProcessor, scanner.Text())
	}
	return scanner.Err()
}
