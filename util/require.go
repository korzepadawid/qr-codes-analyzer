package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func RequireBodyMatchObject[T interface{}](t *testing.T, b *bytes.Buffer, obj T) {
	data, err := ioutil.ReadAll(b)
	require.NoError(t, err)
	var got T
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	require.Equal(t, obj, got)
}
