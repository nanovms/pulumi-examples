import * as pulumi from "@pulumi/pulumi";
import * as resources from "@pulumi/azure-native/resources";
import * as storage from "@pulumi/azure-native/storage";
import * as azure_native from "@pulumi/azure-native";

const config = new pulumi.Config();

// set the following:
const prefix = "myimg";
const location = "WestUS2";
const resourceGroup = "myregion";
const subId = "my-subscription-id-goes-here";
const bucket = "mybucket";
const galleryName = "gallery-name";

// Create a Virtual Network
const virtualNetwork = new azure_native.network.VirtualNetwork("vm-vnet", {
    addressSpace: {
      addressPrefixes: ["10.0.0.0/16"],
    },
    location: location,
    resourceGroupName: resourceGroup,
});

// Create a Subnet
const subnet = new azure_native.network.Subnet("vm-subnet", {
    name: "default",
    resourceGroupName: resourceGroup,
    virtualNetworkName: virtualNetwork.name,
    addressPrefixes: ["10.0.1.0/24"],
});

// Create a Network Interface
const networkInterface = new azure_native.network.NetworkInterface("vm-nic", {
    networkInterfaceName: "vm-nic",
    location: location,
    resourceGroupName: resourceGroup,
    ipConfigurations: [{
        name: "ipconfig1",
        subnet: {
          id: subnet.id,
        },
    }],
});

// Create the Virtual Machine
const virtualMachine = new azure_native.compute.VirtualMachine("vm", {
    vmName: `${prefix}-vm`,
    location: location,
    resourceGroupName: resourceGroup,
    diagnosticsProfile: {
        bootDiagnostics: {
            enabled: true,
            storageUri: "http://" + bucket + ".blob.core.windows.net",
        },
    },
    networkProfile: {
        networkInterfaces: [{
            id: "/subscriptions/" + subId + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Network/networkInterfaces/" + "vm-nic",
            primary: true,
        }],
    },
    hardwareProfile: {
        vmSize: azure_native.compute.VirtualMachineSizeTypes.Standard_D2s_v3,
    },
    storageProfile: {
      imageReference: {
        id: "/subscriptions/" + subId + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Compute/galleries/" + galleryName + "/images/" + prefix,
      },
      osDisk: {
            caching: azure_native.compute.CachingTypes.ReadWrite,
            createOption: azure_native.compute.DiskCreateOptionTypes.FromImage,
            managedDisk: {
                storageAccountType: azure_native.compute.StorageAccountTypes.Premium_LRS,
            },
            name: "myVMosdisk2",
        },
    },
    osProfile: {
        computerName: "myvm",
        adminUsername: "pulumiadmin",
        adminPassword: "SOme-random-pass-that-is-not-used",
    },
});

// if you wish to create a scale set instead uncomment the following and comment out the vm block above
/*
const virtualMachineScaleSet = new azure_native.compute.VirtualMachineScaleSet("virtualMachineScaleSet", {
    location: location,
    overprovision: true,
    resourceGroupName: resourceGroup,
    sku: {
        capacity: 2,
        name: "Standard_B1S",
        tier: "Standard",
    },
    upgradePolicy: {
        mode: azure_native.compute.UpgradeMode.Manual,
    },
    virtualMachineProfile: {
        diagnosticsProfile: {
            bootDiagnostics: {
                enabled: true,
                storageUri: "http://" + bucket + ".blob.core.windows.net",
            },
        },
        networkProfile: {
            networkInterfaceConfigurations: [{
                enableIPForwarding: true,
                ipConfigurations: [{
                    name: "ipconfig1",
                    subnet: {
                        id: subnet.id,
                    },
                }],
                name: "test",
                primary: true,
            }],
        },
        osProfile: {
            computerNamePrefix: "myvm",
            adminUsername: "pulumiadmin",
            adminPassword: "SOme-random-pass-that-is-not-used",
        },
        storageProfile: {
          imageReference: {
            id: "/subscriptions/" + subId + "/resourceGroups/" + resourceGroup + "/providers/Microsoft.Compute/galleries/" + galleryName + "/images/" + prefix,
          },
          osDisk: {
                caching: azure_native.compute.CachingTypes.ReadWrite,
                createOption: azure_native.compute.DiskCreateOptionTypes.FromImage,
                managedDisk: {
                    storageAccountType: azure_native.compute.StorageAccountTypes.Premium_LRS,
                },
            },
        },
    },
    vmScaleSetName: "myss",
});
*/
