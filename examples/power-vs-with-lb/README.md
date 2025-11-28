# IBM Cloud Power Virtual Server with Load Balancer

This example demonstrates how to create an IBM Cloud Power Virtual Server (Power VS) instance with networking and a VPC Load Balancer for public access.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         IBM Cloud                                │
│                                                                   │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                    Resource Group                           │ │
│  │                                                              │ │
│  │  ┌──────────────────────┐    ┌─────────────────────────┐  │ │
│  │  │   Power VS Workspace │    │      VPC Network         │  │ │
│  │  │   (dal12)            │    │      (us-south)          │  │ │
│  │  │                      │    │                          │  │ │
│  │  │  ┌───────────────┐   │    │  ┌────────────────────┐ │  │ │
│  │  │  │ Private Net   │   │    │  │   Public Gateway   │ │  │ │
│  │  │  │ 192.168.100/24│   │    │  │   (Outbound NAT)   │ │  │ │
│  │  │  └───────┬───────┘   │    │  └────────────────────┘ │  │ │
│  │  │          │           │    │           │              │  │ │
│  │  │  ┌───────▼───────┐   │    │  ┌────────▼───────────┐ │  │ │
│  │  │  │  Power VS VM  │   │    │  │  Load Balancer     │ │  │ │
│  │  │  │  RHEL 8.4     │◄──┼────┼──┤  (Public)          │ │  │ │
│  │  │  │  4GB / 0.5vCPU│   │    │  │  - HTTP:80         │ │  │ │
│  │  │  │  192.168.100.x│   │    │  │  - HTTPS:443       │ │  │ │
│  │  │  └───────────────┘   │    │  └────────────────────┘ │  │ │
│  │  │                      │    │           │              │  │ │
│  │  └──────────────────────┘    └───────────┼──────────────┘  │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                              │                   │
└──────────────────────────────────────────────┼───────────────────┘
                                               │
                                         ┌─────▼─────┐
                                         │  Internet │
                                         │   Users   │
                                         └───────────┘
```

## Components Created

### 1. Power Virtual Server Resources
- **Power VS Workspace**: Service instance for Power VS resources
- **Power VS Network**: Private VLAN network (192.168.100.0/24)
- **Power VS Instance**: Virtual machine running RHEL 8.4
  - 4 GB RAM
  - 0.5 shared processors
  - Tier 3 (SSD) storage
  - SSH key authentication

### 2. VPC Resources
- **VPC**: Virtual Private Cloud for load balancer
- **Subnet**: Subnet for load balancer (10.240.0.0/24)
- **Public Gateway**: For outbound internet connectivity
- **Security Group**: Firewall rules for HTTP/HTTPS/SSH

### 3. Load Balancer
- **Load Balancer**: Public Application Load Balancer
- **Backend Pool**: Routes traffic to Power VS instance
- **Health Checks**: HTTP health monitoring
- **Listeners**: HTTP (port 80) and optional HTTPS (port 443)

## Prerequisites

### 1. IBM Cloud Account
- Active IBM Cloud account
- IBM Cloud API key with appropriate permissions

### 2. Tools
- [Pulumi CLI](https://www.pulumi.com/docs/install/) (v3.0+)
- [Go](https://go.dev/dl/) (v1.23+)
- SSH key pair

### 3. Permissions
Your IBM Cloud API key needs:
- **Power VS**: Create and manage Power Virtual Servers
- **VPC**: Create VPC, subnets, load balancers
- **Resource Groups**: Create and manage resource groups

## Configuration

### Required Configuration

1. **IBM Cloud API Key**:
   ```bash
   export IC_API_KEY="your-ibm-cloud-api-key"
   # OR
   pulumi config set ibmcloud:ibmcloudApiKey --secret
   ```

2. **SSH Public Key**:
   ```bash
   pulumi config set sshPublicKey "$(cat ~/.ssh/id_rsa.pub)"
   ```

### Optional Configuration

```bash
# VPC Region (default: us-south)
pulumi config set region eu-de

# VPC Zone (default: us-south-1)
pulumi config set zone eu-de-1

# Power VS Zone (default: dal12)
# Available zones: dal12, dal13, us-south, us-east, eu-de-1, eu-de-2, lon04, lon06, tok04, syd04, syd05
pulumi config set powerVsZone lon06
```

### Power VS Zone Mapping

| VPC Region | Recommended Power VS Zones |
|------------|---------------------------|
| us-south   | dal12, dal13              |
| us-east    | wdc06, wdc07              |
| eu-de      | eu-de-1, eu-de-2, fra04, fra05 |
| eu-gb      | lon04, lon06              |
| jp-tok     | tok04                     |
| au-syd     | syd04, syd05              |

## Deployment

### 1. Initialize Go Module

```bash
cd examples/power-vs-with-lb
go mod init power-vs-with-lb
go mod tidy
```

### 2. Install Dependencies

```bash
go get github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud
go get github.com/pulumi/pulumi/sdk/v3/go/pulumi
```

### 3. Configure Pulumi

```bash
# Create a new stack
pulumi stack init dev

# Set required configuration
pulumi config set ibmcloud:ibmcloudApiKey --secret
# Enter your IBM Cloud API key when prompted

pulumi config set sshPublicKey "$(cat ~/.ssh/id_rsa.pub)"

# Optional: Set regions
pulumi config set region us-south
pulumi config set zone us-south-1
pulumi config set powerVsZone dal12
```

### 4. Preview and Deploy

```bash
# Preview changes
pulumi preview

# Deploy infrastructure
pulumi up
```

Deployment takes approximately **15-20 minutes**:
- Power VS workspace: ~2 minutes
- Power VS instance: ~10-15 minutes
- VPC and Load Balancer: ~3-5 minutes

### 5. Access Your Power VS Instance

After deployment completes:

```bash
# Get the load balancer hostname
LB_HOSTNAME=$(pulumi stack output loadBalancerHostname)

# Access via HTTP
curl http://$LB_HOSTNAME

# SSH to the instance (if SSH listener configured)
ssh root@$LB_HOSTNAME

# Or get the direct command
pulumi stack output accessUrl
pulumi stack output sshCommand
```

## Outputs

The stack exports the following outputs:

| Output | Description |
|--------|-------------|
| `resourceGroupId` | Resource group ID |
| `powerVsWorkspaceId` | Power VS workspace ID |
| `powerVsWorkspaceGuid` | Power VS workspace GUID |
| `powerVsNetworkId` | Power VS network ID |
| `powerVsInstanceId` | Power VS instance ID |
| `powerVsInstanceInternalIP` | Internal IP of Power VS instance |
| `vpcId` | VPC ID |
| `loadBalancerId` | Load balancer ID |
| `loadBalancerHostname` | Load balancer public hostname |
| `loadBalancerPublicIps` | Load balancer public IPs |
| `accessUrl` | HTTP access URL |
| `sshCommand` | SSH command to connect |

## Costs

Estimated monthly costs for this configuration:

| Resource | Specification | Estimated Cost (USD/month) |
|----------|--------------|---------------------------|
| Power VS Instance | 0.5 vCPU, 4GB RAM, Shared | ~$25-40 |
| Power VS Storage | 100GB Tier 3 (SSD) | ~$10-15 |
| Power VS Network | Private VLAN | $0 |
| VPC Load Balancer | Public ALB | ~$60-80 |
| VPC Data Transfer | First 5GB free | Variable |
| **Total** | | **~$95-135/month** |

> **Note**: Costs vary by region. Power VS pricing is based on processor type (shared/dedicated), memory, and storage tier. Check [IBM Cloud Pricing](https://cloud.ibm.com/pricing) for current rates.

## Customization

### Change Power VS Instance Specifications

```go
powerVSInstance, err := ibmcloud.NewPiInstance(ctx, "power-vs-instance", &ibmcloud.PiInstanceArgs{
    // Increase memory
    PiMemory: pulumi.Float64(8),  // 8 GB RAM

    // Increase processors
    PiProcessors: pulumi.Float64(1), // 1 full core

    // Use dedicated processors
    PiProcType: pulumi.String("dedicated"),

    // Use faster storage
    PiStorageType: pulumi.String("tier1"), // NVMe storage

    // Change system type
    PiSysType: pulumi.String("e980"), // or s922, s1022
})
```

### Use Different Operating Systems

Available OS images (check with `ibmcloud pi images`):
- **AIX**: aix-7.2, aix-7.3
- **IBM i**: ibm-i-7.4, ibm-i-7.5
- **Linux**: rhel-8.4, rhel-9.0, sles-15-sp4

```go
// Use AIX
PiImageId: pulumi.String("aix-7.3"),

// Use IBM i
PiImageId: pulumi.String("ibm-i-7.5"),

// Use SLES
PiImageId: pulumi.String("sles-15-sp4"),
```

### Add HTTPS Listener

1. Create certificate in IBM Cloud Certificate Manager
2. Uncomment HTTPS listener code in main.go
3. Add certificate CRN:

```go
_, err = ibmcloud.NewIsLbListener(ctx, "lb-listener-https", &ibmcloud.IsLbListenerArgs{
    Lb:                  loadBalancer.ID(),
    DefaultPool:         lbPool.ID(),
    Port:                pulumi.Int(443),
    Protocol:            pulumi.String("https"),
    CertificateInstance: pulumi.String("crn:v1:bluemix:public:cloudcerts:..."),
})
```

### Add Multiple Power VS Instances

```go
// Create instance 2
powerVSInstance2, err := ibmcloud.NewPiInstance(ctx, "power-vs-instance-2", &ibmcloud.PiInstanceArgs{
    PiInstanceName:    pulumi.String("power-vs-demo-vm-2"),
    // ... same configuration ...
})

// Add to load balancer pool
_, err = ibmcloud.NewIsLbPoolMember(ctx, "lb-pool-member-2", &ibmcloud.IsLbPoolMemberArgs{
    Lb:            loadBalancer.ID(),
    Pool:          lbPool.ID(),
    Port:          pulumi.Int(80),
    TargetAddress: powerVSInstance2.PiNetworks.Index(pulumi.Int(0)).IpAddress(),
    Weight:        pulumi.Int(50),
})
```

## Troubleshooting

### Power VS Instance Creation Fails

**Error**: "No capacity available in zone"
```
Solution: Try a different Power VS zone:
pulumi config set powerVsZone dal13
```

**Error**: "Image not found"
```bash
# List available images
ibmcloud pi images --json

# Use correct image ID
PiImageId: pulumi.String("actual-image-id-from-list")
```

### Load Balancer Not Accessible

1. **Check security group rules**:
   ```bash
   pulumi stack output securityGroupId
   # Verify HTTP/HTTPS rules exist
   ```

2. **Check health status**:
   ```bash
   # Load balancer health checks must pass
   # Ensure your application responds on port 80
   ```

3. **Check Power VS instance status**:
   ```bash
   pulumi stack output powerVsInstanceStatus
   # Should be: ACTIVE
   ```

### SSH Connection Fails

1. **Verify SSH key was added**:
   ```bash
   # Check SSH key in Power VS workspace
   ibmcloud pi keys
   ```

2. **Check Power VS instance is running**:
   ```bash
   pulumi stack output powerVsInstanceStatus
   ```

3. **Add SSH listener to load balancer** (not included by default):
   ```go
   _, err = ibmcloud.NewIsLbListener(ctx, "lb-listener-ssh", &ibmcloud.IsLbListenerArgs{
       Lb:          loadBalancer.ID(),
       DefaultPool: lbPool.ID(),
       Port:        pulumi.Int(22),
       Protocol:    pulumi.String("tcp"),
   })
   ```

## Cleanup

To destroy all resources:

```bash
pulumi destroy
```

> **Warning**: This will permanently delete:
> - Power VS instance and all its data
> - Power VS network
> - VPC and all networking components
> - Load balancer

## Next Steps

- [ ] Add SSL/TLS certificate for HTTPS
- [ ] Configure auto-scaling with multiple instances
- [ ] Add monitoring and alerting
- [ ] Set up backup and disaster recovery
- [ ] Implement infrastructure as code CI/CD
- [ ] Add Cloud Object Storage for backups
- [ ] Configure VPN for private access
- [ ] Add IBM Cloud Databases

## Resources

- [IBM Cloud Power Virtual Server Documentation](https://cloud.ibm.com/docs/power-iaas)
- [IBM Cloud VPC Documentation](https://cloud.ibm.com/docs/vpc)
- [IBM Cloud Load Balancer Documentation](https://cloud.ibm.com/docs/vpc?topic=vpc-nlb-vs-elb)
- [Pulumi IBM Cloud Provider](https://github.com/mapt-oss/pulumi-ibmcloud)
- [IBM Cloud CLI Reference](https://cloud.ibm.com/docs/cli)

## License

Apache 2.0
