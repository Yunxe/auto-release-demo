package gpt

import (
	"auto-release-demo/pkg/prompts"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

var (
	API_KEY = os.Getenv("AZURE_OPENAI_APIKEY")
	HOST    = os.Getenv("AZURE_OPENAI_HOST")
)

func processContent(content string) string {
	lines := strings.Split(content, "\n")
	var processedLines []string

	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			processedLines = append(processedLines, trimmed)
		}
	}

	return strings.Join(processedLines, "\n\n")
}

func NewGPT(releaseNote string) {
	ky, err := azopenai.NewKeyCredential(API_KEY)
	if err != nil {
		log.Fatalf("error new key credential %s", err.Error())
	}

	client, err := azopenai.NewClientWithKeyCredential("https://"+HOST+".openai.azure.com", ky, nil)
	if err != nil {
		log.Fatalf("error new azure client %s", err.Error())
	}

	messages := []azopenai.ChatMessage{
		{
			Role:    to.Ptr(azopenai.ChatRoleUser),
			Content: to.Ptr(prompts.GreetingPrompts),
		},
		{
			Role:    to.Ptr(azopenai.ChatRoleUser),
			Content: to.Ptr(prompts.TemplateENPrompts),
		},
		{
			Role:    to.Ptr(azopenai.ChatRoleUser),
			Content: to.Ptr(prompts.TemplateZHPrompts),
		},
		{
			Role:    to.Ptr(azopenai.ChatRoleUser),
			Content: to.Ptr(prompts.GeneratePrompts),
		},
		{
			Role:    to.Ptr(azopenai.ChatRoleUser),
			Content: to.Ptr(releaseNote),
		},
	}

	gotReply := false

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:    messages,
		Deployment:  "gpt-4",
		Temperature: ToPointer(float32(0)),
	}, nil)
	if err != nil {
		log.Printf("error generate notes %s", err.Error())
	}

	choice := resp.Choices[0]
	gotReply = true
	content := *choice.Message.Content
	fmt.Println(content)

	if gotReply {
		fmt.Fprintf(os.Stderr, "Got chat completions reply\n")
	}
	// Split the content into English and Chinese parts
	parts := strings.Split(content, "--------")
	fmt.Printf("%q\n", parts)
	// if len(parts) < 2 {
	// 	fmt.Println("Invalid content format")
	// 	return
	// }

	englishContent := processContent(parts[0])
	chineseContent := processContent(parts[1])
	englishFilename := "CHANGELOG/CHANGELOG-v0.2.0.md"
	chineseFilename := "CHANGELOG/CHANGELOG-v0.2.0-zh.md"

	err = createMarkdownFile(englishFilename, englishContent)
	if err != nil {
		fmt.Println("Error creating English Markdown file:", err)
		return
	}

	err = createMarkdownFile(chineseFilename, chineseContent)
	if err != nil {
		fmt.Println("Error creating Chinese Markdown file:", err)
		return
	}
}

func ToPointer[T any](v T) *T {
	return &v
}

func createMarkdownFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
