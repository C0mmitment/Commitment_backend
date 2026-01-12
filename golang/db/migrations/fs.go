package migrations

import "embed"

// 同じフォルダにある .sql ファイルをすべて埋め込む
//
//go:embed *.sql
var FS embed.FS
