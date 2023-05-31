/*
Copyright 2022 Weaveworks.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package microvm

// VMSpec holds the configuration for the MicroVM which will be passed
// directy to firecracker.
type VMSpec struct {
	// Provider allows you to specify the name of the microvm provider to use.
	// If this isn't supplied then the default provider will be used.
	// NOTE that the default provider cannot be controlled here: it would have been
	// chosen by the operator configuring Flintlock on the remote host.
	// +optional
	Provider string `json:"provider,omitempty"`

	// VCPU specifies how many vcpu's the microvm will be allocated.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum:=1
	VCPU int64 `json:"vcpu"`

	// MemoryMb is the amount of memory in megabytes that the microvm will be allocated.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum:=1024
	MemoryMb int64 `json:"memoryMb"`

	// RootVolume specifies the volume to use for the root of the microvm.
	// +kubebuilder:validation:Required
	RootVolume Volume `json:"rootVolume"`

	// AdditionalVolumes specifies additional non-root volumes to attach to the microvm.
	// +optional
	AdditionalVolumes []Volume `json:"volumes,omitempty"`

	// Kernel specifies the kernel and its arguments to use.
	// +kubebuilder:validation:Required
	Kernel ContainerFileSource `json:"kernel"`

	// KernelCmdLine are the additional args to use for the kernel cmdline.
	// Each MicroVM provider has its own recommended list, they will be used
	// automatically. This field is for additional values.
	KernelCmdLine map[string]string `json:"kernelCmdline,omitempty"`

	// Initrd is an optional initial ramdisk to use.
	// +optional
	Initrd *ContainerFileSource `json:"initrd,omitempty"`

	// NetworkInterfaces specifies the network interfaces attached to the microvm.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	NetworkInterfaces []NetworkInterface `json:"networkInterfaces"`

	// Labels allow you to include extra data on the Microvm
	// +optional
	Labels map[string]string `json:"labels"`
}

// ContainerFileSource represents a file coming from a container image.
type ContainerFileSource struct {
	// Image is the container image to use.
	// +kubebuilder:validation:Required
	Image string `json:"image"`
	// Filename is the name of the file in the container to use.
	// +optional
	Filename string `json:"filename,omitempty"`
}

// Volume represents a volume to be attached to a microvm.
type Volume struct {
	// ID is a unique identifier for this volume.
	// +kubebuilder:validation:Required
	ID string `json:"id"`
	// Image is the container image to use as the source for the volume.
	// +kubebuilder:validation:Required
	Image string `json:"image"`
	// ReadOnly specifies that the volume is to be mounted readonly.
	// +kubebuilder:default:=false
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`
	// MountPoint specifies the guest mountpoint for the volume.
	// This will only be applied to additional volumes.
	// +optional
	MountPoint string `json:"mountPoint,omitempty"`
}

// IfaceType is a type representing the network interface types.
type IfaceType string

const (
	// IfaceTypeTap is a TAP network interface.
	IfaceTypeTap = "tap"
	// IfaceTypeMacvtap is a MACVTAP network interface.
	IfaceTypeMacvtap = "macvtap"
)

// NetworkInterface represents a network interface for the microvm.
type NetworkInterface struct {
	// GuestDeviceName is the name of the network interface to create in the microvm.
	// +kubebuilder:validation:Required
	GuestDeviceName string `json:"guestDeviceName"`
	// GuestMAC allows the specifying of a specific MAC address to use for the interface. If
	// not supplied a autogenerated MAC address will be used.
	// +optional
	GuestMAC string `json:"guestMac,omitempty"`
	// Type is the type of host network interface type to create to use by the guest.
	// +kubebuilder:validation:Enum=macvtap;tap
	Type IfaceType `json:"type"`
	// Address is an optional IP address to assign to this interface. If not supplied then DHCP will be used.
	// +optional
	Address string `json:"address,omitempty"`
}

// VMState is a type that represents the state of a microvm.
type VMState string

var (
	// VMStatePending indicates the microvm hasn't been started.
	VMStatePending = VMState("pending")
	// VMStateRunning indicates the microvm is running.
	VMStateRunning = VMState("running")
	// VMStateFailed indicates the microvm has failed.
	VMStateFailed = VMState("failed")
	// VMStateDeleted indicates the microvm has been deleted.
	VMStateDeleted = VMState("deleted")
	// VMStateUnknown indicates the microvm is in an state that is unknown/supported by CAPMVM.
	VMStateUnknown = VMState("unknown")
)

type Host struct {
	// Name is an optional name for the host.
	// +optional
	Name string `json:"name,omitempty"`
	// Endpoint is the API endpoint for the microvm service (i.e. flintlock)
	// including the port.
	// +kubebuilder:validation:Required
	Endpoint string `json:"endpoint"`
}

type SSHPublicKey struct {
	// User is the name of the user to add keys for (eg root, ubuntu).
	// +kubebuilder:validation:Required
	User string `json:"user,omitempty"`
	// AuthorizedKeys is a list of public keys to add to the user
	// +kubebuilder:validation:Required
	AuthorizedKeys []string `json:"authorizedKeys,omitempty"`
}
