package gemini

import (
	"context"
	"encoding/json"
	"fmt"
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
func (s *GeminiAIService) GetCompositionAdvice(ctx context.Context, category string, imageBytes []byte, mimeType string) (*model.CompositionAnalysis, error) {

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
		},
		Required: []string{"advice", "category", "visual_cues"},
	}

	// 3. プロンプト設計 (フォーマット指示を削除した軽量版)
	promptTemplate := `あなたはプロの写真家です。画像を見て構図改善のアドバイスをください。
    カテゴリ: %s
    
    ### 判断基準:
    - directionの "forward"/"backward" はズームでも撮影者の移動でも可。
    - adviceは日本語で、専門用語を使わず優しい口調にしてください。`

	prompt := fmt.Sprintf(promptTemplate, category)

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
