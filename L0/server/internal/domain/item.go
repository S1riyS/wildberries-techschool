package domain

type Item struct {
	ChrtID      int    `json:"chrt_id" faker:"boundary_start=1000000, boundary_end=9999999"`
	TrackNumber string `json:"track_number" faker:"len=15"`
	Price       int    `json:"price" faker:"boundary_start=100, boundary_end=10000"`
	RID         string `json:"rid" faker:"uuid_hyphenated"`
	Name        string `json:"name" faker:"word"`
	Sale        int    `json:"sale" faker:"boundary_start=0, boundary_end=70"`
	Size        string `json:"size" faker:"oneof: 0, S, M, L, XL, XXL"`
	TotalPrice  int    `json:"total_price" faker:"boundary_start=50, boundary_end=5000"`
	NMID        int    `json:"nm_id" faker:"boundary_start=1000000, boundary_end=9999999"`
	Brand       string `json:"brand" faker:"oneof: Vivienne Sabo, L'Oreal, Maybelline, NYX, MAC, Fenty"`
	Status      int    `json:"status" faker:"boundary_start=100, boundary_end=299"`
}
