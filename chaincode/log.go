package main

type TranscationLog struct {
	Index     int    `json:"index"`
	FromID    string `json:"from"`
	ToID      string `json:"to"`
	Coin      string `json:"coin"`
	Timestamp string `json:"time"`
}
