package model

type Comparison struct {
	Reason   string `json:"reason"`
	Advice   string `json:"advice"`
	Category string `json:"category"`
}

type CompositionAnalysis struct {
	Reason     string      `json:"reason"`
	Advice     string      `json:"advice"`
	Category   string      `json:"category"`
	VisualCues []VisualCue `json:"visual_cues"`
	Evaluate   string      `json:"evaluate"`
}

type VisualCue struct {
	Target    string `json:"target"`
	Direction string `json:"direction"`
}
