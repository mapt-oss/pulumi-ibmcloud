# Release Guide

## What Happens When You Create a Tag

When you push this repository to GitHub at `github.com/mapt-oss/pulumi-ibmcloud` and create a tag like `v0.0.1`, here's the complete flow:

### 1. Push to GitHub
```bash
git remote add origin https://github.com/mapt-oss/pulumi-ibmcloud.git
git push -u origin main
```

### 2. Create and Push a Tag
```bash
git tag v0.0.1
git push origin v0.0.1
```

### 3. GitHub Actions Workflow Triggers ✅

The `simple-release.yml` workflow automatically runs and:

1. **Builds provider binaries** for 6 platforms:
   - `linux-amd64`
   - `linux-arm64`
   - `darwin-amd64` (macOS Intel)
   - `darwin-arm64` (macOS Apple Silicon)
   - `windows-amd64`
   - `windows-arm64`

2. **Packages binaries**:
   - Linux/macOS: `.tar.gz` archives
   - Windows: `.zip` archives

3. **Creates GitHub Release** at:
   ```
   https://github.com/mapt-oss/pulumi-ibmcloud/releases/tag/v0.0.1
   ```

4. **Attaches artifacts**:
   - `pulumi-resource-ibmcloud-v0.0.1-linux-amd64.tar.gz`
   - `pulumi-resource-ibmcloud-v0.0.1-linux-arm64.tar.gz`
   - `pulumi-resource-ibmcloud-v0.0.1-darwin-amd64.tar.gz`
   - `pulumi-resource-ibmcloud-v0.0.1-darwin-arm64.tar.gz`
   - `pulumi-resource-ibmcloud-v0.0.1-windows-amd64.zip`
   - `pulumi-resource-ibmcloud-v0.0.1-windows-arm64.zip`
   - `checksums.txt` (SHA256 hashes)

### 4. Provider Becomes Usable as a Regular Provider ✅

Once the release is created, users can use it **exactly like any other Pulumi provider**:

## User Installation

### Option 1: Automatic Installation (Recommended)
Pulumi CLI can automatically download the plugin from GitHub releases:

```bash
pulumi plugin install resource ibmcloud v0.0.1 --server github://github.com/mapt-oss
```

Or, when users run a Pulumi program that uses the provider, Pulumi will automatically prompt to install it.

### Option 2: Manual Installation
Users can manually download and install:

```bash
# Download for your platform
wget https://github.com/mapt-oss/pulumi-ibmcloud/releases/download/v0.0.1/pulumi-resource-ibmcloud-v0.0.1-linux-amd64.tar.gz

# Extract
tar xzf pulumi-resource-ibmcloud-v0.0.1-linux-amd64.tar.gz

# Install
mkdir -p ~/.pulumi/bin
mv pulumi-resource-ibmcloud ~/.pulumi/bin/
chmod +x ~/.pulumi/bin/pulumi-resource-ibmcloud
```

## Using the Provider

### Go
```go
import "github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud"

// Go modules will automatically fetch from GitHub
```

```bash
go get github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud@v0.0.1
```

### TypeScript
```typescript
import * as ibmcloud from "@pulumi/ibmcloud";
```

**Installation** (local link since not published to NPM):
```bash
cd /path/to/pulumi-ibmcloud/sdk/nodejs
npm install && npm run build && npm link

# In your project
npm link @pulumi/ibmcloud
```

### Python
```python
import pulumi_ibmcloud as ibmcloud
```

**Installation** (local install since not published to PyPI):
```bash
pip install /path/to/pulumi-ibmcloud/sdk/python
```

## Example Pulumi Program

Create `Pulumi.yaml`:
```yaml
name: my-ibmcloud-project
runtime: go
description: IBM Cloud infrastructure
```

Create `main.go`:
```go
package main

import (
    "github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        rg, err := ibmcloud.NewResourceGroup(ctx, "my-rg", &ibmcloud.ResourceGroupArgs{
            Name: pulumi.String("pulumi-test-rg"),
        })
        if err != nil {
            return err
        }

        ctx.Export("resourceGroupId", rg.ID())
        return nil
    })
}
```

Run:
```bash
# Pulumi will automatically install the plugin from GitHub releases
pulumi up
```

## Verification

After creating the release, verify it works:

1. **Check the release exists**:
   ```bash
   curl -s https://api.github.com/repos/mapt-oss/pulumi-ibmcloud/releases/latest | jq -r '.tag_name'
   ```

2. **Verify plugin download works**:
   ```bash
   pulumi plugin install resource ibmcloud v0.0.1 --server github://github.com/mapt-oss
   ```

3. **Check installed plugin**:
   ```bash
   pulumi plugin ls
   ```

   Should show:
   ```
   NAME      KIND      VERSION  SIZE    INSTALLED
   ibmcloud  resource  0.0.1    274 MB  3 seconds ago
   ```

## Release Checklist

Before creating a release:

- [ ] All tests pass (`simple-build.yml` is green)
- [ ] Version number follows semver (e.g., `v0.0.1`, `v1.0.0`)
- [ ] CHANGELOG.md updated (if you maintain one)
- [ ] Schema is up-to-date (already embedded in the code)
- [ ] Go SDK is committed to the repository

Create the release:

- [ ] Create and push git tag
- [ ] Wait for `simple-release.yml` workflow to complete (~10-15 minutes)
- [ ] Verify release on GitHub has all 6 platform binaries + checksums
- [ ] Test plugin installation: `pulumi plugin install resource ibmcloud v0.0.1 --server github://github.com/mapt-oss`
- [ ] Test a simple Pulumi program

## Release Types

### Stable Release
```bash
git tag v1.0.0
git push origin v1.0.0
```
- Creates a **stable** release (not marked as prerelease)
- Sets as latest release

### Prerelease
```bash
git tag v1.0.0-alpha.1
git push origin v1.0.0-alpha.1
```
- Creates a **prerelease** (marked with "Pre-release" badge on GitHub)
- Not set as latest release

## Summary

✅ **Yes!** After you push to GitHub and create a tag `v0.0.1`:

1. GitHub Actions automatically builds and releases binaries
2. Users can install the provider with `pulumi plugin install resource ibmcloud v0.0.1 --server github://github.com/mapt-oss`
3. Users can use the Go SDK with `go get github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud`
4. Users can use TypeScript/Python SDKs via local installation
5. It works exactly like a regular Pulumi provider

**No secrets required. No external publishing. Just GitHub releases!**
