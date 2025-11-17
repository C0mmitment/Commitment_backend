package gemini

import (
	"context"
	"fmt"
	"log"
	"strings"

	// "github.com/86shin/commit_goback/domain/service"

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
func (s *GeminiAIService) GetCompositionAdvice(ctx context.Context, imageBytes []byte, mimeType string) (string, error) {

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
		return "", fmt.Errorf("geminiクライアントの作成に失敗: %w (apiキーの設定を確認してください)", err)
	}
	defer client.Close()

	// プロンプトは元のコードと同じものを使用
	prompt := "あなたはプロの写真家です。この画像を見て、写真がもっと良くなるためのアドバイスをください。以下のルールを厳守してください。\n1. 専門用語（例：三分割法）は使わない。\n2. 「良い感じ」「もっと素敵に」のような抽象的な表現は使わない。\n3. 「何を」「どうすれば」良くなるか、具体的な行動（例：「もう少し右に寄る」「少し下から撮る」）を指示する。\n4. 「人」や「物」の位置や向きに注目する。\n5. アドバイスは80文字以内。\n6. 最後に、被写体が「人」か「食事」かを判断し、[人]、[飯]、[人,飯]、[x]（どちらでもない場合）のいずれかを必ず付ける。"

	content := []genai.Part{
		genai.ImageData(finalMediaType, imageBytes),
		genai.Text(prompt),
	}

	resp, err := client.GenerativeModel("gemini-2.5-flash").GenerateContent(ctx, content...)
	if err != nil {
		return "", fmt.Errorf("コンテンツ生成リクエストに失敗: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini APIからの応答が空です")
	}

	// 応答解析
	part := resp.Candidates[0].Content.Parts[0]
	textPart, ok := part.(genai.Text)
	if !ok {
		return "", fmt.Errorf("gemini APIからの応答形式が予期されていません (応答がテキストではありません)")
	}

	advice := strings.TrimSpace(string(textPart))
	return advice, nil
}
