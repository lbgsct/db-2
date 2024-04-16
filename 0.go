package main

import (
	"fmt"
	"os"
)

// PoolManager отвечает за управление пулами данных.
type PoolManager struct {
	pools map[string]struct{}
}

func NewPoolManager() *PoolManager {
	return &PoolManager{pools: make(map[string]struct{})}
}

// AddPool добавляет новый пул данных.
func (pm *PoolManager) AddPool(poolName string) {
	if _, ok := pm.pools[poolName]; !ok {
		pm.pools[poolName] = struct{}{}
		fmt.Printf("Pool '%s' добавлен.\n", poolName)
	} else {
		fmt.Printf("Pool '%s' существует.\n", poolName)
	}
}

// RemovePool удаляет существующий пул данных.
func (pm *PoolManager) RemovePool(poolName string) {
	if _, ok := pm.pools[poolName]; ok {
		delete(pm.pools, poolName)
		fmt.Printf("Pool '%s' удален .\n", poolName)
	} else {
		fmt.Printf("Pool '%s' не существует.\n", poolName)
	}
}

// SchemaManager отвечает за управление схемами данных.
type SchemaManager struct {
	schemas map[string]map[string]struct{} // pool -> schema -> {}
}

func NewSchemaManager() *SchemaManager {
	return &SchemaManager{schemas: make(map[string]map[string]struct{})}
}

// AddSchema добавляет новую схему данных в указанный пул данных.
func (sm *SchemaManager) AddSchema(poolName, schemaName string) {
	if _, ok := sm.schemas[poolName]; !ok {
		sm.schemas[poolName] = make(map[string]struct{})
	}
	if _, ok := sm.schemas[poolName][schemaName]; !ok {
		sm.schemas[poolName][schemaName] = struct{}{}
		fmt.Printf("Schema '%s' добавлена в pool '%s'.\n", schemaName, poolName)
	} else {
		fmt.Printf("Schema '%s' уже существует в pool '%s'.\n", schemaName, poolName)
	}
}

// RemoveSchema удаляет существующую схему данных из указанного пула данных.
func (sm *SchemaManager) RemoveSchema(poolName, schemaName string) {
	if schemas, ok := sm.schemas[poolName]; ok {
		if _, ok := schemas[schemaName]; ok {
			delete(sm.schemas[poolName], schemaName)
			fmt.Printf("Schema '%s' удалена из pool '%s'.\n", schemaName, poolName)
			return
		}
	}
	fmt.Printf("Schema '%s' не существует в pool '%s'.\n", schemaName, poolName)
}

// CollectionManager отвечает за управление коллекциями данных.
type CollectionManager struct {
	collections map[string]map[string]map[string]struct{} // pool -> schema -> collection -> {}
}

func NewCollectionManager() *CollectionManager {
	return &CollectionManager{collections: make(map[string]map[string]map[string]struct{})}
}

// AddCollection добавляет новую коллекцию данных в указанную схему данных в указанный пул данных.
func (cm *CollectionManager) AddCollection(poolName, schemaName, collectionName string) {
	if _, ok := cm.collections[poolName]; !ok {
		cm.collections[poolName] = make(map[string]map[string]struct{})
	}
	if _, ok := cm.collections[poolName][schemaName]; !ok {
		cm.collections[poolName][schemaName] = make(map[string]struct{})
	}
	if _, ok := cm.collections[poolName][schemaName][collectionName]; !ok {
		cm.collections[poolName][schemaName][collectionName] = struct{}{}
		fmt.Printf("Collection '%s' добавлена в schema '%s' в pool '%s'.\n", collectionName, schemaName, poolName)
	} else {
		fmt.Printf("Collection '%s' уже существует в schema '%s' в pool '%s'.\n", collectionName, schemaName, poolName)
	}
}

// RemoveCollection удаляет существующую коллекцию данных из указанной схемы данных в указанном пуле данных.
func (cm *CollectionManager) RemoveCollection(poolName, schemaName, collectionName string) {
	if collections, ok := cm.collections[poolName]; ok {
		if _, ok := collections[schemaName]; ok {
			if _, ok := collections[schemaName][collectionName]; ok {
				delete(cm.collections[poolName][schemaName], collectionName)
				fmt.Printf("Collection '%s' удалена из schema '%s' в pool '%s'.\n", collectionName, schemaName, poolName)
				return
			}
		}
	}
	fmt.Printf("Collection '%s' не существует в schema '%s' в pool '%s'.\n", collectionName, schemaName, poolName)
}

// CommandProcessor отвечает за обработку команд из файла.
type CommandProcessor struct {
	poolManager       *PoolManager
	schemaManager     *SchemaManager
	collectionManager *CollectionManager
}

func NewCommandProcessor(pm *PoolManager, sm *SchemaManager, cm *CollectionManager) *CommandProcessor {
	return &CommandProcessor{
		poolManager:       pm,
		schemaManager:     sm,
		collectionManager: cm,
	}
}

// ProcessCommands читает команды из файла и вызывает соответствующие методы управления данными.
func (cp *CommandProcessor) ProcessCommands(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var command, arg1, arg2 string
	for {
		_, err := fmt.Fscanf(file, "%s %s %s\n", &command, &arg1, &arg2)
		if err != nil {
			break // reached end of file
		}
		switch command {
		case "add_pool":
			cp.poolManager.AddPool(arg1)
		case "remove_pool":
			cp.poolManager.RemovePool(arg1)
		case "add_schema":
			cp.schemaManager.AddSchema(arg1, arg2)
		case "remove_schema":
			cp.schemaManager.RemoveSchema(arg1, arg2)
		case "add_collection":
			cp.collectionManager.AddCollection(arg1, arg2, arg2)
		case "remove_collection":
			cp.collectionManager.RemoveCollection(arg1, arg2, arg2)
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
	return nil
}