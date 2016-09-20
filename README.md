# kbot

Yet another [telegram bot]() lib in golang.

*Work in progress*

## Usage

First, get the lib:

  $ go get github.com/diogok/kbot

Import and use it:

```golang
import (
    "github.com/diogok/kbot"
)

func handler(updt kbot.Update, outMsg chan<- kbot.OutMessage, outQuery chan<- kbot.OutQuery) {
	text := updt.Message.Text
	chatId := updt.Message.Chat.Id
  /* Do your logic */
  outMsg <- kbot.OutMessage{Chat_id: chatId, Text: "Hello, world"}
}

func main() {
  token := "Your token here"
  host := "your.domain.com"
	bot := kbot.Bot{Token: token, Host: host, Port: "80", Handler: handler}
	_, done := kbot.Start(bot)
	<-done
}
```

If "Host" is defined at the bot, it will start a server at Port and setup the telegram webhook properly. Otherwise it will listen to the updates api from telegram.

You can take a look at types.go to understand what you will receive and what you can send.

## License

MIT

