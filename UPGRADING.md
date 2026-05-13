# Upgrading to New Terraform Provider Versions

This guide explains how to upgrade the Pulumi IBM Cloud provider when new versions of the upstream Terraform provider are released.

## Quick Upgrade Process

1. **Update the provider dependency** in `provider/go.mod`:
   ```bash
   cd provider
   go get github.com/IBM-Cloud/terraform-provider-ibm@v1.XX.X
   go mod tidy
   ```

2. **Rebuild the schema and provider**:
   ```bash
   cd /home/default/workdir/pulumi-ibmcloud

   # Build tfgen
   cd provider
   go build -o ../bin/pulumi-tfgen-ibmcloud \
     github.com/mapt-oss/pulumi-ibmcloud/provider/cmd/pulumi-tfgen-ibmcloud
   cd ..

   # Regenerate schema
   ./bin/pulumi-tfgen-ibmcloud schema --out provider/cmd/pulumi-resource-ibmcloud

   # Build provider
   cd provider
   go build -o ../bin/pulumi-resource-ibmcloud \
     github.com/mapt-oss/pulumi-ibmcloud/provider/cmd/pulumi-resource-ibmcloud
   cd ..
   ```

3. **Regenerate SDKs**:
   ```bash
   # Clean old SDKs
   rm -rf sdk/go/ibmcloud/*.go
   rm -rf sdk/nodejs/*.ts sdk/nodejs/bin
   rm -rf sdk/python/pulumi_ibmcloud/*.py

   # Generate new SDKs
   ./bin/pulumi-tfgen-ibmcloud go --out sdk/go/
   ./bin/pulumi-tfgen-ibmcloud nodejs --out sdk/nodejs/
   ./bin/pulumi-tfgen-ibmcloud python --out sdk/python/
   ```

4. **Test the build**:
   ```bash
   cd sdk/go
   go build ./...
   ```

## Known Issues & Fixes

### Issue: GetProjectConfigTypeOutput Duplicate Declaration

**Status**: ✅ Fixed in `provider/resources.go`

This issue was caused by a naming collision between two schema types:
- `ibmcloud:index/getProjectConfig:getProjectConfig` (nested type)
- `ibmcloud:index/getProjectConfigOutput:getProjectConfigOutput` (output values)

**Permanent Fix**:
The fix in `provider/resources.go` (lines 185-203) renames the `outputs` field to `OutputValues`, preventing the collision. This configuration will be applied automatically during schema generation.

**Verification**:
After regenerating the Go SDK, verify there's only one `GetProjectConfigTypeOutput` declaration:
```bash
grep -n "^type GetProjectConfigTypeOutput " sdk/go/ibmcloud/pulumiTypes30.go
# Should return exactly 1 line
```

And verify the renamed type exists:
```bash
grep -n "^type GetProjectConfigOutputValue " sdk/go/ibmcloud/pulumiTypes30.go
# Should return exactly 1 line
```

### Issue: C# SDK Filename Too Long

**Status**: ❌ Not fixed (open issue)

The .NET SDK generation fails due to a 255+ character filename for the backup recovery resource.

**Workaround**:
Skip C# SDK generation, or add custom name mappings in `provider/resources.go`:

```go
prov.Resources = map[string]*tfbridge.ResourceInfo{
    "ibm_backup_recovery_protection_sources": {
        Tok: tokens.MakeResource(mainPkg, mainMod, "BackupRecoveryProtectionSources"),
        Fields: map[string]*tfbridge.SchemaInfo{
            // Flatten deeply nested properties here
        },
    },
}
```

## Dependency Management

### Kubernetes Dependencies

The IBM Cloud provider has complex Kubernetes dependencies. If you encounter dependency errors:

1. **Check upstream provider's `go.mod`**:
   ```bash
   curl -s https://raw.githubusercontent.com/IBM-Cloud/terraform-provider-ibm/v1.XX.X/go.mod | grep -A5 "replace ("
   ```

2. **Update our `provider/go.mod`** to match their replace directives

3. **Common issue**: k8s.io version conflicts
   - Solution: Align all k8s.io modules to the same version
   - Check lines 13-45 in `provider/go.mod` for current k8s replacements

### Example go.mod Update

If the upstream provider updates to use k8s.io v0.34.0:

```bash
cd provider
# Update all k8s.io replacements
sed -i 's/v0.33.4/v0.34.0/g' go.mod
go mod tidy
```

## Testing After Upgrade

1. **Schema validation**:
   ```bash
   # Check schema size and resource count
   cat provider/cmd/pulumi-resource-ibmcloud/schema.json | jq '.resources | length'
   cat provider/cmd/pulumi-resource-ibmcloud/schema.json | jq '.functions | length'
   ```

2. **Go SDK compilation**:
   ```bash
   cd sdk/go
   go build ./...
   go test ./... -v
   ```

3. **TypeScript SDK**:
   ```bash
   cd sdk/nodejs
   npm install
   npm run build
   ```

4. **Python SDK**:
   ```bash
   cd sdk/python
   pip install -e .
   python -c "import pulumi_ibmcloud; print('OK')"
   ```

## Version Compatibility Matrix

| Pulumi Provider | Terraform Provider | Go Version | Pulumi CLI |
|-----------------|-------------------|------------|------------|
| v1.0.0          | v1.85.0           | 1.23.4     | 3.208.0    |
| v1.1.0          | v2.1.0            | 1.25.2     | 3.200+     |

## Breaking Changes to Watch For

When upgrading the Terraform provider, check for:

1. **Removed resources** - Check the upstream CHANGELOG
2. **Renamed properties** - May require SDK breaking changes
3. **New required fields** - May break existing programs
4. **Deprecated fields** - Add deprecation notices in schema

## Troubleshooting

### Build fails with "missing go.sum entry"
```bash
cd provider
GOFLAGS="-mod=mod" go mod tidy
go mod download
```

### SDK generation fails
```bash
# Check schema is valid JSON
jq empty provider/cmd/pulumi-resource-ibmcloud/schema.json

# Verify provider builds
cd provider
go build ./...
```

### Type collision errors (similar to GetProjectConfigTypeOutput)
Add field mappings in `provider/resources.go`:
```go
prov.DataSources["ibm_problematic_resource"] = &tfbridge.DataSourceInfo{
    Fields: map[string]*tfbridge.SchemaInfo{
        "conflicting_field": {
            Name: "AlternativeName",
        },
    },
}
```

## Automation Recommendations

Consider creating a GitHub Actions workflow for automated upgrades:

```yaml
name: Check Terraform Provider Updates
on:
  schedule:
    - cron: '0 0 * * 1'  # Weekly on Monday
  workflow_dispatch:

jobs:
  check-updates:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Check for new provider version
        run: |
          LATEST=$(curl -s https://api.github.com/repos/IBM-Cloud/terraform-provider-ibm/releases/latest | jq -r .tag_name)
          CURRENT=$(grep 'github.com/IBM-Cloud/terraform-provider-ibm' provider/go.mod | awk '{print $2}')
          if [ "$LATEST" != "$CURRENT" ]; then
            echo "New version available: $LATEST (current: $CURRENT)"
            # Create issue or PR
          fi
```

## Getting Help

- **Issue tracker**: https://github.com/mapt-oss/pulumi-ibmcloud/issues
- **Upstream provider**: https://github.com/IBM-Cloud/terraform-provider-ibm
- **Pulumi bridge docs**: https://github.com/pulumi/pulumi-terraform-bridge

## References

- Original build documentation: `DEVELOPMENT.md`
- Session notes: `CLAUDE.md`
- Provider configuration: `provider/resources.go`
