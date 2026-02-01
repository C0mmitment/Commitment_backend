package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/86shin/commit_goback/domain/model"
	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiAIService は domain/service.AIConnector インターフェースの実装です。
type GeminiAIService struct {
	Client *genai.Client
}

// NewGeminiAIService は GeminiAIService のコンストラクタです。
func NewGeminiAIService(apiKey string) (*GeminiAIService, error) {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("geminiクライアントの作成に失敗: %v", err)
		return nil, fmt.Errorf("geminiクライアントの作成に失敗: %w (apiキーの設定を確認してください)", err)
	}

	// defer client.Close() は、アプリケーションのメイン関数やDIコンテナのシャットダウン処理で行うべき
	// ここでは、クライアントを構造体に保持して返す
	return &GeminiAIService{Client: client}, nil
}

// GetCompositionAdvice は AIConnector のインターフェースを実装します。
func (s *GeminiAIService) GetCompositionAdvice(ctx context.Context, category string, imageReader io.Reader, mimeType string, prevAnalysis *model.CompositionAnalysis) (*model.CompositionAnalysis, error) {

	// ★ ここでだけバイト化（SDKの制約）
	imageBytes, err := io.ReadAll(imageReader)
	if err != nil {
		return nil, fmt.Errorf("画像ストリームの読み込みに失敗: %w", err)
	}
	if len(imageBytes) == 0 {
		return nil, fmt.Errorf("画像データが空です")
	}

	originalMimeType := strings.ToLower(mimeType)
	finalMediaType := "image/jpeg" // APIの仕様上 "image/" プレフィックス推奨

	// シンプルなMIME判定
	if strings.Contains(originalMimeType, "png") {
		finalMediaType = "image/png"
	}

	// 1. モデル設定 (gemini-2.5-flash)
	modelName := "gemini-2.5-flash"
	genModel := s.Client.GenerativeModel(modelName)

	// 2. JSONモードとスキーマ定義 (これでトークン節約＆高速化)
	genModel.ResponseMIMEType = "application/json"
	genModel.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"reason": {
				Type:        genai.TypeString,
				Description: "なぜそのアドバイスが必要なのか、現状の課題や原因（例：被写体が暗い、水平が取れていない、余白が多すぎる）",
			},
			"advice":   {Type: genai.TypeString, Description: "60文字以内の具体的で優しいアドバイス"},
			"category": {Type: genai.TypeString},
			"visual_cues": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"target": {Type: genai.TypeString, Enum: []string{"camera"}},
						"direction": {
							Type: genai.TypeString,
							// 選択肢を固定して高速化
							Enum: []string{"left", "right", "up", "down", "up_left", "up_right", "down_left", "down_right", "forward", "backward"},
						},
					},
					Required: []string{"target", "direction"},
				},
			},
			"evaluation": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"status": {
						Type:        genai.TypeString,
						Enum:        []string{"improved", "unchanged", "regressed", "first_time"},
						Description: "improved:改善, unchanged:変化なし, regressed:悪化, first_time:初回",
					},
					"comment": {Type: genai.TypeString, Description: "前回の課題(reason)が解決されたかどうかのコメント"},
				},
				Required: []string{"status", "comment"},
			},
		},
		Required: []string{"reason", "advice", "category", "visual_cues", "evaluation"},
	}

	// 3. プロンプト設計 (フォーマット指示を削除した軽量版)
	// プロンプト設計
	basePrompt := `あなたはプロの写真家です。画像を見て構図改善のアドバイスをください。
    カテゴリ: %s
    
    ### ルール:
    - "reason" には「何が良くないか（原因）」を簡潔に書いてください。
    - "advice" には「どう動けばいいか（解決策）」を優しく書いてください。`

	var prompt string

	if prevAnalysis == nil {
		// --- 初回の場合 ---
		prompt = fmt.Sprintf(basePrompt+`
        
        現在は「初回撮影」です。
        evaluation.status は "first_time" に設定してください。
        evaluation.comment は「撮影ありがとうございます！まずは今の状態を分析します」等の挨拶にしてください。`, category)
	} else {
		// --- 比較の場合 ---
		prevJSONBytes, _ := json.Marshal(prevAnalysis)

		// ★ プロンプト強化: 前回の「Reason（課題）」を解消できたかチェックさせる
		prompt = fmt.Sprintf(basePrompt+`

        これは前回アドバイスを受けた後の「修正版」の写真です。
        以下の【前回データ】と比較して評価してください。

        【前回データ】
        %s

        ### 評価のポイント:
        1. 前回の "reason"（課題）が、今回の写真で解消されているか確認してください。
        2. 解消されていれば "status": "improved" とし、"comment" で褒めてください。
        3. まだ解消されていない、あるいは別の問題が出た場合は、新しい "reason" と "advice" を出力してください。
        `, category, string(prevJSONBytes))
	}
	// prompt := fmt.Sprintf(promptTemplate, category)

	// 4. リクエスト実行
	resp, err := genModel.GenerateContent(ctx,
		genai.ImageData(strings.TrimPrefix(finalMediaType, "image/"), imageBytes),
		genai.Text(prompt),
	)
	if err != nil {
		return nil, fmt.Errorf("コンテンツ生成リクエストに失敗: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini APIからの応答が空です")
	}

	// 5. 応答解析 (文字列操作不要で直接Unmarshal可能)
	var result model.CompositionAnalysis

	responseText := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText = string(txt)
			break
		}
	}

	// ここで strings.Trim などの処理は不要になります
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		return nil, fmt.Errorf("分析結果のjsonパースに失敗しました: %w", err)
	}

	// カテゴリの補完
	if result.Category == "" {
		result.Category = category
	}

	return &result, nil
}
