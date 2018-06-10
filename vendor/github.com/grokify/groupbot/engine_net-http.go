package groupbot

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ServeNetHttp(intentRouter IntentRouter) {
	bot := Groupbot{}
	bot.Initialize()
	bot.IntentRouter = intentRouter

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", http.HandlerFunc(bot.HandleNetHTTP))
	mux.HandleFunc("/webhook/", http.HandlerFunc(bot.HandleNetHTTP))

	log.Info(fmt.Sprintf("Starting server on port [%v]", bot.AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", bot.AppConfig.Port), mux))
}
