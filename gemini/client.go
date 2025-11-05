package gemini

import (
	"context"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateContent(apiKey, prompt string) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer client.Close() // self note: defer is used to close a client or a file, it delays the action unless the function returns

	model := client.GenerativeModel("gemini-2.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(textPart), nil
	}

	return "Response not a text", err
}
