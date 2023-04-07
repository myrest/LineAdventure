package function

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	openai "github.com/sashabaranov/go-openai"
)

var Adv_V2 AdventureVersionSetting

func init() {
	Adv_V2 = AdventureVersionSetting{
		SystemPrompt: "你是一個故事接龍系統，以第三人的角度描述故事，每次產生故事的內容限制在100字以內。不要出現重複或類似的場景、對話，輸出的內容要連貫，要符合故事背景，如果場景中的人物在對話，請把對話內容完整輸出来，如果角色與場景中的任何生物、物品互動，請把互動過程詳細描述出來",
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
			MaxTokens:   500,
			Temperature: 1,
			TopP:        1,
		},
	}
}
func Adventure_V2_InitUser(userid string) {
	if _, ok := UsersState[userid]; ok {
		delete(UsersState, userid)
	}

	UsersState[userid] = UserConfig{
		ChatHistory: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: Adv_V2.SystemPrompt,
			}},
		ConversationCount: 0,
		CurrentVersion:    "Adventure_V2",
	}
	InDialogs[userid] = false
}

func Adventure_V2_Command(bot *linebot.Client, event *linebot.Event, userMessage string) {
	if !InDialogs[event.Source.UserID] {
		InDialogs[event.Source.UserID] = true
		go func() {
			userSetting := UsersState[event.Source.UserID]
			StopGame := false
			if userSetting.ConversationCount > 20 {
				userMessage = fmt.Sprint(userMessage, ",請總結對話內容，並結束遊戲。")
				StopGame = true
			}
			fmt.Println("V2:", userMessage)
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
