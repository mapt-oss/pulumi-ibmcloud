# Basic Python Example - IBM Cloud Resource Group and Cloud Object Storage

This example demonstrates how to create an IBM Cloud Resource Group and Cloud Object Storage instance using the Pulumi IBM Cloud provider with Python.

## Prerequisites

1. [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)
2. [Python](https://www.python.org/downloads/) 3.8 or later
3. [IBM Cloud API Key](https://cloud.ibm.com/iam/apikeys)

## Setup

1. **Install the provider plugin**:
   ```bash
   pulumi plugin install resource ibmcloud v0.0.1 --server github://github.com/mapt-oss
   ```

2. **Set IBM Cloud credentials**:
   ```bash
   export IC_API_KEY="your-ibm-cloud-api-key"
   export IC_REGION="global"
   ```

3. **Create a virtual environment and install dependencies**:
   ```bash
   python3 -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   pip install -r requirements.txt

   # Install the IBM Cloud provider SDK locally
   pip install /path/to/pulumi-ibmcloud/sdk/python
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
- 1 Cloud Object Storage instance (Standard plan)

## Outputs

- `resource_group_id` - The unique ID of the resource group
- `resource_group_name` - The name of the resource group
- `cos_instance_id` - The unique ID of the COS instance
- `cos_instance_crn` - The Cloud Resource Name (CRN) of the COS instance

## Cost

- **Resource Group**: Free
- **Cloud Object Storage (Standard plan)**: ~$0.023/GB/month for storage
  - First 25 GB/month free
  - This example only creates the instance, no storage costs until you create buckets and store data

See [IBM Cloud Object Storage pricing](https://www.ibm.com/cloud/object-storage/pricing) for details.
