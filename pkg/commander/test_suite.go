package commander

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCommander(t *testing.T) {
	type args struct {
		name string
		usg  string
	}
	tests := []struct {
		name string
		args args
		want *Commander
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommander(tt.args.name, tt.args.usg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommander() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommander_AddScript(t *testing.T) {
	type fields struct {
		root   *cobra.Command
		script *cobra.Command
	}
	type args struct {
		name   string
		usg    string
		dir    string
		script string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commander{
				root:   tt.fields.root,
				script: tt.fields.script,
			}
			c.AddScript(tt.args.name, tt.args.usg, tt.args.dir, tt.args.script)
		})
	}
}

func TestCommander_AddDescription(t *testing.T) {
	type fields struct {
		root   *cobra.Command
		script *cobra.Command
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commander{
				root:   tt.fields.root,
				script: tt.fields.script,
			}
			c.AddDescription(tt.args.s)
		})
	}
}

func TestCommander_AddVersion(t *testing.T) {
	type fields struct {
		root   *cobra.Command
		script *cobra.Command
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commander{
				root:   tt.fields.root,
				script: tt.fields.script,
			}
			c.AddVersion(tt.args.s)
		})
	}
}

func TestCommander_Execute(t *testing.T) {
	type fields struct {
		root   *cobra.Command
		script *cobra.Command
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commander{
				root:   tt.fields.root,
				script: tt.fields.script,
			}
			if err := c.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Commander.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}