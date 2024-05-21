package main

import (
	"errors"
	"fmt"
	"time"
)




// Команда для сохранения состояния данных
type Command interface {
	Execute()
}

type SaveCommand struct {
	pool       *AllPools
	schemaName string
	collection string
	key        string
	value      interface{}
}

func (c *SaveCommand) Execute() {
	pool, err := c.pool.GetPool("mainPool")
	if err != nil {
		fmt.Println("Ошибка при получении пула:", err)
		return
	}
	schema, err := pool.GetSchema(c.schemaName)
	if err != nil {
		fmt.Println("Ошибка при получении схемы:", err)
		return
	}
	coll, err := schema.GetCollection(c.collection)
	if err != nil {
		fmt.Println("Ошибка при получении коллекции:", err)
		return
	}
	coll.Update(c.key, c.value)
}

// Обработчик времени
type Handler interface {
	SetNext(handler Handler)
	HandleRequest(request string) (interface{}, error)
}

type TimeHandler struct {
	next Handler
	pool *AllPools
	time time.Time
}

func (th *TimeHandler) SetNext(handler Handler) {
	th.next = handler
}

func (th *TimeHandler) HandleRequest(request string) (interface{}, error) {
	if th.time.IsZero() {
		return nil, errors.New("Не указано время")
	}
	pool, err := th.pool.GetPool("mainPool")
	if err != nil {
		return nil, err
	}
	return th.getDataAtTime(pool, th.time)
}

func (th *TimeHandler) getDataAtTime(pool *Pool, t time.Time) (interface{}, error) {
	// Добавьте здесь реальную логику получения данных на указанное время
	return "Current data", nil
}
