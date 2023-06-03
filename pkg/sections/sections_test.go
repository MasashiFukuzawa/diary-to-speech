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

### Advanced English
This is an advanced English section.
`

	section := "Intermediate English"
	expected := "This is an intermediate English section."

	result, err := sections.Extract(source, section, []string{"Simple English", "Intermediate English", "Advanced English"})
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

### Advanced English
This is an advanced English section.
`

	section := "Advanced French"

	_, err := sections.Extract(source, section, []string{"Simple English", "Intermediate English", "Advanced English"})
	assert.Error(t, err)
	assert.EqualError(t, err, "section \"Advanced French\" not found")
}
