package main

import (
	"HomeWork2/internal/lib/api/billing"
	balance "HomeWork2/internal/lib/api/calculated_balance"
	"encoding/json"
	"log"
	"os"
)

func main() {
	jsonFile, err := os.Open("internal/lib/api/json/billing.json")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer jsonFile.Close()

	stat, _ := jsonFile.Stat()
	buf := make([]byte, stat.Size())
	_, err = jsonFile.Read(buf)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	var roots []billing.Root
	if err := json.Unmarshal(buf, &roots); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	companyBalance := balance.CountBalance(roots)
	jsonData, err := json.MarshalIndent(companyBalance, "", "	")
	if err != nil {
		log.Fatalf("failed to marshal JSON: %v", err)
	}

	err = os.WriteFile("internal/lib/api/json/out.json", jsonData, os.ModePerm)
	if err != nil {
		log.Fatal("failed to write file: %v", err)
	}
}
