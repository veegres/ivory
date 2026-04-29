package ssh

import "testing"

func TestNormalizeDockerCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{name: "prefix plain command", command: "ps", expected: "docker ps"},
		{name: "keep docker command", command: "docker ps", expected: "docker ps"},
		{name: "keep sudo docker command", command: "sudo docker ps", expected: "sudo docker ps"},
		{name: "trim spaces", command: "  images  ", expected: "docker images"},
	}

	c := NewClient()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := c.normalizeDockerCommand(test.command)
			if actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
