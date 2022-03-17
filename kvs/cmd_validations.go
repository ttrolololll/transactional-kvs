package kvs

import (
	"errors"
	"strings"
)

const (
	CmdSet      = "SET"
	CmdGet      = "GET"
	CmdDelete   = "DELETE"
	CmdCount    = "COUNT"
	CmdBegin    = "BEGIN"
	CmdCommit   = "COMMIT"
	CmdRollback = "ROLLBACK"
)

var (
	cmdMap = map[string]interface{}{
		CmdSet:      nil,
		CmdGet:      nil,
		CmdDelete:   nil,
		CmdCount:    nil,
		CmdBegin:    nil,
		CmdCommit:   nil,
		CmdRollback: nil,
	}
	cmdValidateFuncMap = map[string]validateFuncType{
		CmdSet:    validateSetCmd,
		CmdGet:    validateGetCmd,
		CmdDelete: validateDeleteCmd,
		CmdCount:  validateCountCmd,
	}
)

type validateFuncType func(inputs []string) error

// ValidateCmdInputs takes in full command input, parse to params and validates them
func ValidateCmdInputs(inputs []string) error {
	if len(inputs) == 0 {
		return errors.New("no command given")
	}

	cmd := strings.ToUpper(inputs[0])

	if _, exists := cmdMap[cmd]; !exists {
		return errors.New("command does not exists")
	}

	// early return for commands without input
	if cmd == CmdBegin || cmd == CmdCommit || cmd == CmdRollback {
		return nil
	}

	validateFunc, validateFuncExists := cmdValidateFuncMap[cmd]
	if !validateFuncExists {
		return errors.New("command is currently unsupported")
	}

	return validateFunc(inputs[1:])
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
