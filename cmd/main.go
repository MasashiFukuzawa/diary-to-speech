package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type ErrSectionNotFound string

func (e ErrSectionNotFound) Error() string {
	return fmt.Sprintf("section %q not found", string(e))
}

func extractSections(source, section string, sections []string) (string, error) {
	nextSection := ""
	for i, sec := range sections {
		if sec == section && i+1 < len(sections) {
			nextSection = sections[i+1]
			break
		}
	}

	regexPattern := ""
	if nextSection == "" {
		regexPattern = `(?s)### ` + regexp.QuoteMeta(section) + `(.*?)$`
	} else {
		regexPattern = `(?s)### ` + regexp.QuoteMeta(section) + `(.*?)### ` + regexp.QuoteMeta(nextSection)
	}

	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(source)
	if len(matches) < 2 {
		return "", ErrSectionNotFound(section)
	}
	return strings.TrimSpace(matches[1]), nil
}

func main() {
	fmt.Println("Starting...")

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	now := time.Now().Local()
	year := now.Format("2006")

	filename := fmt.Sprintf("%s.md", now.Format("01-02"))
	markdownFilepath := filepath.Join(os.Getenv("MARKDOWN_FILEPATH_BASE"), year, filename)
	markdown, err := os.ReadFile(markdownFilepath)
	if err != nil {
		panic(err)
	}
	if len(markdown) == 0 {
		panic("markdown file is empty")
	}

	source := string(markdown)
	sections := []string{
		"Simple English",
		"Sophisticated English",
	}
	results := make(map[string]string)
	for _, section := range sections {
		result, err := extractSections(source, section, sections)
		if err != nil {
			panic(err)
		}
		results[section] = result
	}

	for section, result := range results {
		fmt.Println("SECTION:", section)
		fmt.Println("RESULT:", result)
	}

	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println("Synthesizing native speeches...")

	for section, result := range results {
		req := &texttospeechpb.SynthesizeSpeechRequest{
			Input: &texttospeechpb.SynthesisInput{
				InputSource: &texttospeechpb.SynthesisInput_Text{Text: result},
			},
			// A WaveNet generates speech that sounds more natural than other text-to-speech systems.
			// It synthesizes speech with more human-like emphasis and inflection on syllables, phonemes, and words.
			// see: https://cloud.google.com/text-to-speech/docs/wavenet
			Voice: &texttospeechpb.VoiceSelectionParams{
				LanguageCode: os.Getenv("LANGUAGE_CODE"),
				Name:         "en-US-Wavenet-A",
			},
			AudioConfig: &texttospeechpb.AudioConfig{
				AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			},
		}

		resp, err := client.SynthesizeSpeech(ctx, req)
		if err != nil {
			panic(err)
		}

		filename := fmt.Sprintf("%s_%s.mp3", now.Format("20060102"), strings.ToLower(strings.ReplaceAll(section, " ", "_")))
		outputPath := filepath.Join(os.Getenv("OUTPUT_PATH_BASE"), filename)
		err = os.WriteFile(outputPath, resp.AudioContent, 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("Completed!")
	}
}
