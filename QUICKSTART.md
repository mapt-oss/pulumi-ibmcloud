# Quick Start Guide - Pulumi IBM Cloud Provider

This guide will help you build and test the Pulumi IBM Cloud provider.

## Current Status

✅ **Completed:**
- Provider configuration (`provider/resources.go`)
- Go module setup with IBM Cloud Terraform Provider v1.85.0
- All dependencies resolved
- Organization set to `mapt-oss`
- Upstream provider correctly set to `IBM-Cloud`
- Module paths updated to `github.com/mapt-oss/pulumi-ibmcloud`

⏳ **Next Steps:**
- Generate the provider schema
- Build the provider binary
- Generate language SDKs

## Understanding the SDK Structure

The `sdk/` directory currently contains **placeholder configuration files** only. The actual SDK code will be **automatically generated** when you run the build process. Here's what will happen:

### Before Build (Current State)
```
sdk/
├── go.mod                    # ✅ Updated module path
├── nodejs/
│   └── package.json         # ✅ Updated package name
├── python/
│   └── (placeholder files will be replaced)
└── dotnet/
    └── (placeholder files will be replaced)
```

### After Build (Generated State)
```
sdk/
├── go/
│   └── ibmcloud/           # ← Generated Go SDK
├── nodejs/
│   └── bin/                # ← Generated TypeScript SDK
├── python/
│   └── pulumi_ibmcloud/    # ← Generated Python SDK
└── dotnet/
    └── Pulumi.IBMCloud/    # ← Generated C# SDK
```

## Prerequisites

You need these tools installed:

1. **Go 1.24+** ✅ (already installed at `$HOME/go/bin/go`)
2. **Node.js & npm** - For TypeScript SDK generation
3. **Python 3.8+** - For Python SDK generation
4. **.NET SDK 8.0+** - For C# SDK generation
5. **pulumictl** - Pulumi's build tool
6. **Pulumi CLI** - The Pulumi command-line tool

## Installation Commands

### Install Node.js (if needed)
```bash
# Install nvm first
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc

# Install Node.js LTS
nvm install --lts
nvm use --lts
```

### Install Python (usually pre-installed)
```bash
python3 --version  # Should be 3.8+
```

### Install .NET SDK
```bash
# Download and install .NET 8.0
curl -sSL https://dot.net/v1/dotnet-install.sh | bash /dev/stdin --channel 8.0
echo 'export PATH=$PATH:$HOME/.dotnet' >> ~/.bashrc
source ~/.bashrc
```

### Install pulumictl
```bash
$HOME/go/bin/go install github.com/pulumi/pulumictl@latest
```

### Install Pulumi CLI
```bash
curl -fsSL https://get.pulumi.com | sh
echo 'export PATH=$PATH:$HOME/.pulumi/bin' >> ~/.bashrc
source ~/.bashrc
```

## Build Process

Once all prerequisites are installed, follow these steps:

### Step 1: Generate the Schema
```bash
cd /home/default/workdir/pulumi-ibmcloud
export PATH=$PATH:$HOME/go/bin
make tfgen
```

This will:
- Inspect the IBM Cloud Terraform provider
- Generate `provider/cmd/pulumi-resource-ibmcloud/schema.json`
- Create automatic token mappings for all resources

**Expected output:** You'll see a list of resources and data sources being mapped.

### Step 2: Build the Provider Binary
```bash
make provider
```

This creates:
- `bin/pulumi-resource-ibmcloud` - The provider plugin

**Expected output:** Binary created successfully.

### Step 3: Generate Language SDKs
```bash
make build_sdks
```

This generates all four SDKs:
- `sdk/nodejs/bin/` - TypeScript/JavaScript
- `sdk/python/pulumi_ibmcloud/` - Python
- `sdk/go/ibmcloud/` - Go
- `sdk/dotnet/Pulumi.IBMCloud/` - C#

**Expected output:** SDKs generated for all languages.

### Step 4: Tidy Go SDK Dependencies
```bash
cd sdk
$HOME/go/bin/go mod tidy
cd ..
```

## Testing the Provider Locally

### Install the Provider Binary
```bash
# Create Pulumi plugins directory
mkdir -p ~/.pulumi/bin

# Copy the provider binary
cp bin/pulumi-resource-ibmcloud ~/.pulumi/bin/
```

### Link the Node.js SDK (for testing)
```bash
make install_nodejs_sdk
```

### Create a Test Program
```bash
mkdir -p examples/test-ibmcloud
cd examples/test-ibmcloud
pulumi new typescript --yes
```

### Link to Local SDK
```bash
npm link @pulumi/ibmcloud
```

### Write Test Code
Edit `index.ts`:
```typescript
import * as pulumi from "@pulumi/pulumi";
import * as ibmcloud from "@pulumi/ibmcloud";

// Simple test - just configure the provider
const config = new pulumi.Config("ibmcloud");
export const configured = "IBM Cloud provider loaded successfully!";
```

### Set Up Credentials
```bash
export IC_API_KEY="your-ibm-cloud-api-key"
pulumi config set ibmcloud:region us-south
```

### Test the Stack
```bash
pulumi preview
```

## Troubleshooting

### "make: command not found"
The Makefile needs `make` utility. It should be available, but if not, you can run individual commands from the Makefile directly.

### "go: command not found"
Make sure Go is in your PATH:
```bash
export PATH=$PATH:$HOME/go/bin
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
```

### "pulumi: command not found"
Install the Pulumi CLI (see prerequisites above).

### Schema generation warnings
Some warnings during `make tfgen` are normal. Look for actual errors that stop the build.

### Dependency conflicts
Already resolved! The `go.mod` includes necessary replace directives for Kubernetes dependencies.

## What You Have Right Now

```
✅ Provider code configured (provider/resources.go)
✅ Go dependencies resolved (provider/go.mod)
✅ Command binaries configured (provider/cmd/)
✅ SDK metadata updated (sdk/*/package.json, go.mod)
✅ Documentation created (README.md, DEVELOPMENT.md)
✅ CI/CD config updated (.ci-mgmt.yaml)

⏳ Waiting to generate:
   - Provider schema (make tfgen)
   - Provider binary (make provider)
   - Language SDKs (make build_sdks)
```

## Next Steps After Build

1. **Test with real IBM Cloud resources**
2. **Add example programs** in `examples/`
3. **Run integration tests** with `make test`
4. **Set up GitHub repository** at `github.com/mapt-oss/pulumi-ibmcloud`
5. **Configure GitHub Actions** for CI/CD
6. **Publish to Pulumi Registry** (optional)

## Support

- Check `DEVELOPMENT.md` for detailed developer information
- See `README.md` for user-facing documentation
- File issues on GitHub once repository is created

## Summary

The provider is **fully configured** and ready to build. The `sdk/` folder will contain the actual generated SDK code once you run the build commands. All the configuration is correct and pointing to `mapt-oss/pulumi-ibmcloud`.
