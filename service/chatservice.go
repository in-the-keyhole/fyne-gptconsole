package service

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

var apikey = "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"

var client = openai.NewClient(apikey)

var list []string

func Add(s string) {

	list = append(list, s)

}

func List() []string {

	return list
}

func Prompt(content string) string {
	//client := openai.NewClient(apikey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "Error"
	}

	fmt.Println(resp.Choices[0].Message.Content)

	r := resp.Choices[0].Message.Content

	result := strings.ReplaceAll(r, "```", "")

	return result
}
