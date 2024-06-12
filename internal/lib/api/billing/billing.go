package billing

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"unicode"
)

type Operation struct {
	Type              CustomType      `json:"type,omitempty"`
	Value             CustomValue     `json:"value,omitempty"`
	ID                CustomID        `json:"id,omitempty"`
	CreatedAt         CustomCreatedAt `json:"created_at,omitempty"`
	InvalidOperations []interface{}
	ValidOperations   int
}

type Root struct {
	Company   CustomName      `json:"company"`
	Operation *Operation      `json:"operation,omitempty"`
	Type      CustomType      `json:"type,omitempty"`
	Value     CustomValue     `json:"value,omitempty"`
	ID        CustomID        `json:"id,omitempty"`
	CreatedAt CustomCreatedAt `json:"created_at,omitempty"`
}

var invalidOperations = make(map[string][]interface{}) // companyName : slice
var validOperations = make(map[string]int)             // companyName : slice
var Name string                                        //company name for the same root

type CustomName struct {
	Name string
}
type CustomValue struct {
	Int interface{}
}

type CustomID struct {
	String interface{}
}

type CustomCreatedAt struct {
	CreatedAt interface{}
}

type CustomType struct {
	Type interface{}
}

const (
	Income  = "income"
	Outcome = "outcome"
	Plus    = "+"
	Minus   = "-"
)

// ValidOperationTypes holds all valid operation types
var ValidOperationTypes = map[string]bool{
	Income:  true,
	Outcome: true,
	Plus:    true,
	Minus:   true,
}

func (r *Root) UnmarshalJSON(data []byte) error {
	// Создаем временную структуру без метода UnmarshalJSON
	temp := struct {
		Company   CustomName      `json:"company"`
		Operation *Operation      `json:"operation,omitempty"`
		Type      CustomType      `json:"type,omitempty"`
		Value     CustomValue     `json:"value,omitempty"`
		ID        CustomID        `json:"id,omitempty"`
		CreatedAt CustomCreatedAt `json:"created_at,omitempty"`
	}{}

	// Десериализуем данные в временную структуру
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Копируем данные из временной структуры в оригинальную структуру
	r.Company = temp.Company
	r.Type = temp.Type
	r.Value = temp.Value
	r.ID = temp.ID
	r.CreatedAt = temp.CreatedAt

	if r.Operation == nil {
		r.Operation = &Operation{}
	}
	r.Operation.ValidOperations = validOperations[Name]
	r.Operation.InvalidOperations = invalidOperations[Name]

	// Здесь вы можете добавить дополнительную логику по обработке InvalidOperations и ValidOperations
	// Например, если вам нужно проверить корректность операции и т.д.

	return ValidateOperation(r)
}

func (cv *CustomValue) UnmarshalJSON(data []byte) error {
	if data[0] == '"' && data[len(data)-1] == '"' {
		if _, err := dataParsing(data[1 : len(data)-1]); err == nil {
			sumValidValues()
			if err := json.Unmarshal(data[1:len(data)-1], &cv.Int); err != nil {
				return fmt.Errorf("Failed to unmarshal custom value: %s", err)
			}
		}
	} else if _, err := dataParsing(data); err == nil {
		sumValidValues()
		if err := json.Unmarshal(data, &cv.Int); err != nil {
			return fmt.Errorf("Failed to unmarshal custom value: %s", err)

		}
	}
	return nil
}

func dataParsing(data []byte) (bool, error) {
	for i := 0; i < len(data); i++ {
		if i == 0 && data[i] == '-' {
			continue
		}
		if !unicode.IsDigit(rune(data[i])) {
			collectInvalidValues(string(data))
			return true, fmt.Errorf("Invalid custom value: %s", string(data))
		}
	}
	return checkFloat(string(data)), nil
}

func checkFloat(data string) bool {
	if _, err := strconv.ParseFloat(data, 64); err != nil {
		collectInvalidValues(data)
		return false
	}
	return true
}

func (cid *CustomID) UnmarshalJSON(data []byte) error {
	if data[0] == '"' && data[len(data)-1] == '"' {
		sumValidValues()
		if err := json.Unmarshal(data, &cid.String); err != nil {
			return fmt.Errorf("Failed to unmarshal custom id: %s", err)
		}
	} else {
		flag, err := dataParsing(data)
		if err == nil && flag == false {
			sumValidValues()
			if err := json.Unmarshal(data, &cid.String); err != nil {
				return fmt.Errorf("Failed to unmarshal custom value: %s", err)

			}
		} else if flag == true {
			collectInvalidValues(string(data))
		}
	}

	return nil
}

func (tc *CustomCreatedAt) UnmarshalJSON(data []byte) error {
	_, err := time.Parse(time.RFC3339, string(data[1:len(data)-1]))
	if err != nil {
		collectInvalidValues(string(data))
	} else {
		sumValidValues()
		if err := json.Unmarshal(data, &tc.CreatedAt); err != nil {
			return fmt.Errorf("Failed to unmarshal custom time: %s", err)
		}
	}
	return nil
}

func (ct *CustomType) UnmarshalJSON(data []byte) error {
	if !ValidOperationTypes[string(data[1:len(data)-1])] {
		collectInvalidValues(string(data))
	} else {
		sumValidValues()
		if err := json.Unmarshal(data, &ct.Type); err != nil {
			return fmt.Errorf("Failed to unmarshal custom type: %s", err)
		}
	}
	return nil
}

func (cn *CustomName) UnmarshalJSON(data []byte) error {
	Name = string(data)
	sumValidValues()
	if err := json.Unmarshal(data, &cn.Name); err != nil {
		return fmt.Errorf("failed to unmarshal name: %w", err)
	}
	return nil
}

func collectInvalidValues(data interface{}) {
	_, exists := invalidOperations[Name]
	if !exists {
		invalidOperations[Name] = make([]interface{}, 0)
	}
	invalidOperations[Name] = append(invalidOperations[Name], data)
}

func sumValidValues() {
	_, exists := validOperations[Name]
	if !exists {
		validOperations[Name] = 0
	}
	validOperations[Name]++
}

func ValidateOperation(root *Root) error {
	var operation Operation

	if root.Operation != nil {
		operation = *root.Operation
	}
	//else {
	//	root.Operation = &Operation{}
	//}

	if root.Type.Type != nil {
		operation.Type.Type = root.Type.Type
	}

	if root.Value.Int != nil {
		operation.Value.Int = root.Value.Int
	}

	if root.CreatedAt.CreatedAt != nil {
		operation.CreatedAt.CreatedAt = root.CreatedAt.CreatedAt
	}

	if root.ID.String != nil {
		operation.ID.String = root.ID.String
	}

	root.Operation = &operation
	return nil
}
