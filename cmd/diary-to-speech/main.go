package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/MasashiFukuzawa/diary-to-speech/pkg/sections"
	"github.com/MasashiFukuzawa/diary-to-speech/pkg/speech"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type ErrSectionNotFound string

func (e ErrSectionNotFound) Error() string {
	return fmt.Sprintf("section %q not found", string(e))
}

func loadEnvVariables() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}
	return nil
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

type textToSpeechClient struct {
	client *texttospeech.Client
}

func (c *textToSpeechClient) SynthesizeSpeech(ctx context.Context, in *texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error) {
	return c.client.SynthesizeSpeech(ctx, in)
}

type osFileWriter struct{}

func (f *osFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func main() {
	fmt.Println("Starting...")

	handleError(loadEnvVariables())

	now := time.Now().Local()
	year := now.Format("2006")

	filename := fmt.Sprintf("%s.md", now.Format("01-02"))
	markdownFilepath := filepath.Join(os.Getenv("MARKDOWN_FILEPATH_BASE"), year, filename)
	markdown, err := os.ReadFile(markdownFilepath)
	handleError(err)

	if len(markdown) == 0 {
		fmt.Println("Markdown file is empty")
		os.Exit(1)
	}

	source := string(markdown)
	sectionNames := []string{
		"Simple English",
		"Advanced English",
	}
	results := make(map[string]string)
	for _, section := range sectionNames {
		result, err := sections.Extract(source, section, sectionNames)
		handleError(err)
		results[section] = result
	}

	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	handleError(err)

	defer client.Close()

	textToSpeechClient := &textToSpeechClient{client: client}

	fmt.Println("Synthesizing native speeches...")

	fileWriter := &osFileWriter{}
	for section, result := range results {
		handleError(speech.Synthesize(ctx, textToSpeechClient, fileWriter, now, section, result))
	}

	fmt.Println("Completed!")
}
