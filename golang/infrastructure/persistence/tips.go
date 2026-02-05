package persistence

import (
	"context"
	"database/sql"
	"os"

	"fmt"

	"github.com/86shin/commit_goback/domain/model"
	"github.com/86shin/commit_goback/domain/repository"
	"github.com/86shin/commit_goback/domain/utils"
)

// データベース接続クライアントのラッパー
type TipsListRepositoryImpl struct {
	Client    *sql.DB // 例として標準のsql.DBを使用
	TableName string
}

func NewTipsListRepositoryImpl(db *sql.DB) repository.TipsRepository {
	tableName := os.Getenv("TABLE_NAME_T")
	if tableName == "" {
		panic("TABLE_NAME が環境変数に設定されていません")
	}

	return &TipsListRepositoryImpl{
		Client:    db,
		TableName: tableName,
	}
}

func (i *TipsListRepositoryImpl) TipsList(ctx context.Context) ([]*model.TipsList, error) {
	if !utils.ValidateTableName(i.TableName) {
		return nil, fmt.Errorf("不正なテーブル名です")
	}

	query := fmt.Sprintf(`
    	SELECT tips_id, title, category, content
    	FROM %s
	`, i.TableName)

	rows, err := i.Client.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("tipsの一覧取得に失敗しました。")
	}
	defer rows.Close()

	tips := make([]*model.TipsList, 0)

	for rows.Next() {
		var t model.TipsList
		// Scanの引数はポインタを渡す
		if err := rows.Scan(&t.TipsId, &t.Title, &t.Category, &t.Content); err != nil {
			return nil, fmt.Errorf("tipsの一覧取得に失敗しました。(Scan)")
		}
		tips = append(tips, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("tipsの一覧取得に失敗しました。(Rows)")
	}

	return tips, nil
}
