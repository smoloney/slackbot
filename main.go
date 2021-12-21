package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/slack-go/slack"
)

var token = os.Getenv("SLACK_AUTH_TOKEN")
var appToken = os.Getenv("SLACK_APP_TOKEN")

func main() {
	fmt.Println("Hello world")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello", ServeHTTP)
	router.HandleFunc("/action-complete", actionComplete).Queries("id", "{id:[a-zA-Z0-9]+}")
	router.HandleFunc("/health", healthCheckHandler)

	http.ListenAndServe("0.0.0.0:5000", router)
	fmt.Println("Listening on port :5000")
}

func actionComplete(w http.ResponseWriter, r *http.Request) {
	api := slack.New(token)
	attachment := slack.Attachment{
		Text:       "Foobar i am santa",
		Fallback:   "We don't currently support your client",
		CallbackID: "accept_or_reject",
		Color:      "#3AA3E3",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "accept",
				Text:  "Accept",
				Type:  "button",
				Value: "accept",
			},
			slack.AttachmentAction{
				Name:  "reject",
				Text:  "Reject",
				Type:  "button",
				Value: "reject",
				Style: "danger",
			},
		},
	}

	_, _, err := api.PostMessage("C02MV7GBU22", slack.MsgOptionAttachments(attachment))
	if err != nil {
		fmt.Errorf("failed to post message %w", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func helloMessage() {

	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	attachment := slack.Attachment{}
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: time.Now().String(),
		},
	}

	attachment.Actions = []slack.AttachmentAction{
		{
			Name: "Test button",
			Text: "Cancel",
			Type: "button",
		},
	}

	attachment.Text = fmt.Sprintf("Hello world")
	attachment.Color = "#4af030"

	_, _, err := client.PostMessage("C02MV7GBU22", slack.MsgOptionAttachments(attachment))
	if err != nil {
		fmt.Errorf("failed to post message %w", err)
	}
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r)
	// fmt.Println(r.Body)
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}
	// fmt.Println(payload)
	fmt.Printf("Message button pressed by user %s with value %s", payload.User.Name, payload.Value)
	// fmt.Printf("callback id: %s", payload.CallbackID)
	// fmt.Printf("responseurl: %s", payload.ResponseURL)
	// fmt.Println("Original message:")
	// fmt.Println(payload.OriginalMessage)
	fmt.Printf("\nText:")
	fmt.Println(payload.OriginalMessage.Msg.Text)
	// fmt.Println("action callback")
	// fmt.Println(payload.ActionCallback)

	// var jsonStr = []byte(`{"ref":"main"}`)

	// req, err := http.NewRequest("POST", "https://api.github.com/repos/smoloney-org/hello-world/actions/workflows/15727762/dispatches", bytes.NewBuffer(jsonStr))
	// if err != nil {
	// 	// handle err
	// }
	// req.Header.Set("Accept", "application/vnd.github.v3+json")
	// req.Header.Set("Authorization", "token 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")

	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	// handle err
	// }
	// defer resp.Body.Close()
	// // log.Println(r.Body)
	// // if r.Method != http.MethodPost {
	// // 	log.Printf("[ERROR] Invalid method: %s", r.Method)
	// // 	w.WriteHeader(http.StatusMethodNotAllowed)
	// // 	return
	// // }

	// // buf, err := ioutil.ReadAll(r.Body)
	// // log.Println(buf)
	// // if err != nil {
	// // 	log.Printf("[ERROR] Failed to read request body: %s", err)
	// // 	w.WriteHeader(http.StatusInternalServerError)
	// // 	return
	// // }
	// // jsonStr, err := url.QueryUnescape(string(buf)[8:])

	// // var message slack.InteractionCallback

	// // if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
	// // 	log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
	// // 	w.WriteHeader(http.StatusInternalServerError)
	// // 	return
	// // }
	// // if message.Token != h.verificationToken {
	// // 	log.Printf("[ERROR] Invalid token: %s", message.Token)
	// // 	w.WriteHeader(http.StatusUnauthorized)
	// // 	return
	// // }
	// // log.Println(message)
	// // log.Printf(message.Name)
	// // if message.Name == "productionButton" {
	// // 	handleHelloCommand(w, r)
	// // 	// helloMessage()
	// // 	return
	// // } else {
	// // 	log.Printf("[ERROR] ]Invalid action was submitted: %s", message.Name)
	// // 	w.WriteHeader(http.StatusInternalServerError)
	// // 	return
	// // }

}
