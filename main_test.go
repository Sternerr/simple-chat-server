package main

import "testing"

func testParseFlag(t *testing.T) {
	cases := []struct {
		name string
		args []string
		expected string
		expectedError bool
	}{
		{
			name: "valid mode repl",
			args: []string{"--mode", "repl"},
			expected: "repl",
			expectedError: false,
		}, 
		{
			name: "valid mode tui",
			args: []string{"--mode", "tui"},
			expected: "tui",
			expectedError: false,
		}, 
		{
			name: "valid mode Mixed Case",
			args: []string{"--mode", "RePl"},
			expected: "repl",
			expectedError: false,
		}, 
		{
			name: "no mode flag",
			args: []string{"", "repl"},
			expected: "",
			expectedError: true,
		}, 
		{
			name: "no flag value",
			args: []string{"--mode", ""},
			expected: "",
			expectedError: true,
		}, 
		{
			name: "followed my another flag",
			args: []string{"--mode", "--port"},
			expected: "",
			expectedError: true,
		}, 
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			value, err := parseFlag("--mode", c.args)
			if c.expectedError {
                if err == nil {
                    t.Errorf("expected error, got nil")
                }
                if value != "" {
                    t.Errorf("expected empty string, got %q", value)
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if c.expected != value {
                    t.Errorf("expected %q, got %q", c.expected, value)
                }
            }
		})
	}
}
