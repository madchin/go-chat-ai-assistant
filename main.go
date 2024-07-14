package main

import (
        "context"
        "encoding/json"
        "fmt"
        "io"
	"os"
	"flag"
        "cloud.google.com/go/vertexai/genai"
)

func generateContentFromText(w io.Writer, projectID string) error {
        location := "us-central1"
        modelName := "gemini-1.5-flash-001"

        ctx := context.Background()
        client, err := genai.NewClient(ctx, projectID, location)
        if err != nil {
                return fmt.Errorf("error creating client: %w", err)
        }
        gemini := client.GenerativeModel(modelName)
        prompt := genai.Text(
                "What's a good name for a flower shop that specializes in selling bouquets of dried flowers?")

        resp, err := gemini.GenerateContent(ctx, prompt)
        if err != nil {
                return fmt.Errorf("error generating content: %w", err)
        }
        // See the JSON response in
        // https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerateContentResponse.
        rb, err := json.MarshalIndent(resp, "", "  ")
        if err != nil {
                return fmt.Errorf("json.MarshalIndent: %w", err)
        }
        fmt.Fprintln(w, string(rb))
        return nil
}

func main() {
	googleCloudProjectId := flag.String("project-id", "", "google cloud project id used for communication with Gemini AI model")
	flag.Parse()
	fmt.Printf("%s google cloud projec tid value ", googleCloudProjectId)
	if *googleCloudProjectId == "" {
		fmt.Println("In order to run POC chat model you need to specify, google cloud project id ")
		return
	}
	
	err := generateContentFromText(os.Stdout, *googleCloudProjectId)
	if err != nil {
		fmt.Println("error occured %v", err)
	}
}
