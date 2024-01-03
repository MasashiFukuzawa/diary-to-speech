package sections_test

import (
	"testing"

	"github.com/MasashiFukuzawa/diary-to-speech/pkg/sections"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	source := `
# Diary

### Simple English
This is a simple English section.

### Intermediate English
This is an intermediate English section.

### Colloquial English
This is an Colloquial English section.
`

	section := "Intermediate English"
	expected := "This is an intermediate English section."

	result, err := sections.Extract(source, section, []string{"Simple English", "Intermediate English", "Colloquial English"})
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExtract_NonexistentSection(t *testing.T) {
	source := `
# Diary

### Simple English
This is a simple English section.

### Intermediate English
This is an intermediate English section.

### Colloquial English
This is an Colloquial English section.
`

	section := "Colloquial French"

	_, err := sections.Extract(source, section, []string{"Simple English", "Intermediate English", "Colloquial English"})
	assert.Error(t, err)
	assert.EqualError(t, err, "section \"Colloquial French\" not found")
}
