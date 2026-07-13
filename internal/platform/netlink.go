package platform

import (
	"encoding/binary"
	"fmt"
	"net"
	"syscall"
)

// InterfaceAddress is an address assigned to a Linux network interface.
type InterfaceAddress struct {
	Interface string
	Address   string
}

// AddressSource obtains interface addresses from a native Linux API.
type AddressSource interface {
	InterfaceAddresses() ([]InterfaceAddress, error)
}

type netlinkSource struct{}

// RealNetlink returns the production RTM_GETADDR adapter.
func RealNetlink() AddressSource { return netlinkSource{} }

func (netlinkSource) InterfaceAddresses() ([]InterfaceAddress, error) {
	rib, err := syscall.NetlinkRIB(syscall.RTM_GETADDR, syscall.AF_UNSPEC)
	if err != nil {
		return nil, fmt.Errorf("reading address netlink dump: %w", err)
	}
	messages, err := syscall.ParseNetlinkMessage(rib)
	if err != nil {
		return nil, fmt.Errorf("parsing address netlink dump: %w", err)
	}
	var result []InterfaceAddress
	for i := range messages {
		message := &messages[i]
		if message.Header.Type != syscall.RTM_NEWADDR || len(message.Data) < syscall.SizeofIfAddrmsg {
			continue
		}
		family, prefix := message.Data[0], message.Data[1]
		if family != syscall.AF_INET && family != syscall.AF_INET6 {
			continue
		}
		index := binary.NativeEndian.Uint32(message.Data[4:8])
		iface, err := net.InterfaceByIndex(int(index))
		if err != nil {
			continue
		}
		attrs, err := syscall.ParseNetlinkRouteAttr(message)
		if err != nil {
			return nil, fmt.Errorf("parsing address attributes: %w", err)
		}
		address := addressFromAttrs(attrs, family, prefix)
		if address != "" {
			result = append(result, InterfaceAddress{Interface: iface.Name, Address: address})
		}
	}
	return result, nil
}

func addressFromAttrs(attrs []syscall.NetlinkRouteAttr, family, prefix uint8) string {
	var fallback []byte
	for _, attr := range attrs {
		switch attr.Attr.Type {
		case syscall.IFA_LOCAL:
			return formatAddress(attr.Value, family, prefix)
		case syscall.IFA_ADDRESS:
			fallback = attr.Value
		}
	}
	return formatAddress(fallback, family, prefix)
}

func formatAddress(value []byte, family, prefix uint8) string {
	if len(value) == 0 {
		return ""
	}
	ip := net.IP(value)
	if family == syscall.AF_INET && len(value) != net.IPv4len {
		return ""
	}
	if family == syscall.AF_INET6 && len(value) != net.IPv6len {
		return ""
	}
	return fmt.Sprintf("%s/%d", ip.String(), prefix)
}
