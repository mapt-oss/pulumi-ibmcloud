# IBM Cloud Power11 Virtual Server with Public Access

A comprehensive example demonstrating how to provision an IBM Power11 virtual server with full public internet access using VPC Application Load Balancer.

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              IBM Cloud                                       │
│                                                                               │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                        Resource Group                                   │ │
│  │                                                                          │ │
│  │  ┌──────────────────────────┐     ┌────────────────────────────────┐  │ │
│  │  │   Power VS Workspace     │     │        VPC Network              │  │ │
│  │  │   (Power11 Capable)      │     │        (us-south)               │  │ │
│  │  │                          │     │                                  │  │ │
│  │  │  ┌────────────────────┐  │     │  ┌──────────────────────────┐  │  │ │
│  │  │  │  Private Network   │  │     │  │   Public Gateway         │  │  │ │
│  │  │  │  192.168.50.0/24   │  │     │  │   (NAT for outbound)     │  │  │ │
│  │  │  └─────────┬──────────┘  │     │  └──────────────────────────┘  │  │ │
│  │  │            │              │     │              │                  │  │ │
│  │  │  ┌─────────▼──────────┐  │     │  ┌───────────▼──────────────┐  │  │ │
│  │  │  │   Power11 Instance │  │     │  │  Subnet: 10.240.64.0/24  │  │  │ │
│  │  │  │   ────────────────  │  │     │  │                          │  │  │ │
│  │  │  │   System: s1022    │  │     │  │  ┌────────────────────┐  │  │  │ │
│  │  │  │   (Power10/11)     │◄─┼─────┼──┤  │  Load Balancer     │  │  │  │ │
│  │  │  │   ────────────────  │  │     │  │  │  (Public ALB)      │  │  │  │ │
│  │  │  │   CPU: 2 cores     │  │     │  │  │                    │  │  │  │ │
│  │  │  │   RAM: 16 GB       │  │     │  │  │  Listeners:        │  │  │  │ │
│  │  │  │   Storage: Tier1   │  │     │  │  │  - HTTP :80        │  │  │  │ │
│  │  │  │   (NVMe)           │  │     │  │  │  - HTTPS:443       │  │  │  │ │
│  │  │  │   OS: RHEL 9.2     │  │     │  │  │  - SSH  :2222      │  │  │  │ │
│  │  │  │   ────────────────  │  │     │  │  └────────┬───────────┘  │  │  │ │
│  │  │  │   IP: 192.168.50.x │  │     │  │           │              │  │  │ │
│  │  │  └────────────────────┘  │     │  │  ┌────────▼───────────┐  │  │  │ │
│  │  │                          │     │  │  │  Security Group    │  │  │  │ │
│  │  └──────────────────────────┘     │  │  │  - HTTP  (80)      │  │  │  │ │
│  │                                    │  │  │  - HTTPS (443)     │  │  │  │ │
│  │                                    │  │  │  - SSH   (2222)    │  │  │  │ │
│  │                                    │  │  └────────────────────┘  │  │  │ │
│  │                                    │  └─────────┬────────────────┘  │  │ │
│  │                                    └────────────┼────────────────────┘  │ │
│  └─────────────────────────────────────────────────┼───────────────────────┘ │
│                                                     │                         │
└─────────────────────────────────────────────────────┼─────────────────────────┘
                                                      │
                                              ┌───────▼────────┐
                                              │   Internet     │
                                              │   Users        │
                                              │                │
                                              │  HTTP/HTTPS    │
                                              │  SSH Access    │
                                              └────────────────┘
```

## 📋 What This Creates

### Power Virtual Server Resources
- **Power VS Workspace**: Service instance in a Power11-capable zone
- **SSH Key**: For secure instance access
- **Private Network**: VLAN network (192.168.50.0/24) with DNS
- **Power11 Instance**: Virtual machine with:
  - **System Type**: s1022 (Power10/11 architecture)
  - **CPU**: 2 dedicated cores (configurable)
  - **Memory**: 16 GB RAM (configurable)
  - **Storage**: Tier1 NVMe (highest performance)
  - **OS**: RHEL 9.2 (or AIX/IBM i)

### VPC Network Resources
- **VPC**: Virtual Private Cloud for load balancer
- **Subnet**: 10.240.64.0/24 for load balancer placement
- **Public Gateway**: Provides NAT for outbound internet access
- **Security Group**: Firewall rules for HTTP/HTTPS/SSH traffic

### Load Balancer
- **Application Load Balancer**: Public-facing ALB
- **Backend Pool**: Routes traffic to Power11 instance
- **HTTP Listener**: Port 80 (HTTP traffic)
- **HTTPS Listener**: Port 443 (optional, requires certificate)
- **SSH Listener**: Port 2222 (TCP passthrough for SSH)
- **Health Checks**: Automated health monitoring

## 🚀 Quick Start

### Prerequisites

1. **IBM Cloud Account** with:
   - IBM Cloud API key
   - Appropriate IAM permissions for Power VS and VPC
   - Access to Power11-capable zones

2. **Tools**:
   - [Pulumi CLI](https://www.pulumi.com/docs/install/) v3.0+
   - Language runtime (Node.js 18+, Python 3.9+, or Go 1.23+)
   - SSH key pair

3. **Power11 Availability**:
   Power11 systems (s1022, e1080) are available in select zones:
   - **us-south** (Dallas)
   - **us-east** (Washington DC)
   - **eu-de** (Frankfurt)
   - **eu-gb** (London)
   - **jp-tok** (Tokyo)
   - **au-syd** (Sydney)

### Choose Your Language

This example is available in three languages:

| Language | Directory | Runtime |
|----------|-----------|---------|
| TypeScript | `typescript/` | Node.js 18+ |
| Python | `python/` | Python 3.9+ |
| Go | `go/` | Go 1.23+ |

## 📦 Deployment

### TypeScript Example

```bash
cd typescript

# Install dependencies
npm install

# Initialize Pulumi
pulumi stack init dev

# Configure
export IC_API_KEY="your-ibm-cloud-api-key"
pulumi config set ibmcloud:ibmcloudApiKey --secret
pulumi config set sshPublicKey "$(cat ~/.ssh/id_rsa.pub)"
pulumi config set region us-south
pulumi config set zone us-south-1
pulumi config set powerVsZone us-south
pulumi config set instanceName my-power11-vm

# Deploy
pulumi up

# Get access information
pulumi stack output httpUrl
pulumi stack output sshCommand
```

### Python Example

```bash
cd python

# Create virtual environment
python3 -m venv venv
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Initialize Pulumi
pulumi stack init dev

# Configure
export IC_API_KEY="your-ibm-cloud-api-key"
pulumi config set ibmcloud:ibmcloudApiKey --secret
pulumi config set sshPublicKey "$(cat ~/.ssh/id_rsa.pub)"
pulumi config set region us-south
pulumi config set zone us-south-1
pulumi config set powerVsZone us-south

# Deploy
pulumi up

# Get access information
pulumi stack output httpUrl
pulumi stack output sshCommand
```

### Go Example

```bash
cd go

# Initialize Go module
go mod tidy

# Initialize Pulumi
pulumi stack init dev

# Configure
export IC_API_KEY="your-ibm-cloud-api-key"
pulumi config set ibmcloud:ibmcloudApiKey --secret
pulumi config set sshPublicKey "$(cat ~/.ssh/id_rsa.pub)"
pulumi config set region us-south
pulumi config set zone us-south-1
pulumi config set powerVsZone us-south

# Deploy
pulumi up

# Get access information
pulumi stack output httpUrl
pulumi stack output sshCommand
```

## ⏱️ Deployment Timeline

Total deployment time: **20-30 minutes**

| Component | Time | Notes |
|-----------|------|-------|
| Resource Group | ~30 seconds | Fast |
| Power VS Workspace | ~2 minutes | Service provisioning |
| SSH Key | ~10 seconds | Fast |
| Private Network | ~1 minute | VLAN creation |
| **Power11 Instance** | **15-20 minutes** | Most time-consuming |
| VPC | ~30 seconds | Fast |
| Subnet | ~30 seconds | Fast |
| Public Gateway | ~2 minutes | IP allocation |
| Security Group + Rules | ~30 seconds | Fast |
| Load Balancer | ~5-8 minutes | ALB provisioning |
| Backend Pool + Member | ~1 minute | Configuration |
| Listeners | ~30 seconds | Fast |

> **Note**: Power11 instance provisioning is the longest step. IBM Cloud needs to allocate dedicated Power processors and provision the OS.

## 🔐 Access Your Power11 Instance

After deployment completes:

### Get Access Information

```bash
# View deployment summary
pulumi stack output deploymentSummary

# Get specific outputs
pulumi stack output httpUrl          # HTTP access URL
pulumi stack output httpsUrl         # HTTPS access URL (if configured)
pulumi stack output sshCommand       # SSH command
pulumi stack output loadBalancerHostname  # Load balancer hostname
```

### SSH Access

```bash
# SSH via load balancer (port 2222)
ssh -p 2222 root@<load-balancer-hostname>

# Or use the exported command
$(pulumi stack output sshCommand)
```

### HTTP Access

```bash
# Test HTTP access
curl http://<load-balancer-hostname>

# Open in browser
open $(pulumi stack output httpUrl)
```

### Direct Instance Access

The Power11 instance's internal IP:

```bash
pulumi stack output power11InstanceInternalIp
# Output: 192.168.50.x
```

> **Note**: Direct access to the internal IP requires VPN or Direct Link connection to the Power VS workspace.

## 🔧 Configuration Options

### Power11 System Specifications

Modify in your code to adjust resources:

**TypeScript**:
```typescript
const power11Instance = new ibmcloud.PiInstance("power11-instance", {
    piMemory: 32,              // 32 GB RAM
    piProcessors: 4,           // 4 cores
    piProcType: "dedicated",   // or "shared" for cost savings
    piSysType: "s1022",        // s1022 or e1080
    piStorageType: "tier1",    // tier1 (NVMe), tier3 (SSD)
    // ...
});
```

**Python**:
```python
power11_instance = ibmcloud.PiInstance(
    "power11-instance",
    pi_memory=32.0,           # 32 GB RAM
    pi_processors=4.0,        # 4 cores
    pi_proc_type="dedicated", # or "shared"
    pi_sys_type="s1022",      # s1022 or e1080
    pi_storage_type="tier1",  # tier1 (NVMe), tier3 (SSD)
    # ...
)
```

**Go**:
```go
power11Instance, err := ibmcloud.NewPiInstance(ctx, "power11-instance", &ibmcloud.PiInstanceArgs{
    PiMemory:      pulumi.Float64(32),         // 32 GB RAM
    PiProcessors:  pulumi.Float64(4),          // 4 cores
    PiProcType:    pulumi.String("dedicated"), // or "shared"
    PiSysType:     pulumi.String("s1022"),     // s1022 or e1080
    PiStorageType: pulumi.String("tier1"),     // tier1 (NVMe), tier3 (SSD)
    // ...
})
```

### Available Power11 System Types

| System Type | Architecture | Use Case | CPU Options |
|-------------|--------------|----------|-------------|
| `s1022` | Power10/11 | General purpose, scale-out | Shared or Dedicated |
| `e1080` | Power10/11 | Enterprise, scale-up | Dedicated only |

### Available Operating Systems

Check available images for your zone:

```bash
# Using IBM Cloud CLI
ibmcloud pi images --json

# Common Power11 images
# - rhel-9-2 (Red Hat Enterprise Linux 9.2)
# - rhel-8-6 (Red Hat Enterprise Linux 8.6)
# - aix-7-3 (AIX 7.3)
# - ibm-i-7-5 (IBM i 7.5)
# - sles-15-sp5 (SUSE Linux Enterprise Server 15 SP5)
```

Update the `piImageId` in your code:

```typescript
piImageId: "rhel-9-2"  // or "aix-7-3", "ibm-i-7-5", etc.
```

### Storage Tiers

| Tier | Type | IOPS | Use Case | Cost |
|------|------|------|----------|------|
| tier1 | NVMe | Highest | Production, databases | $$$ |
| tier3 | SSD | Medium | General purpose | $$ |

### Processor Types

| Type | Description | Cost | Use Case |
|------|-------------|------|----------|
| dedicated | Full core allocation | Higher | Production workloads |
| shared | Shared processor pool | Lower | Dev/test, cost-sensitive |

## 💰 Cost Estimate

Estimated monthly costs (US South region):

| Component | Specification | Monthly Cost (USD) |
|-----------|--------------|-------------------|
| Power11 Instance | 2 cores, 16GB, dedicated | $120-160 |
| Storage | 100GB Tier1 NVMe | $15-20 |
| Power VS Network | Private VLAN | $0 |
| VPC | VPC instance | $0 |
| Public Gateway | NAT gateway | $40-50 |
| Load Balancer | Public ALB | $60-80 |
| Data Transfer | First 5GB free | Variable |
| **Total** | | **~$235-310/month** |

> **Cost Savings Tips**:
> - Use **shared** processors: Save 40-60%
> - Use **tier3** storage: Save 30-40%
> - Use smaller instance: 0.5 cores, 4GB RAM saves ~60%
> - Deploy in dev/test for short periods only

**Dev/Test Configuration** (0.5 cores, 4GB, shared, tier3):
- Estimated cost: **$50-70/month**

## 🎯 Use Cases

### 1. **AIX Application Hosting**
- Run legacy AIX applications with modern internet access
- Migrate on-premises AIX workloads to cloud
- Provide web interfaces to AIX applications

### 2. **IBM i Modernization**
- Expose IBM i applications via REST APIs
- Integrate IBM i with cloud-native services
- Provide public access to IBM i web services

### 3. **SAP on Power**
- Host SAP applications on Power11 architecture
- High-performance SAP HANA deployments
- Public-facing SAP Fiori applications

### 4. **Oracle Database on Power**
- Run Oracle databases with Power11 performance
- Public database endpoints for applications
- High-availability Oracle RAC clusters

### 5. **Development & Testing**
- Power architecture development environments
- CI/CD pipelines for Power applications
- Testing new Power11 features and capabilities

## 🔒 Security Considerations

### Network Security

1. **Security Groups**: Restrict inbound traffic
   ```typescript
   // Example: Restrict SSH to specific IP
   remote: "203.0.113.0/24"  // Your office IP range
   ```

2. **Private Networking**: Power11 instance is on private network
   - Only accessible via load balancer
   - Cannot be directly reached from internet

3. **Public Gateway**: Only provides outbound NAT
   - No inbound access through gateway
   - Instances cannot receive unsolicited inbound traffic

### Authentication

1. **SSH Key Authentication**: Password auth disabled by default
2. **Strong SSH Keys**: Use ED25519 or RSA 4096-bit keys
3. **Rotate Keys**: Regularly rotate SSH keys

### Additional Security Measures

```bash
# After deployment, harden the instance

# 1. Update all packages
ssh -p 2222 root@<hostname>
dnf update -y

# 2. Configure firewall
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --reload

# 3. Install fail2ban
dnf install -y fail2ban
systemctl enable --now fail2ban

# 4. Configure SELinux (if using RHEL)
getenforce  # Should return "Enforcing"
```

## 📊 Monitoring & Health Checks

### Load Balancer Health Checks

The load balancer automatically monitors instance health:

- **Protocol**: HTTP
- **URL**: `/`
- **Interval**: 10 seconds
- **Timeout**: 5 seconds
- **Retries**: 3

To customize health checks, modify the backend pool configuration:

```typescript
healthMonitorUrl: "/health",      // Custom health endpoint
healthMonitorPort: 8080,          // Custom port
healthType: "https",              // Use HTTPS
```

### IBM Cloud Monitoring

View metrics in IBM Cloud console:

1. Navigate to **Observability** > **Monitoring**
2. Select your Power VS workspace
3. View CPU, memory, network metrics

## 🔧 Troubleshooting

### Issue: Power11 Instance Provisioning Fails

**Error**: "No capacity available for s1022"

**Solution**:
```bash
# Try alternative Power11 system type
pulumi config set powerVsZone us-east

# Or use e1080 instead of s1022
# Modify piSysType: "e1080" in your code
```

### Issue: Cannot Find Power11 Images

**Error**: "Image rhel-9-2 not found"

**Solution**:
```bash
# List available images in your zone
ibmcloud pi workspace target <workspace-guid>
ibmcloud pi images

# Use exact image ID from the list
# Update piImageId in your code with the actual image ID
```

### Issue: Load Balancer Health Checks Failing

**Symptom**: Load balancer shows "unhealthy" status

**Solution**:
```bash
# SSH into instance
ssh -p 2222 root@<hostname>

# Start a simple web server for testing
python3 -m http.server 80

# Or install nginx
dnf install -y nginx
systemctl enable --now nginx
firewall-cmd --permanent --add-service=http
firewall-cmd --reload
```

### Issue: SSH Connection Refused

**Check**:
1. Instance status: `pulumi stack output power11InstanceStatus` (should be "ACTIVE")
2. Load balancer status: Check IBM Cloud console
3. Security group rules: Verify port 2222 is allowed

### Issue: Deployment Timeout

**Error**: "Timeout waiting for instance to be ready"

**Solution**:
```bash
# Power11 instances can take 20-30 minutes
# Increase timeout in your code:

# TypeScript/Python
customTimeouts: { create: "45m" }

# Go
pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "45m"})
```

## 🧹 Cleanup

To destroy all resources:

```bash
pulumi destroy
```

> **Warning**: This permanently deletes:
> - Power11 instance and all its data
> - All networks and networking components
> - Load balancer
> - VPC resources
>
> **Ensure you have backups before destroying!**

## 📚 Next Steps

### Enhancements

1. **Add HTTPS Support**
   - Obtain SSL certificate from IBM Certificate Manager
   - Uncomment HTTPS listener in code
   - Configure certificate CRN

2. **Add Multiple Instances**
   - Create additional Power11 instances
   - Add to load balancer pool for high availability
   - Configure session affinity if needed

3. **Add Cloud Object Storage**
   - Create COS instance for backups
   - Configure automated backups
   - Set up lifecycle policies

4. **Add Monitoring**
   - Deploy IBM Cloud Monitoring agent
   - Create custom dashboards
   - Configure alerts

5. **Add VPN Access**
   - Deploy VPN gateway for private access
   - Alternative to load balancer for management
   - More secure for administrative tasks

6. **Add Direct Link**
   - For production deployments
   - Direct connection from on-premises
   - Lower latency, higher bandwidth

## 📖 Additional Resources

### IBM Cloud Documentation
- [Power Virtual Server Documentation](https://cloud.ibm.com/docs/power-iaas)
- [Power10/11 Specifications](https://www.ibm.com/products/power-systems)
- [VPC Load Balancer Documentation](https://cloud.ibm.com/docs/vpc?topic=vpc-nlb-vs-elb)
- [IBM Cloud Regions and Zones](https://cloud.ibm.com/docs/overview?topic=overview-locations)

### Pulumi Documentation
- [Pulumi IBM Cloud Provider](https://github.com/mapt-oss/pulumi-ibmcloud)
- [Pulumi Getting Started](https://www.pulumi.com/docs/get-started/)
- [Pulumi Best Practices](https://www.pulumi.com/docs/guides/best-practices/)

### Power Systems Resources
- [IBM Power Systems Community](https://community.ibm.com/community/user/power)
- [AIX Documentation](https://www.ibm.com/docs/en/aix)
- [IBM i Documentation](https://www.ibm.com/docs/en/i)

## 🤝 Contributing

Found an issue or want to improve this example?

1. Open an issue at [github.com/mapt-oss/pulumi-ibmcloud/issues](https://github.com/mapt-oss/pulumi-ibmcloud/issues)
2. Submit a pull request
3. Share your Power11 use cases!

## 📄 License

Apache 2.0 - See [LICENSE](../../LICENSE) for details.

---

**Built with** ❤️ **by the mapt-oss community**

For questions or support, please open an issue on GitHub.
