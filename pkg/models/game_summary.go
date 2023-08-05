package models

type GameSummary struct {
	BestReportedTier string  `json:"bestReportedTier"`
	Confidence       string  `json:"confidence"`
	Score            float64 `json:"score"`
	Tier             string  `json:"tier"`
	Total            int     `json:"total"`
	TrendingTier     string  `json:"trendingTier"`
}
