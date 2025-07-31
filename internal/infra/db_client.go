package infra

import (
	"context"
	"fmt"
	"time"

	"QuickBrick/internal/domain/ent"

	_ "github.com/go-sql-driver/mysql"
)

func NewEntClient(dsn string) (*ent.Client, error) {
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to mysql: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // 防止 context 泄漏

	fmt.Println("正在同步数据库 schema...") // 调试输出

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	fmt.Println("数据库 schema 同步完成") // 成功提示
	return client, nil

}
