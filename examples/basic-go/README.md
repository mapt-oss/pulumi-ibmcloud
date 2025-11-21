# Basic Go Example - IBM Cloud Resource Group

This example demonstrates how to create an IBM Cloud Resource Group using the Pulumi IBM Cloud provider with Go.

## Prerequisites

1. [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)
2. [Go](https://golang.org/dl/) 1.20 or later
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

3. **Initialize dependencies**:
   ```bash
   go mod init example
   go get github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud
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

## Outputs

- `resourceGroupId` - The unique ID of the resource group
- `resourceGroupName` - The name of the resource group

## Cost

Resource Groups are free on IBM Cloud.
