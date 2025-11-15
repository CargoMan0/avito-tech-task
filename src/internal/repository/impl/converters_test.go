package impl

import (
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"reflect"
	"testing"
)

func Test_statusFromDomainToEnum(t *testing.T) {
	type args struct {
		pullRequestStatus domain.PullRequestStatus
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantPanic bool
	}{
		{
			name: "domain.PullRequestStatusOpen argument should result in string \"OPEN\"",
			args: args{
				pullRequestStatus: domain.PullRequestStatusOpen,
			},
			want: "OPEN",
		},
		{
			name: "domain.PullRequestStatusMerged argument should result in string \"MERGED\"",
			args: args{
				pullRequestStatus: domain.PullRequestStatusMerged,
			},
			want: "MERGED",
		},
		{
			name: "unknown domain.PullRequestStatus argument should result panic",
			args: args{
				pullRequestStatus: domain.PullRequestStatus(0),
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("statusFromDomainToEnum() expected panic, but did not panic")
					}
				}()
			}

			if got := statusFromDomainToEnum(tt.args.pullRequestStatus); got != tt.want {
				t.Errorf("statusFromDomainToEnum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prStatusScanner_Scan(t *testing.T) {
	type args struct {
		src any
	}
	tests := []struct {
		name    string
		args    args
		want    domain.PullRequestStatus
		wantErr bool
	}{
		{
			name: "\"OPEN\" string argument should result in domain.PullRequestStatusOpen",
			args: args{
				src: "OPEN",
			},
			want: domain.PullRequestStatusOpen,
		},
		{
			name: "\"MERGED\" string argument should result in domain.PullRequestStatusMerged",
			args: args{
				src: "MERGED",
			},
			want: domain.PullRequestStatusMerged,
		},
		{
			name: "\"OPEN\" []byte argument should result in domain.PullRequestStatusOpen",
			args: args{
				src: []byte("OPEN"),
			},
			want: domain.PullRequestStatusOpen,
		},
		{
			name: "\"MERGED\" []byte argument should result in domain.PullRequestStatusMerged",
			args: args{
				src: []byte("MERGED"),
			},
			want: domain.PullRequestStatusMerged,
		},
		{
			name: "incorrect string argument should result in a non-nil error",
			args: args{
				src: "incorrect_string",
			},
			wantErr: true,
		},
		{
			name: "incorrect []byte argument should result in a non-nil error",
			args: args{
				src: []byte("incorrect_byte"),
			},
			wantErr: true,
		},
		{
			name: "empty string argument should result in a non-nil error",
			args: args{
				src: "",
			},
			wantErr: true,
		},
		{
			name: "empty []byte argument should result in a non-nil error",
			args: args{
				src: []byte(""),
			},
			wantErr: true,
		},
		{
			name: "nil argument should result in a non-nil error",
			args: args{
				src: nil,
			},
			wantErr: true,
		},
		{
			name: "incorrect data type argument should result in a non-nil error",
			args: args{
				src: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &prStatusScanner{}
			err := a.Scan(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v", err)
			}
			if got := a.Status; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scan() got = %v, want %v", got, tt.want)
			}
		})
	}
}
