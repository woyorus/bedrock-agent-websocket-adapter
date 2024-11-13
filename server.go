package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var (
	once                      sync.Once
	bedrockagentRuntimeClient *bedrockagentruntime.Client
	sdkConfig                 aws.Config
	initError                 error
	upgrader                  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func initAwsBedrockAgentRuntimeClient() (*bedrockagentruntime.Client, error) {
	once.Do(func() {
		ctx := context.TODO()
		sdkConfig, initError = config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"), config.WithSharedConfigProfile(os.Getenv("AWS_SHARED_CONFIG_PROFLE")))
		if initError != nil {
			log.Printf("Error loading default config: %v", initError)
			return
		}

		bedrockagentRuntimeClient = bedrockagentruntime.NewFromConfig(sdkConfig)
	})
	return bedrockagentRuntimeClient, nil
}

func processBedrock(prompt string) (string, error) {
	client, err := initAwsBedrockAgentRuntimeClient()
	if err != nil {
		return "", err
	}

	params := &bedrockagentruntime.InvokeAgentInput{
		AgentId:      aws.String(os.Getenv("AGENT_ID")),
		AgentAliasId: aws.String(os.Getenv("AGENT_ALIAS_ID")),
		InputText:    aws.String(prompt),
		SessionId:    aws.String("agent-test-session"),
	}
	invoke, err := client.InvokeAgent(context.Background(), params)

	if err != nil {
		return "", err
	}

	var result string
	for event := range invoke.GetStream().Events() {
		switch e := event.(type) {
		case *types.ResponseStreamMemberChunk:
			result = string(e.Value.Bytes)
		}
	}

	return result, nil
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

	for {
		_, prompt, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("error reading message ", err)
			break
		}

		fmt.Println("Mesage received: ", string(prompt))
		response, err := processBedrock(string(prompt))
		if err != nil {
			log.Println("unable process bedrock ", err)
			break
		}

		err = ws.WriteMessage(websocket.TextMessage, []byte(response))
		if err != nil {
			log.Println("error writing message ", err)
			break
		}
	}

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Printf("Listening for WebSocket connections at localhost:7001/ws")

	http.HandleFunc("/ws", handleConnection)
	err = http.ListenAndServe(":7001", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
