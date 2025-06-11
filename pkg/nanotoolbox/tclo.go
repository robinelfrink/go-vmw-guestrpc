// SPDX-FileCopyrightText: Copyright (c) 2020 Oliver Kuckertz, Siderolabs and Equinix
// SPDX-License-Identifier: Apache-2.0

package nanotoolbox

import (
	"log/slog"

	"github.com/equinix-ms/go-vmw-guestrpc/pkg/hypercall"
)

// TCLO represents a "TCLO" communication interface with the hypervisor.
type TCLO struct {
	channel *hypercall.Channel
	logger  *slog.Logger
}

// TCLOCallBack is a callback function.
type TCLOCallBack func(command string) (reply string, err error)

// NewTCLO creates a new TCLO interface.
func NewTCLO(log *slog.Logger) (*TCLO, error) {
	l := log.With("module", "nanotoolbox.TCLO")

	l.Debug("initializing")

	return &TCLO{
		channel: nil,
		logger:  l,
	}, nil
}

// Start starts the TCLO.
func (t *TCLO) Start() error {
	t.logger.Debug("starting")

	ch, err := hypercall.NewChannel(hypercall.TCLOProto, t.logger.With("module", "hypercall.Channel"))
	if err != nil {
		return err
	}

	t.channel = ch

	return nil
}

// Stop stops the TCLO.
func (t *TCLO) Stop() error {
	t.logger.Debug("closing")

	return t.channel.Close()
}

// Send sends data over TCLO.
func (t *TCLO) Send(data []byte) error {
	return t.channel.Send(data)
}

// Receive receives data over TCLO.
func (t *TCLO) Receive() ([]byte, error) {
	return t.channel.Receive()
}
