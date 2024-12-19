package calc

import (
	"errors"
	"log"
	"math"
	"slices"
	"strconv"
	"unicode"
)

var (
	// Errors
	// ErrDivisionByZero                     = errors.New("division by zero is not allowed")
	// ErrEmptyExpression                    = errors.New("empty expression")
	// ErrUnsupportedCharacters              = errors.New(`unsupported characters. expression can only contain numbers, operations (^, *, /, +, -) and float points`)
	// ErrTwoOperatorsInARow                 = errors.New("two operators in a row")
	// ErrExpressionStartsOrEndsWithOperator = errors.New("expression cant start or end with an operator")
	// ErrUnclosedParenthesises              = errors.New("unclosed parenthesises")
	// ErrOperatorsCount                     = errors.New("operators count must be less than numbers count")

	Err422 = errors.New("expression is not valid")
	Err500 = errors.New("internal server error")
)

func Calc(expression string) (float64, error) {
	res, err := calc(expression)
	if err != nil {
		log.Printf("calculation resulted in error: %v\n", err)
	} else {
		log.Println("calculation successful")
	}
	return res, err
}

func calc(expression string) (float64, error) {
	if expression == "" {
		return 0, Err422 //ErrEmptyExpression
	}

	var result float64 = 0

	out_arr := make([]interface{}, 0)
	operationsStack := make([]rune, 0)

	Priority := map[rune]int{
		'^': 3,
		'*': 2,
		'/': 2,
		'+': 1,
		'-': 1,
	}

	i, j := 0, 0

	// удаляет все пробелы из строки

	newExpression := ""
	for _, r := range expression {
		if r != ' ' {
			newExpression += string(r)
		}
	}

	openParenthesises := 0
	operatorsCount := 0
	numbersCount := 0
	wasLastRuneOperator := false

	for ii, r := range newExpression {
		if !unicode.IsDigit(r) && r != '+' && r != '-' && r != '*' && r != '/' && r != '^' && r != '(' && r != ')' && r != '.' {
			return 0, Err422 //ErrUnsupportedCharacters
		}
		// накапливаем цифры для преобразования в число
		if unicode.IsDigit(r) || r == '.' {
			wasLastRuneOperator = false
			j++
			// если последняя цифра в выражении - её нужно ниже добавить в out_arr
			if j < len(newExpression) {
				continue
			}
		}
		if i != j {
			num, err := strconv.ParseFloat(newExpression[i:j], 64)
			if err != nil {
				return result, Err500
			}
			out_arr = append(out_arr, num)
			numbersCount++
		}

		// накапливаем операции
		stack_len := len(operationsStack)
		if r == '+' || r == '-' || r == '*' || r == '/' || r == '^' {
			if wasLastRuneOperator {
				return 0, Err422 //ErrTwoOperatorsInARow
			}
			wasLastRuneOperator = true
			if ii == 0 || ii == len(newExpression)-1 {
				return 0, Err422 //ErrExpressionStartsOrEndsWithOperator
			}
			if stack_len == 0 {
				operationsStack = append(operationsStack, r)
			} else if Priority[r] <= Priority[operationsStack[stack_len-1]] {
				out_arr = append(out_arr, operationsStack[stack_len-1])
				operationsStack[stack_len-1] = r
			} else {
				operationsStack = append(operationsStack, r)
			}
			operatorsCount++
		} else if r == '(' {
			operationsStack = append(operationsStack, r)
			openParenthesises++
		} else if r == ')' {
			openParenthesises--
			k := stack_len
			for operationsStack[k-1] != '(' {
				k--
				out_arr = append(out_arr, operationsStack[k])
				operationsStack = slices.Delete(operationsStack, k, k+1)
			}
		}

		// пропускаем операторы и скобки для преобразования следующих чисел
		j++
		i = j
	}

	if operatorsCount >= numbersCount {
		return 0, Err422 //ErrOperatorsCount
	}

	if openParenthesises != 0 {
		return 0, Err422 //ErrUnclosedParenthesises
	}

	for i := range operationsStack {
		out_arr = append(out_arr, operationsStack[len(operationsStack)-1-i])
	}

	foundNum := false

	for _, token := range out_arr {
		_, ok := token.(float64)
		if ok {
			foundNum = true
			break
		}
		foundNum = false
	}

	if !foundNum {
		return out_arr[0].(float64), nil
	}

	// считаем
	left := []interface{}{}
	for _, token := range out_arr {
		if token != '(' && token != ')' {
			left = append(left, token)
		}
	}
	right := make([]interface{}, 0)

	for _, rr := range left {
		res, ok := rr.(float64)

		if ok {
			if len(right) > 1 {
				// если последний элемент - число, добавляем текущее число либо знак
				lastElem := right[len(right)-1]
				if _, ok := lastElem.(float64); ok {
					right = append(right, res)
					// если же последний элемент это знак - считаем последний элемент, который является числом
				} else {
					if lastElem == '+' {
						result = right[len(right)-2].(float64) + res
						right = slices.Delete(right, len(right)-2, len(right))
						right = append(right, result)

					} else if lastElem == '-' {
						result = right[len(right)-2].(float64) - res
						right = slices.Delete(right, len(right)-2, len(right))
						right = append(right, result)

					} else if lastElem == '*' {
						result = right[len(right)-2].(float64) * res
						right = slices.Delete(right, len(right)-2, len(right))
						right = append(right, result)

					} else if lastElem == '/' {
						if res == 0 {
							return 0, Err422 //ErrDivisionByZero
						}
						result = right[len(right)-2].(float64) / res
						right = slices.Delete(right, len(right)-2, len(right))
						right = append(right, result)

					} else if lastElem == '^' {
						result = math.Pow(right[len(right)-2].(float64), res)
						right = slices.Delete(right, len(right)-2, len(right))
						right = append(right, result)
					}
				}
			} else {
				//иначе добавляем число либо знак.
				right = append(right, res)
			}
			continue
		}
		// если есть 2 или более числа, то считаем 2 последних элемента. результат сохраняем в предпоследний элемент
		if len(right) >= 2 {
			if rr == '+' {
				result = right[len(right)-2].(float64) + right[len(right)-1].(float64)
				right = slices.Delete(right, len(right)-2, len(right))
				right = append(right, result)

			} else if rr == '-' {
				result = right[len(right)-2].(float64) - right[len(right)-1].(float64)
				right = slices.Delete(right, len(right)-2, len(right))
				right = append(right, result)

			} else if rr == '*' {
				result = right[len(right)-2].(float64) * right[len(right)-1].(float64)
				right = slices.Delete(right, len(right)-2, len(right))
				right = append(right, result)

			} else if rr == '/' {
				if right[len(right)-1].(float64) == 0 {
					return 0, Err422 //ErrDivisionByZero
				}
				result = right[len(right)-2].(float64) / right[len(right)-1].(float64)
				right = slices.Delete(right, len(right)-2, len(right))
				right = append(right, result)

			} else if rr == '^' {
				result = math.Pow(right[len(right)-2].(float64), right[len(right)-1].(float64))
				right = slices.Delete(right, len(right)-2, len(right))
				right = append(right, result)
			}
		} else {
			// иначе добавляем число
			right = append(right, rr)
		}
	}

	return result, nil
}
