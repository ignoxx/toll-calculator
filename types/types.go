package types

type ObuData struct {
	ObuID     int     `json:"obuId"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	RequestID int     `json:"requestId"`
}
