package function

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	openai "github.com/sashabaranov/go-openai"
)

type AdventureVersionSetting struct {
	SystemPrompt                string
	ChatCompletionRequestFirst  openai.ChatCompletionRequest
	ChatCompletionRequestNormal openai.ChatCompletionRequest
}

type UserConfig struct {
	ChatHistory       []openai.ChatCompletionMessage
	ConversationCount int
	CurrentVersion    string
}

var InDialogs = make(map[string]bool)
var UsersState = make(map[string]UserConfig)

func Adventure(w http.ResponseWriter, req *http.Request) {
	//開發版
	LINE_CHANNEL_SECRET := os.Getenv("LINE_CHANNEL_SECRET")
	LINE_ACCESS_TOKEN := os.Getenv("LINE_ACCESS_TOKEN")
	bot, err := linebot.New(
		LINE_CHANNEL_SECRET,
		LINE_ACCESS_TOKEN,
	)
	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	// 解析POST請求中的JSON資料，以獲取Line訊息的內容
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				inputText := strings.ToLower(message.Text)
				//取得使用者資料
				switch inputText {
				case "help", "?":
					SendLineTextMessageWithQuickReply(bot, event, Help())
					break
				case "自由的文字冒險":
					Adventure_V1_InitUser(event.Source.UserID)
					Adventure_V1_Command(bot, event, "請提供故事背景及場景")
				case "開始故事接龍":
					Adventure_V2_InitUser(event.Source.UserID)
					Adventure_V2_Command(bot, event, "請提供故事開頭")
				case "開始文字冒險":
					Adventure_V3_InitUser(event.Source.UserID)
					Adventure_V3_Command(bot, event, "請提供故事背景並提供可以執行動作的選項")
				default:
					userSetting := UsersState[event.Source.UserID]
					if _, ok := UsersState[event.Source.UserID]; ok {
						switch userSetting.CurrentVersion {
						case "Adventure_V1":
							Adventure_V1_Command(bot, event, message.Text)
							break
						case "Adventure_V2":
							Adventure_V2_Command(bot, event, message.Text)
							break
						case "Adventure_V3":
							Adventure_V3_Command(bot, event, message.Text)
							break
						}
					} else {
						SendLineTextMessage(bot, event, "請輸入「Help」或是「?」，以取得使用說明。")
					}
				}
			}
		}
	}
}

func SendLineTextMessage(bot *linebot.Client, event *linebot.Event, reply string) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		log.Print(err)
	}
}

func SendLineTextMessageWithQuickReply(bot *linebot.Client, event *linebot.Event, adv AdvanceResponse) {
	var quickReplyButtons []*linebot.QuickReplyButton
	if len(adv.Items) > 0 {
		for _, item := range adv.Items {
			quickReplyButtons = append(quickReplyButtons, linebot.NewQuickReplyButton("", linebot.NewMessageAction(item, item)))
		}

		message := linebot.NewTextMessage(adv.Body).
			WithQuickReplies(linebot.NewQuickReplyItems(quickReplyButtons...))

		if _, err := bot.ReplyMessage(event.ReplyToken, message).Do(); err != nil {
			log.Print(err)
		}
	} else {
		SendLineTextMessage(bot, event, adv.Body)
	}
}
