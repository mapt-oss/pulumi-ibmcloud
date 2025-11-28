# Pulumi IBM Cloud Provider Examples

This directory contains examples demonstrating various IBM Cloud resources using Pulumi.

## Available Examples

### Power Virtual Server

| Example | Description | Complexity | Monthly Cost |
|---------|-------------|------------|--------------|
| [power-vs-basic](./power-vs-basic) | Basic Power VS instance with private networking | Beginner | ~$20-30 |
| [power-vs-with-lb](./power-vs-with-lb) | Power VS with VPC Load Balancer for public access | Intermediate | ~$95-135 |
| [power11-public-access](./power11-public-access) | **Power11 (s1022)** instance with full public access, ALB, and multi-language examples | Advanced | ~$235-310 |

### CI/CD & DevOps

| Example | Description | Complexity | Monthly Cost |
|---------|-------------|------------|--------------|
| [gitlab-runner-multi-arch](./gitlab-runner-multi-arch) | Multi-architecture GitLab runners (s390x + ppc64) on VPC | Intermediate | ~$480 |

### Coming Soon

- **vpc-kubernetes**: VPC with IBM Kubernetes Service cluster
- **vpc-web-app**: Complete web application with VPC, VSI, and databases
- **cos-cdn**: Cloud Object Storage with CDN
- **multizone-ha**: High-availability multi-zone deployment
- **hybrid-cloud**: Hybrid cloud with Direct Link and VPN

## Prerequisites

All examples require:

1. **IBM Cloud Account** with appropriate permissions
2. **IBM Cloud API Key**:
   ```bash
   export IC_API_KEY="your-api-key"
   ```
3. **Pulumi CLI**: [Installation guide](https://www.pulumi.com/docs/install/)
4. **Go**: Version 1.23 or later

## Quick Start

```bash
# 1. Choose an example
cd power-vs-basic

# 2. Initialize
pulumi stack init dev

# 3. Configure
pulumi config set ibmcloud:ibmcloudApiKey --secret
pulumi config set sshPublicKey "$(cat ~/.ssh/id_rsa.pub)"

# 4. Deploy
pulumi up

# 5. Clean up when done
pulumi destroy
```

## Configuration

### Global IBM Cloud Settings

```bash
# Set API key (required)
pulumi config set ibmcloud:ibmcloudApiKey --secret

# Set default region (optional)
pulumi config set ibmcloud:region us-south

# Set default resource group (optional)
pulumi config set ibmcloud:resourceGroup default
```

### Example-Specific Settings

Each example has its own configuration options. See the README in each example directory for details.

## Example Structure

Each example follows this structure:

```
example-name/
├── main.go           # Main Pulumi program
├── Pulumi.yaml       # Project definition
├── Pulumi.dev.yaml   # Sample configuration
├── go.mod            # Go module file
└── README.md         # Example documentation
```

## Cost Estimation

All examples include estimated monthly costs. Actual costs may vary based on:
- Region selected
- Resource utilization
- Data transfer
- Additional services

Use the [IBM Cloud Cost Estimator](https://cloud.ibm.com/estimator) for precise estimates.

## Getting Help

- **Pulumi Docs**: https://www.pulumi.com/docs/
- **IBM Cloud Docs**: https://cloud.ibm.com/docs
- **Provider Issues**: https://github.com/mapt-oss/pulumi-ibmcloud/issues
- **Pulumi Community**: https://slack.pulumi.com/

## Contributing

Have an example to share? Submit a PR!

Guidelines:
- Follow existing structure
- Include comprehensive README
- Add cost estimates
- Test thoroughly
- Use descriptive variable names
- Add comments for complex logic

## License

Apache 2.0
