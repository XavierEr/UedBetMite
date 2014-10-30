package uedBetDataJson

import (
	"encoding/json"
	"fmt"
	"github.com/XavierEr/UedBetMite/Model"
)

func Parse(uedBetJson []byte) (uedBetData model.UedBetData) {
	var uedBetData model.UedBetData

	err := json.Unmarshal(uedBetJson, &uedBetData)
	if err != nil {
		fmt.Println(err)
	}
	return
}
