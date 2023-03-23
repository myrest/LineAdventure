
# 用Line + ChatGPT玩文字冒險遊戲

目前處於開發階段，遇到Bug是正常的，請以平常心面對。在我的ChatGPT API，或是Google免費額度用罄之前，系統會持續運作。



## 系統運作

 - [Cloud Functions](https://cloud.google.com/functions)
 - [Go](https://go.dev/)
 - [OpenAI](https://openai.com/)
 
## 環境變數設定

在本機開發時，請先設定好以下Line及OpenAI環境變數

`LINE_CHANNEL_SECRET` 
`LINE_ACCESS_TOKEN`
可以從[Line Console](https://developers.line.biz/console/)取得

`OPEN_API_authToken`可以從[Open AI](https://platform.openai.com/account/api-keys)取得


## Run Locally

Clone the project

```bash
  git clone https://github.com/myrest/LineAdventure
```

Go to the project directory

```bash
  cd LineAdventure
```

Install dependencies

```bash
  go mod tidy
```

Start the server

```bash
  go run main.go
```


## Authors

- [Roy Tai](https://www.facebook.com/roy.tai.58)

