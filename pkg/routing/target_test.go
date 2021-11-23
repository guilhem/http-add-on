package routing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTargetServiceURL(t *testing.T) {
	r := require.New(t)

	target := Target{
		Host:       "example.com",
		Service:    "testsvc",
		Port:       8081,
		Deployment: "testdeploy",
	}
	svcURL, err := target.ServiceURL()
	r.NoError(err)
	r.Equal(
		fmt.Sprintf("%s:%d", target.Host, target.Port),
		svcURL.Host,
	)
}
