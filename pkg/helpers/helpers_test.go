package helpers_test

import (
	"os"
	"ozonLinkShortener/pkg/helpers"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	config := helpers.GetConfig()
	os.Setenv("SHORT_LINK_BASE", "http://localhost:15001/")
	os.Setenv("POSTGRES_CONN_STRING", "postgres://user:password@localhost:5432/linkShortener")
	require.Equal(t, "http://localhost:15001/", config.ShortLinkBase)
	require.Equal(t, "postgres://user:password@localhost:5432/linkShortener", config.DbConnString)
}
