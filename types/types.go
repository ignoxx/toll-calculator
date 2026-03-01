package types

type ObuData struct {
	ObuID     int     `json:"obuId"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	RequestID int     `json:"requestId"`
}

type Distance struct {
	Value     float64 `json:"value"`
	ObuID     int     `json:"obuId"`
	Timestamp int64   `json:"timestamp"`
}
