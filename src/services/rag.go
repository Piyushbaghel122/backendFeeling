package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pinecone-io/go-pinecone/v6/pinecone"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

// RunRAGPipeline sets up a Chroma vector store, adds documents, and runs a RetrievalQA chain
func RunRAGPipeline(prompt string) (string, error) {
	ctx := context.Background()

	// Initialize LLM
	llm, err := googleai.New(ctx, 
		googleai.WithAPIKey(os.Getenv("API_GEMINI")),
		googleai.WithDefaultEmbeddingModel("text-embedding-004"),
	)
	if err != nil {
		log.Println("Error creating LLM client:", err)
		return "", err
	}
	// Wrap the Gemini LLM in an Embedder so we can pass it to Chroma
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Println("Error creating embedder:", err)
		return "", err
	}
	
	// Create a Vector store
	store, err := chroma.New(
		chroma.WithChromaURL("http://localhost:8080"),
		chroma.WithEmbedder(embedder),
	)
	if err != nil {
		log.Println("Error creating Chroma store:", err)
		return "", err
	}

	// Add sample documents
	documents := []schema.Document{
		{PageContent: "Paris is the capital of France"},
		{PageContent: "London is the capital of England"},
	}

	_, err = store.AddDocuments(ctx, documents)
	if err != nil {
		log.Println("Error adding documents to Chroma:", err)
		return "", err
	}

	// Set up the Retrieval QA chain
	chain := chains.NewRetrievalQAFromLLM(llm, vectorstores.ToRetriever(store, 3))

	// Run the chain
	answer, err := chains.Run(ctx, chain, prompt)
	if err != nil {
		log.Println("Error running RetrievalQA chain:", err)
		return "", err
	}

	return answer, nil
}

// GenerateResponse calls the LLM with the given prompt
func GenerateResponse(prompt string) (string, error) {
	ctx := context.Background()
	
	llm, err := googleai.New(ctx, 
		googleai.WithAPIKey(os.Getenv("API_GEMINI")),
		googleai.WithDefaultEmbeddingModel("text-embedding-004"),
	)
	if err != nil {
		log.Println("Error creating LLM client:", err)
		return "", err
	}

	completion, err := llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	})
	if err != nil {
		log.Println("Error generating content:", err)
		return "", err
	}
 
	if len(completion.Choices) > 0 {
		return completion.Choices[0].Content, nil
	}

	return "", fmt.Errorf("no completion choices returned")
}

// QueryPinecone embeds a user's query and searches the Pinecone index "cohort-2-rag"
func QueryPinecone(queryText string) (string, error) {
	ctx := context.Background()

	// 1. Initialize Gemini to create embeddings (using your existing Gemini setup)
	llm, err := googleai.New(ctx, 
		googleai.WithAPIKey(os.Getenv("API_GEMINI")),
		googleai.WithDefaultEmbeddingModel("text-embedding-004"),
	)
	if err != nil {
		return "", fmt.Errorf("error creating LLM client: %v", err)
	}
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return "", fmt.Errorf("error creating embedder: %v", err)
	}

	// 2. Embed the query text
	queryEmbeddings, err := embedder.EmbedQuery(ctx, queryText)
	if err != nil {
		return "", fmt.Errorf("error embedding query: %v", err)
	}

	// 3. Initialize the Pinecone client locally
	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: os.Getenv("PINECONE_API_KEY"),
	})
	if err != nil {
		return "", fmt.Errorf("error creating pinecone client: %v", err)
	}

	// 4. Look up the index host
	indexDetails, err := pc.DescribeIndex(ctx, "cohort-2-rag")
	if err != nil {
		return "", fmt.Errorf("error describing index: %v", err)
	}

	// 5. Connect to the specific Pinecone index using the Host
	idx, err := pc.Index(pinecone.NewIndexConnParams{
		Host: indexDetails.Host,
	})
	if err != nil {
		return "", fmt.Errorf("error connecting to index: %v", err)
	}

	// 4. Query the index
	res, err := idx.QueryByVectorValues(ctx, &pinecone.QueryByVectorValuesRequest{
		Vector:          queryEmbeddings,
		TopK:            2,
		IncludeMetadata: true,
	})
	if err != nil {
		return "", fmt.Errorf("error querying pinecone: %v", err)
	}

	// Return the result as a formatted string (you can parse this in your controller!)
	return fmt.Sprintf("Found %d results: %+v", len(res.Matches), res.Matches), nil
}
