package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/tools"
)

// InternetSearchTool implements the tools.Tool interface for Langchaingo
type InternetSearchTool struct{}

func (t InternetSearchTool) Name() string {
	return "searchInternet"
}

func (t InternetSearchTool) Description() string {
	return "Use this tool to get the latest information from the internet. Input should be a search query."
}

func (t InternetSearchTool) Call(ctx context.Context, input string) (string, error) {
	log.Println("Agent is searching the internet for:", input)
	return "Simulated search results for: " + input, nil
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GenerateAgentResponse generates an agent-based response with tools
func GenerateAgentResponse(messages []Message) (string, error) {
	ctx := context.Background()

	llm, err := googleai.New(ctx,
		googleai.WithAPIKey(os.Getenv("API_GEMINI")),
		googleai.WithDefaultModel("gemini-1.5-flash"),
	)
	if err != nil {
		return "", fmt.Errorf("error creating LLM client: %v", err)
	}

	searchTool := InternetSearchTool{}

	executor, err := agents.Initialize(
		llm,
		[]tools.Tool{searchTool},
		agents.ZeroShotReactDescription,
		agents.WithMaxIterations(3),
	)
	if err != nil {
		return "", fmt.Errorf("error initializing agent: %v", err)
	}

	prompt := "You are a helpful and precise assistant for answering questions. If you don't know the answer, say you don't know. If the question requires up-to-date information, use the 'searchInternet' tool to get the latest information from the internet and then answer based on the search results.\n\n"
	for _, msg := range messages {
		prompt += msg.Role + ": " + msg.Content + "\n"
	}
	prompt += "ai:"

	res, err := chains.Run(ctx, executor, prompt)
	if err != nil {
		return "", err
	}

	return res, nil
}

// GenerateChatTitle generates a concise title for a chat
func GenerateChatTitle(message string) (string, error) {
	ctx := context.Background()
	llm, err := googleai.New(ctx,
		googleai.WithAPIKey(os.Getenv("API_GEMINI")),
		googleai.WithDefaultModel("gemini-1.5-flash"),
	)
	if err != nil {
		return "", fmt.Errorf("error creating LLM client: %v", err)
	}

	sysPrompt := "You are a helpful assistant that generates concise and descriptive titles for chat conversations. Generate a title that captures the essence of the conversation in 2-4 words. The title should be clear, relevant, and engaging."
	userPrompt := fmt.Sprintf("Generate a title for a chat conversation based on the following first message:\n\"%s\"", message)

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, sysPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userPrompt),
	}

	completion, err := llm.GenerateContent(ctx, content)
	if err != nil {
		return "", err
	}

	if len(completion.Choices) > 0 {
		return completion.Choices[0].Content, nil
	}

	return "", fmt.Errorf("no title generated")
}