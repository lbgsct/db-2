package main

import (
	"errors"
	"fmt"
	"time"
)

// Интерфейс для ассоциативного контейнера который производит операции над коллекцией.
type Tree interface {
	Insert(key string, value interface{}) error
	Get(key string) (interface{}, error)
	GetRange(minValue, maxValue string) ([]string, error)
	Update(key string, value interface{}) error
	Remove(key string) error
}

type TreeCollection struct {
	tree Tree
}

func NewTreeCollection(treeType string) *TreeCollection {
	var tree Tree
	switch treeType {
	case "avl":
		tree = Tree(NewAVLTree())
	/*case "redblack":
		tree = NewRedBlackTree()*/
	default:
		tree = NewAVLTree() // По умолчанию используем AVL-дерево
	}
	return &TreeCollection{tree: tree}
}

func (tc *TreeCollection) Insert(key string, value interface{}) error {
	return tc.tree.Insert(key, value)
}

func (tc *TreeCollection) Get(key string) (interface{}, error) {
	return tc.tree.Get(key)
}

func (tc *TreeCollection) GetRange(minValue, maxValue string) ([]string, error) {
	return tc.tree.GetRange(minValue, maxValue)
}

func (tc *TreeCollection) Update(key string, value interface{}) error {
	return tc.tree.Update(key, value)
}

func (tc *TreeCollection) Remove(key string) error {
	return tc.tree.Remove(key)
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

/*func (mc *MapCollection) Get(key string) (interface{}, error) {
	value, exists := mc.data[key]
	if !exists {
		return nil, errors.New("Элемент не найден!")
	}
	return value, nil
}*/
func (mc *MapCollection) Get(key string, time time.Time, pools *AllPools) (interface{}, error) {
	timeHandler := &TimeHandler{pool: pools, time: time}
	return timeHandler.HandleRequest(key)
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

func (mc *MapCollection) Update(key string, value interface{}, pools *AllPools) error {
	saveCmd := &SaveCommand{
		pool:       pools,
		schemaName: "mainSchema",
		collection: "mainCollection",
		key:        key,
		value:      value,
	}
	saveCmd.Execute()
	
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
	if pool, exists := pools.pools[name]; exists {
		// Удаляем все схемы и их коллекции
		for schemaName := range pool.schema {
			schema := pool.schema[schemaName]
			// Удаляем все коллекции в схеме
			for collectionName := range schema.collection {
				schema.RemoveCollection(collectionName)
			}
			// Удаляем схему
			pool.RemoveSchema(schemaName)
		}
		// Удаляем пул
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
    fmt.Print("Схема с именем ", name, " добавлена в пул ")
}


func (pool *Pool) RemoveSchema(name string) {
	if schema, exists := pool.schema[name]; exists {
		// Удаляем все коллекции в схеме
		for collectionName := range schema.collection {
			schema.RemoveCollection(collectionName)
		}
		// Удаляем схему из пула
		delete(pool.schema, name)
		fmt.Println("Схема с именем", name, "удалена из пула")
	} else {
		fmt.Println("Схема с именем", name, "не найдена в пуле")
	}
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
	fmt.Print("Коллекция с именем ", name, " добавлена в схему ")
	return nil
}

func (schema *Schema) RemoveCollection(name string) {
	delete(schema.collection, name)
	fmt.Print("Коллекция с именем ", name, " удалена из схемы ")
}
