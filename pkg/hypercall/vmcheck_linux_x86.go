// SPDX-FileCopyrightText: Copyright (c) 2020 Oliver Kuckertz, Siderolabs and Equinix
// SPDX-License-Identifier: Apache-2.0

//go:build linux && (amd64 || i386)

package hypercall

import (
	"syscall"

	"github.com/klauspost/cpuid/v2"
)

func hypercallPreCheck() error {
	// is this a VM according to CPUID?
	if !cpuid.CPU.VM() {
		return ErrCpuIdMismatch
	}

	// .. of vendor VMWare?
	if cpuid.CPU.HypervisorVendorID != cpuid.VMware {
		return ErrHypervisorMismatch
	}

	// try to change I/O privilege level to 3. If this succeeds, we are (probably) a VM. If not,
	// we should not try to knock the backdoor port, causing a SEGV
	if err := syscall.Iopl(3); err != nil {
		return ErrSetPivilegeLevel
	}

	return nil
}
