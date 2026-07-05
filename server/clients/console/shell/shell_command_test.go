package shell

import (
	"strings"
	"testing"
)

func TestCommand_Id(t *testing.T) {
	cmd1 := &Command{Name: "ls", Args: []string{"-l", "/tmp"}}
	cmd2 := &Command{Name: "ls", Args: []string{"-l", "/tmp"}}
	cmd3 := &Command{Name: "ls", Args: []string{"-a", "/tmp"}}
	cmd4 := &Command{Name: "ls", Args: []string{"-l", "/var"}}

	if cmd1.Id() != cmd2.Id() {
		t.Errorf("Expected same ID for identical commands, got %s and %s", cmd1.Id(), cmd2.Id())
	}
	if cmd1.Id() == cmd3.Id() {
		t.Errorf("Expected different ID for different args")
	}
	if cmd1.Id() == cmd4.Id() {
		t.Errorf("Expected different ID for different paths")
	}

	customID := "custom-job-123"
	cmdCustom := &Command{Name: "ls", JobID: customID}
	if cmdCustom.Id() != customID {
		t.Errorf("Expected custom JobID %s, got %s", customID, cmdCustom.Id())
	}
}

func TestCommand_Execute(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *Command
		want    []string
		wantErr bool
	}{
		{
			name: "successful command",
			cmd:  &Command{Name: "echo", Args: []string{"hello\nworld"}},
			want: []string{"hello", "world"},
		},
		{
			name:    "failing command",
			cmd:     &Command{Name: "ls", Args: []string{"/non-existent-directory-ivory"}},
			wantErr: true,
		},
		{
			name: "command with empty output",
			cmd:  &Command{Name: "true"},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("Execute() got %v, want %v", got, tt.want)
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("Execute() got[%d] = %v, want %v", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestCommand_Abort(t *testing.T) {
	// Use a command that sleeps to test abort
	cmd := &Command{Name: "sleep", Args: []string{"10"}}

	_, err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	errAbort := cmd.Abort()
	if errAbort != nil {
		t.Errorf("Abort() error = %v", errAbort)
	}

	errWait := cmd.Wait()
	if errWait == nil {
		t.Error("Expected error from Wait() after Abort(), got nil")
	} else if !strings.Contains(errWait.Error(), "killed") && !strings.Contains(errWait.Error(), "signal") {
		t.Errorf("Expected killed/signal error, got: %v", errWait)
	}
}
