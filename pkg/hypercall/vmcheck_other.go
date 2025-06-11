// SPDX-FileCopyrightText: Copyright (c) 2020 Oliver Kuckertz, Siderolabs and Equinix
// SPDX-License-Identifier: Apache-2.0

//go:build !(linux && (amd64 || i386))

package hypercall

func hypercallPreCheck() bool {
	// we have not found a way to do a generic check on non-x86/non-linux
	// so we assume all is fine (and risk segfaults)
	return true
}
