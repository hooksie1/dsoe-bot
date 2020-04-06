package p

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"math/rand"

	"github.com/mailgun/mailgun-go/v3"
	"gitlab.com/hooksie1/excuses"
)

var token = os.Getenv("TOKEN")
var port = os.Getenv("PORT")
var url = os.Getenv("URL")

func buzzword() string {
        firstColumn := []string{"integrated","total","systematized","parallel","functional","responsive","optimal","synchronized","compatible","balanced"}
        secondColumn := []string{"management","organizational","monitored","reciprocal","digital","logistical","transitional","incremental","third-generation","policy"}
        thirdColumn := []string{"options","flexibility","capability","mobility","programming","concept","time-phase","projection","hardware","contingency"}

        first := rand.Intn(len(firstColumn))
        second := rand.Intn(len(secondColumn))
        third := rand.Intn(len(thirdColumn))

        f := firstColumn[first]
        s := secondColumn[second]
        t := thirdColumn[third]

        return fmt.Sprintf("%s %s %s", f, s, t) 
}

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

func buildExcuse(m Message, e excuses.Excuse) string {
	return fmt.Sprintf("%s %s %s", m.Message.From.FirstName, m.Message.From.LastName, e.Message)
}

func SendManager(excuse string) (string, error) {
	mg := mailgun.NewMailgun(os.Getenv("DOMAIN"), os.Getenv("API_KEY"))
	toAddress := os.Getenv("SEND_TO")
	m := mg.NewMessage(
		"john@hooks.technology",
		"Absense",
		excuse,
		toAddress,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}

func Bot(w http.ResponseWriter, r *http.Request) {
	
	rand.Seed(time.Now().UnixNano())

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

	if message.Message.Text == "/excuse send" || message.Message.Text == "/excuse@dsoebot send" {
		excuse := excuses.NewExcuse()
		note := buildExcuse(message, excuse)
		sendMessage(message, note, "")
		id, err := SendManager(note)
		if err != nil {
			log.Println(err)
		}
		log.Println(id)
	}
	
	if message.Message.Text == "/buzzword please" || message.Message.Text == "/buzzword@dsoebot please" {
		phrase := buzzword()
		sendMessage(message, phrase, "")
	}

}
