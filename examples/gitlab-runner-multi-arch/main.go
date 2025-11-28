package main

import (
	"github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a resource group for the GitLab runners
		rg, err := ibmcloud.NewResourceGroup(ctx, "gitlab-runner-rg", &ibmcloud.ResourceGroupArgs{
			Name: pulumi.String("gitlab-runner-multiarch-rg"),
			Tags: pulumi.StringArray{
				pulumi.String("gitlab"),
				pulumi.String("ci-cd"),
				pulumi.String("multi-arch"),
			},
		})
		if err != nil {
			return err
		}

		// Create VPC
		vpc, err := ibmcloud.NewIsVpc(ctx, "gitlab-runner-vpc", &ibmcloud.IsVpcArgs{
			Name:          pulumi.String("gitlab-runner-vpc"),
			ResourceGroup: rg.ID(),
			Tags: pulumi.StringArray{
				pulumi.String("gitlab"),
				pulumi.String("network"),
			},
		})
		if err != nil {
			return err
		}

		// Create subnet for s390x runner (IBM Z)
		subnetS390x, err := ibmcloud.NewIsSubnet(ctx, "gitlab-runner-subnet-s390x", &ibmcloud.IsSubnetArgs{
			Name:          pulumi.String("gitlab-runner-subnet-s390x"),
			Vpc:           vpc.ID(),
			Zone:          pulumi.String("us-east-1"), // s390x available in us-east
			Ipv4CidrBlock: pulumi.String("10.240.0.0/24"),
			ResourceGroup: rg.ID(),
		})
		if err != nil {
			return err
		}

		// Create subnet for ppc64 runner (IBM Power)
		subnetPpc64, err := ibmcloud.NewIsSubnet(ctx, "gitlab-runner-subnet-ppc64", &ibmcloud.IsSubnetArgs{
			Name:          pulumi.String("gitlab-runner-subnet-ppc64"),
			Vpc:           vpc.ID(),
			Zone:          pulumi.String("us-south-2"), // ppc64 available in us-south
			Ipv4CidrBlock: pulumi.String("10.240.1.0/24"),
			ResourceGroup: rg.ID(),
		})
		if err != nil {
			return err
		}

		// Create Public Gateway for s390x subnet (provides outbound internet for GitLab.com)
		pgwS390x, err := ibmcloud.NewIsPublicGateway(ctx, "gitlab-runner-pgw-s390x", &ibmcloud.IsPublicGatewayArgs{
			Name:          pulumi.String("gitlab-runner-pgw-s390x"),
			Vpc:           vpc.ID(),
			Zone:          pulumi.String("us-east-1"),
			ResourceGroup: rg.ID(),
			Tags: pulumi.StringArray{
				pulumi.String("gitlab"),
				pulumi.String("s390x"),
			},
		})
		if err != nil {
			return err
		}

		// Create Public Gateway for ppc64 subnet (provides outbound internet for GitLab.com)
		pgwPpc64, err := ibmcloud.NewIsPublicGateway(ctx, "gitlab-runner-pgw-ppc64", &ibmcloud.IsPublicGatewayArgs{
			Name:          pulumi.String("gitlab-runner-pgw-ppc64"),
			Vpc:           vpc.ID(),
			Zone:          pulumi.String("us-south-2"),
			ResourceGroup: rg.ID(),
			Tags: pulumi.StringArray{
				pulumi.String("gitlab"),
				pulumi.String("ppc64"),
			},
		})
		if err != nil {
			return err
		}

		// Attach Public Gateways to subnets for outbound connectivity
		_, err = ibmcloud.NewIsSubnetPublicGatewayAttachment(ctx, "pgw-attachment-s390x", &ibmcloud.IsSubnetPublicGatewayAttachmentArgs{
			SubnetId:        subnetS390x.ID(),
			PublicGatewayId: pgwS390x.ID(),
		})
		if err != nil {
			return err
		}

		_, err = ibmcloud.NewIsSubnetPublicGatewayAttachment(ctx, "pgw-attachment-ppc64", &ibmcloud.IsSubnetPublicGatewayAttachmentArgs{
			SubnetId:        subnetPpc64.ID(),
			PublicGatewayId: pgwPpc64.ID(),
		})
		if err != nil {
			return err
		}

		// Create security group for GitLab runners
		securityGroup, err := ibmcloud.NewIsSecurityGroup(ctx, "gitlab-runner-sg", &ibmcloud.IsSecurityGroupArgs{
			Name:          pulumi.String("gitlab-runner-sg"),
			Vpc:           vpc.ID(),
			ResourceGroup: rg.ID(),
		})
		if err != nil {
			return err
		}

		// Allow SSH inbound
		_, err = ibmcloud.NewIsSecurityGroupRule(ctx, "allow-ssh", &ibmcloud.IsSecurityGroupRuleArgs{
			Group:     securityGroup.ID(),
			Direction: pulumi.String("inbound"),
			Remote:    pulumi.String("0.0.0.0/0"), // CHANGE THIS to your IP range for better security
			Tcp: &ibmcloud.IsSecurityGroupRuleTcpArgs{
				PortMin: pulumi.Int(22),
				PortMax: pulumi.Int(22),
			},
		})
		if err != nil {
			return err
		}

		// Allow all outbound traffic (needed for GitLab.com, Docker registries, package repos)
		_, err = ibmcloud.NewIsSecurityGroupRule(ctx, "allow-all-outbound", &ibmcloud.IsSecurityGroupRuleArgs{
			Group:     securityGroup.ID(),
			Direction: pulumi.String("outbound"),
			Remote:    pulumi.String("0.0.0.0/0"),
		})
		if err != nil {
			return err
		}

		// Create SSH key (REPLACE with your actual public key)
		sshKey, err := ibmcloud.NewIsSshKey(ctx, "gitlab-runner-key", &ibmcloud.IsSshKeyArgs{
			Name:          pulumi.String("gitlab-runner-key"),
			PublicKey:     pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC... your-public-key-here"),
			ResourceGroup: rg.ID(),
		})
		if err != nil {
			return err
		}

		// Fedora installation script for GitLab Runner
		fedoraUserData := `#!/bin/bash
# Update system
dnf update -y

# Install Docker
dnf install -y dnf-plugins-core
dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
dnf install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
systemctl enable docker
systemctl start docker

# Add default user to docker group
usermod -aG docker fedora

# Install GitLab Runner
curl -L "https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.rpm.sh" | bash
dnf install -y gitlab-runner

# Create GitLab Runner registration script
cat > /home/fedora/register-runner.sh << 'EOF'
#!/bin/bash
# Run this script with your GitLab registration token:
# sudo ./register-runner.sh YOUR_REGISTRATION_TOKEN [TAGS]

if [ -z "$1" ]; then
  echo "Usage: $0 <registration-token> [tags]"
  echo "Example: $0 glrt-xxxxxxxxxxxx s390x,docker"
  exit 1
fi

ARCH=$(uname -m)
TAGS="${2:-$ARCH,docker}"

gitlab-runner register \
  --non-interactive \
  --url "https://gitlab.com/" \
  --registration-token "$1" \
  --executor "docker" \
  --docker-image "fedora:latest" \
  --description "IBM Cloud Runner - $ARCH" \
  --tag-list "$TAGS" \
  --run-untagged="true" \
  --locked="false" \
  --access-level="not_protected"

echo "GitLab Runner registered successfully!"
echo "Architecture: $ARCH"
echo "Tags: $TAGS"
EOF

chmod +x /home/fedora/register-runner.sh
chown fedora:fedora /home/fedora/register-runner.sh

# Enable GitLab Runner service
systemctl enable gitlab-runner
systemctl start gitlab-runner

echo "GitLab Runner installation complete!"
echo "Architecture: $(uname -m)"
echo "To register the runner, SSH to this instance and run:"
echo "  sudo /home/fedora/register-runner.sh YOUR_TOKEN"
`

		// Create s390x (IBM Z) GitLab runner instance
		// NOTE: Replace with actual Fedora s390x image ID
		// Find it with: ibmcloud is images --visibility public | grep -i fedora | grep s390x
		imageS390x := pulumi.String("r014-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx") // REPLACE with actual image ID

		instanceS390x, err := ibmcloud.NewIsInstance(ctx, "gitlab-runner-s390x", &ibmcloud.IsInstanceArgs{
			Name:          pulumi.String("gitlab-runner-s390x"),
			Image:         imageS390x,
			Profile:       pulumi.String("bz2-2x8"), // s390x profile: 2 vCPU, 8 GB RAM
			Vpc:           vpc.ID(),
			Zone:          pulumi.String("us-east-1"),
			ResourceGroup: rg.ID(),
			Keys:          pulumi.StringArray{sshKey.ID()},
			PrimaryNetworkInterface: &ibmcloud.IsInstancePrimaryNetworkInterfaceArgs{
				Subnet: subnetS390x.ID(),
				SecurityGroups: pulumi.StringArray{
					securityGroup.ID(),
				},
			},
			UserData: pulumi.String(fedoraUserData),
			Tags: pulumi.StringArray{
				pulumi.String("gitlab-runner"),
				pulumi.String("s390x"),
				pulumi.String("ibm-z"),
			},
		})
		if err != nil {
			return err
		}

		// Create ppc64 (IBM Power) GitLab runner instance
		// NOTE: Replace with actual Fedora ppc64 image ID
		// Find it with: ibmcloud is images --visibility public | grep -i fedora | grep ppc
		imagePpc64 := pulumi.String("r006-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx") // REPLACE with actual image ID

		instancePpc64, err := ibmcloud.NewIsInstance(ctx, "gitlab-runner-ppc64", &ibmcloud.IsInstanceArgs{
			Name:          pulumi.String("gitlab-runner-ppc64"),
			Image:         imagePpc64,
			Profile:       pulumi.String("bp2-2x8"), // ppc64 profile: 2 vCPU, 8 GB RAM
			Vpc:           vpc.ID(),
			Zone:          pulumi.String("us-south-2"),
			ResourceGroup: rg.ID(),
			Keys:          pulumi.StringArray{sshKey.ID()},
			PrimaryNetworkInterface: &ibmcloud.IsInstancePrimaryNetworkInterfaceArgs{
				Subnet: subnetPpc64.ID(),
				SecurityGroups: pulumi.StringArray{
					securityGroup.ID(),
				},
			},
			UserData: pulumi.String(fedoraUserData),
			Tags: pulumi.StringArray{
				pulumi.String("gitlab-runner"),
				pulumi.String("ppc64"),
				pulumi.String("ibm-power"),
			},
		})
		if err != nil {
			return err
		}

		// Create floating IPs for SSH access
		floatingIpS390x, err := ibmcloud.NewIsFloatingIp(ctx, "gitlab-runner-fip-s390x", &ibmcloud.IsFloatingIpArgs{
			Name: pulumi.String("gitlab-runner-fip-s390x"),
			Target: instanceS390x.PrimaryNetworkInterface.ApplyT(func(ni ibmcloud.IsInstancePrimaryNetworkInterface) string {
				return *ni.Id
			}).(pulumi.StringOutput),
			Zone:          pulumi.String("us-east-1"),
			ResourceGroup: rg.ID(),
		})
		if err != nil {
			return err
		}

		floatingIpPpc64, err := ibmcloud.NewIsFloatingIp(ctx, "gitlab-runner-fip-ppc64", &ibmcloud.IsFloatingIpArgs{
			Name: pulumi.String("gitlab-runner-fip-ppc64"),
			Target: instancePpc64.PrimaryNetworkInterface.ApplyT(func(ni ibmcloud.IsInstancePrimaryNetworkInterface) string {
				return *ni.Id
			}).(pulumi.StringOutput),
			Zone:          pulumi.String("us-south-2"),
			ResourceGroup: rg.ID(),
		})
		if err != nil {
			return err
		}

		// Export important values
		ctx.Export("resourceGroupId", rg.ID())
		ctx.Export("vpcId", vpc.ID())

		// s390x exports
		ctx.Export("s390x_instanceId", instanceS390x.ID())
		ctx.Export("s390x_privateIP", instanceS390x.PrimaryNetworkInterface.ApplyT(func(ni ibmcloud.IsInstancePrimaryNetworkInterface) string {
			return *ni.PrimaryIpv4Address
		}))
		ctx.Export("s390x_floatingIP", floatingIpS390x.Address)
		ctx.Export("s390x_sshCommand", pulumi.Sprintf("ssh fedora@%s", floatingIpS390x.Address))

		// ppc64 exports
		ctx.Export("ppc64_instanceId", instancePpc64.ID())
		ctx.Export("ppc64_privateIP", instancePpc64.PrimaryNetworkInterface.ApplyT(func(ni ibmcloud.IsInstancePrimaryNetworkInterface) string {
			return *ni.PrimaryIpv4Address
		}))
		ctx.Export("ppc64_floatingIP", floatingIpPpc64.Address)
		ctx.Export("ppc64_sshCommand", pulumi.Sprintf("ssh fedora@%s", floatingIpPpc64.Address))

		return nil
	})
}
