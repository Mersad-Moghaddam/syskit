package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDev(t *testing.T) {
	interfaces, err := ParseDev([]byte("Inter-|\n face |bytes\n eth0: 10 2 3 4 0 0 0 0 20 5 6 7 0 0 0 0\n"))
	assert.NoError(t, err)
	assert.Len(t, interfaces, 1)
	assert.Equal(t, "eth0", interfaces[0].Name)
	assert.Equal(t, uint64(20), interfaces[0].TXBytes)
}
