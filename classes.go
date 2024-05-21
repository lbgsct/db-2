package main

import (
	"errors"
	"fmt"
)

// Интерфейс для ассоциативного контейнера который производит операции над коллекцией.
type Collection interface {
	Insert(key string, value interface{}) error
	Get(key string) (interface{}, error)
	GetRange(minValue, maxValue string) ([]string, error)
	Update(key string, value interface{}) error
	Remove(key string) error
}

// Пример реализации интерфейса Collection на основе map.
type MapCollection struct {
	data map[string]interface{}
}

func NewMapCollection() *MapCollection {
	return &MapCollection{
		data: make(map[string]interface{}),
	}
}

func (mc *MapCollection) Insert(key string, value interface{}) error {
	if _, exists := mc.data[key]; exists {
		return errors.New("Элемент с таким ключом уже существует!")
	}
	mc.data[key] = value
	fmt.Println("Элемент успешно добавлен с ключом", key)
	return nil
}

func (mc *MapCollection) Get(key string) (interface{}, error) {
	value, exists := mc.data[key]
	if !exists {
		return nil, errors.New("Элемент не найден!")
	}
	return value, nil
}

func (mc *MapCollection) GetRange(minValue, maxValue string) ([]string, error) {
	var result []string
	for key := range mc.data {
		if key >= minValue && key <= maxValue {
			result = append(result, key)
		}
	}
	return result, nil
}

func (mc *MapCollection) Update(key string, value interface{}) error {
	if _, exists := mc.data[key]; !exists {
		return errors.New("Элемент не найден!")
	}
	mc.data[key] = value
	fmt.Println("Значение элемента с ключом", key, "успешно обновлено.")
	return nil
}

func (mc *MapCollection) Remove(key string) error {
	if _, exists := mc.data[key]; !exists {
		return errors.New("Элемент не найден!")
	}
	delete(mc.data, key)
	return nil
}

type Pool struct {
	schema map[string]*Schema
}

type AllPools struct {
	pools map[string]*Pool
}

func InitPool() *AllPools {
	return &AllPools{
		pools: make(map[string]*Pool),
	}
}

func (pools *AllPools) AddPool(name string) {
    if _, exists := pools.pools[name]; exists {
        fmt.Println("Пул с именем", name, "уже существует.")
    } else {
        pools.pools[name] = NewPool()
		fmt.Println("Добавлен пул с именем", name)
    }
}

func (pools *AllPools) RemovePool(name string) {
    if _, exists := pools.pools[name]; exists {
        delete(pools.pools, name)
        fmt.Println("Пул с именем", name, "удален.")
    } else {
        fmt.Println("Пул с именем", name, "не существует.")
    }
}


func (pools *AllPools) GetPool(name string) (*Pool, error) {
	returnEl, ok := pools.pools[name]
	if !ok {
		return nil, errors.New("Элемент не найден!")
	}
	return returnEl, nil
}

func NewPool() *Pool {
	return &Pool{
		schema: make(map[string]*Schema),
	}
}

func (pool *Pool) GetSchema(schemaName string) (*Schema, error) {
	returnEl, ok := pool.schema[schemaName]
	if !ok {
		return nil, errors.New("Элемент не найден!")
	}
	return returnEl, nil
}

func (pool *Pool) AddSchema(name string) {

    pool.schema[name] = InitSchema()
    fmt.Println("Схема с именем", name, "добавлена в пул.")
}


func (pool *Pool) PopSchema(name string) {
	delete(pool.schema, name)
}

type Schema struct {
	collection map[string]Collection
}

func InitSchema() *Schema {
	return &Schema{
		collection: make(map[string]Collection),
	}
}

func (schema *Schema) GetCollection(name string) (Collection, error) {
	returnEl, ok := schema.collection[name]
	if !ok {
		return nil, errors.New("Элемент не найден!")
	}
	return returnEl, nil
}

func (schema *Schema) AddCollection(name string, collection Collection) error {
	if _, exists := schema.collection[name]; exists {
		return errors.New("Коллекция с таким именем уже существует!")
	}
	schema.collection[name] = collection
	fmt.Println("Коллекция с именем", name, "добавлена в схему.")
	return nil
}

func (schema *Schema) PopCollection(name string) {
	delete(schema.collection, name)
}
