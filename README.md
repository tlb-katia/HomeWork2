# Financial Operations Balance Calculator

This project processes financial operations data for companies, calculates their balances, and outputs the results in a formatted JSON file.

## Project Overview

The main goal of this project is to read financial operations data from a file, environment variable, or standard input, validate the data, calculate the balance for each company, and output the results in a specified JSON format.

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [File Structure](#file-structure)
- [Custom Deserialization](#custom-deserialization)
- [Validations](#validations)
- [Output](#output)
- [Implementation Details](#implementation-details)


## Getting Started

### Prerequisites

- Go 1.16 or later

### Installation

Clone the repository:

```sh
git clone https://github.com/tlb-katia/HomeWork2/tree/hw2
cd HomeWork2
```

### Usage
```sh
go run main.go --file=path/to/input.json
```

### File Structure
The input file should be a JSON file with the following structure:

```sh
[
    {
        "company": "horns",
        "operation": {
            "type": "income",
            "value": 123,
            "id": 1,
            "created_at": "2021-09-09T12:55:00Z"
        }
    },
    {
        "company": "hoofs",
        "type": "-",
        "value": "123",
        "id": "abcd-123-iydc",
        "created_at": "2021-09-09T12:55:00Z"
    }
]
```

### Custom Deserialization

The project uses custom deserialization methods for each field in the Operation struct to handle various data formats and ensure validity. The custom types and their deserialization logic are defined in the billing package.

### Validations

- operation.type can be one of income, outcome, +, -.
- operation.value must be an integer, float (always integral), or a string (always integral).
- operation.id must be an integer or string.
- Missing or invalid values for type, value, or id are considered invalid operations.
- Invalid operations are collected and recorded for each company.

## Implementation Details
### calculated_balance Package
#### Structures
- CompanyBalance: Represents the balance information for a company, including valid operations count, balance, and invalid operations.
  Functions
- CountBalance: Processes the financial operations data, validates the operations, and calculates the balance for each company.
  doSum: Adds the operation value to the company's balance if the value is valid.
- checkIntValue: Checks and converts the operation value to an integer.
- assignValidInvalidOperations: Assigns valid and invalid operations to the company.
- recordAndSortFinalData: Collects the results, sorts them alphabetically by company name, and returns the sorted data.

### Output
```sh
[
    {
        "company": "hoofs",
        "valid_operations_count": 123,
        "balance": -25,
        "invalid_operations": [
            "abc",
            3,
            "cdr"
        ]
    },
    {
        "company": "horns",
        "valid_operations_count": 123,
        "balance": -25,
        "invalid_operations": [
            "abc",
            3,
            "cdr"
        ]
    }
]

```
The results are sorted alphabetically by company name and formatted with tabs for indentation. If a company has no invalid operations, the invalid_operations field is omitted.