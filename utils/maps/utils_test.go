package maps

import (
	"github.com/grasp-labs/ds-boilerplate-api-go/utils/strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetKeysFromMap(t *testing.T) {
	a := make(map[string][]byte)
	a["a"] = []byte("x")
	a["b"] = []byte("y")
	a["c"] = []byte("z")

	assert.True(t, strings.ArrayEqualElements([]string{"a", "b", "c"}, GetKeysFromMap(a)))
}
