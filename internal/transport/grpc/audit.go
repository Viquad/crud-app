package grpc

import (
	"github.com/Viquad/crud-audit-service/pkg/domain/audit"
	"google.golang.org/grpc"
)

type AuditClient struct {
	audit.AuditServiceClient
}

func NewAuditClient(addr string) (*AuditClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &AuditClient{audit.NewAuditServiceClient(conn)}, nil
}
