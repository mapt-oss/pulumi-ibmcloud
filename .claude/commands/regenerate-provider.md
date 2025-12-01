# Regenerate Pulumi IBM Cloud Provider

This command regenerates the complete Pulumi IBM Cloud provider, including:
- Schema generation
- Provider binaries
- All SDK packages (Go, TypeScript, Python)

## Prerequisites

Before running this regeneration, ensure:
1. Go 1.24.10+ is installed at `$HOME/go/bin/go`
2. Pulumi CLI is installed at `$HOME/.pulumi/bin/pulumi`
3. You are in the `/home/default/workdir` directory
4. All provider source code is up to date

## Full Regeneration Steps

Execute the following steps in order:

### 1. Clean Previous Artifacts
```bash
cd /home/default/workdir
rm -rf provider/bin/* sdk/go/ibmcloud/* sdk/nodejs/*.ts sdk/nodejs/bin sdk/nodejs/types sdk/python/pulumi_ibmcloud/* provider/cmd/pulumi-resource-ibmcloud/schema.json
```

### 2. Build tfgen Binary
```bash
cd provider
GOFLAGS="-mod=mod" $HOME/go/bin/go build \
  -o ./bin/pulumi-tfgen-ibmcloud \
  -ldflags "-X github.com/mapt-oss/pulumi-ibmcloud/provider/pkg/version.Version=0.0.9" \
  ./cmd/pulumi-tfgen-ibmcloud
```

Expected output: `provider/bin/pulumi-tfgen-ibmcloud` (~241MB)

### 3. Generate Schema
```bash
cd /home/default/workdir
./provider/bin/pulumi-tfgen-ibmcloud schema --out provider/cmd/pulumi-resource-ibmcloud
```

Expected output: `provider/cmd/pulumi-resource-ibmcloud/schema.json` (~5-41MB)

### 4. Build Provider Binary
```bash
cd provider
GOFLAGS="-mod=mod" $HOME/go/bin/go build \
  -o ./bin/pulumi-resource-ibmcloud \
  -ldflags "-X github.com/mapt-oss/pulumi-ibmcloud/provider/pkg/version.Version=0.0.9" \
  ./cmd/pulumi-resource-ibmcloud
```

Expected output: `provider/bin/pulumi-resource-ibmcloud` (~228-274MB)

### 5. Set Environment Variables
```bash
export PULUMI_HOME=/home/default/workdir/.pulumi
export PULUMI_CONVERT=1
export PATH=$PATH:$HOME/.pulumi/bin:$HOME/go/bin
```

### 6. Generate Go SDK
```bash
cd /home/default/workdir
./provider/bin/pulumi-tfgen-ibmcloud go --out sdk/go/
```

Expected output: ~538+ Go files in `sdk/go/ibmcloud/`

### 7. Generate TypeScript SDK
```bash
./provider/bin/pulumi-tfgen-ibmcloud nodejs --out sdk/nodejs/
```

Expected output: ~535+ TypeScript files in `sdk/nodejs/`

### 8. Generate Python SDK
```bash
./provider/bin/pulumi-tfgen-ibmcloud python --out sdk/python/
```

Expected output: ~534+ Python files in `sdk/python/pulumi_ibmcloud/`

## Quick Verification

After regeneration, verify the build:

```bash
# Check binaries
ls -lh provider/bin/

# Check schema
ls -lh provider/cmd/pulumi-resource-ibmcloud/schema.json

# Count SDK files
find sdk/go -name "*.go" | wc -l
find sdk/nodejs -name "*.ts" | wc -l
find sdk/python -name "*.py" | wc -l
```

## Troubleshooting

### Go version mismatch
If you get "go.mod requires go >= 1.24.7", install the latest Go:
```bash
rm -rf $HOME/go
curl -L https://go.dev/dl/go1.24.10.linux-amd64.tar.gz -o /tmp/go.tar.gz
tar -C $HOME -xzf /tmp/go.tar.gz
```

### Pulumi not found
Install Pulumi CLI:
```bash
curl -fsSL https://get.pulumi.com | sh
```

### Missing go.sum entries
Use `GOFLAGS="-mod=mod"` to allow go.sum updates during build.

## Notes

- Build time: ~10-20 minutes total
- Expected warnings about GOPATH/GOROOT are normal
- Schema generation warnings about upstream docs are expected
- C# SDK generation is currently disabled due to filename length issues
