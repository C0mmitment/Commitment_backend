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

type GeminiAIService struct {
	Client *genai.Client
}

func NewGeminiAIService(apiKey string) (*GeminiAIService, error) {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("geminiクライアントの作成に失敗: %v", err)
		return nil, fmt.Errorf("geminiクライアントの作成に失敗: %w (apiキーの設定を確認してください)", err)
	}

	return &GeminiAIService{Client: client}, nil
}

func (s *GeminiAIService) GetCompositionAdvice(ctx context.Context, category string, imageReader io.Reader, mimeType string, prevAnalysis *model.Comparison) (*model.CompositionAnalysis, error) {

	imageBytes, err := io.ReadAll(imageReader)
	if err != nil {
		return nil, fmt.Errorf("画像ストリームの読み込みに失敗: %w", err)
	}
	if len(imageBytes) == 0 {
		return nil, fmt.Errorf("画像データが空です")
	}

	originalMimeType := strings.ToLower(mimeType)
	finalMediaType := "image/jpeg"

	if strings.Contains(originalMimeType, "png") {
		finalMediaType = "image/png"
	}

	genModel := s.Client.GenerativeModel("gemini-2.5-flash")

	genModel.SafetySettings = []*genai.SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockNone},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockNone},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockNone},
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockNone},
	}
	genModel.SetTemperature(0.4)
	genModel.SetMaxOutputTokens(3000)

	genModel.ResponseMIMEType = "application/json"
	genModel.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"reason": {
				Type:        genai.TypeString,
				Description: "35文字以内。簡潔に。",
			},
			"advice": {
				Type:        genai.TypeString,
				Description: "50文字以内。画面の線を使った具体的指示。",
			},
			"category": {Type: genai.TypeString},
			"visual_cues": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"target": {Type: genai.TypeString, Enum: []string{"camera"}},
						"direction": {Type: genai.TypeString, Enum: []string{
							"left", "right", "up", "down", "forward", "backward",
							"up-left", "up-right", "down-left", "down-right",
						}},
					},
					Required: []string{"target", "direction"},
				},
			},
			"evaluate": {
				Type: genai.TypeString,
				Enum: []string{"improved", "unchanged", "regressed", "first_time"},
			},
		},
		Required: []string{"reason", "advice", "category", "visual_cues", "evaluate"},
	}

	const systemInst = `役割: 初心者向けプロ写真コーチ
		UI環境: ユーザーの画面には【三分割法のグリッド線(縦2本・横2本)】が表示されている。

		# 用語定義 (これ以外使うな)
		- 「右の縦線」「左の縦線」
		- 「上の横線」「下の横線」
		- 「交点」: 線が交わる4つの点
		- 「画面の中央」: 線と線の間のスペース
		- ★禁止用語: 「真ん中の線」(存在しない)、「グリッド」「三分割法」(専門用語)

		# 診断フロー (NGがあれば即回答)

		1. 【光】チェック (最優先・絶対基準)
		- 暗い / 顔に影 / 逆光 / のっぺり ならNG。
		- NGなら構図の話はせず、光の修正だけを指示する。
		- 例: 「顔が暗いです。明るい方(窓の方)を向いて」

		2. 【構図】チェック (光が合格点の場合のみ)
		- 基準: 被写体を「線の上」か「交点」、または「画面の中央」のいずれか最適な場所に配置する。
		- 基本的に「日の丸構図(ど真ん中)」は避け、「左右の線」や「交点」を優先して検討すること。
		- ただし、左右対称や集合写真など、中央が最適な場合のみ「画面の中央」を指示する。

		【指示のルール】
		- 画面に見えている「線」や「点」の位置関係で具体的に指示する。
		- 「画面に収める」だけの指示は禁止（配置のバランスを指摘すること）。
        1. 位置(X/Y軸)の指示:
           - 「右」かつ「上」の修正が必要な場合は、別々に出さず「右上(up-right)」のように斜めの指示を1つ出すこと。
           - 常に主役を「交点」や「線」に最短距離で導く方向を選ぶこと。

        2. 距離(Z軸)の指示:
           - 被写体が遠すぎる/近すぎる場合は、「前(forward)」「後ろ(backward)」を追加してよい。
           - 例: 位置合わせで「右(right)」＋ サイズ調整で「前(forward)」＝ 2つのcueを出力。

        3. 出力配列(visual_cues)の制限:
           - 最大2つまで（方向1つ ＋ 距離1つ）。
           - 矢印だらけにしてユーザーを混乱させないこと。

		【良いアドバイスの例】
		- ○ 「右の縦線に顔が重なるように動いて」
		- ○ 「地平線を下側の横線に合わせてみて」
		- ○ 「画面の中央に、二人がバランスよく収まるように」

		# 出力制約
		- reason: 35文字以内
		- advice: 50文字以内 (「線に合わせて」等、直感的に)
		- カテゴリ: %s`

	var prompt string
	if prevAnalysis == nil {
		prompt = fmt.Sprintf(systemInst+`
        # 状況: 初回撮影
        evaluate: "first_time" を選択せよ。`, category)
	} else {
		prompt = fmt.Sprintf(systemInst+`
        # 状況: 再撮影(前回比較)
        - 前回課題: %s
        - 前回助言: %s
        - 前回カテゴリ: %s
        
        # 判定ルール:
        1. 改善なら "improved"、変化なしなら "unchanged"、悪化なら "regressed"。
        2. 前回のアドバイス通りに動けているか厳しく判定せよ。`,
			category,
			prevAnalysis.Reason,
			prevAnalysis.Advice,
			prevAnalysis.Category,
		)
	}

	resp, err := genModel.GenerateContent(ctx,
		genai.ImageData(strings.TrimPrefix(finalMediaType, "image/"), imageBytes),
		genai.Text(prompt),
	)
	if err != nil {
		return nil, fmt.Errorf("APIリクエスト失敗")
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("候補(Candidates)がゼロです")
	}

	candidate := resp.Candidates[0]

	if candidate.FinishReason != genai.FinishReasonStop {
		log.Printf("生成が中断されました。")
		if candidate.FinishReason == genai.FinishReasonSafety {
			return nil, fmt.Errorf("セーフティフィルタにより生成がブロックされました")
		}
	}

	responseText := ""
	for _, part := range candidate.Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText = string(txt)
			break
		}
	}

	if strings.TrimSpace(responseText) == "" {
		return nil, fmt.Errorf("AIからの応答テキストが空でした")
	}
	responseText = strings.ReplaceAll(responseText, "```json", "")
	responseText = strings.ReplaceAll(responseText, "```", "")

	var result model.CompositionAnalysis
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		log.Printf("JSON Parse Error! Raw: %s", responseText)
		return nil, fmt.Errorf("JSONパース失敗")
	}

	if result.Category == "" {
		result.Category = category
	}

	return &result, nil
}
