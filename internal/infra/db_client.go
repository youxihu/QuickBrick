package infra

import (
	"context"
	"fmt"

	"QuickBrick/internal/domain/ent"

	_ "github.com/go-sql-driver/mysql"
)

func NewEntClient(dsn string) (*ent.Client, error) {
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to mysql: %w", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return client, nil
}
