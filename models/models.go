package models

type WalletReq struct {
	UUID      string `json:"uuid"`
	Operation string `json:"operation_type"`
	Amount    string `json:"amount"`
}

type WalletResp struct {
	Amount string `json:"amount"`
}
