package model

type ResultEvaluation struct {
	Status  string `json:"status"`  // "improved", "unchanged", "regressed", "first_time"
	Comment string `json:"comment"` // 評価コメント
}

type CompositionAnalysis struct {
	Reason     string           `json:"reason"`
	Advice     string           `json:"advice"`
	Category   string           `json:"category"`
	VisualCues []VisualCue      `json:"visual_cues"`
	Evaluation ResultEvaluation `json:"evaluation"`
}

type VisualCue struct {
	Target    string `json:"target"`
	Direction string `json:"direction"`
}
