# diary-to-speech

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Diary-to-Speech is a Go-based CLI tool that converts English markdown diaries into mp3 files using Google's Text-to-Speech API for language learning.

# Prerequisites

- go version: 1.20.4
- Must have been issued an API key for the Google Text-To-Speech API
  - details: https://cloud.google.com/text-to-speech/docs/before-you-begin

# How To Use

```sh
$ git clone git@github.com:MasashiFukuzawa/diary-to-speech.git
$ cd diary-to-speech

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
$ go run cmd/diary-to-speech/main.go

# You can specify date.
$ go run cmd/diary-to-speech/main.go -date 2024-01-01
```
