package main

type AddUserOffer struct {
	PublickKey string `json:"public_key"`
	Time       int64  `json:"timestamp"`
}

type UpgradeUserOffser struct {
	CallID string `json:"callid"`
	UserID string `json:"id"`
	Limit  int    `json:"limit"`
	Time   int64  `json:"timestamp"`
}

type GetUserOffser struct {
	CallID string `json:"callid"`
	UserID string `json:"id"`
	Time   int64  `json:"timestamp"`
}

type SendTranscationOffser struct {
	CallID string `json:"callid"`
	ToUser string `json:"to_id"`
	Coin   int    `json:"coin"`
	Time   int64  `json:"timestamp"`
}
