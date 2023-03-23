package function

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	openai "github.com/sashabaranov/go-openai"
)

func appendChatGPTResponse(userid string, userMessage string, responseStr string) {
	userSetting := UsersState[userid]
	userSetting.ConversationCount++
	userSetting.ChatHistory = append(userSetting.ChatHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})
	userSetting.ChatHistory = append(userSetting.ChatHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: responseStr,
	})
	if len(userSetting.ChatHistory) > 24 {
		userSetting.ChatHistory = append(userSetting.ChatHistory[:1], userSetting.ChatHistory[2:]...)
	}
	UsersState[userid] = userSetting
}

func DoAnAdventure(userid string, userMessage string, ChatCompletionRequest openai.ChatCompletionRequest) (adv AdvanceResponse, err error) {
	userSetting := UsersState[userid]
	OPEN_API_authToken := "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"
	if os.Getenv("OPEN_API_authToken") != "" {
		OPEN_API_authToken = os.Getenv("OPEN_API_authToken")
	}
	client := openai.NewClient(OPEN_API_authToken)
	ChatCompletionMessage := append(userSetting.ChatHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})
	ChatCompletionRequest.Messages = ChatCompletionMessage
	resp, err := client.CreateChatCompletion(
		context.Background(),
		ChatCompletionRequest,
	)

	if err != nil {
		return adv, err
	}
	story := resp.Choices[0].Message.Content
	adv, err = ConvertJsonToAdv(story)
	if err != nil {
		adv.Body = story
		err = nil
	}
	appendChatGPTResponse(userid, userMessage, resp.Choices[0].Message.Content)
	return adv, err
}

func DoAnAdventureV2(userid string, userMessage string, ChatCompletionRequest openai.ChatCompletionRequest) (adv AdvanceResponse, err error) {
	userSetting := UsersState[userid]
	userSetting.ConversationCount++
	OPEN_API_authToken := os.Getenv("OPEN_API_authToken")
	client := openai.NewClient(OPEN_API_authToken)
	ChatCompletionMessage := append(userSetting.ChatHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})
	ChatCompletionRequest.Messages = ChatCompletionMessage
	resp, err := client.CreateChatCompletion(
		context.Background(),
		ChatCompletionRequest,
	)

	if err != nil {
		return adv, err
	}
	story := resp.Choices[0].Message.Content
	fmt.Println(story)
	//取得動作選項
	actions, fistOption := extractStringsAfterDot(story)
	if len(actions) > 0 {
		story = splitString(story, fistOption)
	}

	adv.Body = story
	adv.Items = actions
	if len(adv.Body) < 1 {
		adv.Body = "(請重新選擇一次)"
	}
	appendChatGPTResponse(userid, userMessage, resp.Choices[0].Message.Content)
	return adv, nil
}

func splitString(originalString string, splitBy string) string {
	// 拆解字串成切片
	parts := strings.Split(originalString, splitBy)
	return strings.Trim(parts[0], "")
}

func extractStringsAfterDot(input string) ([]string, string) {
	var result []string
	firstOptionLine := ""
	re := regexp.MustCompile(`\d+\..*`)
	matches := re.FindAllString(input, -1)
	for _, match := range matches {
		if firstOptionLine == "" {
			firstOptionLine = match
		}
		inLineAction := keepFirst20Chars(strings.Trim(match[strings.Index(match, ".")+1:], " "))
		if len(inLineAction) > 0 {
			result = append(result, inLineAction)
		}
	}
	return result, firstOptionLine
}

func keepFirst20Chars(str string) string {
	if utf8.RuneCountInString(str) <= 20 {
		return str
	}
	length := 0
	for i, r := range str {
		if length+utf8.RuneLen(r) > 40 {
			return str[:i]
		}
		length += utf8.RuneLen(r)
	}
	return ""
}
