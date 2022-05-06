package root

import (
	"fmt"
	"io"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	c := Command()
	c.SetOut(io.Discard)
	c.SetArgs([]string{
		"version",
		fmt.Sprintf("--%s", colorFlag),
	})
	err := c.Execute()
	require.NoError(t, err)
	assert.False(t, color.NoColor)
}
