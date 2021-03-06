package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	n, err := w.Write([]byte("test\n"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 5, n)
	assert.Zero(t, w.ind)
	assert.True(t, w.atBeginOfLine)
}

func TestWriteWithoutln(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	n, err := w.Write([]byte("test"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 4, n)
	assert.Zero(t, w.ind)
	assert.False(t, w.atBeginOfLine)
}

func TestWriteWithoutlnInNewLine(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)
	w.atBeginOfLine = true
	w.ind = 1

	// when
	n, err := w.Write([]byte("test"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, 5, w.body.Len())
	assert.Equal(t, 1, w.ind)
	assert.False(t, w.atBeginOfLine)
}

func TestWriteMultiLines(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)
	w.atBeginOfLine = true
	w.ind = 1

	// when
	n, err := w.Write([]byte("line1\nline2"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, 13, w.body.Len())
	assert.Equal(t, 1, w.ind)
	assert.False(t, w.atBeginOfLine)
}

func TestWriteString(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	n, err := w.WriteString("test\n")

	// then
	assert.Nil(t, err)
	assert.Equal(t, 5, n)
	assert.Zero(t, w.ind)
	assert.True(t, w.atBeginOfLine)
}

func TestWriteByte(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	err := w.WriteByte('a')

	// then
	assert.Nil(t, err)
	assert.Equal(t, 1, w.body.Len())
	assert.False(t, w.atBeginOfLine)
}

func TestWriteByteln(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	err := w.WriteByte('\n')

	// then
	assert.Nil(t, err)
	assert.Equal(t, 1, w.body.Len())
	assert.True(t, w.atBeginOfLine)
}

func TestWriteHead(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	_, err := w.WriteHead("head")
	require.NoError(t, err)

	// then
	assert.Equal(t, 5, w.head.Len())
}

func TestIndent(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	w.Indent()

	// then
	assert.Equal(t, 1, w.ind)
}

func TestUnindent(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)

	assert.Panics(t, func() {
		w.Unindent()
	})

	w.Indent()
	assert.Equal(t, 1, w.ind)

	w.Unindent()
	assert.Zero(t, w.ind)

	assert.Panics(t, func() {
		w.Unindent()
	})
}

func TestActivate(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)

	w.Activate("a")
	assert.Equal(t, 11, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{"a": 1}, w.active)

	w.Activate("a")
	assert.Equal(t, 22, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{"a": 2}, w.active)
}

func TestActivated(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)

	d := w.Activated("a", false)
	assert.Equal(t, 11, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{"a": 1}, w.active)

	d()
	assert.Equal(t, 24, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)

	d()
	assert.Equal(t, 24, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)
}

func TestActivatedWithSuppressed(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)

	d := w.Activated("a", true)
	assert.Zero(t, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)

	d()
	assert.Zero(t, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)

	d()
	assert.Zero(t, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)
}

func TestDeactivate(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)

	w.Activate("a")
	assert.Equal(t, 11, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{"a": 1}, w.active)

	w.Deactivate("a")
	assert.Equal(t, 24, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)

	w.Deactivate("a")
	assert.Equal(t, 24, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)

	w.Deactivate("b")
	assert.Equal(t, 24, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)
}

func TestWriteIndent(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)
	w.atBeginOfLine = true
	w.ind = 1

	w.writeIndent()
	assert.Equal(t, 1, w.body.Len())
	assert.False(t, w.atBeginOfLine)

	w.writeIndent()
	assert.Equal(t, 1, w.body.Len())
	assert.False(t, w.atBeginOfLine)
}

func TestStringer(t *testing.T) {
	t.Parallel()

	// Given
	w := MakeSequenceDiagramWriter(false)
	_, err := w.WriteHead("head")
	require.NoError(t, err)
	_, err = w.WriteString("body\n")
	require.NoError(t, err)

	// When
	s := w.String()

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "@startuml\nhead\nbody\n@enduml\n", s)
}

func TestStringerEmpty(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(false)
	_, err := w.WriteString("body\n")

	s := w.String()
	expected := ""

	assert.Nil(t, err)
	assert.Equal(t, expected, s)
}

func TestStringerWithAutogen(t *testing.T) {
	t.Parallel()

	// Given
	w := MakeSequenceDiagramWriter(true)
	_, err := w.WriteHead("head")
	require.NoError(t, err)
	_, err = w.WriteString("body\n")
	require.NoError(t, err)

	// When
	s := w.String()

	// Then
	expected := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
head
body
@enduml
`
	assert.Nil(t, err)
	assert.Equal(t, expected, s)
}

func TestWriteTo(t *testing.T) {
	t.Parallel()

	// Given
	w := MakeSequenceDiagramWriter(false, "properties 1")
	_, err := w.WriteHead("head")
	require.NoError(t, err)
	_, err = w.WriteString("body\n")
	require.NoError(t, err)

	// When
	var b bytes.Buffer
	n, err := w.WriteTo(&b)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, int64(41), n)
	assert.Equal(t, "@startuml\nhead\nproperties 1\nbody\n@enduml\n", b.String())
}
