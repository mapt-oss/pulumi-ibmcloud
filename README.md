# Pulumi IBM Cloud Provider

This repository contains a Pulumi provider for IBM Cloud, bridged from the [terraform-provider-ibm](https://github.com/IBM-Cloud/terraform-provider-ibm).

The IBM Cloud Resource Provider lets you manage IBM Cloud resources including VPC infrastructure, Cloud Object Storage, Watson services, and more.

## Installing

This package is available for several languages/platforms:

### Node.js (JavaScript/TypeScript)

To use from JavaScript or TypeScript in Node.js, install using either `npm`:

```bash
npm install @pulumi/ibmcloud
```

or `yarn`:

```bash
yarn add @pulumi/ibmcloud
```

### Python

To use from Python, install using `pip`:

```bash
pip install pulumi_ibmcloud
```

### Go

To use from Go, use `go get` to grab the latest version of the library:

```bash
go get github.com/mapt-oss/pulumi-ibmcloud/sdk/go/...
```

### .NET

To use from .NET, install using `dotnet add package`:

```bash
dotnet add package Pulumi.IBMCloud
```

## Configuration

The IBM Cloud provider supports configuration through environment variables or Pulumi configuration. Key configuration options include:

- `ibmcloud:region` - The IBM Cloud region (e.g., `us-south`, `eu-gb`)
- `ibmcloud:iaas_classic_username` - Classic Infrastructure username (optional)
- `ibmcloud:iaas_classic_api_key` - Classic Infrastructure API key (optional)
- `ibmcloud:ibmcloud_api_key` (environment: `IC_API_KEY` or `IBMCLOUD_API_KEY`) - IBM Cloud API key for authentication
- `ibmcloud:resource_group` - Default resource group for resources

For a complete list of configuration options, see the [IBM Cloud Terraform Provider documentation](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs).

## Example Usage

### TypeScript

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as ibmcloud from "@pulumi/ibmcloud";

// Create a resource group
const resourceGroup = new ibmcloud.ResourceGroup("my-resource-group", {
    name: "my-resource-group",
});

// Create a VPC
const vpc = new ibmcloud.IsVpc("my-vpc", {
    name: "my-vpc",
    resourceGroup: resourceGroup.id,
});
```

### Python

```python
import pulumi
import pulumi_ibmcloud as ibmcloud

# Create a resource group
resource_group = ibmcloud.ResourceGroup("my-resource-group",
    name="my-resource-group"
)

# Create a VPC
vpc = ibmcloud.IsVpc("my-vpc",
    name="my-vpc",
    resource_group=resource_group.id
)
```

### Go

```go
package main

import (
	"github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a resource group
		resourceGroup, err := ibmcloud.NewResourceGroup(ctx, "my-resource-group", &ibmcloud.ResourceGroupArgs{
			Name: pulumi.String("my-resource-group"),
		})
		if err != nil {
			return err
		}

		// Create a VPC
		_, err = ibmcloud.NewIsVpc(ctx, "my-vpc", &ibmcloud.IsVpcArgs{
			Name:          pulumi.String("my-vpc"),
			ResourceGroup: resourceGroup.ID(),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
```

### C# (.NET)

```csharp
using Pulumi;
using Pulumi.IBMCloud;

class Program
{
    static Task<int> Main() => Deployment.RunAsync(() =>
    {
        // Create a resource group
        var resourceGroup = new ResourceGroup("my-resource-group", new ResourceGroupArgs
        {
            Name = "my-resource-group",
        });

        // Create a VPC
        var vpc = new IsVpc("my-vpc", new IsVpcArgs
        {
            Name = "my-vpc",
            ResourceGroup = resourceGroup.Id,
        });
    });
}
```

## Building and Development

### Prerequisites

- Go 1.24 or later
- Node.js (Active LTS)
- Python 3
- .NET SDK
- [pulumictl](https://github.com/pulumi/pulumictl)

### Build the Provider

1. Clone this repository
2. Navigate to the `provider` directory
3. Install dependencies:

```bash
cd provider
go mod download
```

4. Build the provider and generate SDKs:

```bash
make build_sdks
```

### Testing

Run the tests:

```bash
make test
```

## Reference

For detailed reference documentation on all available resources and functions, please visit the [Pulumi Registry](https://www.pulumi.com/registry/packages/ibmcloud/api-docs/) (once published).

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This provider is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

## Acknowledgments

This provider is built using the [Pulumi Terraform Bridge](https://github.com/pulumi/pulumi-terraform-bridge) and wraps the [IBM Cloud Terraform Provider](https://github.com/IBM-Cloud/terraform-provider-ibm).
