# IBM Cloud Provider Examples

This directory contains example Pulumi programs demonstrating how to use the IBM Cloud provider.

## Prerequisites

### 1. Install Pulumi
```bash
curl -fsSL https://get.pulumi.com | sh
```

### 2. Install the Provider Plugin
```bash
# After a release is created
pulumi plugin install resource ibmcloud v0.0.1 --server github://github.com/mapt-oss

# Or build locally
cd /path/to/pulumi-ibmcloud
make provider
cp bin/pulumi-resource-ibmcloud ~/.pulumi/bin/
```

### 3. IBM Cloud Credentials

Get your IBM Cloud API key from: https://cloud.ibm.com/iam/apikeys

```bash
export IC_API_KEY="your-ibm-cloud-api-key"
export IC_REGION="us-south"  # or your preferred region
```

## Available Examples

### Basic Examples (Getting Started)

#### 1. `basic-go` - Resource Group (Go)
**What it does**: Creates an IBM Cloud Resource Group

**Resources**: 1 resource group

**Run it**:
```bash
cd basic-go
go mod init example
go get github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud
pulumi stack init dev
pulumi up
```

#### 2. `basic-ts` - Resource Group + VPC (TypeScript)
**What it does**: Creates a Resource Group and VPC

**Resources**: 1 resource group, 1 VPC

**Run it**:
```bash
cd basic-ts
npm install
pulumi stack init dev
pulumi up
```

#### 3. `basic-py` - Resource Group + Cloud Object Storage (Python)
**What it does**: Creates a Resource Group and COS instance

**Resources**: 1 resource group, 1 COS instance

**Run it**:
```bash
cd basic-py
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
pulumi stack init dev
pulumi up
```

### Advanced Examples

#### 4. `vpc-with-subnet` - Complete VPC Setup
**What it does**: Creates a full VPC environment with subnets, security groups, and gateways

**Resources**:
- 1 Resource Group
- 1 VPC
- 3 Subnets (across 3 availability zones)
- 1 Public Gateway
- 1 Security Group with rules
- 1 Network ACL

**Languages**: Go, TypeScript, Python

#### 5. `kubernetes-cluster` - IKS Cluster
**What it does**: Deploys a Kubernetes cluster on IBM Cloud

**Resources**:
- 1 Resource Group
- 1 VPC
- 3 Subnets
- 1 IKS cluster with worker pool

**Languages**: Go, TypeScript

#### 6. `object-storage-bucket` - Cloud Object Storage
**What it does**: Creates a COS instance with multiple buckets

**Resources**:
- 1 Resource Group
- 1 COS instance
- 3 Buckets (standard, cold, vault)
- Bucket lifecycle policies

**Languages**: Python, TypeScript

#### 7. `databases` - Cloud Databases
**What it does**: Deploys PostgreSQL and Redis databases

**Resources**:
- 1 Resource Group
- 1 PostgreSQL database
- 1 Redis database
- Service credentials

**Languages**: Go, Python

### Real-World Examples

#### 8. `three-tier-app` - Three-Tier Application
**What it does**: Complete infrastructure for a web application

**Resources**:
- VPC with subnets
- Load balancer
- IKS cluster
- PostgreSQL database
- COS bucket for static assets
- IAM policies

**Languages**: Go, TypeScript

#### 9. `serverless-api` - Serverless API
**What it does**: Serverless API with Functions and API Gateway

**Resources**:
- Cloud Functions
- API Gateway
- COS buckets
- Cloudant database

**Languages**: TypeScript, Python

## Example Structure

Each example follows this structure:
```
example-name/
├── README.md           # Detailed description and instructions
├── Pulumi.yaml        # Pulumi project configuration
├── main.go            # Go example (if applicable)
├── index.ts           # TypeScript example (if applicable)
├── __main__.py        # Python example (if applicable)
├── package.json       # TypeScript dependencies
├── requirements.txt   # Python dependencies
├── go.mod            # Go dependencies
└── .gitignore        # Ignored files
```

## Creating Your Own Example

### 1. Initialize a new Pulumi project
```bash
mkdir my-ibmcloud-example
cd my-ibmcloud-example
pulumi new go  # or typescript, python
```

### 2. Import the provider
**Go**:
```go
import "github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud"
```

**TypeScript**:
```typescript
import * as ibmcloud from "@pulumi/ibmcloud";
```

**Python**:
```python
import pulumi_ibmcloud as ibmcloud
```

### 3. Create resources
```go
rg, err := ibmcloud.NewResourceGroup(ctx, "my-rg", &ibmcloud.ResourceGroupArgs{
    Name: pulumi.String("my-resource-group"),
})
```

### 4. Run
```bash
pulumi up
```

## Testing Examples

Examples include integration tests that can be run with:

```bash
cd examples
go test -v -timeout 2h
```

Individual example tests:
```bash
go test -v -run TestGoExample
go test -v -run TestTypeScriptExample
go test -v -run TestPythonExample
```

## Cost Considerations

⚠️ **Important**: Most IBM Cloud resources incur costs. Always run `pulumi destroy` after testing to avoid unexpected charges.

### Estimated Costs (as of 2025)

| Example | Monthly Cost (USD) | Hourly Cost (USD) |
|---------|-------------------|-------------------|
| basic-* | Free | Free |
| vpc-with-subnet | ~$10 | ~$0.01 |
| kubernetes-cluster | ~$180 | ~$0.25 |
| object-storage-bucket | ~$1 | ~$0.001 |
| databases | ~$60 | ~$0.08 |
| three-tier-app | ~$300 | ~$0.42 |

Always check current IBM Cloud pricing: https://cloud.ibm.com/pricing

## Cleaning Up

After testing, destroy resources to avoid charges:

```bash
cd example-directory
pulumi destroy
pulumi stack rm dev  # Remove the stack
```

## Common Issues

### Issue: Plugin not found
```
error: could not load plugin for ibmcloud provider
```
**Solution**: Install the plugin:
```bash
pulumi plugin install resource ibmcloud v0.0.1 --server github://github.com/mapt-oss
```

### Issue: Authentication failed
```
Error: could not configure provider: failed to initialize client
```
**Solution**: Set your IBM Cloud API key:
```bash
export IC_API_KEY="your-api-key"
```

### Issue: Resource already exists
```
Error: resource with name already exists
```
**Solution**: Use unique names or configure auto-naming:
```go
Name: pulumi.String(fmt.Sprintf("my-resource-%s", ctx.Stack())),
```

## Contributing Examples

Want to contribute an example?

1. Fork the repository
2. Create a new example directory
3. Include:
   - Working code in at least one language
   - README.md with description and instructions
   - Clean resource naming
   - Proper error handling
4. Test with `pulumi up` and `pulumi destroy`
5. Submit a pull request

## Next Steps

After running the examples:

1. Explore the [IBM Cloud documentation](https://cloud.ibm.com/docs)
2. Check the [Provider documentation](../README.md)
3. View all 600 available resources in the [schema](../provider/cmd/pulumi-resource-ibmcloud/schema.json)
4. Join the [Pulumi Community Slack](https://slack.pulumi.com/)

## Resources

- [Pulumi Documentation](https://www.pulumi.com/docs/)
- [IBM Cloud Documentation](https://cloud.ibm.com/docs)
- [IBM Cloud Terraform Provider Docs](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs)
- [Pulumi IBM Cloud Provider GitHub](https://github.com/mapt-oss/pulumi-ibmcloud)
