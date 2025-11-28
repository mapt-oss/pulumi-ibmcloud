# GitLab Runner Multi-Architecture Example

This example demonstrates how to deploy GitLab runners on IBM Cloud for multiple architectures:
- **s390x** (IBM Z/LinuxONE)
- **ppc64** (IBM Power Systems)

Both runners use Fedora and connect to GitLab.com for CI/CD workloads.

## Architecture Overview

```
IBM Cloud VPC: gitlab-runner-vpc (10.240.0.0/16)
│
├── Subnet 1 (us-east-1): 10.240.0.0/24
│   ├── Public Gateway (for outbound to GitLab.com)
│   └── s390x Instance (IBM Z)
│       ├── Fedora OS
│       ├── Docker
│       └── GitLab Runner
│
└── Subnet 2 (us-south-2): 10.240.1.0/24
    ├── Public Gateway (for outbound to GitLab.com)
    └── ppc64 Instance (IBM Power)
        ├── Fedora OS
        ├── Docker
        └── GitLab Runner
```

## Why Separate Subnets?

- **Subnets are zone-specific**: Each subnet exists in exactly one availability zone
- **Different architectures require different zones**:
  - s390x is available in zones like `us-east-1`, `us-east-2`
  - ppc64 is available in zones like `us-south-2`, `us-south-3`
- **Both subnets are in the same VPC**: They can communicate over the private network

## Prerequisites

1. **IBM Cloud Account** with appropriate permissions
2. **IBM Cloud API Key**: Set as environment variable
   ```bash
   export IC_API_KEY="your-api-key"
   ```
3. **SSH Public Key**: Update in `main.go` at line ~136
4. **Fedora Image IDs**: Find the correct image IDs for each architecture

### Finding Fedora Image IDs

```bash
# For s390x (IBM Z)
ibmcloud is images --visibility public | grep -i fedora | grep s390x

# For ppc64 (IBM Power)
ibmcloud is images --visibility public | grep -i fedora | grep ppc
```

Update the image IDs in `main.go`:
- Line ~245: `imageS390x` for s390x
- Line ~271: `imagePpc64` for ppc64

## Deployment

### Step 1: Update Configuration

Edit `main.go` and replace:
1. **SSH Public Key** (line ~136):
   ```go
   PublicKey: pulumi.String("ssh-rsa AAAAB3... your-actual-key"),
   ```

2. **s390x Image ID** (line ~245):
   ```go
   imageS390x := pulumi.String("r014-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
   ```

3. **ppc64 Image ID** (line ~271):
   ```go
   imagePpc64 := pulumi.String("r006-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
   ```

### Step 2: Initialize Pulumi

```bash
cd examples/gitlab-runner-multi-arch
pulumi stack init dev
pulumi config set ibmcloud:region us-south
```

### Step 3: Deploy

```bash
pulumi up
```

Review the resources to be created and confirm.

### Step 4: Get Connection Information

```bash
pulumi stack output
```

Example output:
```
Current stack outputs (10):
    OUTPUT                VALUE
    ppc64_floatingIP      169.xx.xx.xx
    ppc64_instanceId      0717_xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    ppc64_privateIP       10.240.1.4
    ppc64_sshCommand      ssh fedora@169.xx.xx.xx
    resourceGroupId       xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    s390x_floatingIP      150.xx.xx.xx
    s390x_instanceId      0717_xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    s390x_privateIP       10.240.0.4
    s390x_sshCommand      ssh fedora@150.xx.xx.xx
    vpcId                 r006-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

## Register GitLab Runners

### Step 1: Get Your GitLab Registration Token

1. Go to your GitLab project/group
2. Navigate to **Settings > CI/CD > Runners**
3. Expand the **Runners** section
4. Copy the registration token (starts with `glrt-`)

### Step 2: Register s390x Runner

```bash
# SSH to the s390x instance
ssh fedora@$(pulumi stack output s390x_floatingIP)

# Register the runner with s390x tag
sudo /home/fedora/register-runner.sh glrt-YOUR_TOKEN s390x,docker,ibm-z

# Verify installation
sudo gitlab-runner list
uname -m  # Should show s390x
```

### Step 3: Register ppc64 Runner

```bash
# SSH to the ppc64 instance
ssh fedora@$(pulumi stack output ppc64_floatingIP)

# Register the runner with ppc64 tag
sudo /home/fedora/register-runner.sh glrt-YOUR_TOKEN ppc64,docker,ibm-power

# Verify installation
sudo gitlab-runner list
uname -m  # Should show ppc64le
```

## Using the Runners in GitLab CI/CD

### Example `.gitlab-ci.yml`

```yaml
# Build on s390x
build-s390x:
  tags:
    - s390x
  script:
    - echo "Building on s390x architecture"
    - uname -m
    - docker build -t myapp:s390x .

# Build on ppc64
build-ppc64:
  tags:
    - ppc64
  script:
    - echo "Building on ppc64 architecture"
    - uname -m
    - docker build -t myapp:ppc64 .

# Multi-arch build
build-multiarch:
  parallel:
    matrix:
      - ARCH: [s390x, ppc64]
  tags:
    - ${ARCH}
  script:
    - echo "Building on ${ARCH}"
    - docker build -t myapp:${ARCH} .
```

## Network Details

### Outbound Connectivity (via Public Gateway)

Both runners can reach:
- ✅ gitlab.com (port 443) - Job polling and artifact upload
- ✅ Docker Hub - Container image pulls
- ✅ Package repositories (dnf, pip, npm, etc.)
- ✅ Any internet resources

### Inbound Connectivity

- ❌ No inbound internet access by default (secure)
- ✅ SSH access via Floating IP (for administration)
- ✅ Private network communication between runners within VPC

### Security

The security group allows:
- **Inbound**: SSH (port 22) from anywhere (⚠️ consider restricting to your IP)
- **Outbound**: All traffic (required for GitLab.com and dependencies)

To restrict SSH access, update `main.go` line ~145:
```go
Remote: pulumi.String("YOUR_IP_ADDRESS/32"), // Your specific IP
```

## Cost Estimation

Approximate monthly costs (us-south/us-east regions):

| Resource | Quantity | Cost/Month |
|----------|----------|------------|
| Virtual Server Instance (bz2-2x8) s390x | 1 | ~$180 |
| Virtual Server Instance (bp2-2x8) ppc64 | 1 | ~$160 |
| Public Gateway | 2 | ~$130 |
| Floating IP | 2 | ~$10 |
| VPC (free tier) | 1 | $0 |
| **Total** | | **~$480/month** |

**Note**: Costs vary by region and resource usage. Check [IBM Cloud Pricing](https://cloud.ibm.com/vpc-ext/provision/vs) for current rates.

## Cleanup

To destroy all resources:

```bash
pulumi destroy
```

Confirm the destruction when prompted.

## Troubleshooting

### Runner Cannot Connect to GitLab.com

1. **Check Public Gateway**: Ensure the subnet has a public gateway attached
   ```bash
   ibmcloud is subnets
   ibmcloud is subnet <subnet-id>
   ```

2. **Check Security Group**: Verify outbound rules allow HTTPS
   ```bash
   ibmcloud is security-group <sg-id>
   ```

3. **Test from instance**:
   ```bash
   ssh fedora@<floating-ip>
   curl -I https://gitlab.com
   ```

### SSH Connection Fails

1. **Verify Floating IP**: Check that the floating IP is attached
   ```bash
   pulumi stack output s390x_floatingIP
   ```

2. **Check SSH key**: Ensure your private key matches the public key in the code
   ```bash
   ssh -i ~/.ssh/your_key fedora@<floating-ip>
   ```

3. **Security Group**: Verify SSH rule allows your IP

### Docker Permission Denied

The `fedora` user should be in the `docker` group automatically. If not:
```bash
sudo usermod -aG docker fedora
# Log out and back in
```

### GitLab Runner Not Starting

Check the service status:
```bash
sudo systemctl status gitlab-runner
sudo journalctl -u gitlab-runner -f
```

## Architecture-Specific Notes

### s390x (IBM Z)

- **Zones**: us-east-1, us-east-2, eu-de-1, eu-de-2
- **Profiles**: bz2-* (Balanced Z2)
- **Use cases**: Enterprise workloads, cryptography, compliance
- **Container platforms**: Docker, Podman (native support)

### ppc64 (IBM Power)

- **Zones**: us-south-1, us-south-2, us-south-3
- **Profiles**: bp2-* (Balanced Power)
- **Use cases**: AI/ML workloads, databases, analytics
- **Container platforms**: Docker, Podman (native support)

## Additional Resources

- [IBM Cloud VPC Documentation](https://cloud.ibm.com/docs/vpc)
- [GitLab Runner Documentation](https://docs.gitlab.com/runner/)
- [Docker on Multi-Architecture](https://docs.docker.com/build/building/multi-platform/)
- [IBM Z Container Platform](https://www.ibm.com/z/resources/containers)
- [IBM Power Systems Containers](https://www.ibm.com/power/solutions/containers)

## License

This example is provided under the Apache 2.0 license.
