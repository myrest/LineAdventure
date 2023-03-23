package function

import (
	"encoding/json"
	"fmt"
)

func Help() AdvanceResponse {
	var rtn AdvanceResponse
	rtn.Body = fmt.Sprint("請參考以下關鍵字：\n--------\n",
		"【自由的文字冒險】\n開始一場華麗的文字冒險遊戲，幾乎是全自由的發揮，故事的好、壞都由你來決定。\n",
		"\n【開始故事接龍】\n系統將提供一個故事的開頭，一起來創作故事\n",
		"\n【開始文字冒險】\n開始一場文字冒險遊戲，線索就在指尖上\n",
		"--------\n因為是免費版的，所以ChatGPT回覆慢是正常的，尤其是晚上、下班後使用的人多時，會更慢。請保持耐心。")
	rtn.Items = append(rtn.Items, "自由的文字冒險", "開始文字冒險", "開始故事接龍")
	return rtn
}

type AdvanceResponse struct {
	Body  string   `json:"Body,omitempty"`
	Items []string `json:"Items,omitempty"`
}

func ConvertJsonToAdv(str string) (AdvanceResponse, error) {
	fmt.Println(str)
	var adv AdvanceResponse
	err := json.Unmarshal([]byte(str), &adv)
	return adv, err
}
