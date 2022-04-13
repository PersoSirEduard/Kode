package kode

import (
	"errors"
	"strconv"
)

type Function struct {
	Arguments map[string]Variable
	Variables map[string]Variable
	Code      string
}

/**
 * Create a new function with scope.
 * @param arguments : map[string]Variable - The arguments of the function.
 * @param variables : map[string]Variable - The variables of the function.
 * @param code : string - The code of the function.
 * @return Function - The new function.
 */
func CreateFunction(arguments map[string]Variable, variables map[string]Variable, code string) Function {
	return Function{
		Arguments: arguments,
		Variables: variables,
		Code:      code,
	}
}

func (scope *Function) Run(arguments interface{}) (interface{}, error) {

	// Split the code into lines.
	lines := LineParse((*scope).Code)

	// Loop through the lines.
	for currentLine, line := range lines {

		// Skip empty lines.
		if line == "" || line == "\r" {
			continue
		}

		// Get the command tokens of the line.
		tokens := Queue{}
		tokensArr := InlineParse(line, " \t\r=#")
		for _, token := range tokensArr {
			tokens.Push(token)
		}

		for !tokens.IsEmpty() {
			command, _ := tokens.Pop()

			switch command {

			// Comment
			case "#":
				// Ignore the rest of the line.
				tokens.Clear()
				continue

			// Variable creation
			case "val":

				// Get variable name
				name, nameProvided := tokens.Pop()

				// Check if the name for the variable was provided
				if !nameProvided {
					return nil, errors.New("Error on line " + strconv.Itoa(currentLine) + ": Missing variable name.")
				}

				// Check if the variable name is valid
				if !HasValidVariableName(name.(string)) {
					return nil, errors.New("Error on line " + strconv.Itoa(currentLine) + ": Invalid variable name. The name must be alphanumeric and start with a letter.")
				}

				// Check if the variable name is already in use
				if (*scope).VariableExists(name.(string)) {
					return nil, errors.New("Error on line " + strconv.Itoa(currentLine) + ": Variable was already defined.")
				}

				equal, assign := tokens.Pop()

				// Check if the variable has an assignments
				if !assign || equal.(string) != "=" {
					return nil, errors.New("Error on line " + strconv.Itoa(currentLine) + ": Missing variable assignment \"=\".")
				}

				// Get the rest of the line tokens and join them to feed the variable.
				value := ""
				for !tokens.IsEmpty() {
					tokenValue, _ := tokens.Pop()

					// Exit if a comment is found.
					if tokenValue.(string) == "#" {
						break
					}

					value += tokenValue.(string)
				}

				// Make sure the variable value is valid
				if value == "" {
					return nil, errors.New("Error on line " + strconv.Itoa(currentLine) + ": Missing variable value.")
				}

				// Create the variable.
				evaluatedValue, err := EvaluateExpression(scope, value)
				if err != nil {
					return nil, err
				}

				// Create the variable in the current scope.

				(*scope).Variables[name.(string)] = evaluatedValue
				println("Created variable " + name.(string) + "(" + evaluatedValue.Type + ").")

				break
			default:
				// Command not found
				return nil, errors.New("Error: Unknown command.")
			}

		}

	}

	return nil, nil
}
