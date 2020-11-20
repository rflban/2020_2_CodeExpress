package models

type Search struct {
	Albums  []*Album  `json:"albums"`
	Artists []*Artist `json:"artists"`
	Tracks  []*Track  `json:"tracks"`
}
