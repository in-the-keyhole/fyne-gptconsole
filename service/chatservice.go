package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Chat struct {
	Context  string
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

	apiKey = decryptKey(string(s[:]))
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
	_, err = file.Write([]byte(encryptKey(k)))
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
		return "Error, make sure you have a valid ChatGPT API Key"
	}

	fmt.Println(resp.Choices[0].Message.Content)

	r := resp.Choices[0].Message.Content

	//result := strings.ReplaceAll(r, "```", "")

	return r
}

func encryptKey(k string) string {

	user, _ := user.Current()
	key := []byte(padKey(user.Username)) // 128-bit AES key

	// Encrypt the plaintext using AES
	ciphertext, err := encrypt(key, k)
	if err != nil {
		panic(err)
	}

	return string(ciphertext)

}

func encrypt(key []byte, plaintext string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Pad plaintext to a multiple of the block size
	paddedPlaintext := pad(plaintext, block.BlockSize())

	// Create a new AES cipher block mode encryption
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCTR(block, iv)

	// Encrypt the padded plaintext
	ciphertext := make([]byte, len(paddedPlaintext))
	stream.XORKeyStream(ciphertext, []byte(paddedPlaintext))

	return ciphertext, nil
}

func decryptKey(k string) string {

	user, _ := user.Current()
	key := []byte(padKey(user.Username)) // 128-bit AES key

	// Decrypt the ciphertext using AES
	decrypted, err := decrypt(key, []byte(k))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted message: %s\n", decrypted)

	return decrypted

}

func decrypt(key []byte, ciphertext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block mode decryption
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCTR(block, iv)

	// Decrypt the ciphertext
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	// Remove padding from plaintext
	unpaddedPlaintext := unpad(string(plaintext), block.BlockSize())

	return unpaddedPlaintext, nil
}

func pad(plaintext string, blockSize int) string {
	padding := blockSize - len(plaintext)%blockSize
	padText := fmt.Sprintf("%s%s", plaintext, strings.Repeat(string(byte(padding)), padding))
	return padText
}

func unpad(paddedText string, blockSize int) string {
	lastByte := paddedText[len(paddedText)-1]
	padding := int(lastByte)
	if padding > blockSize || padding > len(paddedText) {
		return ""
	}
	unpaddedText := paddedText[:len(paddedText)-padding]
	return unpaddedText
}

func padKey(k string) string {

	padded := strings.Repeat("0", 16-len(k)) + k

	return padded

}
