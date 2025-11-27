package model

// Geminiの解析結果を表すドメインモデル
type CompositionAnalysis struct {
	Advice     string      `json:"advice"`
	Category   string      `json:"category"`
	VisualCues []VisualCue `json:"visual_cues"`
}

type VisualCue struct {
	Target    string `json:"target"`
	Direction string `json:"direction"`
}
