package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/slack-go/slack"
)

var token = os.Getenv("SLACK_AUTH_TOKEN")
var appToken = os.Getenv("SLACK_APP_TOKEN")

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello", ServeHTTP)
	router.HandleFunc("/action-complete", actionComplete).Queries("id", "{id:[a-zA-Z0-9_/-]+}", "sha", "{sha:[a-zA-Z0-9]+}", "lastSuccessSha", "{lastSuccessSha:[a-zA-Z0-9]+}")
	router.HandleFunc("/health", healthCheckHandler)

	http.ListenAndServe("0.0.0.0:5000", router)
	fmt.Println("Listening on port :5000")
}

func queryParser(str string) map[string]string {
	values := strings.Split(str, "&")

	generateMap := make(map[string]string)
	for _, e := range values {
		parts := strings.Split(e, "=")
		generateMap[parts[0]] = parts[1]
	}

	return generateMap
}

func actionComplete(w http.ResponseWriter, r *http.Request) {
	parsedQueries := queryParser(r.URL.RawQuery)
	textText := fmt.Sprintf("Repo: %s SHA: %s", parsedQueries["id"], parsedQueries["sha"])
	titleText := fmt.Sprintf("Deployment alert for %s", parsedQueries["id"])
	fallBackText := fmt.Sprintf("Deployment to %s", parsedQueries["id"])
	api := slack.New(token)

	attachment := slack.Attachment{
		Title:      titleText,
		Text:       textText,
		Fallback:   fallBackText,
		CallbackID: "deployment",
		Color:      "#3AA3E3",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "deploy",
				Text:  "Deploy",
				Type:  "button",
				Value: "deploy",
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

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}
	// log.Println("Printing payload")
	// log.Println(payload)
	// fmt.Println("hello world1")
	// fmt.Println(payload.OriginalMessage.Msg.Attachments[0].Text)
	// fmt.Println(payload.ActionCallback.AttachmentActions.Text)

	outputText := payload.OriginalMessage.Msg.Attachments[0].Text
	fmt.Println(outputText)
	fmt.Println(strings.Split(outputText, ":"))
	// deployInfoMap := make(map[string]string)

	// generateMap := make(map[string]string)
	// for _, e := range outputText {
	// 	e := string(e)
	// 	parts := strings.Split(e, ":")
	// 	generateMap[parts[0]] = parts[1]
	// }
	// fmt.Println(deployInfoMap)

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

type JsonStruct struct {
	Commit struct {
		Author struct {
			Date string `json:"date"`
		}
	}
}

func parseJson() {

	jsonData := `{
		"url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		"node_id": "MDY6Q29tbWl0NmRjYjA5YjViNTc4NzVmMzM0ZjYxYWViZWQ2OTVlMmU0MTkzZGI1ZQ==",
		"html_url": "https://github.com/octocat/Hello-World/commit/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		"comments_url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e/comments",
		"commit": {
		  "url": "https://api.github.com/repos/octocat/Hello-World/git/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		  "author": {
			"name": "Monalisa Octocat",
			"email": "mona@github.com",
			"date": "2011-04-14T16:00:49Z"
		  },
		  "committer": {
			"name": "Monalisa Octocat",
			"email": "mona@github.com",
			"date": "2011-04-14T16:00:49Z"
		  },
		  "message": "Fix all the bugs",
		  "tree": {
			"url": "https://api.github.com/repos/octocat/Hello-World/tree/6dcb09b5b57875f334f61aebed695e2e4193db5e",
			"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		  },
		  "comment_count": 0,
		  "verification": {
			"verified": false,
			"reason": "unsigned",
			"signature": null,
			"payload": null
		  }
		}
	  }
	`

	var jsonValues JsonStruct

	json.Unmarshal([]byte(jsonData), &jsonValues)
	dateTime := jsonValues.Commit.Author.Date
	layout := "2010-01-01T01:00:00z"
	t, err := time.Parse(layout, dateTime)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
}
