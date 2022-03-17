package kvs

import (
	"errors"
	"strings"
)

const (
	CmdSet    = "SET"
	CmdGet    = "GET"
	CmdDelete = "DELETE"
	CmdCount  = "COUNT"
)

var (
	cmdMap = map[string]interface{}{
		CmdSet:    nil,
		CmdGet:    nil,
		CmdDelete: nil,
		CmdCount:  nil,
	}
	cmdValidateFuncMap = map[string]validateFuncType{
		CmdSet:    validateSetCmd,
		CmdGet:    validateGetCmd,
		CmdDelete: validateDeleteCmd,
		CmdCount:  validateCountCmd,
	}
)

type validateFuncType func(inputs []string) error

func ValidateCmdInputs(inputs []string) error {
	if len(inputs) == 0 {
		return errors.New("no command given")
	}

	cmd := strings.ToUpper(inputs[0])

	if _, exists := cmdMap[cmd]; !exists {
		return errors.New("command does not exists")
	}

	validateFunc, validateFuncExists := cmdValidateFuncMap[cmd]
	if !validateFuncExists {
		return errors.New("command is currently unsupported")
	}

	return validateFunc(inputs)
}

var validateSetCmd = func(inputs []string) error {
	if len(inputs) != 2 {
		return errors.New("invalid number of inputs, need 2")
	}
	return nil
}

var validateGetCmd = func(inputs []string) error {
	if len(inputs) != 1 {
		return errors.New("invalid number of inputs, need 1")
	}
	return nil
}

var validateDeleteCmd = func(inputs []string) error {
	if len(inputs) != 1 {
		return errors.New("invalid number of inputs, need 1")
	}
	return nil
}

var validateCountCmd = func(inputs []string) error {
	if len(inputs) != 1 {
		return errors.New("invalid number of inputs, need 1")
	}
	return nil
}
