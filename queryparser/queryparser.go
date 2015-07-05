// Package queryparser provides tools to query json values.
package queryparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/jsonq"
)

func New(jsonBytes []byte) *QueryParser {
	qp := &QueryParser{}
	qp.JSONBytes = jsonBytes
	return qp
}

type QueryParser struct {
	JSONBytes []byte
}

func (qp *QueryParser) JSONQuery() ([]interface{}, error) {
	var data []interface{}

	err := json.Unmarshal(qp.JSONBytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (qp *QueryParser) EvalSingleValue() (bool, error) {
	jsonQuery, err := qp.JSONQuery()
	if err != nil {
		return false, err
	}

	if len(jsonQuery) == 1 {
		if value, ok := jsonQuery[0].(bool); ok {
			return value, nil
		}
	}
	return false, nil
}

// EvalSingleExpression given full dataset from storage.
func (qp *QueryParser) EvalSingleExpression(data map[string][]byte) (bool, error) {
	jsonQuery, err := qp.JSONQuery()
	if err != nil {
		return false, err
	}

	if len(jsonQuery) == 3 {
		operator, ok := jsonQuery[0].(string)
		if !ok {
			return false, errors.New(fmt.Sprintf("Operator %v must be of type string", jsonQuery[0]))
		}

		datapointMap := jsonQuery[1].(map[string]interface{})

		var dataJsonBytes []byte
		var dataJson map[string]interface{}
		var jsonSelector string

		for key, jsonSelectorInterface := range datapointMap {
			jsonSelector, ok = jsonSelectorInterface.(string)
			if !ok {
				return false, errors.New(fmt.Sprintf("JSON selector %v must be of type string", jsonSelectorInterface))
			}
			dataJsonBytes = data[key]
			break
		}

		err = json.Unmarshal(dataJsonBytes, &dataJson)
		if err != nil {
			return false, err
		}

		jq := jsonq.NewQuery(dataJson)

		jsonSelectorChunks := strings.Split(jsonSelector, ".")
		jsonSelectorChunks = append([]string{"Data"}, jsonSelectorChunks...) // Always query from "Data" sub-structure.

		queryValueInteface := jsonQuery[2]
		switch queryValue := queryValueInteface.(type) {
		case int:
			dataValue, err := jq.Int(jsonSelectorChunks...)
			if err != nil {
				return false, err
			}

			if operator == "==" {
				return dataValue == queryValue, nil

			} else if operator == "!=" {
				return dataValue != queryValue, nil

			} else if operator == ">" {
				return dataValue > queryValue, nil

			} else if operator == ">=" {
				return dataValue >= queryValue, nil

			} else if operator == "<" {
				return dataValue < queryValue, nil

			} else if operator == "<=" {
				return dataValue <= queryValue, nil
			}

		case int64:
			dataValueInterface, err := jq.Interface(jsonSelectorChunks...)
			if err != nil {
				return false, err
			}
			dataValue, ok := dataValueInterface.(int64)
			if !ok {
				return false, errors.New("Query value is of type int64 but data value is not of type int64")
			}

			if operator == "==" {
				return dataValue == queryValue, nil

			} else if operator == "!=" {
				return dataValue != queryValue, nil

			} else if operator == ">" {
				return dataValue > queryValue, nil

			} else if operator == ">=" {
				return dataValue >= queryValue, nil

			} else if operator == "<" {
				return dataValue < queryValue, nil

			} else if operator == "<=" {
				return dataValue <= queryValue, nil
			}

		case float32:
			dataValueFloat64, err := jq.Float(jsonSelectorChunks...)
			if err != nil {
				return false, err
			}

			dataValue := float32(dataValueFloat64)

			if operator == "==" {
				return dataValue == queryValue, nil

			} else if operator == "!=" {
				return dataValue != queryValue, nil

			} else if operator == ">" {
				return dataValue > queryValue, nil

			} else if operator == ">=" {
				return dataValue >= queryValue, nil

			} else if operator == "<" {
				return dataValue < queryValue, nil

			} else if operator == "<=" {
				return dataValue <= queryValue, nil
			}

		case float64:
			dataValue, err := jq.Float(jsonSelectorChunks...)
			if err != nil {
				return false, err
			}

			if operator == "==" {
				return dataValue == queryValue, nil

			} else if operator == "!=" {
				return dataValue != queryValue, nil

			} else if operator == ">" {
				return dataValue > queryValue, nil

			} else if operator == ">=" {
				return dataValue >= queryValue, nil

			} else if operator == "<" {
				return dataValue < queryValue, nil

			} else if operator == "<=" {
				return dataValue <= queryValue, nil
			}

		case string:
			dataValue, err := jq.String(jsonSelectorChunks...)
			if err != nil {
				return false, err
			}

			if operator == "==" {
				return dataValue == queryValue, nil

			} else if operator == "!=" {
				return dataValue != queryValue, nil
			}

		default:
			return false, errors.New("Supported types are: int, int64, float32, float64, and string")
		}
	}
	return false, nil
}
