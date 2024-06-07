package calculated_balance

import (
	"HomeWork2/internal/lib/api/billing"
	"fmt"
	"reflect"
	"sort"
)

type CompanyBalance struct {
	Company           string        `json:"company"`
	ValidOperations   int           `json:"valid_operations_count,omitempty"`
	Balance           int           `json:"balance,omitempty"`
	InvalidOperations []interface{} `json:"invalid_operations,omitempty"`
}

// CountBalance вычисляет баланс каждой компании на основе списка Root.
// Возвращает список структур CompanyBalance.
func CountBalance(roots []billing.Root) []CompanyBalance {
	companies := make(map[string]CompanyBalance) // name = struct
	for _, root := range roots {
		t := reflect.TypeOf(root)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Name == "Value" {
				err := sumBalanceValues(&companies, root.Company, &root.Operation.Value)
				if err == nil {
					countValidOperations(&companies, root.Company)
				}
			}
		}
		distributeFields(reflect.ValueOf(root), &companies, root.Company)
	}
	ret := make([]CompanyBalance, 0, len(companies))
	recordAndSortFinalData(&ret, &companies)
	return ret
}

// checkValue проверяет наличие значения true в мапе Value и возвращает это значение как int.
// Если значение не найдено, возвращает ошибку.
func checkValue(m *map[interface{}]bool) (int, error) {
	for value, v := range *m {
		if v == true {
			return value.(int), nil
		}
	}
	return 0, fmt.Errorf("value not found in company map")
}

// sumBalanceValues суммирует значения в мапе с данными Values и добавляет их к балансу компании.
// Если компания не существует в мапе, создает новую запись.
func sumBalanceValues(companies *map[string]CompanyBalance, key string, m *map[interface{}]bool) error {
	valInt, err := checkValue(m)
	if err != nil {
		return fmt.Errorf("value not found in company map ", err)
	}
	companyStruct, exists := (*companies)[key]

	if !exists {
		companyStruct = CompanyBalance{Company: key, Balance: 0}
	}
	companyStruct.Balance += valInt
	(*companies)[key] = companyStruct

	return nil
}

// addInvalidOperations добавляет невалидные операции в список невалиндных операций компании.
// Если компания не существует в мапе, создает новую запись.
func addInvalidOperations(companies *map[string]CompanyBalance, key string, value interface{}) {
	companyStruct, exists := (*companies)[key]

	if !exists {
		companyStruct = CompanyBalance{Company: key, Balance: 0, InvalidOperations: make([]interface{}, 0)}
	}
	companyStruct.InvalidOperations = append(companyStruct.InvalidOperations, value)
	(*companies)[key] = companyStruct

}

// countValidOperations сумирует валидные операции
// Если компания не существует в мапе, создает новую запись.
func countValidOperations(companies *map[string]CompanyBalance, company string) {
	companyBalance, exists := (*companies)[company]
	if !exists {
		companyBalance = CompanyBalance{Company: company, Balance: 0, InvalidOperations: make([]interface{}, 0)}
	}
	companyBalance.ValidOperations += 1
	(*companies)[company] = companyBalance
}

// distributeFields распределяет значения мап в структуре v на соответствующие списки в map[string]CompanyBalance.
func distributeFields(v reflect.Value, companies *map[string]CompanyBalance, company string) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.Kind() == reflect.Map {
			for _, key := range field.MapKeys() {
				value := field.MapIndex(key)
				if value.Bool() == true {
					countValidOperations(companies, company)
				} else {
					if keyValue := key.Interface(); keyValue != nil && keyValue != "" {
						addInvalidOperations(companies, company, keyValue)
					}
				}
			}
		}
	}
}

// recordAndSortFinalData записывает данные из map[string]CompanyBalance
// в переданный список CompanyBalance и сортирует его по названию компании.
func recordAndSortFinalData(ret *[]CompanyBalance, companies *map[string]CompanyBalance) {
	for _, company := range *companies {
		*ret = append(*ret, company)
	}

	sort.SliceStable(*ret, func(i, j int) bool {
		return (*ret)[i].Company < (*ret)[j].Company
	})
}
