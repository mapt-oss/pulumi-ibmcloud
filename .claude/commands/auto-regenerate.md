# Auto-Regenerate Provider (One Command)

Automatically regenerate the entire Pulumi IBM Cloud provider with a single command.

Run this bash script to regenerate everything:

```bash
#!/bin/bash
set -e

echo "🧹 Cleaning previous artifacts..."
cd /home/default/workdir
rm -rf provider/bin/* sdk/go/ibmcloud/* sdk/nodejs/*.ts sdk/nodejs/bin sdk/nodejs/types sdk/python/pulumi_ibmcloud/* provider/cmd/pulumi-resource-ibmcloud/schema.json

echo "🔨 Building tfgen binary..."
cd provider
GOFLAGS="-mod=mod" $HOME/go/bin/go build \
  -o ./bin/pulumi-tfgen-ibmcloud \
  -ldflags "-X github.com/mapt-oss/pulumi-ibmcloud/provider/pkg/version.Version=0.0.9" \
  ./cmd/pulumi-tfgen-ibmcloud

echo "📋 Generating schema..."
cd /home/default/workdir
./provider/bin/pulumi-tfgen-ibmcloud schema --out provider/cmd/pulumi-resource-ibmcloud

echo "🔨 Building provider binary..."
cd provider
GOFLAGS="-mod=mod" $HOME/go/bin/go build \
  -o ./bin/pulumi-resource-ibmcloud \
  -ldflags "-X github.com/mapt-oss/pulumi-ibmcloud/provider/pkg/version.Version=0.0.9" \
  ./cmd/pulumi-resource-ibmcloud

echo "📦 Setting up environment..."
export PULUMI_HOME=/home/default/workdir/.pulumi
export PULUMI_CONVERT=1
export PATH=$PATH:$HOME/.pulumi/bin:$HOME/go/bin

echo "🐹 Generating Go SDK..."
cd /home/default/workdir
./provider/bin/pulumi-tfgen-ibmcloud go --out sdk/go/

echo "📘 Generating TypeScript SDK..."
./provider/bin/pulumi-tfgen-ibmcloud nodejs --out sdk/nodejs/

echo "🐍 Generating Python SDK..."
./provider/bin/pulumi-tfgen-ibmcloud python --out sdk/python/

echo "✅ Regeneration complete!"
echo ""
echo "📊 Summary:"
echo "  - tfgen binary: $(ls -lh provider/bin/pulumi-tfgen-ibmcloud | awk '{print $5}')"
echo "  - provider binary: $(ls -lh provider/bin/pulumi-resource-ibmcloud | awk '{print $5}')"
echo "  - schema: $(ls -lh provider/cmd/pulumi-resource-ibmcloud/schema.json | awk '{print $5}')"
echo "  - Go files: $(find sdk/go -name "*.go" | wc -l)"
echo "  - TypeScript files: $(find sdk/nodejs -name "*.ts" | wc -l)"
echo "  - Python files: $(find sdk/python -name "*.py" | wc -l)"
```

To use this command, just ask Claude to:
"Run the /auto-regenerate command"
