package repoconvert

import (
	"errors"
	"fmt"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
)

func StatusFromDomainToEnum(s domain.PullRequestStatus) string {
	switch s {
	case domain.PullRequestStatusOpen:
		return "OPEN"
	case domain.PullRequestStatusMerged:
		return "MERGED"
	default:
		panic(fmt.Sprintf("unknown domain.PullRequestStatus: (%v)", s))
	}
}

type PRStatusScanner struct {
	Status domain.PullRequestStatus
}

func (p *PRStatusScanner) Scan(src any) error {
	if src == nil {
		return errors.New("src is nil")
	}

	switch v := src.(type) {
	case string:
		return p.scanString(v)
	case []byte:
		return p.scanString(string(v))
	default:
		return fmt.Errorf("cannot convert type %T to domain.PullRequestStatus", src)
	}
}

func (p *PRStatusScanner) scanString(v string) error {
	switch v {
	case "OPEN":
		p.Status = domain.PullRequestStatusOpen
		return nil
	case "MERGED":
		p.Status = domain.PullRequestStatusMerged
		return nil
	default:
		return fmt.Errorf("unknown enum value %s", v)
	}
}
