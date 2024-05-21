package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunCommand(pools *AllPools, command string) error {
	args := strings.Fields(command)
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	switch args[0] {
	case "add-pool":
		if len(args) < 2 {
			return fmt.Errorf("Недостаточно аргументов для команды add-pool.")
		}
		pools.AddPool(args[1])
	case "remove-pool":
		if len(args) < 2 {
			return fmt.Errorf("Недостаточно аргументов для команды remove-pool.")
		}
		pools.RemovePool(args[1])
	case "add-schema":
		if len(args) < 3 {
			return fmt.Errorf("Недостаточно аргументов для команды add-schema.")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		pool.AddSchema(args[2])
	case "remove-schema":
		if len(args) < 3 {
			return fmt.Errorf("Недостаточно аргументов для команды remove-schema.")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		pool.PopSchema(args[2])
	case "add-collection":
		if len(args) < 4 {
			return fmt.Errorf("Недостаточно аргументов для команды add-collection.")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		schema, err := pool.GetSchema(args[2])
		if err != nil {
			return err
		}
		collection := NewMapCollection()
		if err = schema.AddCollection(args[3], collection); err != nil {
			return err
		}
	case "remove-collection":
		if len(args) < 4 {
			return fmt.Errorf("Недостаточно аргументов для команды remove-collection.")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		schema, err := pool.GetSchema(args[2])
		if err != nil {
			return err
		}
		schema.PopCollection(args[3])
	case "add-record":
		if len(args) < 5 {
			return fmt.Errorf("Недостаточно аргументов для команды add-record.")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		schema, err := pool.GetSchema(args[2])
		if err != nil {
			return err
		}
		collection, err := schema.GetCollection(args[3])
		if err != nil {
			return err
		}
		if err := collection.Insert(args[4], args[5]); err != nil {
			return err
		}
	case "update-record":
		if len(args) < 6 {
			return fmt.Errorf("Недостаточно аргументов для команды update-record.")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		schema, err := pool.GetSchema(args[2])
		if err != nil {
			return err
		}
		collection, err := schema.GetCollection(args[3])
		if err != nil {
			return err
		}
		if err := collection.Update(args[4], args[5]); err != nil {
			return err
		}
	case "read-record":
		if len(args) < 5 {
			return fmt.Errorf("Недостаточно аргументов для команды read-record.")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		schema, err := pool.GetSchema(args[2])
		if err != nil {
			return err
		}
		collection, err := schema.GetCollection(args[3])
		if err != nil {
			return err
		}
		result, err := collection.Get(args[4])
		if err != nil {
			return err
		}
		fmt.Printf("key: %v, value: %v\n", args[4], result)
	case "delete-record":
		if len(args) < 5 {
			return fmt.Errorf("Недостаточно аргументов для команды delete-record.")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		schema, err := pool.GetSchema(args[2])
		if err != nil {
			return err
		}
		collection, err := schema.GetCollection(args[3])
		if err != nil {
			return err
		}
		if err := collection.Remove(args[4]); err != nil {
			return err
		}
	case "exit":
		return nil
	default:
		return fmt.Errorf("Неизвестная команда.")
	}
	return nil
}


func main() {

	cm := InitPool()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите команду:")
	for scanner.Scan() {
		command := scanner.Text()
		if command == "exit" {
			break
		} else if strings.Contains(command, ".txt") {
			file, err := os.Open(command)
			if err != nil {
				fmt.Println("Ошибка открытия файла:", err)
				continue
			}
			defer file.Close()

			fileScanner := bufio.NewScanner(file)
			for fileScanner.Scan() {
				cmd := fileScanner.Text()
				if err := RunCommand(cm, cmd); err != nil {
					fmt.Println("Ошибка выполнения команды:", err)
				}
			}

			if err := fileScanner.Err(); err != nil {
				fmt.Println("Ошибка чтения файла:", err)
			}
		} else {
			if err := RunCommand(cm, command); err != nil {
				fmt.Println("Ошибка выполнения команды:", err)
			}
		}
		fmt.Println("Введите команду:")
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка ввода:", err)
	}
}
