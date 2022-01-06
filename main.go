package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"github.com/slack-go/slack"
)

var token = os.Getenv("SLACK_AUTH_TOKEN")
var appToken = os.Getenv("SLACK_APP_TOKEN")
var gitToken = os.Getenv("GIT_TOKEN")

type data struct {
	Repo string
	Sha  string
}

func main() {
	outputText := strings.ReplaceAll(fmt.Sprintf("Repo: smoloney/slack-go \nSHA: 6989ffd533e41844ee19bfbc72cfe6916511790e"), " ", "")
	splitString := strings.FieldsFunc(outputText, func(r rune) bool { return strings.ContainsRune("\n:", r) })
	deployInfoMap := make(map[string]string)

	for i := 0; i <= len(splitString)-1; i += 2 {
		deployInfoMap[splitString[i]] = splitString[i+1]
	}

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
	textText := fmt.Sprintf("Repo: %s \nSHA: %s", parsedQueries["id"], parsedQueries["sha"])
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

type GithubInput struct {
	ref    string
	inputs struct {
		SHA  string
		Repo string
	}
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}

	callBackText := strings.ReplaceAll(payload.OriginalMessage.Msg.Attachments[0].Text, " ", "")
	callBackSplit := strings.FieldsFunc(callBackText, func(r rune) bool { return strings.ContainsRune("\n:", r) })
	deployInfoMap := make(map[string]string)

	for i := 0; i <= len(callBackSplit)-1; i += 2 {
		deployInfoMap[callBackSplit[i]] = callBackSplit[i+1]
	}

	githubStruct := GithubInput{
		ref: "main",
	}

	githubStruct.inputs.SHA = deployInfoMap["SHA"]
	githubStruct.inputs.Repo = deployInfoMap["Repo"]
	fmt.Println(githubStruct)
	b, err := json.Marshal(githubStruct)
	if err != nil {
		fmt.Println(err)
	}

	// var jsonStr = []byte(fmt.Sprintf("{\"ref\":\"main\", \"inputs\":{\"SHA:%s}}"")
	inputValues := fmt.Sprintf("{\"SHA\": %s, \"Repo\": %s}", deployInfoMap["SHA"], deployInfoMap["Repo"])
	requestUrl := fmt.Sprintf("https://api.github.com/repos/%s/actions/workflows/echo.yml/dispatches", deployInfoMap["Repo"])
	tokenString := fmt.Sprintf("token %s", gitToken)
	form := url.Values{}
	form.Add("ref", "main")
	form.Add("inputs", inputValues)

	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", tokenString)

	http.DefaultClient.Do(req)

}
