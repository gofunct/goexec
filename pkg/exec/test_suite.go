package exec

import (
	"bytes"
	"context"
	"io"
	osexec "os/exec"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Interface
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_executor_Command(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name     string
		executor *executor
		args     args
		want     Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &executor{}
			if got := executor.Command(tt.args.script); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("executor.Command() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_executor_CommandContext(t *testing.T) {
	type args struct {
		ctx    context.Context
		script string
	}
	tests := []struct {
		name     string
		executor *executor
		args     args
		want     Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &executor{}
			if got := executor.CommandContext(tt.args.ctx, tt.args.script); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("executor.CommandContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_executor_LookPath(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name     string
		executor *executor
		args     args
		want     string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &executor{}
			got, err := executor.LookPath(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("executor.LookPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("executor.LookPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdWrapper_SetDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		cmd  *cmdWrapper
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			cmd.SetDir(tt.args.dir)
		})
	}
}

func Test_cmdWrapper_SetStdin(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name string
		cmd  *cmdWrapper
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			cmd.SetStdin(tt.args.in)
		})
	}
}

func Test_cmdWrapper_SetStdout(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cmdWrapper
		wantOut string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			out := &bytes.Buffer{}
			cmd.SetStdout(out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("cmdWrapper.SetStdout() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_cmdWrapper_SetStderr(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cmdWrapper
		wantOut string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			out := &bytes.Buffer{}
			cmd.SetStderr(out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("cmdWrapper.SetStderr() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_cmdWrapper_SetEnv(t *testing.T) {
	type args struct {
		env []string
	}
	tests := []struct {
		name string
		cmd  *cmdWrapper
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			cmd.SetEnv(tt.args.env)
		})
	}
}

func Test_cmdWrapper_StdoutPipe(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cmdWrapper
		want    io.ReadCloser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			got, err := cmd.StdoutPipe()
			if (err != nil) != tt.wantErr {
				t.Errorf("cmdWrapper.StdoutPipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmdWrapper.StdoutPipe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdWrapper_StderrPipe(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cmdWrapper
		want    io.ReadCloser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			got, err := cmd.StderrPipe()
			if (err != nil) != tt.wantErr {
				t.Errorf("cmdWrapper.StderrPipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmdWrapper.StderrPipe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdWrapper_Start(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cmdWrapper
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			if err := cmd.Start(); (err != nil) != tt.wantErr {
				t.Errorf("cmdWrapper.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cmdWrapper_Wait(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cmdWrapper
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			if err := cmd.Wait(); (err != nil) != tt.wantErr {
				t.Errorf("cmdWrapper.Wait() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cmdWrapper_Run(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cmdWrapper
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			if err := cmd.Run(); (err != nil) != tt.wantErr {
				t.Errorf("cmdWrapper.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cmdWrapper_CombinedOutput(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cmdWrapper
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			got, err := cmd.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Errorf("cmdWrapper.CombinedOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmdWrapper.CombinedOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdWrapper_Output(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *cmdWrapper
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			got, err := cmd.Output()
			if (err != nil) != tt.wantErr {
				t.Errorf("cmdWrapper.Output() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmdWrapper.Output() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdWrapper_Stop(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cmdWrapper
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmdWrapper{}
			cmd.Stop()
		})
	}
}

func TestExitErrorWrapper_ExitStatus(t *testing.T) {
	type fields struct {
		ExitError *osexec.ExitError
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eew := ExitErrorWrapper{
				ExitError: tt.fields.ExitError,
			}
			if got := eew.ExitStatus(); got != tt.want {
				t.Errorf("ExitErrorWrapper.ExitStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodeExitError_Error(t *testing.T) {
	type fields struct {
		Err  error
		Code int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := CodeExitError{
				Err:  tt.fields.Err,
				Code: tt.fields.Code,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("CodeExitError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodeExitError_String(t *testing.T) {
	type fields struct {
		Err  error
		Code int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := CodeExitError{
				Err:  tt.fields.Err,
				Code: tt.fields.Code,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("CodeExitError.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodeExitError_Exited(t *testing.T) {
	type fields struct {
		Err  error
		Code int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := CodeExitError{
				Err:  tt.fields.Err,
				Code: tt.fields.Code,
			}
			if got := e.Exited(); got != tt.want {
				t.Errorf("CodeExitError.Exited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodeExitError_ExitStatus(t *testing.T) {
	type fields struct {
		Err  error
		Code int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := CodeExitError{
				Err:  tt.fields.Err,
				Code: tt.fields.Code,
			}
			if got := e.ExitStatus(); got != tt.want {
				t.Errorf("CodeExitError.ExitStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
