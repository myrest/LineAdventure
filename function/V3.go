package function

import (
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	openai "github.com/sashabaranov/go-openai"
)

var Adv_V3 AdventureVersionSetting

func init() {
	Adv_V3 = AdventureVersionSetting{
		SystemPrompt: "你是一個使用繁體中文的文字冒險遊戲系統，符合以下功能。\n 1.描述故事內容\n 2.提供動作選項\n 3.每個選項在10個中文字以內\n 4.每次提供的選項數量不超過3個\n 5.選項以阿拉伯數字開頭\n 6.故事內容或提供的動作選項符合故事背景\n 7.故事要連貫、要曲折離奇、高潮迭起\n 8.玩家所以使用的任何物品，必需在遊戲中取得才能使用\n 9.在達成故事的目標後，結束遊戲。\n 10.玩家的動作若不符合故事背景，直接拒絕。\n",
		ChatCompletionRequestFirst: openai.ChatCompletionRequest{
			Model:            openai.GPT3Dot5Turbo0301,
			MaxTokens:        500,
			Temperature:      1,
			TopP:             1,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
		},
		ChatCompletionRequestNormal: openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo0301,
			MaxTokens:   400,
			Temperature: 0.7,
			TopP:        0.5,
		},
	}
}

func Adventure_V3_InitUser(userid string) {
	if _, ok := UsersState[userid]; ok {
		delete(UsersState, userid)
	}

	UsersState[userid] = UserConfig{
		ChatHistory: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: Adv_V1.SystemPrompt,
			}},
		ConversationCount: 0,
		CurrentVersion:    "Adventure_V3",
	}
	InDialogs[userid] = false
}

func Adventure_V3_Command(bot *linebot.Client, event *linebot.Event, userMessage string) {
	var adv AdvanceResponse
	var err error
	if !InDialogs[event.Source.UserID] {
		InDialogs[event.Source.UserID] = true
		go func() {
			userSetting := UsersState[event.Source.UserID]
			StopGame := false
			if userSetting.ConversationCount > 15 {
				userMessage = fmt.Sprint(userMessage, ",請總結遊戲內容，並結束遊戲。")
				StopGame = true
			}
			fmt.Println("V3:", userMessage)
			ChatCompletionRequest := Adv_V1.ChatCompletionRequestFirst
			if userSetting.ConversationCount > 0 {
				ChatCompletionRequest = Adv_V1.ChatCompletionRequestNormal
			}

			timeout := 2 * time.Minute // 設置2分鐘超時時間
			done := make(chan bool)

			go func() {
				// 執行你的goroutine任務
				adv, err = DoAnAdventureV2(event.Source.UserID, userMessage, ChatCompletionRequest)
				if len(adv.Items) < 1 {
					adv.Body = adv.Body + "\n(遊戲真的結束了)"
					StopGame = true
				}
				if err != nil {
					SendLineTextMessage(bot, event, fmt.Sprintln("出錯了，請再試一次。", err.Error()))
				} else {
					SendLineTextMessageWithQuickReply(bot, event, adv)
				}
				if StopGame {
					delete(UsersState, event.Source.UserID)
				}
				InDialogs[event.Source.UserID] = false
				done <- true // 任務完成後向done通道發送一個true值
			}()
			select {
			case <-done:
				// 任務在超時之前完成
				InDialogs[event.Source.UserID] = false
			case <-time.After(timeout):
				adv = AdvanceResponse{
					Body:  "ChatGPT想太久了。。。。。\n你可以選擇再送一次答案，或是重新開始。",
					Items: []string{userMessage, "開始文字冒險"},
				}
				InDialogs[event.Source.UserID] = false
			}

		}()
	} else {
		fmt.Println("ChatGPT還沒回覆呢！命令：", userMessage, "未處理")
	}
}
