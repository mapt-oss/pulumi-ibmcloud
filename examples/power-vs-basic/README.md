# IBM Cloud Power Virtual Server - Basic Example

A minimal example showing how to create an IBM Cloud Power Virtual Server instance with private networking.

## What This Creates

- **Resource Group**: Container for all resources
- **Power VS Workspace**: Power Virtual Server service instance
- **Private Network**: VLAN network (192.168.200.0/24)
- **SSH Key**: For secure access
- **Power VS Instance**: Virtual machine running RHEL
  - 0.25 shared processors
  - 4 GB RAM
  - Tier 3 (SSD) storage

## Quick Start

```bash
# 1. Configure
pulumi config set ibmcloud:ibmcloudApiKey --secret
pulumi config set sshPublicKey "$(cat ~/.ssh/id_rsa.pub)"
pulumi config set powerVsZone dal12  # Optional

# 2. Deploy
pulumi up

# 3. Get instance details
pulumi stack output instanceIP
pulumi stack output instanceStatus
```

## Cost

Estimated cost: **~$20-30/month**
- Power VS: 0.25 vCPU, 4GB RAM, 100GB storage
- Private network: Free

## Access

This example creates a private instance with no public access. To connect:

**Option 1: VPN Connection**
```bash
# Set up IBM Cloud VPN and connect to private network
ssh root@$(pulumi stack output instanceIP)
```

**Option 2: Add Load Balancer**
See the [power-vs-with-lb](../power-vs-with-lb) example for public access.

**Option 3: Jump Host**
Create a VPC VSI as a bastion host in the same region.

## Cleanup

```bash
pulumi destroy
```

## Next Steps

- Add [Load Balancer](../power-vs-with-lb) for public access
- Add storage volumes
- Add multiple instances
- Configure backups
