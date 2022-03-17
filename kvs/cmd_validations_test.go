package kvs

import "testing"

func TestValidateCmdInputs(t *testing.T) {
	testcases := []struct {
		tName     string
		input     []string
		expectErr bool
	}{
		{
			tName:     "happy - set",
			input:     []string{"set", "k", "v"},
			expectErr: false,
		},
		{
			tName:     "happy - get",
			input:     []string{"get", "k"},
			expectErr: false,
		},
		{
			tName:     "happy - delete",
			input:     []string{"delete", "k"},
			expectErr: false,
		},
		{
			tName:     "happy - count",
			input:     []string{"count", "v"},
			expectErr: false,
		},
		{
			tName:     "happy - begin",
			input:     []string{"begin"},
			expectErr: false,
		},
		{
			tName:     "happy - commit",
			input:     []string{"commit"},
			expectErr: false,
		},
		{
			tName:     "happy - rollback",
			input:     []string{"rollback"},
			expectErr: false,
		},
		{
			tName:     "err - empty input",
			input:     nil,
			expectErr: true,
		},
		{
			tName:     "err - non-existent validate func",
			input:     []string{"someCmd"},
			expectErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.tName, func(t *testing.T) {
			err := ValidateCmdInputs(tc.input)
			if tc.expectErr && err == nil {
				t.Errorf("expected error")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("did not expected error")
			}
		})
	}
}
