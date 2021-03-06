package main

import (
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestAppDependency_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "AppA:EndptA:AppB:EndptB", (&AppDependency{
		Self:      AppElement{Name: "AppA", Endpoint: "EndptA"},
		Target:    AppElement{Name: "AppB", Endpoint: "EndptB"},
		Statement: &sysl.Statement{},
	}).String())
}
