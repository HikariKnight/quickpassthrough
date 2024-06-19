package command

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_processCmdString(t *testing.T) {
	type args struct {
		cmd string
	}

	kernel_args := "intel_iommu=on iommu=pt vfio-pci.ids=10de:1c03,10de:10f1"

	tests := []struct {
		name  string
		args  args
		want  string
		want1 []string
	}{
		{
			name: "ls -l",
			args: args{
				cmd: "ls -l",
			},
			want: "ls",
			want1: []string{
				"-l",
			},
		},
		{
			name: "ls -l -a",
			args: args{
				cmd: "ls -l -a",
			},
			want: "ls",
			want1: []string{
				"-l",
				"-a",
			},
		},
		{
			name: "rm -v \"file.txt\"",
			args: args{
				cmd: "rm -v \"file.txt\"",
			},
			want: "rm",
			want1: []string{
				"-v",
				"\"file.txt\"",
			},
		},
		{
			name: "rm -v \"file.txt\" -f",
			args: args{
				cmd: "rm -v \"file.txt\" -f",
			},
			want: "rm",
			want1: []string{
				"-v",
				"\"file.txt\"",
				"-f",
			},
		},
		{
			name: fmt.Sprintf("kernelstub -a \"%s\"", kernel_args),
			args: args{
				cmd: fmt.Sprintf("kernelstub -a \"%s\"", kernel_args),
			},
			want: "kernelstub",
			want1: []string{
				"-a",
				fmt.Sprintf("\"%s\"", kernel_args),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := processCmdString(tt.args.cmd)
			if got != tt.want {
				t.Errorf("processCmdString() got = %v, want %v", got, tt.want)
			}
			t.Logf("got: %v", got)
			t.Logf("got1: %v", got1)
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("processCmdString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
