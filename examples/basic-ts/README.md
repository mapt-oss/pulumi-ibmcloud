# Basic TypeScript Example - IBM Cloud Resource Group and VPC

This example demonstrates how to create an IBM Cloud Resource Group, VPC, and Subnet using the Pulumi IBM Cloud provider with TypeScript.

## Prerequisites

1. [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)
2. [Node.js](https://nodejs.org/) 18 or later
3. [IBM Cloud API Key](https://cloud.ibm.com/iam/apikeys)

## Setup

1. **Install the provider plugin**:
   ```bash
   pulumi plugin install resource ibmcloud v0.0.1 --server github://github.com/mapt-oss
   ```

2. **Set IBM Cloud credentials**:
   ```bash
   export IC_API_KEY="your-ibm-cloud-api-key"
   export IC_REGION="us-south"
   ```

3. **Install dependencies**:
   ```bash
   npm install

   # Link the IBM Cloud provider SDK locally
   cd /path/to/pulumi-ibmcloud/sdk/nodejs
   npm install && npm run build && npm link
   cd -
   npm link @pulumi/ibmcloud
   ```

## Deploy

1. **Initialize a new stack**:
   ```bash
   pulumi stack init dev
   ```

2. **Preview the deployment**:
   ```bash
   pulumi preview
   ```

3. **Deploy the resources**:
   ```bash
   pulumi up
   ```

4. **View outputs**:
   ```bash
   pulumi stack output
   ```

## Clean Up

To destroy the resources:

```bash
pulumi destroy
pulumi stack rm dev
```

## Resources Created

- 1 IBM Cloud Resource Group
- 1 VPC (Virtual Private Cloud)
- 1 Subnet (in zone us-south-1)

## Outputs

- `resourceGroupId` - The unique ID of the resource group
- `resourceGroupName` - The name of the resource group
- `vpcId` - The unique ID of the VPC
- `vpcName` - The name of the VPC
- `subnetId` - The unique ID of the subnet
- `subnetCidr` - The CIDR block of the subnet

## Cost

- **Resource Group**: Free
- **VPC**: Free for the VPC itself
- **Subnet**: Free
- **Public Gateway** (not included in this example): ~$0.05/hour (~$36/month)

VPCs on IBM Cloud have no cost, but resources within them (compute instances, load balancers, etc.) do incur charges.

See [IBM Cloud VPC pricing](https://www.ibm.com/cloud/vpc/pricing) for details.

## Network Configuration

This example creates a subnet with CIDR block `10.240.0.0/24` in zone `us-south-1`, providing 254 usable IP addresses for resources.
