package ai

import (
	"context"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func MustNewClient(ctx context.Context, key string) *genai.Client {
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		log.Fatalf("failed to init genai client: %s", err)
	}
	return client
}

type contentGenerator struct {
	model *genai.GenerativeModel
}

type ContentGenerator interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}

func NewContentGenerator(client *genai.Client, modelName string, defaultPrompt string) ContentGenerator {
	model := client.GenerativeModel(modelName)
	model.SystemInstruction = genai.NewUserContent(genai.Text(defaultPrompt))
	return &contentGenerator{
		model: model,
	}
}

func (c *contentGenerator) GenerateContent(ctx context.Context, prompt string) (string, error) {
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, part := range resp.Candidates {
		if part.Content != nil {
			for _, msgPart := range part.Content.Parts {
				if text, ok := msgPart.(genai.Text); ok {
					sb.WriteString(string(text))
					sb.WriteByte('\n')
				}
			}
		}
	}

	return strings.TrimSpace(sb.String()), nil
}
