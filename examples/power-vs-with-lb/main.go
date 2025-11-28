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
		region := cfg.Get("region")
		if region == "" {
			region = "us-south" // Default region
		}

		zone := cfg.Get("zone")
		if zone == "" {
			zone = "us-south-1" // Default zone
		}

		powerVSZone := cfg.Get("powerVsZone")
		if powerVSZone == "" {
			powerVSZone = "dal12" // Default Power VS zone
		}

		sshPublicKey := cfg.Require("sshPublicKey")

		// Step 1: Create a Resource Group
		resourceGroup, err := ibmcloud.NewResourceGroup(ctx, "power-vs-rg", &ibmcloud.ResourceGroupArgs{
			Name: pulumi.String("power-vs-demo-rg"),
			Tags: pulumi.StringArray{
				pulumi.String("environment:demo"),
				pulumi.String("managed-by:pulumi"),
			},
		})
		if err != nil {
			return err
		}

		// Step 2: Create Power VS Workspace (Service Instance)
		powerVSWorkspace, err := ibmcloud.NewResourceInstance(ctx, "power-vs-workspace", &ibmcloud.ResourceInstanceArgs{
			Name:              pulumi.String("power-vs-workspace"),
			Service:           pulumi.String("power-iaas"),
			Plan:              pulumi.String("power-virtual-server-group"),
			Location:          pulumi.String(powerVSZone),
			ResourceGroupId:   resourceGroup.ID(),
			Tags: pulumi.StringArray{
				pulumi.String("environment:demo"),
			},
		})
		if err != nil {
			return err
		}

		// Step 3: Create SSH Key for Power VS
		sshKey, err := ibmcloud.NewPiKey(ctx, "power-vs-ssh-key", &ibmcloud.PiKeyArgs{
			PiKeyName:          pulumi.String("power-vs-demo-key"),
			PiSshKey:           pulumi.String(sshPublicKey),
			PiCloudInstanceId:  powerVSWorkspace.Guid,
		})
		if err != nil {
			return err
		}

		// Step 4: Create Power VS Private Network (Subnet)
		powerVSNetwork, err := ibmcloud.NewPiNetwork(ctx, "power-vs-network", &ibmcloud.PiNetworkArgs{
			PiNetworkName:      pulumi.String("power-vs-private-network"),
			PiCloudInstanceId:  powerVSWorkspace.Guid,
			PiNetworkType:      pulumi.String("vlan"), // vlan or pub-vlan
			PiCidr:             pulumi.String("192.168.100.0/24"),
			PiGateway:          pulumi.String("192.168.100.1"),
			PiDnsServers: pulumi.StringArray{
				pulumi.String("9.9.9.9"),
				pulumi.String("1.1.1.1"),
			},
		})
		if err != nil {
			return err
		}

		// Step 5: Get available Power VS images (AIX, IBM i, or Linux)
		// Note: In production, you'd use GetPiImages to list available images
		// For this example, we'll use a common RHEL image name pattern

		// Step 6: Create Power VS Instance (Virtual Machine)
		powerVSInstance, err := ibmcloud.NewPiInstance(ctx, "power-vs-instance", &ibmcloud.PiInstanceArgs{
			PiInstanceName:    pulumi.String("power-vs-demo-vm"),
			PiCloudInstanceId: powerVSWorkspace.Guid,

			// System configuration
			PiMemory:          pulumi.Float64(4),  // 4 GB RAM
			PiProcessors:      pulumi.Float64(0.5), // 0.5 cores
			PiProcType:        pulumi.String("shared"), // shared or dedicated
			PiSysType:         pulumi.String("s922"), // System type: s922, e980, etc.

			// Image - Use RHEL or other available OS
			PiImageId:         pulumi.String("rhel-8-4"), // Replace with actual image ID from GetPiImages

			// Network configuration
			PiNetworks: ibmcloud.PiInstancePiNetworkArray{
				&ibmcloud.PiInstancePiNetworkArgs{
					NetworkId: powerVSNetwork.NetworkId,
				},
			},

			// SSH Key
			PiKeyPairName: sshKey.PiKeyName,

			// Storage
			PiStorageType: pulumi.String("tier3"), // tier1 (NVMe), tier3 (SSD)

			// Health check
			PiHealthStatus: pulumi.String("OK"),
		})
		if err != nil {
			return err
		}

		// Step 7: Create VPC for Load Balancer (Power VS instances need VPC LB for public access)
		vpc, err := ibmcloud.NewIsVpc(ctx, "power-vs-vpc", &ibmcloud.IsVpcArgs{
			Name:                  pulumi.String("power-vs-demo-vpc"),
			ResourceGroup:         resourceGroup.ID(),
			AddressPrefixManagement: pulumi.String("auto"),
			Tags: pulumi.StringArray{
				pulumi.String("environment:demo"),
			},
		})
		if err != nil {
			return err
		}

		// Step 8: Create VPC Subnet for Load Balancer
		subnet, err := ibmcloud.NewIsSubnet(ctx, "lb-subnet", &ibmcloud.IsSubnetArgs{
			Name:           pulumi.String("power-vs-lb-subnet"),
			Vpc:            vpc.ID(),
			Zone:           pulumi.String(zone),
			Ipv4CidrBlock:  pulumi.String("10.240.0.0/24"),
			ResourceGroup:  resourceGroup.ID(),
		})
		if err != nil {
			return err
		}

		// Step 9: Create Public Gateway for outbound internet access
		publicGateway, err := ibmcloud.NewIsPublicGateway(ctx, "public-gateway", &ibmcloud.IsPublicGatewayArgs{
			Name:          pulumi.String("power-vs-pgw"),
			Vpc:           vpc.ID(),
			Zone:          pulumi.String(zone),
			ResourceGroup: resourceGroup.ID(),
		})
		if err != nil {
			return err
		}

		// Attach public gateway to subnet
		_, err = ibmcloud.NewIsSubnetPublicGatewayAttachment(ctx, "subnet-pgw-attachment", &ibmcloud.IsSubnetPublicGatewayAttachmentArgs{
			Subnet:        subnet.ID(),
			PublicGateway: publicGateway.ID(),
		})
		if err != nil {
			return err
		}

		// Step 10: Create Security Group for Load Balancer
		securityGroup, err := ibmcloud.NewIsSecurityGroup(ctx, "lb-security-group", &ibmcloud.IsSecurityGroupArgs{
			Name:          pulumi.String("power-vs-lb-sg"),
			Vpc:           vpc.ID(),
			ResourceGroup: resourceGroup.ID(),
		})
		if err != nil {
			return err
		}

		// Allow HTTP inbound traffic
		_, err = ibmcloud.NewIsSecurityGroupRule(ctx, "allow-http-inbound", &ibmcloud.IsSecurityGroupRuleArgs{
			Group:     securityGroup.ID(),
			Direction: pulumi.String("inbound"),
			Remote:    pulumi.String("0.0.0.0/0"),
			Tcp: &ibmcloud.IsSecurityGroupRuleTcpArgs{
				PortMin: pulumi.Int(80),
				PortMax: pulumi.Int(80),
			},
		})
		if err != nil {
			return err
		}

		// Allow HTTPS inbound traffic
		_, err = ibmcloud.NewIsSecurityGroupRule(ctx, "allow-https-inbound", &ibmcloud.IsSecurityGroupRuleArgs{
			Group:     securityGroup.ID(),
			Direction: pulumi.String("inbound"),
			Remote:    pulumi.String("0.0.0.0/0"),
			Tcp: &ibmcloud.IsSecurityGroupRuleTcpArgs{
				PortMin: pulumi.Int(443),
				PortMax: pulumi.Int(443),
			},
		})
		if err != nil {
			return err
		}

		// Allow SSH inbound traffic
		_, err = ibmcloud.NewIsSecurityGroupRule(ctx, "allow-ssh-inbound", &ibmcloud.IsSecurityGroupRuleArgs{
			Group:     securityGroup.ID(),
			Direction: pulumi.String("inbound"),
			Remote:    pulumi.String("0.0.0.0/0"),
			Tcp: &ibmcloud.IsSecurityGroupRuleTcpArgs{
				PortMin: pulumi.Int(22),
				PortMax: pulumi.Int(22),
			},
		})
		if err != nil {
			return err
		}

		// Allow all outbound traffic
		_, err = ibmcloud.NewIsSecurityGroupRule(ctx, "allow-all-outbound", &ibmcloud.IsSecurityGroupRuleArgs{
			Group:     securityGroup.ID(),
			Direction: pulumi.String("outbound"),
			Remote:    pulumi.String("0.0.0.0/0"),
		})
		if err != nil {
			return err
		}

		// Step 11: Create Load Balancer
		loadBalancer, err := ibmcloud.NewIsLb(ctx, "power-vs-lb", &ibmcloud.IsLbArgs{
			Name:          pulumi.String("power-vs-demo-lb"),
			Subnets: pulumi.StringArray{
				subnet.ID(),
			},
			Type:          pulumi.String("public"), // public or private
			ResourceGroup: resourceGroup.ID(),
			Tags: pulumi.StringArray{
				pulumi.String("environment:demo"),
			},
		})
		if err != nil {
			return err
		}

		// Step 12: Create Load Balancer Pool
		lbPool, err := ibmcloud.NewIsLbPool(ctx, "lb-pool", &ibmcloud.IsLbPoolArgs{
			Lb:                     loadBalancer.ID(),
			Name:                   pulumi.String("power-vs-pool"),
			Algorithm:              pulumi.String("round_robin"), // round_robin, weighted_round_robin, least_connections
			Protocol:               pulumi.String("http"),
			HealthDelay:            pulumi.Int(5),
			HealthRetries:          pulumi.Int(2),
			HealthTimeout:          pulumi.Int(2),
			HealthType:             pulumi.String("http"),
			HealthMonitorUrl:       pulumi.String("/"),
			HealthMonitorPort:      pulumi.Int(80),
		})
		if err != nil {
			return err
		}

		// Step 13: Add Power VS Instance to Load Balancer Pool
		// Note: You'll need the internal IP of the Power VS instance
		_, err = ibmcloud.NewIsLbPoolMember(ctx, "lb-pool-member", &ibmcloud.IsLbPoolMemberArgs{
			Lb:         loadBalancer.ID(),
			Pool:       lbPool.ID(),
			Port:       pulumi.Int(80),
			TargetAddress: powerVSInstance.PiNetworks.Index(pulumi.Int(0)).IpAddress(), // Internal IP
			Weight:     pulumi.Int(50),
		})
		if err != nil {
			return err
		}

		// Step 14: Create Load Balancer Listener
		_, err = ibmcloud.NewIsLbListener(ctx, "lb-listener-http", &ibmcloud.IsLbListenerArgs{
			Lb:           loadBalancer.ID(),
			DefaultPool:  lbPool.ID(),
			Port:         pulumi.Int(80),
			Protocol:     pulumi.String("http"),
		})
		if err != nil {
			return err
		}

		// Optional: HTTPS Listener (requires SSL certificate)
		// _, err = ibmcloud.NewIsLbListener(ctx, "lb-listener-https", &ibmcloud.IsLbListenerArgs{
		// 	Lb:           loadBalancer.ID(),
		// 	DefaultPool:  lbPool.ID(),
		// 	Port:         pulumi.Int(443),
		// 	Protocol:     pulumi.String("https"),
		// 	CertificateInstance: pulumi.String("crn:v1:..."), // Certificate Manager CRN
		// })
		// if err != nil {
		// 	return err
		// }

		// Export important values
		ctx.Export("resourceGroupId", resourceGroup.ID())
		ctx.Export("resourceGroupName", resourceGroup.Name)

		ctx.Export("powerVsWorkspaceId", powerVSWorkspace.ID())
		ctx.Export("powerVsWorkspaceGuid", powerVSWorkspace.Guid)

		ctx.Export("powerVsNetworkId", powerVSNetwork.NetworkId)
		ctx.Export("powerVsNetworkName", powerVSNetwork.PiNetworkName)

		ctx.Export("powerVsInstanceId", powerVSInstance.PiInstanceId)
		ctx.Export("powerVsInstanceName", powerVSInstance.PiInstanceName)
		ctx.Export("powerVsInstanceStatus", powerVSInstance.Status)
		ctx.Export("powerVsInstanceInternalIP", powerVSInstance.PiNetworks.Index(pulumi.Int(0)).IpAddress())

		ctx.Export("vpcId", vpc.ID())
		ctx.Export("vpcName", vpc.Name)

		ctx.Export("loadBalancerId", loadBalancer.ID())
		ctx.Export("loadBalancerHostname", loadBalancer.Hostname)
		ctx.Export("loadBalancerPublicIps", loadBalancer.PublicIps)

		// Export access information
		ctx.Export("accessUrl", pulumi.Sprintf("http://%s", loadBalancer.Hostname))
		ctx.Export("sshCommand", pulumi.Sprintf("ssh root@%s", loadBalancer.Hostname))

		return nil
	})
}
