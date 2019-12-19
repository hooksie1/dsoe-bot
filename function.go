package p

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gitlab.com/hooksie1/excuses"
)

var token = os.Getenv("TOKEN")
var port = os.Getenv("PORT")
var url = os.Getenv("URL")

func sendMessage(m Message, s string, p string) {
	log.Println("setting up message")
	var response Response

	response.ChatID = m.Message.Chat.ID

	botURL := "https://api.telegram.org/bot" + token + "/sendMessage"

	response.Text = s

	response.ParseMode = p

	var body []byte

	body, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(body))

	req, err := http.Post(botURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}

	defer req.Body.Close()

}

func buildExcuse(m Message, e Excuse) string {
	excuse := fmt.Sprintf("%s %s %s", m.Message.From.FirstName, m.Message.From.LastName, e.Message)
}

func Bot(w http.ResponseWriter, r *http.Request) {

	log.Println("started")

	var message Message

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(body, &message)

	log.Println(string(body))
	log.Println(message)

	if message.Message.Text == "/excuse" || message.Message.Text == "/excuse@dsoebot" {
		excuse := excuses.NewExcuse()
		note := buildExcuse(message, excuse)
		sendMessage(message, note, "")
	}

}
