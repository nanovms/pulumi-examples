package main

import (
	compute "github.com/pulumi/pulumi-azure-native-sdk/compute/v3"
	network "github.com/pulumi/pulumi-azure-native-sdk/network/v3"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		location := "WestUS2"
		resourceGroupNm := "resourceGroupName"
		subId := "my-sub-id-goes-here"
		bucket := "my-image-bucket"
		galleryName := "my-image-gallery"
		prefix := "my-image-name"
		nicName := "test-nic"

		vnet, err := network.NewVirtualNetwork(ctx, "virtualNetwork", &network.VirtualNetworkArgs{
			AddressSpace: &network.AddressSpaceArgs{
				AddressPrefixes: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"),
				},
			},
			FlowTimeoutInMinutes: pulumi.Int(10),
			Location:             pulumi.String(location),
			ResourceGroupName:    pulumi.String(resourceGroupNm),
			VirtualNetworkName:   pulumi.String("test-vnet"),
		})
		if err != nil {
			return err
		}

		subnet, err := network.NewSubnet(ctx, "subnet", &network.SubnetArgs{
			AddressPrefix:      pulumi.String("10.0.0.0/16"),
			ResourceGroupName:  pulumi.String(resourceGroupNm),
			SubnetName:         pulumi.String("subnet1"),
			VirtualNetworkName: vnet.Name,
		})
		if err != nil {
			return err
		}

		nic, err := network.NewNetworkInterface(ctx, nicName, &network.NetworkInterfaceArgs{
			IpConfigurations: network.NetworkInterfaceIPConfigurationArray{
				&network.NetworkInterfaceIPConfigurationArgs{
					Name: pulumi.String("ipconfig1"),
					Subnet: &network.SubnetTypeArgs{
						Id: subnet.ID(),
					},
				},
			},
			Location:             pulumi.String(location),
			NetworkInterfaceName: pulumi.String(nicName),
			ResourceGroupName:    pulumi.String(resourceGroupNm),
		})
		if err != nil {
			return err
		}

		// uncomment to create a single vm

		/*
				_, err = compute.NewVirtualMachine(ctx, "virtualMachine", &compute.VirtualMachineArgs{
					HardwareProfile: &compute.HardwareProfileArgs{
						VmSize: pulumi.String(compute.VirtualMachineSizeTypes_Standard_D1_v2),
					},
					Location: pulumi.String(location),
					NetworkProfile: &compute.NetworkProfileArgs{
						NetworkInterfaces: compute.NetworkInterfaceReferenceArray{
							&compute.NetworkInterfaceReferenceArgs{
								Id:      pulumi.String("/subscriptions/" + subId + "/resourceGroups/" + resourceGroupNm + "/providers/Microsoft.Network/networkInterfaces/" + nicName),
								Primary: pulumi.Bool(true),
							},
						},
					},
					OsProfile: &compute.OSProfileArgs{
						AdminPassword: pulumi.String("SOme-random-pass-that-is-not-used"),
						AdminUsername: pulumi.String("pulumiadmin"),
						ComputerName:  pulumi.String("myVM"),
					},
					ResourceGroupName: pulumi.String(resourceGroupNm),
					DiagnosticsProfile: &compute.DiagnosticsProfileArgs{
						BootDiagnostics: &compute.BootDiagnosticsArgs{
							Enabled:    pulumi.Bool(true),
							StorageUri: pulumi.String("http://" + bucket + ".blob.core.windows.net"),
						},
					},
					StorageProfile: &compute.StorageProfileArgs{
						ImageReference: &compute.ImageReferenceArgs{
							Id: pulumi.String("/subscriptions/" + subId + "/resourceGroups/" + resourceGroupNm + "/providers/Microsoft.Compute/galleries/" + galleryName + "/images/" + prefix),
						},
						OsDisk: &compute.OSDiskArgs{
							Caching:      compute.CachingTypesReadWrite,
							CreateOption: pulumi.String(compute.DiskCreateOptionTypesFromImage),
							ManagedDisk: &compute.ManagedDiskParametersArgs{
								StorageAccountType: pulumi.String(compute.StorageAccountTypes_Standard_LRS),
							},
							Name: pulumi.String("myVMosdisk"),
						},
					},
					VmName: pulumi.String("myVM"),
				}, pulumi.DependsOn([]pulumi.Resource{nic}))
				if err != nil {
					return err
				}

				return nil

			})
		*/

		_, err = compute.NewVirtualMachineScaleSet(ctx, "virtualMachineScaleSet", &compute.VirtualMachineScaleSetArgs{
			Location:          pulumi.String(location),
			Overprovision:     pulumi.Bool(true),
			ResourceGroupName: pulumi.String(resourceGroupNm),
			Sku: &compute.SkuArgs{
				Capacity: pulumi.Float64(2),
				Name:     pulumi.String("Standard_B1S"),
				Tier:     pulumi.String("Standard"),
			},
			UpgradePolicy: &compute.UpgradePolicyArgs{
				Mode: compute.UpgradeModeManual,
			},
			VirtualMachineProfile: &compute.VirtualMachineScaleSetVMProfileArgs{
				DiagnosticsProfile: &compute.DiagnosticsProfileArgs{
					BootDiagnostics: &compute.BootDiagnosticsArgs{
						Enabled:    pulumi.Bool(true),
						StorageUri: pulumi.String("http://" + bucket + ".blob.core.windows.net"),
					},
				},
				NetworkProfile: &compute.VirtualMachineScaleSetNetworkProfileArgs{
					NetworkInterfaceConfigurations: compute.VirtualMachineScaleSetNetworkConfigurationArray{
						&compute.VirtualMachineScaleSetNetworkConfigurationArgs{
							EnableIPForwarding: pulumi.Bool(true),
							IpConfigurations: compute.VirtualMachineScaleSetIPConfigurationArray{
								&compute.VirtualMachineScaleSetIPConfigurationArgs{
									Name: pulumi.String("vmss-name"),
									Subnet: &compute.ApiEntityReferenceArgs{
										Id: subnet.ID(),
									},
								},
							},
							Name:    pulumi.String("vmss-name"),
							Primary: pulumi.Bool(true),
						},
					},
				},
				OsProfile: &compute.VirtualMachineScaleSetOSProfileArgs{
					AdminPassword:      pulumi.String("SOme-random-pass-that-is-not-used"),
					AdminUsername:      pulumi.String("pulumiadmin"),
					ComputerNamePrefix: pulumi.String("myVM"),
				},
				StorageProfile: &compute.VirtualMachineScaleSetStorageProfileArgs{
					ImageReference: &compute.ImageReferenceArgs{
						Id: pulumi.String("/subscriptions/" + subId + "/resourceGroups/" + resourceGroupNm + "/providers/Microsoft.Compute/galleries/" + galleryName + "/images/" + prefix),
					},
					OsDisk: &compute.VirtualMachineScaleSetOSDiskArgs{
						Caching:      compute.CachingTypesReadWrite,
						CreateOption: pulumi.String(compute.DiskCreateOptionTypesFromImage),
						ManagedDisk: &compute.VirtualMachineScaleSetManagedDiskParametersArgs{
							StorageAccountType: pulumi.String(compute.StorageAccountTypes_Standard_LRS),
						},
					},
				},
			},
			VmScaleSetName: pulumi.String("myss"),
		}, pulumi.DependsOn([]pulumi.Resource{nic}))
		if err != nil {
			return err
		}

		return nil

	})

}
