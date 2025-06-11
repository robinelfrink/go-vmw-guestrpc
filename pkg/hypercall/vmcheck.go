// SPDX-FileCopyrightText: Copyright (c) 2020 Oliver Kuckertz, Siderolabs and Equinix
// SPDX-License-Identifier: Apache-2.0

package hypercall

import "fmt"

// IsVMWareVM determines if this is a VM in VMWare. It indirectly also verifies if we can connect with hypercall/backdoor.
func IsVMWareVM() bool {
	// IsVirtual() can be a bit violent on non-vms: accessing privileged instructions may cause SEGV.
	// so we try to rule out VMs using less destructive tests, ie. with CPU flags etc.
	if !hypercallPreCheck() {
		fmt.Println("fails precheck")
		return false
	}

	// knock the backdoor and see if we can read VMWare version
	return IsVirtual()
}
