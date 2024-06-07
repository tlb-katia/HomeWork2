package billing

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Entity interface {
	SetType(interface{}, bool)
	SetValue(interface{}, bool)
	SetCreatedAt(interface{}, bool)
	SetID(interface{}, bool)
}

type Operation struct {
	Type      map[interface{}]bool `json:"type,omitempty"`
	Value     map[interface{}]bool `json:"value,omitempty"`
	ID        map[interface{}]bool `json:"id,omitempty"`
	CreatedAt map[interface{}]bool `json:"created_at,omitempty"`
}

type Root struct {
	Company   string               `json:"company"`
	Operation *Operation           `json:"operation,omitempty"`
	Type      map[interface{}]bool `json:"type,omitempty"`
	Value     map[interface{}]bool `json:"value,omitempty"`
	ID        map[interface{}]bool `json:"id,omitempty"`
	CreatedAt map[interface{}]bool `json:"created_at,omitempty"`
}

type OperationType string

const (
	Income  OperationType = "Income"
	Outcome OperationType = "Outcome"
	Plus    OperationType = "+"
	Minus   OperationType = "-"
)

// ValidOperationTypes holds all valid operation types
var ValidOperationTypes = map[OperationType]bool{
	Income:  true,
	Outcome: true,
	Plus:    true,
	Minus:   true,
}

// UnmarshalJSON десериализует JSON-данные в структуру Root.
// В случае некорректных данных возвращает ошибку.
func (r *Root) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("Root: invalid JSON: %w", err)
	}

	if company, ok := raw["company"].(string); ok {
		r.Company = company
	} else {
		return errors.New("invalid or missing company")
	}
	initializeMaps(r)

	if op, ok := raw["operation"].(map[string]interface{}); ok {
		parseEntity(r.Operation, &op)
	}
	parseEntity(r, &raw)

	return validateOperation(r)
}

// initializeMaps инициализирует все мапы в структуре Root и вложенной структуре Operation, если они не были инициализированы ранее.
// Если переданное значение r является nil, создается новая структура Root.
func initializeMaps(r *Root) {
	if r == nil {
		r = &Root{}
	}
	initializeStructMaps(reflect.ValueOf(r).Elem())

	if r.Operation == nil {
		r.Operation = &Operation{}
	}
	initializeStructMaps(reflect.ValueOf(r.Operation).Elem())
}

// initializeStructMaps инициализирует все мапы в структуре, представленной значением val.
// Функция используется для рекурсивной инициализации мап во всех вложенных структурах.
func initializeStructMaps(val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := field.Type()

		if field.Kind() == reflect.Map && field.IsNil() {
			field.Set(reflect.MakeMap(fieldType))
		}
	}
}

// parseEntity анализирует данные из map, содержащих информацию о сущности,
// и устанавливает соответствующие значения в переданном объекте Entity.
// Если данные отсутствуют или имеют неправильный формат, используются значения по умолчанию.
func parseEntity(entity Entity, data *map[string]interface{}) {
	if t, ok := (*data)["type"].(string); ok {
		parseType(&entity, t)
	} else {
		entity.SetType((*data)["type"], false)
	}

	if v, ok := (*data)["value"]; ok {
		parseValue(&entity, v)
	}

	if t, ok := (*data)["created_at"]; ok {
		parseCreatedAt(&entity, t)
	}

	if v, ok := (*data)["id"]; ok {
		parseID(&entity, v)
	}
}

// parseType парсит тип операции из строки и устанавливает его в объекте Entity.
// Если тип операции недопустим, устанавливает значение по умолчанию.
func parseType(entity *Entity, data string) {
	operationType := OperationType(data)
	if !ValidOperationTypes[operationType] {
		(*entity).SetType(operationType, false)
	} else {
		(*entity).SetType(operationType, true)
	}
}

// parseValue парсит тип операции из строки и устанавливает его в объекте Entity.
// Если тип операции недопустим, устанавливает значение по умолчанию.
func parseValue(entity *Entity, data interface{}) {
	switch data.(type) {
	case float64:
		(*entity).SetValue(int(data.(float64)), true)
	case string:
		i, err := strconv.Atoi(data.(string))
		if err == nil {
			(*entity).SetValue(i, true)
		} else {
			(*entity).SetValue(i, false)
		}
	case int:
		(*entity).SetValue(data, true)
	default:
		(*entity).SetValue(data, false)
	}
}

// parseCreatedAt парсит дату создания из интерфейса и устанавливает ее в объекте Entity.
// Если дата некорректна, устанавливает значение по умолчанию.
func parseCreatedAt(entity *Entity, data interface{}) {
	createdAt, err := time.Parse(time.RFC3339, data.(string))
	if err != nil {
		(*entity).SetCreatedAt(data.(string), false)
	} else {
		(*entity).SetCreatedAt(createdAt, true)
	}
}

// parseID парсит идентификатор из интерфейса и устанавливает его в объекте Entity.
// Если идентификатор некорректен, устанавливает значение по умолчанию.
func parseID(entity *Entity, data interface{}) {
	switch data.(type) {
	case string:
		(*entity).SetID(data, true)
	case int:
		(*entity).SetID(data, true)
	default:
		(*entity).SetID(data, true)
	}
}

// validateOperation проверяет корректность операции и обновляет структуру Root соответствующим образом.
// Если операция некорректна, возвращает ошибку.
func validateOperation(root *Root) error {
	var operation Operation

	if root.Operation != nil {
		operation = *root.Operation
	}

	for k, v := range root.Type {
		if v == true {
			operation.Type = root.Type
		} else if k != nil {
			operation.Type = root.Type
		}
	}
	for k, v := range root.Value {
		if v == true {
			operation.Value = root.Value
		} else if k != nil {
			operation.Value = root.Value
		}
	}

	for k, v := range root.CreatedAt {
		if v == true {
			operation.CreatedAt = root.CreatedAt
		} else if k != nil {
			operation.CreatedAt = root.CreatedAt
		}
	}
	for k, v := range root.ID {
		if v == true {
			operation.ID = root.ID
		} else if k != nil {
			operation.ID = root.ID
		}
	}

	root.Operation = &operation
	return nil
}

func (o *Operation) SetType(t interface{}, flag bool) {
	o.Type[t] = flag
}

func (o *Operation) SetValue(v interface{}, flag bool) {
	o.Value[v] = flag
}

func (o *Operation) SetCreatedAt(c interface{}, flag bool) {
	o.CreatedAt[c] = flag
}

func (o *Operation) SetID(id interface{}, flag bool) {
	o.ID[id] = flag
}

func (r *Root) SetType(t interface{}, flag bool) {
	r.Type[t] = flag
}

func (r *Root) SetValue(v interface{}, flag bool) {
	r.Value[v] = flag
}

func (r *Root) SetCreatedAt(c interface{}, flag bool) {
	r.CreatedAt[c] = flag
}

func (r *Root) SetID(id interface{}, flag bool) {
	r.ID[id] = flag
}
