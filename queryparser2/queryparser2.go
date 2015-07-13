package queryparser2

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/jsonq"
)

func NewExpression() *Expression {
	ex := &Expression{}
	ex.ID = time.Now().UnixNano()

	return ex
}

type Expression struct {
	ID  int64
	Str string
}

func (exp *Expression) IDString() string {
	return fmt.Sprintf("EXP:%v", exp.ID)
}

func (exp *Expression) EvalParsedExpression(left, right interface{}, op string) (bool, error) {
	switch rightValue := right.(type) {
	case int:
		if op == "==" {
			return left.(int) == rightValue, nil

		} else if op == "!=" {
			return left.(int) != rightValue, nil

		} else if op == ">" {
			return left.(int) > rightValue, nil

		} else if op == ">=" {
			return left.(int) >= rightValue, nil

		} else if op == "<" {
			return left.(int) < rightValue, nil

		} else if op == "<=" {
			return left.(int) <= rightValue, nil
		}

	case int64:
		if op == "==" {
			return left.(int64) == rightValue, nil

		} else if op == "!=" {
			return left.(int64) != rightValue, nil

		} else if op == ">" {
			return left.(int64) > rightValue, nil

		} else if op == ">=" {
			return left.(int64) >= rightValue, nil

		} else if op == "<" {
			return left.(int64) < rightValue, nil

		} else if op == "<=" {
			return left.(int64) <= rightValue, nil
		}

	case float32:
		if op == "==" {
			return left.(float32) == rightValue, nil

		} else if op == "!=" {
			return left.(float32) != rightValue, nil

		} else if op == ">" {
			return left.(float32) > rightValue, nil

		} else if op == ">=" {
			return left.(float32) >= rightValue, nil

		} else if op == "<" {
			return left.(float32) < rightValue, nil

		} else if op == "<=" {
			return left.(float32) <= rightValue, nil
		}

	case float64:
		if op == "==" {
			return left.(float64) == rightValue, nil

		} else if op == "!=" {
			return left.(float64) != rightValue, nil

		} else if op == ">" {
			return left.(float64) > rightValue, nil

		} else if op == ">=" {
			return left.(float64) >= rightValue, nil

		} else if op == "<" {
			return left.(float64) < rightValue, nil

		} else if op == "<=" {
			return left.(float64) <= rightValue, nil
		}

	case string:
		if op == "==" {
			return left.(string) == rightValue, nil

		} else if op == "!=" {
			return left.(string) != rightValue, nil
		} else {
			return false, errors.New("Supported operators are: ==, !=")
		}

	default:
		return false, errors.New("Supported types are: int, int64, float32, float64, and string")
	}

	return false, errors.New("Supported operators are: ==, !=, >, >=, <, <=")
}

func (exp *Expression) Eval(expressions string) (bool, error) {
	results := make([]bool, 0)

	if expressions == "" {
		expressions = exp.Str
	}

	for _, boolOp := range []string{"&&", "||"} {

		for _, chunk := range strings.Split(expressions, boolOp) {
			chunk = strings.TrimSpace(chunk)

			// Found boolean value
			if strings.ToLower(chunk) == "true" {
				results = append(results, true)
				continue
			}
			if strings.ToLower(chunk) == "false" {
				results = append(results, false)
				continue
			}

			// Found expression value
			if strings.HasPrefix(chunk, "EXP:") {

			}

			// Found sub boolean expression
			if strings.Contains(chunk, "||") || strings.Contains(chunk, "&&") {
				result, err := exp.Eval(chunk)
				if err != nil {
					return false, err
				}

				results = append(results, result)
				continue
			}

			// Found arithmetic expression
			if !strings.Contains(chunk, "||") && !strings.Contains(chunk, "&&") {

				var op string

				// Find what kind of operator this expression has.
				for _, op = range []string{"==", "!=", ">", ">=", "<", "<="} {
					if strings.Contains(chunk, op) {
						break
					}
				}

				if op == "" {
					continue
				}

				leftOpRight := strings.Split(chunk, op)
				leftStr := leftOpRight[0]
				rightStr := leftOpRight[len(leftOpRight)-1]

				// Try to parse string value to int64
				leftInt64, leftErr := strconv.ParseInt(leftStr, 10, 64)
				rightInt64, rightErr := strconv.ParseInt(rightStr, 10, 64)

				if leftErr == nil && rightErr == nil {
					result, err := NewExpression().EvalParsedExpression(leftInt64, rightInt64, op)
					if err == nil {
						results = append(results, result)
					}
					continue
				}

				// Try to parse string value to float64
				leftFloat64, leftErr := strconv.ParseFloat(leftStr, 64)
				rightFloat64, rightErr := strconv.ParseFloat(rightStr, 64)

				if leftErr == nil && rightErr == nil {
					result, err := NewExpression().EvalParsedExpression(leftFloat64, rightFloat64, op)
					if err == nil {
						results = append(results, result)
					}
					continue
				}

				// Try to parse string value to string
				result, err := NewExpression().EvalParsedExpression(leftStr, rightStr, op)
				if err == nil {
					results = append(results, result)
				}
				continue

			}
		}

		// Reduce results
		for len(results) > 1 {
			left := results[0]
			right := results[1]

			if boolOp == "&&" {
				results = append(results, left && right)
			} else if boolOp == "||" {
				results = append(results, left || right)
			}

			results = results[2:]
		}
	}

	if len(results) == 0 {
		return false, nil
	}

	return results[0], nil
}

func New(data map[string][]byte) (*QueryParser, error) {
	qp := &QueryParser{}
	qp.data = data
	qp.expressions = make(map[int64]*Expression)

	return qp, nil
}

type QueryParser struct {
	data        map[string][]byte
	expressions map[int64]*Expression
}

func (qp *QueryParser) DataValue(datapath, jsonSelector string) (interface{}, error) {
	dataJsonBytes := qp.data[datapath]
	var dataJson map[string]interface{}

	err := json.Unmarshal(dataJsonBytes, &dataJson)
	if err != nil {
		return nil, err
	}

	jq := jsonq.NewQuery(dataJson)

	jsonSelectorChunks := strings.Split(jsonSelector, ".")
	jsonSelectorChunks = append([]string{"Data"}, jsonSelectorChunks...) // Always query from "Data" sub-structure.

	return jq.Interface(jsonSelectorChunks...)
}

func (qp *QueryParser) NewExpression() *Expression {
	exp := NewExpression()
	qp.expressions[exp.ID] = exp
	return exp
}

func (qp *QueryParser) Parse(query string) (bool, error) {
	query, err := qp.replaceDataPathWithValue(query)
	if err != nil {
		return false, err
	}

	query, err = qp.replaceParensWithExpressions(query)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (qp *QueryParser) replaceDataPathWithValue(query string) (string, error) {
	queryChunks := strings.Fields(query)

	for i, chunk := range queryChunks {
		var chunkWithOpenParen string
		var chunkWithoutOpenParen string

		if strings.HasPrefix(chunk, "(/r/") || strings.HasPrefix(chunk, "(/w/") || strings.HasPrefix(chunk, "(/x/") {
			chunkWithOpenParen = chunk
			chunkWithoutOpenParen = strings.Replace(chunk, "(", "", -1)

		} else if strings.HasPrefix(chunk, "/r/") || strings.HasPrefix(chunk, "/w/") || strings.HasPrefix(chunk, "/x/") {
			chunkWithoutOpenParen = chunk
		}

		if chunkWithoutOpenParen != "" {
			dataPathAndJsonSelectorChunks := strings.Split(chunkWithoutOpenParen, ".")

			dataPath := dataPathAndJsonSelectorChunks[0]
			jsonSelectorChunks := dataPathAndJsonSelectorChunks[1:]
			jsonSelector := strings.Join(jsonSelectorChunks, ".")

			valueInterface, err := qp.DataValue(dataPath, jsonSelector)
			if err != nil {
				return "", err
			}

			valueBytes, err := json.Marshal(valueInterface)
			if err != nil {
				return "", err
			}

			if chunkWithOpenParen == "" {
				queryChunks[i] = string(valueBytes)

			} else {
				queryChunks[i] = "(" + string(valueBytes)
			}
		}
	}

	return strings.Join(queryChunks, " "), nil
}

func (qp *QueryParser) replaceParensWithExpressions(query string) (string, error) {
	leftParenIndex := -1
	RightParenIndex := -1

	for {
		if !strings.Contains(query, "(") {
			break
		}

		for index, rn := range query {
			if string(rn) == "(" {
				leftParenIndex = index
			}
			if string(rn) == ")" {
				RightParenIndex = index
			}

			if leftParenIndex > -1 && RightParenIndex > -1 {
				expStringWithParens := query[leftParenIndex : RightParenIndex+1]
				expString := query[leftParenIndex+1 : RightParenIndex]
				exp := qp.NewExpression()
				exp.Str = expString

				query = strings.Replace(query, expStringWithParens, exp.IDString(), 1)
				leftParenIndex = -1
				RightParenIndex = -1

				break
			}
		}
	}

	return query, nil
}
