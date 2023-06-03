package speech

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	texttospeechpb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

type TextToSpeechClient interface {
	SynthesizeSpeech(ctx context.Context, in *texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error)
}

type FileWriter interface {
	WriteFile(filename string, data []byte, perm fs.FileMode) error
}

func Synthesize(ctx context.Context, client TextToSpeechClient, writer FileWriter, now time.Time, section, result string) error {
	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: result},
		},
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
		return err
	}

	filename := fmt.Sprintf("%s_%s.mp3", now.Format("20060102"), strings.ToLower(strings.ReplaceAll(section, " ", "_")))
	outputPath := filepath.Join(os.Getenv("OUTPUT_PATH_BASE"), filename)
	err = writer.WriteFile(outputPath, resp.AudioContent, 0644)
	if err != nil {
		return err
	}

	return nil
}
