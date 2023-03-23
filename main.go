package main

import (
	"log"
	"net/http"
	"os"

	fff "rest.com.tw/ChatGPTBot/function"
)

func main() {
	//確認環境變數
	if os.Getenv("LINE_CHANNEL_SECRET") == "" || os.Getenv("LINE_ACCESS_TOKEN") == "" || os.Getenv("OPEN_API_authToken") == "" {
		log.Fatal("環境變數未設定。")
	} else {
		http.HandleFunc("/callback", fff.Adventure)

		// 啟動HTTP伺服器以處理Line的訊息
		if err := http.ListenAndServe(":9998", nil); err != nil {
			log.Fatal(err)
		}
	}
}
