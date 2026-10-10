package migrations

import "embed"

// Files 将 migrations 目录中的 SQL 文件嵌入 Go 程序。
//
//go:embed *.sql
var Files embed.FS
