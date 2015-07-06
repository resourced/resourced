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

// EvalExpressions evaluates nested expressions and returns boolean given full dataset from storage.
// Each slice tree has three components: [operator, (expression || data point map), (expression || value)].
// There is an exception where a slice tree only has one value. That value must be a boolean type.
// Operator && and || only work when both left and right sides are expressions.
func (qp *QueryParser) EvalExpressions(data map[string][]byte, jsonQuery []interface{}) (bool, error) {
	if jsonQuery == nil {
		var err error

		jsonQuery, err = qp.JSONQuery()
		if err != nil {
			return false, err
		}
	}

	var operator string
	var isBooleanOperator bool

	var leftSideDatapointMap map[string]interface{}
	var isLeftSideDatapointMap bool

	var rightSideValue interface{}

	// Special case, when length of jsonQuery is 1.
	if len(jsonQuery) == 1 {
		value, ok := jsonQuery[0].(bool)
		if !ok {
			return false, errors.New("Single value query must always be boolean")
		}

		return value, nil
	}

	// Step 1: For each of the 3 parts, gather all the variables.
	for i, jsonQueryPart := range jsonQuery {
		// Part 1 is always the operator
		if i == 0 {
			var ok bool

			operator, ok = jsonQueryPart.(string)
			if !ok {
				return false, errors.New(fmt.Sprintf("Operator %v must be of type string", jsonQueryPart))
			}

			if operator == "&&" || operator == "||" {
				isBooleanOperator = true
			} else {
				isBooleanOperator = false
			}
		}

		// Part 2 is the left side of the tree.
		// It can either be another expression or a data point map.
		// Data point map is used for traversing the JSON data.
		// If an expression is found, evaulate immediately.
		if i == 1 {
			leftSideExpression, isLeftSideExpression := jsonQueryPart.([]interface{})

			if isLeftSideExpression {
				evaluated, err := qp.EvalExpressions(data, leftSideExpression)
				if err != nil {
					return false, err
				}

				jsonQuery[i] = evaluated

			} else {
				leftSideDatapointMap, isLeftSideDatapointMap = jsonQueryPart.(map[string]interface{})
			}
		}

		if i == 2 {
			rightSideExpression, isRightSideExpression := jsonQueryPart.([]interface{})
			if isRightSideExpression {
				evaluated, err := qp.EvalExpressions(data, rightSideExpression)
				if err != nil {
					return false, err
				}

				jsonQuery[i] = evaluated

			} else {
				rightSideValue = jsonQueryPart
			}
		}
	}

	if isBooleanOperator {
		leftAsBool, ok := jsonQuery[1].(bool)
		if !ok {
			return false, errors.New(fmt.Sprintf("Boolean operator: %v cannot operate on non boolean", operator))
		}

		rightAsBool, ok := jsonQuery[2].(bool)
		if !ok {
			return false, errors.New(fmt.Sprintf("Boolean operator: %v cannot operate on non boolean", operator))
		}

		if operator == "&&" {
			return leftAsBool && rightAsBool, nil
		} else if operator == "||" {
			return leftAsBool || rightAsBool, nil
		}
	}

	if isLeftSideDatapointMap {
		// Pulls out the right data
		var dataJsonBytes []byte
		var dataJson map[string]interface{}
		var jsonSelector string

		for key, jsonSelectorInterface := range leftSideDatapointMap {
			var ok bool
			jsonSelector, ok = jsonSelectorInterface.(string)
			if !ok {
				return false, errors.New(fmt.Sprintf("JSON selector %v must be of type string", jsonSelectorInterface))
			}
			dataJsonBytes = data[key]
			break
		}

		err := json.Unmarshal(dataJsonBytes, &dataJson)
		if err != nil {
			return false, err
		}

		jq := jsonq.NewQuery(dataJson)

		jsonSelectorChunks := strings.Split(jsonSelector, ".")
		jsonSelectorChunks = append([]string{"Data"}, jsonSelectorChunks...) // Always query from "Data" sub-structure.

		switch queryValue := rightSideValue.(type) {
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
