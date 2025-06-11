# go-vmw-guestrpc

This is an implementation of the mechanism that VMWare guests use to
interface with the hypervisor through the RPC/VMX "backdoor". It has been
inspired by the following projects:

* [open-vm-tools](https://github.com/vmware/open-vm-tools)
* [vmw-guestinfo](https://github.com/vmware-archive/vmw-guestinfo)
* [vmw_backdoor](https://docs.rs/vmw_backdoor)

It implements the backdoor itself, for i386, amd64 and arm64, as well as
RPCI, TCLO and a minimal toolbox to provide services.