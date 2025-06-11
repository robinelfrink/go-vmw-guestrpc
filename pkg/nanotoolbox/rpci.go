// SPDX-FileCopyrightText: Copyright (c) 2020 Oliver Kuckertz, Siderolabs and Equinix
// SPDX-License-Identifier: Apache-2.0

package nanotoolbox

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/equinix-ms/go-vmw-guestrpc/internal/util"
	"github.com/equinix-ms/go-vmw-guestrpc/pkg/hypercall"
)

// GuestInfoID represents the type of guest information.
type GuestInfoID int

const (
	// GuestInfoError represents an error.
	GuestInfoError GuestInfoID = iota
	// GuestInfoDNSName is the guest info kind for the DNS name.
	GuestInfoDNSName
	// GuestInfoIPAddress is the IP of the guest.
	GuestInfoIPAddress
	// GuestInfoDiskFreeSpace is the amount of free disk space.
	GuestInfoDiskFreeSpace
	// GuestInfoBuildNumber is the build number.
	GuestInfoBuildNumber
	// GuestInfoOSNameFull is the guest info kind for the full OS name.
	GuestInfoOSNameFull
	// GuestInfoOSName is the guest info kind for the OS name.
	GuestInfoOSName
	// GuestInfoUptime is the guest uptime in 100s of seconds.
	GuestInfoUptime
	// GuestInfoMemory is the amount of memory.
	GuestInfoMemory
	// GuestInfoIPAddressV2 is the IP of the guest.
	GuestInfoIPAddressV2
	// GuestInfoIPAddressV3 is the guest info kind for the IP address.
	GuestInfoIPAddressV3
	// GuestInfoOSDetailed is detailed information about the OS.
	GuestInfoOSDetailed
)

var guestInfos = map[GuestInfoID]string{
	GuestInfoError:         "Error",
	GuestInfoDNSName:       "DNS name",
	GuestInfoIPAddress:     "IP address (V1 format)",
	GuestInfoDiskFreeSpace: "Disk free space",
	GuestInfoBuildNumber:   "Build number",
	GuestInfoOSNameFull:    "Full OS name",
	GuestInfoOSName:        "Short OS name",
	GuestInfoUptime:        "Uptime",
	GuestInfoMemory:        "Memory",
	GuestInfoIPAddressV2:   "IP address (V2 format)",
	GuestInfoIPAddressV3:   "IP address (V3 format)",
	GuestInfoOSDetailed:    "Detailed OS info",
}

func (g GuestInfoID) String() string {
	value := "UNKNOWN"

	if v, ok := guestInfos[g]; ok {
		value = v
	}

	return value
}

// RPCI models a RPCI communication interface.
type RPCI struct {
	channel *hypercall.Channel
	logger  *slog.Logger
}

var (
	// RpciOK is the return code for a successful RPCI request.
	RpciOK = []byte{'1', ' '}
	// RpciERR is the return code for a failed RPCI request.
	RpciERR = []byte{'0', ' '}
)

// NewRPCI creates a new RPCI instance.
func NewRPCI(log *slog.Logger) (*RPCI, error) {
	l := log.With("module", "nanotoolbox.RPCI")

	l.Debug("initializing")

	return &RPCI{
		channel: nil,
		logger:  l,
	}, nil
}

// Request sends an RPC command to the vmx and checks the return code for success or error.
func (r *RPCI) Request(request []byte) ([]byte, bool, error) {
	if r.channel == nil {
		return nil, false, fmt.Errorf("no channel available for request %q", request)
	}

	if err := r.channel.Send(request); err != nil {
		return nil, false, err
	}

	reply, err := r.channel.Receive()
	if err != nil {
		return nil, false, err
	}

	if bytes.HasPrefix(reply, RpciOK) {
		return reply[2:], true, nil
	}

	if bytes.HasPrefix(reply, RpciERR) {
		return reply[2:], false, nil
	}

	return nil, false, fmt.Errorf("failed request %q: %q", request, reply)
}

// InfoGet retrieves a property, for instance ExtraConfig (guestinfo).
func (r *RPCI) InfoGet(key string, defaultValue string) (string, error) {
	l := r.logger.With("cmd", "info-get", "key", key)

	reply, ok, err := r.Request([]byte(fmt.Sprintf("info-get %s", key)))
	l.Debug("requested", "reply", reply, "ok", ok, "err", err)

	if err != nil {
		return "", fmt.Errorf("error requesting info: %w", err)
	}

	if !ok {
		return defaultValue, nil
	}

	return string(reply), nil
}

// InfoSet sets a property.
func (r *RPCI) InfoSet(key string, value string) (bool, error) {
	l := r.logger.With("cmd", "info-set", "key", key)

	_, ok, err := r.Request([]byte(fmt.Sprintf("info-set %s %s", key, value)))
	l.Debug("requested", "ok", ok, "err", err)

	if err != nil {
		return false, fmt.Errorf("error requesting info: %w", err)
	}

	return ok, err
}

// SetGuestInfo sets a guest info property. Similar to InfoSet.
func (r *RPCI) SetGuestInfo(kind GuestInfoID, data []byte) {
	l := r.logger.With("cmd", "SetGuestInfo", "kind", kind.String())
	msg := append([]byte(fmt.Sprintf("SetGuestInfo  %d ", kind)), data...)
	util.TraceLog(l, "setting", "msg", string(msg))

	if _, ok, err := r.Request(msg); err != nil {
		l.Error("error sending guestinfo", "err", err)
	} else {
		if !ok {
			l.Warn("request returned not OK", "err", err)
		}
	}
}

// Start starts the RPCI, creating a channel.
func (r *RPCI) Start() error {
	r.logger.Debug("starting")

	ch, err := hypercall.NewChannel(hypercall.RPCIProto, r.logger.With("module", "hypercall.Channel"))
	if err != nil {
		return err
	}

	r.channel = ch

	return nil
}

// Stop closes the RPCI.
func (r *RPCI) Stop() error {
	r.logger.Debug("closing")

	return r.channel.Close()
}

// Send sends data over RPCI.
func (r *RPCI) Send(data []byte) error {
	return r.channel.Send(data)
}

// Receive receives data over RPCI.
func (r *RPCI) Receive() ([]byte, error) {
	return r.channel.Receive()
}
