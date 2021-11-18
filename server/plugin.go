package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"

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
type Body struct {
	Messages []string `json:"messages"`
}
type Output struct {
	Sequence string    `json:"sequence"`
	Labels   []string  `json:"labels"`
	Scores   []float64 `json:"scores"`
}

func Predict(message string, candidates string) (string, float64) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials("AKIA4H5E5ZK6CZIFSBEN", "iXx0GqU+O8MnupcB8aZaeuZRT6e7AC+pJwTC+77u", ""),
	})
	if err != nil {
		fmt.Println(err)
	}

	svc := sagemakerruntime.New(sess)

	body := `{
    "inputs": "%s",
    "parameters": {
        "candidate_labels": [%w]
    }
}`
	body = strings.Replace(body, "%s", message, 1)
	body = strings.Replace(body, "%w", candidates, 1)

	params := sagemakerruntime.InvokeEndpointInput{}
	// params.SetAccept("application/json")
	params.SetContentType("application/json")
	params.SetBody([]byte(body))
	params.SetEndpointName("huggingface-pytorch-inference-2021-11-02-19-38-58-852")

	req, out := svc.InvokeEndpointRequest(&params)

	if err := req.Send(); err != nil {
		// process error
		panic(err)
	}

	var output Output

	if err := json.Unmarshal(out.Body, &output); err != nil {
		panic(err)
	}

	return output.Labels[0], output.Scores[0]
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		_, errors := w.Write([]byte("Received a GET request\n"))
		if errors != nil {
			fmt.Println(errors)
		}
	case "POST":
		reqBody, err := ioutil.ReadAll(r.Body)
		fmt.Println(string(reqBody))
		if err != nil {
			fmt.Println(err)
		}
		var messages Body

		if errors := json.Unmarshal(reqBody, &messages); err != nil {
			panic(errors)
		}

		topics := make(map[string][]string)
		feilds := `"pets", "coding", "sports", "databases"`
		for _, post := range messages.Messages {
			label, score := Predict(post, feilds)
			if score > 0.5 {
				topics[label] = append(topics[label], post)
			}
		}
		topicsJSON, err := json.Marshal(topics)
		if err != nil {
			fmt.Println(err)
		}
		_, err = w.Write(topicsJSON)
		if err != nil {
			fmt.Println(err)
		}
	}
}
