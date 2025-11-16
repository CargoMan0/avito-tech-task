package handlers

import (
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"testing"
)

func Test_convertPRStatusFromDomain(t *testing.T) {
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
						t.Errorf("convertPRStatusFromDomain() expected panic, but did not panic")
					}
				}()
			}

			if got := convertPRStatusFromDomain(tt.args.pullRequestStatus); got != tt.want {
				t.Errorf("convertPRStatusFromDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
