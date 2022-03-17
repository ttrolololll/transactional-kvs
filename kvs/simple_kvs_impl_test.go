package kvs

import (
	"testing"
)

func TestNewSimpleKvs(t *testing.T) {
	simpleKvs := NewSimpleKvs()
	if _, ok := simpleKvs.(*simpleKvsImpl); !ok {
		t.Errorf("is not of type simpleKvsImpl")
	}
}

func TestSimpleKvsImpl_CommandExecutor(t *testing.T) {
	testcases := []struct {
		tName          string
		input          []string
		preSteps       [][]string
		expectedResult string
		expectErr      bool
	}{
		{
			tName:          "err - cmd input validation failure",
			input:          []string{"someCmd"},
			preSteps:       [][]string{},
			expectedResult: "",
			expectErr:      true,
		},
		{
			tName:          "happy - set cmd",
			input:          []string{"set", "k", "v"},
			preSteps:       [][]string{},
			expectedResult: "",
			expectErr:      false,
		},
		{
			tName:          "happy - get cmd",
			input:          []string{"get", "k"},
			preSteps:       [][]string{{"set", "k", "v"}},
			expectedResult: "v",
			expectErr:      false,
		},
		{
			tName:          "happy - delete cmd",
			input:          []string{"delete", "k"},
			preSteps:       [][]string{},
			expectedResult: "",
			expectErr:      false,
		},
		{
			tName:          "happy - count cmd",
			input:          []string{"count", "v"},
			preSteps:       [][]string{{"set", "k1", "v"}, {"set", "k2", "v"}, {"set", "k3", "v"}},
			expectedResult: "3",
			expectErr:      false,
		},
		{
			tName:          "happy - begin cmd",
			input:          []string{"begin"},
			preSteps:       [][]string{},
			expectedResult: "",
			expectErr:      false,
		},
		{
			tName:          "happy - commit cmd",
			input:          []string{"commit"},
			preSteps:       [][]string{{"begin"}},
			expectedResult: "",
			expectErr:      false,
		},
		{
			tName:          "err - commit cmd",
			input:          []string{"commit"},
			preSteps:       [][]string{},
			expectedResult: "",
			expectErr:      true,
		},
		{
			tName:          "happy - rollback cmd",
			input:          []string{"rollback"},
			preSteps:       [][]string{{"begin"}, {"set", "k", "v"}},
			expectedResult: "",
			expectErr:      false,
		},
		{
			tName:          "err - rollback cmd",
			input:          []string{"rollback"},
			preSteps:       [][]string{},
			expectedResult: "",
			expectErr:      true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.tName, func(t *testing.T) {
			kvsInst := NewSimpleKvs()

			if len(tc.preSteps) > 0 {
				for _, cmd := range tc.preSteps {
					_, _ = kvsInst.CommandExecutor(cmd)
				}
			}

			res, err := kvsInst.CommandExecutor(tc.input)
			if res != tc.expectedResult {
				t.Errorf("result does not match")
			}
			if tc.expectErr && err == nil {
				t.Errorf("expected error")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("did not expected error")
			}
		})
	}
}

func TestSimpleKvsImpl_Set(t *testing.T) {
	testcases := []struct {
		tName         string
		key           string
		value         string
		expectedCount string
	}{
		{
			tName:         "happy",
			key:           "k",
			value:         "v",
			expectedCount: "1",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.tName, func(t *testing.T) {
			kvsInst := NewSimpleKvs()

			kvsInst.Set(tc.key, tc.value)
			if tc.expectedCount != kvsInst.Count(tc.value) {
				t.Errorf("count does not match")
			}
		})
	}
}

func TestSimpleKvsImpl_Get(t *testing.T) {
	testcases := []struct {
		tName          string
		key            string
		preSteps       [][]string
		expectedResult string
		expectErr      bool
	}{
		{
			tName:          "happy",
			key:            "k",
			preSteps:       [][]string{{"set", "k", "v"}},
			expectedResult: "v",
			expectErr:      false,
		},
		{
			tName:          "err - key does not exist",
			key:            "k",
			preSteps:       [][]string{},
			expectedResult: "",
			expectErr:      true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.tName, func(t *testing.T) {
			kvsInst := NewSimpleKvs()

			if len(tc.preSteps) > 0 {
				for _, cmd := range tc.preSteps {
					_, _ = kvsInst.CommandExecutor(cmd)
				}
			}

			res, err := kvsInst.Get(tc.key)
			if res != tc.expectedResult {
				t.Errorf("result does not match")
			}
			if tc.expectErr && err == nil {
				t.Errorf("expected error")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("did not expected error")
			}
		})
	}
}

func TestSimpleKvsImpl_Delete(t *testing.T) {
	kvsInst := NewSimpleKvs()
	kvsInst.Set("k", "v")
	kvsInst.Delete("k")
	res, err := kvsInst.Get("k")
	if err == nil {
		t.Errorf("expected err")
	}
	if res != "" {
		t.Errorf("expected empty result")
	}
}

func TestSimpleKvsImpl_Count(t *testing.T) {
	kvsInst := NewSimpleKvs()
	kvsInst.Set("k1", "v")
	kvsInst.Set("k2", "v")
	kvsInst.Set("k3", "v")
	res := kvsInst.Count("v")
	if res != "3" {
		t.Errorf("expected 3")
	}
}

func TestSimpleKvsImpl_Begin(t *testing.T) {
	kvsInst := &simpleKvsImpl{
		kvMap:           map[string]string{},
		vCountMap:       map[string]uint64{},
		activeKvsInst:   nil,
		sessionKvsStack: []*simpleKvsImpl{},
	}
	kvsInst.Begin()
	if len(kvsInst.sessionKvsStack) != 1 {
		t.Errorf("stack len should be 1")
	}
}

func TestSimpleKvsImpl_Commit(t *testing.T) {
	kvsInst := &simpleKvsImpl{
		kvMap:           map[string]string{},
		vCountMap:       map[string]uint64{},
		activeKvsInst:   nil,
		sessionKvsStack: []*simpleKvsImpl{},
	}

	err := kvsInst.Commit()
	if err == nil {
		t.Errorf("expected err")
	}

	kvsInst.Begin()
	err = kvsInst.Commit()
	if err != nil {
		t.Errorf("expected no err")
	}
}

func TestSimpleKvsImpl_Rollback(t *testing.T) {
	kvsInst := &simpleKvsImpl{
		kvMap:           map[string]string{},
		vCountMap:       map[string]uint64{},
		activeKvsInst:   nil,
		sessionKvsStack: []*simpleKvsImpl{},
	}

	err := kvsInst.Rollback()
	if err == nil {
		t.Errorf("expected err")
	}

	kvsInst.Begin()
	err = kvsInst.Rollback()
	if err != nil {
		t.Errorf("expected no err")
	}
}
