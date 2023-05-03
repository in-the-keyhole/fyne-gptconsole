package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type Chat struct {
	Prompt   string
	Response string
}

var chatList []Chat

//var apikey = "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"

var apiKey string = ""

//var client = openai.NewClient(apikey)

var list []string

func ApiKeyExists() bool {

	return !(apiKey == "")
}

func ApiKey() string {

	return apiKey

}

func Save(key string) {

}

func Add(s string) {

	list = append(list, s)

}

func exists(filename string) bool {

	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
		fmt.Printf("File '%s' does not exist\n", filename)
	}

	return true

}

func List() []string {

	return list
}

func Write(l []Chat) {

	// convert the person struct to JSON
	jsonBytes, err := json.Marshal(l)
	if err != nil {
		fmt.Println(err)
		return
	}

	// open a file to write to
	file, err := os.Create("gptconsole.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// write the JSON data to the file
	_, err = file.Write(jsonBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	// close the file
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Read() []Chat {

	if !exists("gptconsole.json") {

		return make([]Chat, 0)

	}

	// open the file
	jsonFile, err := ioutil.ReadFile("gptconsole.json")
	if err != nil {
		fmt.Println(err)
		return nil

	}

	// parse the JSON data into a new person struct
	var l []Chat
	err = json.Unmarshal(jsonFile, &l)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// print out the person's name and age
	//fmt.Println("Name:", p.Name)
	//fmt.Println("Age:", p.Age)

	return l

}

func ReadKey() string {

	if !exists("gptconsole.key") {

		return ""

	}

	// open the file
	s, err := ioutil.ReadFile("gptconsole.key")
	if err != nil {
		fmt.Println(err)
		return ""

	}

	apiKey = string(s[:])
	return apiKey

}

func WriteKey(k string) {

	// open a file to write to
	file, err := os.Create("gptconsole.key")
	if err != nil {
		fmt.Println(err)
		return
	}

	// write the JSON data to the file
	_, err = file.Write([]byte(k))
	if err != nil {
		fmt.Println(err)
		return
	}

	// close the file
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	apiKey = k
}

func Prompt(content string) string {
	client := openai.NewClient(apiKey)
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

	//result := strings.ReplaceAll(r, "```", "")

	return r
}
