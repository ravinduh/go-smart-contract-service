package datastruct

import (
	"encoding/json"
	"fmt"
)

type Request struct {
}

type ReceiptRequest struct {
	Data ReceiptData `json:"data"`
}

type ReceiptData struct {
	NRIC          string `json:"nric"`
	WalletAddress string `json:"wallet_address"`
}

func (r ReceiptRequest) String() string {
	byteArray, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	return string(byteArray)
}
