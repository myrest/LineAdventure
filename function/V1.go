package function

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	openai "github.com/sashabaranov/go-openai"
)

var Adv_V1 AdventureVersionSetting

func init() {
	Adv_V1 = AdventureVersionSetting{
		SystemPrompt: "你是一個文字冒險遊戲系統，每次故事的內容限制在100字以內。由你來描述遊戲場景，等待玩家執行要採取的動怍。輸出的內容要連貫，要符合故事背景，在達成故事的目標後，結束遊戲。如果玩家的動作不符合故事背景或想要直接解開詸題，直接拒絕。",
		ChatCompletionRequestFirst: openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo0301,
			MaxTokens:   500,
			Temperature: 1,
			TopP:        1,
			// FrequencyPenalty: 2,
			// PresencePenalty:  2,
		},
		ChatCompletionRequestNormal: openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo0301,
			MaxTokens:   400,
			Temperature: 1,
			TopP:        1,
		},
	}
}

func Adventure_V1_InitUser(userid string) {
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
		CurrentVersion:    "Adventure_V1",
	}
	InDialogs[userid] = false
}

func Adventure_V1_Command(bot *linebot.Client, event *linebot.Event, userMessage string) {
	if !InDialogs[event.Source.UserID] {
		InDialogs[event.Source.UserID] = true
		go func() {
			userSetting := UsersState[event.Source.UserID]
			StopGame := false
			if userSetting.ConversationCount >= 15 {
				userMessage = fmt.Sprint(userMessage, ",請總結遊戲內容，並結束遊戲。")
				StopGame = true
			}
			fmt.Println("V1:", userMessage)
			ChatCompletionRequest := Adv_V1.ChatCompletionRequestFirst
			if userSetting.ConversationCount > 0 {
				ChatCompletionRequest = Adv_V1.ChatCompletionRequestNormal
			}
			adv, err := DoAnAdventure(event.Source.UserID, userMessage, ChatCompletionRequest)
			if err != nil {
				SendLineTextMessage(bot, event, fmt.Sprintln("出錯了，請再試一次。", err.Error()))
			} else {
				SendLineTextMessageWithQuickReply(bot, event, adv)
			}
			if StopGame {
				delete(UsersState, event.Source.UserID)
			}
			InDialogs[event.Source.UserID] = false
		}()
	} else {
		fmt.Println("ChatGPT還沒回覆呢！命令：", userMessage, "未處理")
	}
}
