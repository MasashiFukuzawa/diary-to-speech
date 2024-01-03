package main

import (
	"context"
	"flag"
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

func main() {
	fmt.Println("Starting...")

	if err := loadEnvVariables(); err != nil {
		fmt.Printf("Error loading env variables: %v\n", err)
		os.Exit(1)
	}

	date, err := parseDateFlag()
	if err != nil {
		fmt.Printf("Error parsing date flag: %v\n", err)
		os.Exit(1)
	}

	markdown, err := readMarkdownFile(date)
	if err != nil {
		fmt.Printf("Error reading markdown file: %v\n", err)
		os.Exit(1)
	}

	if len(markdown) == 0 {
		fmt.Println("Markdown file is empty")
		os.Exit(1)
	}

	source := string(markdown)
	sectionNames := []string{"Simple English", "Colloquial English"}
	results, err := extractSections(source, sectionNames)
	if err != nil {
		fmt.Printf("Error extracting sections: %v\n", err)
		os.Exit(1)
	}

	if err := synthesizeSpeech(results, date); err != nil {
		fmt.Printf("Error in speech synthesis: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Completed!")
}

func loadEnvVariables() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}
	return nil
}

func parseDateFlag() (time.Time, error) {
	var inputDate string
	flag.StringVar(&inputDate, "date", "", "Date in YYYY-MM-DD format")
	flag.Parse()

	if inputDate != "" {
		return time.Parse("2006-01-02", inputDate)
	}
	return time.Now().Local(), nil
}

func readMarkdownFile(date time.Time) ([]byte, error) {
	year := date.Format("2006")
	filename := fmt.Sprintf("%s.md", date.Format("01-02"))
	filepath := filepath.Join(os.Getenv("MARKDOWN_FILEPATH_BASE"), year, filename)
	return os.ReadFile(filepath)
}

func extractSections(source string, sectionNames []string) (map[string]string, error) {
	results := make(map[string]string)
	for _, section := range sectionNames {
		result, err := sections.Extract(source, section, sectionNames)
		if err != nil {
			return nil, err
		}
		results[section] = result
	}
	return results, nil
}

func synthesizeSpeech(results map[string]string, date time.Time) error {
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		return err
	}
	defer client.Close()

	ttsClient := &textToSpeechClient{client: client}
	fileWriter := &osFileWriter{}

	for section, result := range results {
		if err := speech.Synthesize(ctx, ttsClient, fileWriter, date, section, result); err != nil {
			return err
		}
	}
	return nil
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
