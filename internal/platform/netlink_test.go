package platform

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressFromAttrs(t *testing.T) {
	attrs := []syscall.NetlinkRouteAttr{{Attr: syscall.RtAttr{Type: syscall.IFA_ADDRESS}, Value: []byte{192, 0, 2, 10}}}
	assert.Equal(t, "192.0.2.10/24", addressFromAttrs(attrs, syscall.AF_INET, 24))
}
