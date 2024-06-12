package calculated_balance

import (
	"HomeWork2/internal/lib/api/billing"
	"errors"
	"sort"
	"strconv"
)

type CompanyBalance struct {
	Company           string        `json:"company"`
	ValidOperations   int           `json:"valid_operations_count,omitempty"`
	Balance           int           `json:"balance,omitempty"`
	InvalidOperations []interface{} `json:"invalid_operations,omitempty"`
}

var companies = make(map[string]CompanyBalance)

func CountBalance(roots []billing.Root) []CompanyBalance {
	for _, r := range roots {
		if r.Operation.Value.Int != nil {
			doSum(r.Operation.Value.Int, r.Company.Name)
		}

		assignValidInvalidOperations(r)
	}

	return recordAndSortFinalData()
}

func doSum(data interface{}, name string) CompanyBalance {
	companyStruct, ok := companies[name]
	if !ok {
		companyStruct = CompanyBalance{Company: name, Balance: 0}
	}
	if num, err := checkIntValue(data); err == nil {
		companyStruct.Balance += num
		companies[name] = companyStruct
	}
	return companies[name]
}

func checkIntValue(data interface{}) (int, error) {
	switch data.(type) {
	case int:
		return data.(int), nil
	case float64:
		return int(data.(float64)), nil
	case string:
		return strconv.Atoi(data.(string))
	}
	return 0, errors.New("invalid type")
}

func assignValidInvalidOperations(r billing.Root) {
	_, exists := companies[r.Company.Name]
	if !exists {
		companies[r.Company.Name] = CompanyBalance{Company: r.Company.Name}
	}
	compStruct := companies[r.Company.Name]
	compStruct.InvalidOperations = r.Operation.InvalidOperations
	compStruct.ValidOperations = r.Operation.ValidOperations
	companies[r.Company.Name] = compStruct
}

func recordAndSortFinalData() []CompanyBalance {
	var ret []CompanyBalance
	for _, company := range companies {
		ret = append(ret, company)
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return (ret)[i].Company < (ret)[j].Company
	})
	return ret
}
