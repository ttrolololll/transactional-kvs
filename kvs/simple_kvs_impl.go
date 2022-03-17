package kvs

import (
	"errors"
	"strconv"
	"strings"
)

type simpleKvsImpl struct {
	kvMap           map[string]string
	vCountMap       map[string]uint64
	activeKvsInst   *simpleKvsImpl
	sessionKvsStack []*simpleKvsImpl
}

// NewSimpleKvs creates a new instance of simple KVS
func NewSimpleKvs() IKvs {
	return &simpleKvsImpl{
		kvMap:           map[string]string{},
		vCountMap:       map[string]uint64{},
		activeKvsInst:   nil,
		sessionKvsStack: []*simpleKvsImpl{},
	}
}

// CommandExecutor takes in cmd line inputs and executes the correct command
func (kvs *simpleKvsImpl) CommandExecutor(inputs []string) (string, error) {
	err := ValidateCmdInputs(inputs)
	if err != nil {
		return "", err
	}

	kvs.determineInstance()

	cmd := strings.ToUpper(inputs[0])
	result := ""
	var executeErr error

	switch cmd {
	case CmdSet:
		kvs.Set(inputs[1], inputs[2])
	case CmdGet:
		result, executeErr = kvs.Get(inputs[1])
	case CmdDelete:
		kvs.Delete(inputs[1])
	case CmdCount:
		result = kvs.Count(inputs[1])
	case CmdBegin:
		kvs.Begin()
	case CmdCommit:
		executeErr = kvs.Commit()
	case CmdRollback:
		executeErr = kvs.Rollback()
	default:
		return "", errors.New("command execution failed: command not supported")
	}

	return result, executeErr
}

// Set creates a transaction and commit the value
func (kvs *simpleKvsImpl) Set(k, v string) {
	kvs.activeKvsInst.kvMap[k] = v
	if _, exists := kvs.activeKvsInst.vCountMap[v]; exists {
		kvs.activeKvsInst.vCountMap[v] += 1
	} else {
		kvs.activeKvsInst.vCountMap[v] = 1
	}
}

// Get retrieves a value by key
func (kvs *simpleKvsImpl) Get(k string) (string, error) {
	v, exists := kvs.activeKvsInst.kvMap[k]
	if !exists {
		return "", errors.New("key not set")
	}
	return v, nil
}

// Delete removes a value by key
func (kvs *simpleKvsImpl) Delete(k string) {
	v, _ := kvs.activeKvsInst.Get(k)
	if v == "" {
		return
	}

	delete(kvs.activeKvsInst.kvMap, k)
	kvs.activeKvsInst.vCountMap[v] -= 1
}

// Count returns the number of value occurrences
func (kvs *simpleKvsImpl) Count(v string) string {
	if _, exists := kvs.activeKvsInst.vCountMap[v]; exists {
		return strconv.FormatUint(kvs.activeKvsInst.vCountMap[v], 10)
	}
	return "0"
}

// Begin initiates session KVS
func (kvs *simpleKvsImpl) Begin() {
	simpleKvs := NewSimpleKvs()
	sessionKvs, _ := simpleKvs.(*simpleKvsImpl) // ignore err as conversion will be valid
	kvs.sessionKvsStack = append(kvs.sessionKvsStack, sessionKvs)
}

// Commit finalises operations in session KVS at top of stack, by applying it to the KVS next in stack
func (kvs *simpleKvsImpl) Commit() error {
	stackLen := len(kvs.sessionKvsStack)
	if stackLen == 0 {
		return errors.New("no transaction")
	}

	currentKvMap := map[string]string{}
	for k, v := range kvs.activeKvsInst.kvMap {
		currentKvMap[k] = v
	}

	// pop from stack and set active instance to next in stack
	kvs.sessionKvsStack = kvs.sessionKvsStack[:stackLen-1]
	kvs.determineInstance()

	for k, v := range currentKvMap {
		kvs.activeKvsInst.kvMap[k] = v
	}

	return nil
}

// Rollback destroys session KVS at top of stack
func (kvs *simpleKvsImpl) Rollback() error {
	stackLen := len(kvs.sessionKvsStack)

	if stackLen == 0 {
		return errors.New("no transaction")
	}

	kvs.sessionKvsStack = kvs.sessionKvsStack[:stackLen-1]

	return nil
}

// determineInstance sets the current active store,
// it could be the root store, or the deepest child store
func (kvs *simpleKvsImpl) determineInstance() {
	stackLen := len(kvs.sessionKvsStack)
	if stackLen == 0 {
		kvs.activeKvsInst = kvs
	} else {
		kvs.activeKvsInst = kvs.sessionKvsStack[stackLen-1]
	}
}
