package speech_test

import (
	"context"
	"io/fs"
	"testing"
	"time"

	texttospeechpb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/MasashiFukuzawa/diary-to-speech/pkg/speech"
	"github.com/stretchr/testify/mock"
)

type MockTextToSpeechClient struct {
	mock.Mock
}

func (m *MockTextToSpeechClient) SynthesizeSpeech(ctx context.Context, in *texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*texttospeechpb.SynthesizeSpeechResponse), args.Error(1)
}

type MockFileWriter struct {
	mock.Mock
}

func (m *MockFileWriter) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	args := m.Called(filename, data, perm)
	return args.Error(0)
}

func TestSynthesize(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	section := "section"
	result := "result"

	client := &MockTextToSpeechClient{
		Mock: mock.Mock{},
	}
	writer := &MockFileWriter{
		Mock: mock.Mock{},
	}

	client.On("SynthesizeSpeech", mock.Anything, mock.Anything).Return(&texttospeechpb.SynthesizeSpeechResponse{}, nil)
	writer.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := speech.Synthesize(ctx, client, writer, now, section, result)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	client.AssertCalled(t, "SynthesizeSpeech", ctx, mock.Anything)
	writer.AssertCalled(t, "WriteFile", mock.Anything, mock.Anything, mock.Anything)
}
