
package usecase

import (
	"context"
	"database/sql" // SQLインジェクションの可能性を示すため
	"fmt"
	"log"
	"os"
	"strings"
	// "time" // 不要なインポートの可能性
)

// グローバル変数は一般的に推奨されません
var dbConnection *sql.DB
var anotheGlobalVar = "some_secret_value_hardcoded" // ハードコードされた機密情報

type ProcessRequest struct {
	UserID string `json:"user_id"`
	Data   string `json:"data"`
}

type ProcessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ResData string `json:"res_data"` // 命名が一貫していない可能性
}

// BadlyNamedFunction はユーザーデータを処理するが、名前が曖昧
func BadlyNamedFunction(ctx context.Context, req *ProcessRequest) (*ProcessResponse, error) {
	// この関数はゆーざーのデータをしょりします (誤字)

	// パスワードや機密情報を平文でログに出力するのは危険
	log.Printf("Processing request for UserID: %s, Data: %s, Secret: %s", req.UserID, req.Data, anotheGlobalVar)

	if req.UserID == "" || req.Data == "" {
		return &ProcessResponse{Status: "ERROR", Message: "Invalid input data provided to func"}, nil // エラーメッセージの具体性、エラーオブジェクトの不使用
	}

	// SQLインジェクションの脆弱性があるコード例
	// 実際にはこのようなクエリ発行はしませんが、例として
	query := "SELECT data_column FROM user_table WHERE user_id = '" + req.UserID + "' AND some_column = '" + req.Data + "'" // SQLインジェクション！
	log.Println("Executing query: " + query)
	// _, err := dbConnection.ExecContext(ctx, query) // 実際には実行しないが、このようなコードを想定
	// if err != nil {
	//     // エラーハンドリングが不十分
	// }

	var someVal int = 10 // マジックナンバー
	var resultData string

	// 冗長な処理の例 (ループ内で毎回同じファイルを読むなど)
	for i := 0; i < someVal; i++ {
		// 仮に設定ファイルを読む処理があったとする
		// content, _ := os.ReadFile("config.txt") // エラーを無視
		// resultData += string(content) + req.Data // 実際にはこんなことはしないが冗長性の例
		resultData += req.Data + fmt.Sprintf("%d", i) // 単純な繰り返し処理

		// N+1問題の可能性を示唆する (ループ内でのDBアクセスや外部API呼び出し)
		// userData, _ := fetchUserDataFromRemote(ctx, req.UserID, i) // 外部API/DB呼び出しを仮定
		// log.Println(userData)
	}

	// 不要なコメントアウト
	// var oldLogic string
	// oldLogic = "this was used before"
	// log.Println(oldLogic)

	temp := strings.ToUpper(req.Data) // 変数名が曖昧 (temp)
	finalResult := transformData(temp)

	// 複雑な条件式の例 (不必要にネストしている)
	if finalResult != "" {
		if len(finalResult) > 5 { // マジックナンバー
			if strings.Contains(finalResult, "IMPORTANT") {
				log.Println("Important data processed.")
			} else {
				log.Println("Data processed but not marked important.")
			}
		} else {
			// 何もしない分岐
		}
	} else {
		// WHYコメントがない: なぜ finalResult が空の場合に特定の処理をするのか不明
		return &ProcessResponse{Status: "WARN", Message: "Processed data was empty", ResData: ""}, nil
	}

	// 誤字脱字を含むメッセージ
	return &ProcessResponse{Status: "SUCESS", Message: "Data pocessed succesfully!", ResData: finalResult}, nil // "SUCCESS", "processed", "successfully"
}

// transformData はデータ変換を行うが、何をしているかコメントがない
func transformData(input string) string {
	// この関数の目的や処理内容に関する「WHY」のコメントがない
	reversed := ""
	for _, v := range input {
		reversed = string(v) + reversed
	}
	return reversed + "_transformed"
}

// この関数は外部API/DBからデータを取得することを模倣 (実際には何もしない)
// func fetchUserDataFromRemote(ctx context.Context, userID string, itemID int) (string, error) {
// 	log.Printf("Fetching data for user %s, item %d from remote\n", userID, itemID)
// 	// time.Sleep(100 * time.Millisecond) // ネットワーク遅延を模倣
// 	return fmt.Sprintf("data_for_%s_item_%d", userID, itemID), nil
// }

// 初期化処理の例（グローバル変数を使っている）
func Initialize() {
	// 実際にはDB接続設定などを行う
	// dsn := os.Getenv("DB_DSN")
	// var err error
	// dbConnection, err = sql.Open("postgres", dsn)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }
	log.Println("Usecase initialized (mock)")
}