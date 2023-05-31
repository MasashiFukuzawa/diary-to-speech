# gpt-language-tutor

An AI-powered language learning CLI tool leveraging GPT for translation and speech-to-text APIs for native pronunciation. Write in your native language, translate it into a foreign language, and practice listening to the translated text. Enhance your writing and listening skills with gpt-language-tutor.

# Prerequisites

- go version: 1.20.4
- Must have been issued an API key for the Google Text-To-Speech API
  - details: https://cloud.google.com/text-to-speech/docs/before-you-begin

# How To Use

```sh
$ git clone git@github.com:MasashiFukuzawa/gpt-language-tutor.git
$ cd gpt-language-tutor

# for tutorial
# =====================================================
$ cat << EOF > .env
GOOGLE_APPLICATION_CREDENTIALS=$PWD/secrets/credentials.json
LANGUAGE_CODE=en-US
MARKDOWN_FILEPATH_BASE=$PWD/example
OUTPUT_PATH_BASE=$PWD/example
EOF

$ year=$(date +%Y)
$ date=$(date +%m-%d)
$ mkdir example/$year
$ cp example/2023/04-30.md example/$year/$date.md
# =====================================================

$ go mod tidy
$ go run cmd/main.go
```
