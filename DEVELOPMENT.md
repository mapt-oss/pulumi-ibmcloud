# Development Guide for Pulumi IBM Cloud Provider

This guide provides detailed instructions for developing, building, and testing the Pulumi IBM Cloud provider.

## Architecture

This provider is a Pulumi bridge over the [terraform-provider-ibm](https://github.com/IBM-Cloud/terraform-provider-ibm). The bridge:

1. Wraps the Terraform provider using the Pulumi Terraform Bridge
2. Generates language-specific SDKs (TypeScript/JavaScript, Python, Go, C#)
3. Provides idiomatic Pulumi resource interfaces

## Project Structure

```
pulumi-ibmcloud/
├── provider/               # Provider implementation
│   ├── cmd/
│   │   ├── pulumi-resource-ibmcloud/  # Main provider binary
│   │   └── pulumi-tfgen-ibmcloud/     # Schema generation tool
│   ├── resources.go        # Provider configuration and mappings
│   ├── go.mod             # Go dependencies
│   └── pkg/
│       └── version/       # Version information
├── sdk/                   # Generated SDKs (after build)
│   ├── nodejs/           # TypeScript/JavaScript SDK
│   ├── python/           # Python SDK
│   ├── go/               # Go SDK
│   └── dotnet/           # C# SDK
├── examples/             # Example programs
└── Makefile             # Build automation
```

## Prerequisites

Before you begin, ensure you have the following installed:

1. **Go 1.24 or later**
   ```bash
   # Verify installation
   go version
   ```

2. **Node.js** (Active LTS version)
   ```bash
   # Verify installation
   node --version
   npm --version
   ```

3. **Python 3.8 or later**
   ```bash
   # Verify installation
   python3 --version
   ```

4. **.NET SDK 8.0 or later**
   ```bash
   # Verify installation
   dotnet --version
   ```

5. **pulumictl**
   ```bash
   # Install pulumictl
   go install github.com/pulumi/pulumictl@latest
   ```

6. **Pulumi CLI**
   ```bash
   # Install from https://www.pulumi.com/docs/install/
   pulumi version
   ```

## Building the Provider

### Step 1: Download Dependencies

Navigate to the provider directory and download Go dependencies:

```bash
cd provider
go mod download
go mod tidy
```

### Step 2: Generate the Schema

Generate the Pulumi schema from the Terraform provider:

```bash
cd ..
make tfgen
```

This command:
- Inspects the Terraform provider schema
- Generates `provider/cmd/pulumi-resource-ibmcloud/schema.json`
- Creates token mappings for resources and data sources

**Note:** You may see warnings during schema generation. Review them to ensure all resources are properly mapped.

### Step 3: Build the Provider Binary

Build the provider plugin:

```bash
make provider
```

This creates `bin/pulumi-resource-ibmcloud`, the main provider binary.

### Step 4: Generate Language SDKs

Generate SDKs for all supported languages:

```bash
make build_sdks
```

This generates:
- `sdk/nodejs/` - TypeScript/JavaScript SDK
- `sdk/python/` - Python SDK
- `sdk/go/` - Go SDK
- `sdk/dotnet/` - C# SDK

### Step 5: Install the Provider Locally

To test the provider locally:

```bash
# Copy the binary to your PATH
cp bin/pulumi-resource-ibmcloud ~/.pulumi/bin/

# Install the Node.js SDK locally
make install_nodejs_sdk
```

## Testing the Provider

### Manual Testing

1. **Create a test program**:

```bash
mkdir -p examples/quickstart-ts
cd examples/quickstart-ts
pulumi new typescript --yes
```

2. **Link the local SDK** (for TypeScript):

```bash
npm link ../../sdk/nodejs/bin
```

3. **Write test code** in `index.ts`:

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as ibmcloud from "@pulumi/ibmcloud";

// Test by creating a simple resource
const config = new pulumi.Config("ibmcloud");
const region = config.get("region") || "us-south";

export const providerRegion = region;
```

4. **Set up IBM Cloud credentials**:

```bash
export IC_API_KEY="your-ibm-cloud-api-key"
pulumi config set ibmcloud:region us-south
```

5. **Preview the stack**:

```bash
pulumi preview
```

### Automated Testing

Run the integration tests:

```bash
cd examples
go test -v -tags=all
```

## Common Development Tasks

### Updating the Terraform Provider Version

1. Update the version in `provider/go.mod`:

```go
require (
    github.com/IBM-Cloud/terraform-provider-ibm v1.XX.X
    ...
)
```

2. Update dependencies:

```bash
cd provider
go mod tidy
```

3. Rebuild the provider and SDKs:

```bash
make build_sdks
```

### Adding Custom Resource Mappings

If automatic token mapping doesn't work for some resources, add manual mappings in `provider/resources.go`:

```go
prov := tfbridge.ProviderInfo{
    Resources: map[string]*tfbridge.ResourceInfo{
        "ibm_resource_instance": {
            Tok: tfbridge.MakeResource(mainPkg, mainMod, "ResourceInstance"),
            // Add custom configuration here
        },
    },
}
```

### Customizing Provider Configuration

Override provider configuration defaults in `provider/resources.go`:

```go
Config: map[string]*tfbridge.SchemaInfo{
    "region": {
        Default: &tfbridge.DefaultInfo{
            EnvVars: []string{"IC_REGION", "IBMCLOUD_REGION"},
            Value:   "us-south",
        },
    },
},
```

## Troubleshooting

### Schema Generation Issues

If `make tfgen` fails:

1. Check that the upstream provider import is correct in `provider/resources.go`
2. Verify Go dependencies with `go mod tidy`
3. Look for breaking changes in the Terraform provider

### SDK Generation Issues

If `make build_sdks` fails:

1. Ensure the schema was generated successfully
2. Check that all language toolchains are installed
3. Review error messages for missing dependencies

### Dependency Conflicts

The IBM Cloud provider has complex Kubernetes dependencies. If you encounter issues:

1. Check the `replace` directives in `provider/go.mod`
2. Update to compatible versions
3. Run `go mod tidy` to resolve conflicts

## Next Steps

1. **Add Examples**: Create example programs in `examples/` for common use cases
2. **Write Tests**: Add integration tests in `examples/*_test.go`
3. **Documentation**: Document resources and configuration options
4. **CI/CD**: Set up GitHub Actions for automated builds and releases
5. **Publish**: Publish to the Pulumi Registry for public use

## Resources

- [Pulumi Terraform Bridge Documentation](https://github.com/pulumi/pulumi-terraform-bridge)
- [IBM Cloud Terraform Provider](https://github.com/IBM-Cloud/terraform-provider-ibm)
- [Pulumi Provider Authoring Guide](https://www.pulumi.com/docs/iac/packages-and-automation/pulumi-packages/authoring/)
- [IBM Cloud API Documentation](https://cloud.ibm.com/docs)

## Getting Help

- File issues in this repository
- Ask questions in the [Pulumi Community Slack](https://slack.pulumi.com)
- Consult the [Pulumi documentation](https://www.pulumi.com/docs/)
