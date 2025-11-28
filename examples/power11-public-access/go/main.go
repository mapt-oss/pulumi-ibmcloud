package main

import (
	"github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Load configuration
		cfg := config.New(ctx, "")
		region := cfg.Get("region")
		if region == "" {
			region = "us-south"
		}
		zone := cfg.Get("zone")
		if zone == "" {
			zone = "us-south-1"
		}
		powerVsZone := cfg.Get("powerVsZone")
		if powerVsZone == "" {
			powerVsZone = "us-south"
		}
		sshPublicKey := cfg.Require("sshPublicKey")
		instanceName := cfg.Get("instanceName")
		if instanceName == "" {
			instanceName = "power11-vm"
		}

		// Step 1: Create Resource Group
		resourceGroup, err := ibmcloud.NewResourceGroup(ctx, "power11-rg", &ibmcloud.ResourceGroupArgs{
			Name: pulumi.String("power11-public-access-rg"),
			Tags: pulumi.StringArray{
				pulumi.String("environment:production"),
				pulumi.String("managed-by:pulumi"),
				pulumi.String("power-architecture:power11"),
			},
		})
		if err != nil {
			return err
		}

		// Step 2: Create Power VS Workspace for Power11
		powerVsWorkspace, err := ibmcloud.NewResourceInstance(ctx, "power11-workspace", &ibmcloud.ResourceInstanceArgs{
			Name:            pulumi.String("power11-workspace"),
			Service:         pulumi.String("power-iaas"),
			Plan:            pulumi.String("power-virtual-server-group"),
			Location:        pulumi.String(powerVsZone),
			ResourceGroupId: resourceGroup.ID(),
			Tags: pulumi.StringArray{
				pulumi.String("power11"),
				pulumi.String("production"),
			},
		})
		if err != nil {
			return err
		}

		// Step 3: Create SSH Key
		sshKey, err := ibmcloud.NewPiKey(ctx, "power11-ssh-key", &ibmcloud.PiKeyArgs{
			PiKeyName:         pulumi.String("power11-access-key"),
			PiSshKey:          pulumi.String(sshPublicKey),
			PiCloudInstanceId: powerVsWorkspace.Guid,
		})
		if err != nil {
			return err
		}

		// Step 4: Create Private Network for Power VS
		powerVsNetwork, err := ibmcloud.NewPiNetwork(ctx, "power11-network", &ibmcloud.PiNetworkArgs{
			PiNetworkName:     pulumi.String("power11-private-net"),
			PiCloudInstanceId: powerVsWorkspace.Guid,
			PiNetworkType:     pulumi.String("vlan"),
			PiCidr:            pulumi.String("192.168.50.0/24"),
			PiGateway:         pulumi.String("192.168.50.1"),
			PiDnsServers: pulumi.StringArray{
				pulumi.String("9.9.9.9"),
				pulumi.String("1.1.1.1"),
			},
		})
		if err != nil {
			return err
		}

		// Step 5: Create Power11 Virtual Server Instance
		// Power11 uses system type s1022 (Power10/11 architecture)
		power11Instance, err := ibmcloud.NewPiInstance(ctx, "power11-instance", &ibmcloud.PiInstanceArgs{
			PiInstanceName:    pulumi.String(instanceName),
			PiCloudInstanceId: powerVsWorkspace.Guid,

			// Power11 Specifications
			PiMemory:     pulumi.Float64(16),         // 16 GB RAM
			PiProcessors: pulumi.Float64(2),          // 2 cores
			PiProcType:   pulumi.String("dedicated"), // dedicated or shared
			PiSysType:    pulumi.String("s1022"),     // Power10/11 system type

			// Operating System - RHEL 9 for Power11
			PiImageId: pulumi.String("rhel-9-2"), // Replace with actual image ID

			// Network Configuration
			PiNetworks: ibmcloud.PiInstancePiNetworkArray{
				&ibmcloud.PiInstancePiNetworkArgs{
					NetworkId: powerVsNetwork.NetworkId,
				},
			},

			// SSH Access
			PiKeyPairName: sshKey.PiKeyName,

			// Storage - tier1 (NVMe) for best performance
			PiStorageType: pulumi.String("tier1"),
			PiStoragePool: pulumi.String("Tier1-Flash"),

			// Health
			PiHealthStatus: pulumi.String("OK"),
		}, pulumi.Timeouts(&pulumi.CustomTimeouts{
			Create: "30m",
			Update: "30m",
			Delete: "30m",
		}))
		if err != nil {
			return err
		}

		// Step 6: Create VPC for Load Balancer
		vpc, err := ibmcloud.NewIsVpc(ctx, "power11-vpc", &ibmcloud.IsVpcArgs{
			Name:                    pulumi.String("power11-vpc"),
			ResourceGroup:           resourceGroup.ID(),
			AddressPrefixManagement: pulumi.String("auto"),
			Tags: pulumi.StringArray{
				pulumi.String("power11"),
				pulumi.String("loadbalancer"),
			},
		})
		if err != nil {
			return err
		}

		// Step 7: Create Subnet for Load Balancer
		lbSubnet, err := ibmcloud.NewIsSubnet(ctx, "lb-subnet", &ibmcloud.IsSubnetArgs{
			Name:          pulumi.String("power11-lb-subnet"),
			Vpc:           vpc.ID(),
			Zone:          pulumi.String(zone),
			Ipv4CidrBlock: pulumi.String("10.240.64.0/24"),
			ResourceGroup: resourceGroup.ID(),
		})
		if err != nil {
			return err
		}

		// Step 8: Create Public Gateway
		publicGateway, err := ibmcloud.NewIsPublicGateway(ctx, "public-gateway", &ibmcloud.IsPublicGatewayArgs{
			Name:          pulumi.String("power11-pgw"),
			Vpc:           vpc.ID(),
			Zone:          pulumi.String(zone),
			ResourceGroup: resourceGroup.ID(),
		})
		if err != nil {
			return err
		}

		// Attach Public Gateway to Subnet
		_, err = ibmcloud.NewIsSubnetPublicGatewayAttachment(ctx, "pgw-attachment", &ibmcloud.IsSubnetPublicGatewayAttachmentArgs{
			Subnet:        lbSubnet.ID(),
			PublicGateway: publicGateway.ID(),
		})
		if err != nil {
			return err
		}

		// Step 9: Create Security Group
		securityGroup, err := ibmcloud.NewIsSecurityGroup(ctx, "power11-sg", &ibmcloud.IsSecurityGroupArgs{
			Name:          pulumi.String("power11-security-group"),
			Vpc:           vpc.ID(),
			ResourceGroup: resourceGroup.ID(),
		})
		if err != nil {
			return err
		}

		// Security Rules - HTTP
		_, err = ibmcloud.NewIsSecurityGroupRule(ctx, "allow-http", &ibmcloud.IsSecurityGroupRuleArgs{
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

		// Security Rules - HTTPS
		_, err = ibmcloud.NewIsSecurityGroupRule(ctx, "allow-https", &ibmcloud.IsSecurityGroupRuleArgs{
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

		// Security Rules - SSH (custom port 2222)
		_, err = ibmcloud.NewIsSecurityGroupRule(ctx, "allow-ssh", &ibmcloud.IsSecurityGroupRuleArgs{
			Group:     securityGroup.ID(),
			Direction: pulumi.String("inbound"),
			Remote:    pulumi.String("0.0.0.0/0"),
			Tcp: &ibmcloud.IsSecurityGroupRuleTcpArgs{
				PortMin: pulumi.Int(2222),
				PortMax: pulumi.Int(2222),
			},
		})
		if err != nil {
			return err
		}

		// Security Rules - Allow all outbound
		_, err = ibmcloud.NewIsSecurityGroupRule(ctx, "allow-outbound", &ibmcloud.IsSecurityGroupRuleArgs{
			Group:     securityGroup.ID(),
			Direction: pulumi.String("outbound"),
			Remote:    pulumi.String("0.0.0.0/0"),
		})
		if err != nil {
			return err
		}

		// Step 10: Create Application Load Balancer
		loadBalancer, err := ibmcloud.NewIsLb(ctx, "power11-lb", &ibmcloud.IsLbArgs{
			Name: pulumi.String("power11-alb"),
			Subnets: pulumi.StringArray{
				lbSubnet.ID(),
			},
			Type:          pulumi.String("public"),
			ResourceGroup: resourceGroup.ID(),
			SecurityGroups: pulumi.StringArray{
				securityGroup.ID(),
			},
			Tags: pulumi.StringArray{
				pulumi.String("power11"),
				pulumi.String("public-access"),
			},
		}, pulumi.Timeouts(&pulumi.CustomTimeouts{
			Create: "15m",
			Update: "15m",
			Delete: "15m",
		}))
		if err != nil {
			return err
		}

		// Step 11: Create Backend Pool
		backendPool, err := ibmcloud.NewIsLbPool(ctx, "backend-pool", &ibmcloud.IsLbPoolArgs{
			Lb:               loadBalancer.ID(),
			Name:             pulumi.String("power11-backend-pool"),
			Algorithm:        pulumi.String("round_robin"),
			Protocol:         pulumi.String("http"),
			HealthDelay:      pulumi.Int(10),
			HealthRetries:    pulumi.Int(3),
			HealthTimeout:    pulumi.Int(5),
			HealthType:       pulumi.String("http"),
			HealthMonitorUrl: pulumi.String("/"),
			HealthMonitorPort: pulumi.Int(80),
		})
		if err != nil {
			return err
		}

		// Step 12: Add Power11 Instance to Backend Pool
		_, err = ibmcloud.NewIsLbPoolMember(ctx, "power11-member", &ibmcloud.IsLbPoolMemberArgs{
			Lb:            loadBalancer.ID(),
			Pool:          backendPool.ID(),
			Port:          pulumi.Int(80),
			TargetAddress: power11Instance.PiNetworks.Index(pulumi.Int(0)).IpAddress(),
			Weight:        pulumi.Int(100),
		})
		if err != nil {
			return err
		}

		// Step 13: Create HTTP Listener
		_, err = ibmcloud.NewIsLbListener(ctx, "http-listener", &ibmcloud.IsLbListenerArgs{
			Lb:              loadBalancer.ID(),
			DefaultPool:     backendPool.ID(),
			Port:            pulumi.Int(80),
			Protocol:        pulumi.String("http"),
			ConnectionLimit: pulumi.Int(2000),
		})
		if err != nil {
			return err
		}

		// Step 14: Create SSH Listener (TCP passthrough on port 2222)
		_, err = ibmcloud.NewIsLbListener(ctx, "ssh-listener", &ibmcloud.IsLbListenerArgs{
			Lb:              loadBalancer.ID(),
			DefaultPool:     backendPool.ID(),
			Port:            pulumi.Int(2222),
			Protocol:        pulumi.String("tcp"),
			ConnectionLimit: pulumi.Int(100),
		})
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("resourceGroupId", resourceGroup.ID())
		ctx.Export("resourceGroupName", resourceGroup.Name)

		ctx.Export("powerVsWorkspaceId", powerVsWorkspace.ID())
		ctx.Export("powerVsWorkspaceGuid", powerVsWorkspace.Guid)

		ctx.Export("powerVsNetworkId", powerVsNetwork.NetworkId)
		ctx.Export("powerVsNetworkCidr", powerVsNetwork.PiCidr)

		ctx.Export("power11InstanceId", power11Instance.PiInstanceId)
		ctx.Export("power11InstanceName", power11Instance.PiInstanceName)
		ctx.Export("power11InstanceStatus", power11Instance.Status)
		ctx.Export("power11InstanceInternalIp", power11Instance.PiNetworks.Index(pulumi.Int(0)).IpAddress())
		ctx.Export("power11InstanceMacAddress", power11Instance.PiNetworks.Index(pulumi.Int(0)).MacAddress())
		ctx.Export("power11SystemType", power11Instance.PiSysType)
		ctx.Export("power11Memory", power11Instance.PiMemory)
		ctx.Export("power11Processors", power11Instance.PiProcessors)

		ctx.Export("vpcId", vpc.ID())
		ctx.Export("vpcCrn", vpc.Crn)

		ctx.Export("loadBalancerId", loadBalancer.ID())
		ctx.Export("loadBalancerHostname", loadBalancer.Hostname)
		ctx.Export("loadBalancerPublicIps", loadBalancer.PublicIps)

		// Access Information
		ctx.Export("httpUrl", pulumi.Sprintf("http://%s", loadBalancer.Hostname))
		ctx.Export("httpsUrl", pulumi.Sprintf("https://%s", loadBalancer.Hostname))
		ctx.Export("sshCommand", pulumi.Sprintf("ssh -p 2222 root@%s", loadBalancer.Hostname))

		// Summary
		ctx.Export("deploymentSummary", pulumi.All(
			loadBalancer.Hostname,
			power11Instance.PiNetworks.Index(pulumi.Int(0)).IpAddress(),
			power11Instance.Status,
		).ApplyT(func(args []interface{}) map[string]interface{} {
			hostname := args[0].(string)
			ip := args[1].(string)
			status := args[2].(string)

			return map[string]interface{}{
				"message": "Power11 instance deployed successfully!",
				"publicAccess": map[string]string{
					"http":  "http://" + hostname,
					"https": "https://" + hostname,
					"ssh":   "ssh -p 2222 root@" + hostname,
				},
				"instance": map[string]string{
					"internalIp": ip,
					"status":     status,
					"systemType": "s1022 (Power10/11)",
				},
			}
		}))

		return nil
	})
}
