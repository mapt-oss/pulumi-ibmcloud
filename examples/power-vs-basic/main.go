package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Load configuration
		cfg := config.New(ctx, "")
		powerVSZone := cfg.Get("powerVsZone")
		if powerVSZone == "" {
			powerVSZone = "dal12"
		}
		sshPublicKey := cfg.Require("sshPublicKey")

		// Create Resource Group
		rg, err := ibmcloud.NewResourceGroup(ctx, "power-vs-rg", &ibmcloud.ResourceGroupArgs{
			Name: pulumi.String("power-vs-basic-rg"),
			Tags: pulumi.StringArray{
				pulumi.String("environment:dev"),
				pulumi.String("managed-by:pulumi"),
			},
		})
		if err != nil {
			return err
		}

		// Create Power VS Workspace
		workspace, err := ibmcloud.NewResourceInstance(ctx, "power-vs-workspace", &ibmcloud.ResourceInstanceArgs{
			Name:            pulumi.String("power-vs-basic-workspace"),
			Service:         pulumi.String("power-iaas"),
			Plan:            pulumi.String("power-virtual-server-group"),
			Location:        pulumi.String(powerVSZone),
			ResourceGroupId: rg.ID(),
			Tags: pulumi.StringArray{
				pulumi.String("environment:dev"),
			},
		})
		if err != nil {
			return err
		}

		// Create SSH Key
		sshKey, err := ibmcloud.NewPiKey(ctx, "ssh-key", &ibmcloud.PiKeyArgs{
			PiKeyName:         pulumi.String("power-vs-key"),
			PiSshKey:          pulumi.String(sshPublicKey),
			PiCloudInstanceId: workspace.Guid,
		})
		if err != nil {
			return err
		}

		// Create Private Network
		network, err := ibmcloud.NewPiNetwork(ctx, "private-network", &ibmcloud.PiNetworkArgs{
			PiNetworkName:     pulumi.String("power-vs-network"),
			PiCloudInstanceId: workspace.Guid,
			PiNetworkType:     pulumi.String("vlan"),
			PiCidr:            pulumi.String("192.168.200.0/24"),
			PiGateway:         pulumi.String("192.168.200.1"),
			PiDnsServers: pulumi.StringArray{
				pulumi.String("9.9.9.9"),
				pulumi.String("1.1.1.1"),
			},
		})
		if err != nil {
			return err
		}

		// Create Power VS Instance
		instance, err := ibmcloud.NewPiInstance(ctx, "power-vs-vm", &ibmcloud.PiInstanceArgs{
			PiInstanceName:    pulumi.String("basic-rhel-vm"),
			PiCloudInstanceId: workspace.Guid,
			PiMemory:          pulumi.Float64(4),
			PiProcessors:      pulumi.Float64(0.25),
			PiProcType:        pulumi.String("shared"),
			PiSysType:         pulumi.String("s922"),
			PiImageId:         pulumi.String("rhel-8-4"), // Replace with actual image ID
			PiNetworks: ibmcloud.PiInstancePiNetworkArray{
				&ibmcloud.PiInstancePiNetworkArgs{
					NetworkId: network.NetworkId,
				},
			},
			PiKeyPairName: sshKey.PiKeyName,
			PiStorageType: pulumi.String("tier3"),
			PiHealthStatus: pulumi.String("OK"),
		})
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("resourceGroupId", rg.ID())
		ctx.Export("workspaceGuid", workspace.Guid)
		ctx.Export("networkId", network.NetworkId)
		ctx.Export("instanceId", instance.PiInstanceId)
		ctx.Export("instanceName", instance.PiInstanceName)
		ctx.Export("instanceIP", instance.PiNetworks.Index(pulumi.Int(0)).IpAddress())
		ctx.Export("instanceStatus", instance.Status)

		return nil
	})
}
