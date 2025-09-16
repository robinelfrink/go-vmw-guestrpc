# go-vmw-guestrpc

This is an implementation of the mechanism that VMWare guests use to
interface with the hypervisor through the RPC/VMX "backdoor". It has been
inspired by the following projects:

* [open-vm-tools](https://github.com/vmware/open-vm-tools)
* [vmw-guestinfo](https://github.com/vmware-archive/vmw-guestinfo)
* [vmw_backdoor](https://docs.rs/vmw_backdoor)

It implements the backdoor itself, for i386, amd64 and arm64, as well as
RPCI, TCLO and a minimal toolbox to provide services.

Using the RPC/VMX backdoor involves executing a privileged CPU instruction,
which is trapped by the hypervisor. If you use this backdoor on a non-VM, it
will cause a `SEGFAULT`, which is why we recommend testing if you are running
in ESX first. Strange as it may sound, using the backdoor (given that you are
running in ESXi) does not require any privileges (non-root, no caps is OK). The
function `IsVMWareVM()`, however, does require at least `CAP_SYS_RAWIO` to
determine if it's safe to use the backdoor.