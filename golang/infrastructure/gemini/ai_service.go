package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	// "github.com/86shin/commit_goback/domain/service"

	"github.com/86shin/commit_goback/domain/model"
	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiAIService は domain/service.AIConnector インターフェースの実装です。
type GeminiAIService struct {
	APIKey string
}

// NewGeminiAIService は GeminiAIService のコンストラクタです。
func NewGeminiAIService(apiKey string) *GeminiAIService {
	return &GeminiAIService{APIKey: apiKey}
}

// GetCompositionAdvice は AIConnector のインターフェースを実装します。
func (s *GeminiAIService) GetCompositionAdvice(ctx context.Context, imageBytes []byte, mimeType string) (*model.CompositionAnalysis, error) {

	originalMimeType := strings.ToLower(mimeType)
	finalMediaType := "jpeg" // デフォルト設定

	if strings.Contains(originalMimeType, "jpeg") || strings.Contains(originalMimeType, "jpg") {
		finalMediaType = "jpeg"
	} else if strings.Contains(originalMimeType, "png") {
		finalMediaType = "png"
	} else {
		log.Printf("[MIME CRITICAL FIX] Unexpected type found: %s. Forcing MediaType to 'jpeg'.", originalMimeType)
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(s.APIKey))
	if err != nil {
		return &model.CompositionAnalysis{}, fmt.Errorf("geminiクライアントの作成に失敗: %w (apiキーの設定を確認してください)", err)
	}
	defer client.Close()

	// プロンプトは元のコードと同じものを使用
	prompt := `あなたはプロの写真家です。画像を見て構図改善のアドバイスをください。
        以下のJSONフォーマットのみを出力してください（Markdownのバッククォートは不要です）。
        {
            "advice": "80文字以内の具体的なアドバイス（専門用語禁止、命令形ではなく提案）",
            "category": "被写体の種類 ('person', 'food', 'scenery', 'other')", 
            "visual_cues": [
                {
                "target": "camera",
                "direction": "動かす方向 ('left', 'right', 'up', 'down', 'forward', 'backward')"
                }
            ]
        }
        ### 制約事項・判断基準:
        1. targetは必ず "camera" に固定してください。（被写体を動かす指示は禁止）
        2. direction は以下から選択してください：
           - "left", "right", "up", "down": カメラを上下左右に平行移動すべき場合
           - "forward": 被写体に寄るべき場合
           - "backward": 被写体から離れるべき場合
        3. 複数の指示がある場合は配列に追加してください（例: 右に移動して、寄る）。
        4. "advice"は日本語で記述し、優しい口調にしてください。
        `

	content := []genai.Part{
		genai.ImageData(finalMediaType, imageBytes),
		genai.Text(prompt),
	}

	resp, err := client.GenerativeModel("gemini-2.5-flash").GenerateContent(ctx, content...)
	if err != nil {
		return &model.CompositionAnalysis{}, fmt.Errorf("コンテンツ生成リクエストに失敗: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return &model.CompositionAnalysis{}, fmt.Errorf("gemini APIからの応答が空です")
	}

	// 応答解析
	part := resp.Candidates[0].Content.Parts[0]
	textPart, ok := part.(genai.Text)
	if !ok {
		return &model.CompositionAnalysis{}, fmt.Errorf("gemini APIからの応答形式が予期されていません (応答がテキストではありません)")
	}
	rawJSON := string(textPart)

	cleanJSON := strings.TrimSpace(rawJSON)
	cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
	cleanJSON = strings.TrimPrefix(cleanJSON, "```")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")

	// 5. 構造体へ変換 (Unmarshal)
	var result model.CompositionAnalysis
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		// 解析失敗時はエラーを返す
		return nil, fmt.Errorf("分析結果のjsonパースに失敗しました: %w, raw: %s", err, rawJSON)
	}

	// 6. 綺麗な構造体を返す
	return &result, nil
}
