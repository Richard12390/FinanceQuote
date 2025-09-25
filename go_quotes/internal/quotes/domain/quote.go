package domain

type QuoteNorm struct {
	Source       string         `json:"source"`
	AssetType    string         `json:"assetType"`
	Symbol       string         `json:"symbol"`
	SymbolOrigin string         `json:"symbolOrigin,omitempty"`
	Exchange     string         `json:"exchange,omitempty"`
	Currency     string         `json:"currency,omitempty"`
	Price        float64        `json:"price"`
	Open         float64        `json:"open,omitempty"`
	High         float64        `json:"high,omitempty"`
	Low          float64        `json:"low,omitempty"`
	PrevClose    float64        `json:"prevClose,omitempty"`
	Bid          float64        `json:"bid,omitempty"`
	BidSize      float64        `json:"bidSize,omitempty"`
	Ask          float64        `json:"ask,omitempty"`
	AskSize      float64        `json:"askSize,omitempty"`
	Volume       float64        `json:"volume,omitempty"`
	VolumeKind   string         `json:"volumeKind,omitempty"`
	MarketState  string         `json:"marketState,omitempty"`
	TsUnixMs     int64          `json:"tsUnixMs"`
	TsISO        string         `json:"tsIso"`
	TsLocal      string         `json:"tsLocal"`
	Raw          map[string]any `json:"raw,omitempty"`
}
