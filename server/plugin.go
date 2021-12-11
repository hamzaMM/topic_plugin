package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	configuration *configuration
	// configurationLock synchronizes access to the configuration.
	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	plugin.MattermostPlugin

	configurationLock sync.RWMutex
}

type Output struct {
	Sequence string    `json:"sequence"`
	Labels   []string  `json:"labels"`
	Scores   []float64 `json:"scores"`
}
type User struct {
	User   string `json:"user"`
	Labels string `json:"labels"`
}

var kvList []string

func (p *Plugin) Predict(message string, labels string) (string, float64) {
	body := `{
		"inputs": "%s",
		"parameters": {
			"candidate_labels": [%w],
			"multi_label": "True"
		}
	}`
	body = strings.Replace(body, "%s", message, 1)
	body = strings.Replace(body, "%w", labels, 1)

	var jsonData = []byte(body)

	req, _ := http.NewRequest("POST", "https://api-inference.huggingface.co/models/facebook/bart-large-mnli", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer api_xKgQCGwfjYCyYzeKkdxiZwhxpBwWNjYuaq")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var output Output

	_ = json.Unmarshal(response, &output)
	fmt.Println(output.Labels[0])

	return output.Labels[0], output.Scores[0]
}

func (p *Plugin) FilterPost(post *model.Post, labels string, email string) (*model.Post, string) {
	postMessageWithoutAccents := post.Message

	words := strings.Fields(postMessageWithoutAccents)

	postMessageWithoutAccents = strings.Join(words, " ")
	label, score := p.Predict(postMessageWithoutAccents, labels)
	fmt.Println("label: ")
	fmt.Println(label)

	// If class score more than 0.5 it belongs to the topic. This limit can be adjusted
	if score > 0.5 {
		message := label + ": " + post.Message
		err := p.API.SendMail(email, "Topic Found", message)
		if err != nil {
			fmt.Print(err)
		}
	}

	return nil, ""
}

func (p *Plugin) OnActive() error {
	kvList, _ = p.API.KVList(0, 10000)
	return nil
}

func (p *Plugin) MessageHasBeenPosted(_ *plugin.Context, post *model.Post) {
	fmt.Println(kvList)
	for _, id := range kvList {
		user, err := p.API.GetUser(id)
		if err != nil {
			fmt.Println(err)
		}
		labels, _ := p.API.KVGet(id)
		if post.UserId != user.Id {
			p.FilterPost(post, string(labels), user.Email)
		}
	}
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Println(kvList)
	w.Header().Set("Content-Type", "application/json")
	switch r.URL.Path {
	// gets topics user has subscribed too
	case "/topics":
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		kv, _ := p.API.KVGet(string(reqBody))
		tempJSON, err := json.Marshal(string(kv))
		if err != nil {
			fmt.Println(err)
		}
		_, _ = w.Write(tempJSON)

	// User has added or changed their topics
	case "/add_topics":
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("received: ")
		fmt.Println(string(reqBody))
		var user User
		_ = json.Unmarshal(reqBody, &user)
		fmt.Println(user.Labels)
		_ = p.API.KVSet(user.User, []byte(user.Labels))

		kvList, _ = p.API.KVList(0, 10000)
	}
}
