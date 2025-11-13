package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	// "github.com/joho/godotenv"

	"google.golang.org/api/option"

	genai "github.com/google/generative-ai-go/genai"
	"github.com/gorilla/mux"
)

// Node.jsから受け取るリクエストボディ
type ImageRequest struct {
	Base64Image string `json:"image_data_base64"`
	MimeType    string `json:"mime_type"`
}

// Node.jsへ返すレスポンスボディ
// 🚨 修正: Vision AIの Objects を削除
type AnalysisResponse struct {
	Status   string `json:"status"`
	Analysis struct {
		// Objects           []string `json:"objects"` // 削除
		CompositionAdvice string `json:"compositionAdvice"`
	} `json:"analysis"`
}

// グローバルなAPIキー変数を定義 (Gemini用)
var geminiAPIKey string

func main() {

	geminiAPIKey = os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatalf("Fatal: GEMINI_API_KEY is not set.")
	}

	r := mux.NewRouter()
	r.HandleFunc("/analyze", analyzeImageHandler).Methods("POST")

	port := os.Getenv("GO_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Goサーバーがポート %s で起動しました。", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func analyzeImageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req ImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "無効なリクエストフォーマット", http.StatusBadRequest)
		return
	}
	log.Println(req.MimeType)

	imageBytes, err := base64.StdEncoding.DecodeString(req.Base64Image)
	if err != nil {
		http.Error(w, "Base64デコードエラー", http.StatusBadRequest)
		return
	}

	// --- 1. Vision AI の呼び出しを削除 ---

	// --- 2. Gemini APIによる構図アドバイスを同期実行 ---
	// 🚨 修正: Vision AIのラベル(visionLabels)を引数から削除
	compositionAdvice, err := runGeminiAdviceSync(ctx, imageBytes, req.MimeType)
	if err != nil {
		log.Printf("[Gemini API] エラー: %v", err)
		compositionAdvice = "写真の構図に関するアドバイスを取得できませんでした。"
	}

	// --- 3. Node.jsへのレスポンス ---
	w.Header().Set("Content-Type", "application/json")
	res := AnalysisResponse{
		Status: "success",
		Analysis: struct {
			// Objects           []string `json:"objects"` // 削除
			CompositionAdvice string `json:"compositionAdvice"`
		}{
			// 🚨 修正: Objectsをレスポンスから削除
			// Objects:           visionLabels,
			CompositionAdvice: compositionAdvice,
		},
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "JSONエンコードエラー", http.StatusInternalServerError)
	}
}

// 🚨 削除: runVisionAnalysis 関数は不要になりました

// 修正後の Geminiクライアント関数 (Vision AIのラベル依存を排除)
func runGeminiAdviceSync(ctx context.Context, imageBytes []byte, mimeType string) (string, error) {
	// 🚨 修正: labels []string 引数を削除

	originalMimeType := strings.ToLower(mimeType)

	// 1. MIMEタイプを、ライブラリに渡す「拡張子部分」のみに絞り込む
	finalMediaType := ""

	// どのMIMEタイプが来ても、確実な拡張子部分のみを抽出する
	if strings.Contains(originalMimeType, "jpeg") || strings.Contains(originalMimeType, "jpg") {
		finalMediaType = "jpeg"
	} else if strings.Contains(originalMimeType, "png") {
		finalMediaType = "png"
	} else {
		// サポート外の場合はログを出力し、'jpeg'に強制フォールバック
		log.Printf("[MIME CRITICAL FIX] Unexpected type found: %s. Forcing MediaType to 'jpeg'.", originalMimeType)
		finalMediaType = "jpeg"
	}

	// デバッグログ
	log.Printf("[Gemini FINAL MIME] Sending MediaType: %s", finalMediaType)

	// 認証処理 (Gemini APIはAPIキー認証)
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		return "", fmt.Errorf("geminiクライアントの作成に失敗: %w (apiキーの設定を確認してください)", err)
	}
	defer client.Close()

	// 🚨 修正: プロンプトをさらに具体的に、抽象的な表現を禁止するよう変更
	prompt := "あなたはプロの写真家です。この画像を見て、写真がもっと良くなるためのアドバイスをください。以下のルールを厳守してください。\n1. 専門用語（例：三分割法）は使わない。\n2. 「良い感じ」「もっと素敵に」のような抽象的な表現は使わない。\n3. 「何を」「どうすれば」良くなるか、具体的な行動（例：「もう少し右に寄る」「少し下から撮る」）を指示する。\n4. 「人」や「物」の位置や向きに注目する。\n5. アドバイスは80文字以内。\n6. 最後に、被写体が「人」か「食事」かを判断し、[人]、[飯]、[人,飯]、[x]（どちらでもない場合）のいずれかを必ず付ける。"
	log.Println(prompt)

	content := []genai.Part{
		// 修正点: 拡張子部分のみの finalMediaType を genai.ImageData に渡す
		genai.ImageData(finalMediaType, imageBytes),
		genai.Text(prompt),
	}

	// 🚨 変更なし: ご要望通り gemini-2.5-flash を使用
	resp, err := client.GenerativeModel("gemini-2.5-flash").GenerateContent(ctx, content...)
	if err != nil {
		return "", fmt.Errorf("コンテンツ生成リクエストに失敗: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini APIからの応答が空です")
	}

	// genai.Text に型アサーションし、そこから string に変換
	part := resp.Candidates[0].Content.Parts[0]
	textPart, ok := part.(genai.Text)

	if !ok {
		return "", fmt.Errorf("gemini APIからの応答形式が予期されていません (応答がテキストではありません)")
	}

	advice := string(textPart)
	advice = strings.TrimSpace(advice)
	return advice, nil
}
