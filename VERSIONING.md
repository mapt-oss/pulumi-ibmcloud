# Versioning Guide for IBM Cloud Pulumi Provider

This document describes the versioning system for the IBM Cloud Pulumi provider and how to create releases.

## Table of Contents

- [Overview](#overview)
- [Version Detection](#version-detection)
- [Release Process](#release-process)
- [Local Development](#local-development)
- [CI/CD Workflows](#cicd-workflows)
- [Troubleshooting](#troubleshooting)

## Overview

The IBM Cloud Pulumi provider uses **git tags** as the source of truth for versioning. The version is automatically detected from git tags and injected into:

- Provider binaries
- Generated SDKs (Go, Node.js/TypeScript, Python)
- Package metadata files
- Release artifacts

### Version Format

Versions follow [Semantic Versioning 2.0.0](https://semver.org/):

- **Release versions**: `v1.2.3` (tag) → `1.2.3` (in code)
- **Pre-release versions**: `v1.2.3-alpha.1`, `v1.2.3-beta.2`, `v1.2.3-rc.1`
- **Development versions**: `1.2.3+dev.5.g1a2b3c4` (auto-generated between releases)

### Tag Structure

The repository uses two types of tags:

1. **Provider tags**: `v1.2.3` - Main provider release tags
2. **SDK tags**: `sdk/v1.2.3` - SDK-specific tags for Go module versioning

## Version Detection

### Automatic Detection Script

The version is automatically detected by `scripts/get-version.sh`:

```bash
./scripts/get-version.sh
# Output: 0.0.6 (if on tag v0.0.6)
# Output: 0.0.6+dev.3.g1a2b3c4 (if 3 commits after v0.0.6)
```

### Detection Logic

1. **Exact tag match**: If current commit has a `v*.*.*` tag → use that version
2. **Between tags**: If commits exist after latest tag → append dev suffix
3. **No tags**: Fall back to `0.0.0+dev.{sha}`

### Integration with Build System

The Makefile automatically uses the version detection script:

```makefile
PROVIDER_VERSION ?= $(shell ./scripts/get-version.sh 2>/dev/null || echo "1.0.0-alpha.0+dev")
```

You can override the version:

```bash
make build PROVIDER_VERSION=1.2.3
```

## Release Process

### Step 1: Prepare the Release

1. **Update CHANGELOG** (if you maintain one):
   ```bash
   vim CHANGELOG.md
   # Add release notes for new version
   ```

2. **Ensure clean working directory**:
   ```bash
   git status
   # Should show clean working tree
   ```

3. **Ensure you're on main branch**:
   ```bash
   git checkout main
   git pull origin main
   ```

### Step 2: Create and Push Tag

1. **Create annotated tag**:
   ```bash
   VERSION=v1.2.3
   git tag -a "$VERSION" -m "Release $VERSION"
   ```

2. **Push the tag**:
   ```bash
   git push origin "$VERSION"
   ```

### Step 3: Automated Release (GitHub Actions)

Once the tag is pushed, GitHub Actions automatically:

1. **Build Provider Binaries** (`.github/workflows/simple-release.yml`):
   - Builds for multiple platforms (Linux, macOS, Windows)
   - Creates GitHub release with provider binaries
   - Creates `sdk/v1.2.3` tag

2. **Generate and Publish SDKs** (`.github/workflows/publish-sdks.yml`):
   - Generates SDKs for Go, Node.js, Python
   - Builds and packages each SDK
   - Uploads SDK packages to GitHub release
   - Creates SDK-specific release under `sdk/v1.2.3` tag

### Step 4: Verify Release

1. Check GitHub Releases page:
   ```
   https://github.com/mapt-oss/pulumi-ibmcloud/releases
   ```

2. Verify two releases were created:
   - `v1.2.3` - Provider binaries
   - `sdk/v1.2.3` - SDK packages

3. Test the release:
   ```bash
   # Go
   go get github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud@sdk/v1.2.3

   # Node.js
   npm install @pulumi/ibmcloud@1.2.3

   # Python
   pip install pulumi-ibmcloud==1.2.3
   ```

## Local Development

### Building with Auto-Detected Version

When building locally without a tag, the version automatically includes commit information:

```bash
# On commit 3 commits after v0.0.6
make provider
# Version will be: 0.0.6+dev.3.g1a2b3c4
```

### Building with Custom Version

Override the version for testing:

```bash
# Build provider with custom version
make provider PROVIDER_VERSION=1.2.3-test

# Generate SDKs with custom version
make generate_sdks PROVIDER_VERSION=1.2.3-test

# Full build with custom version
make build PROVIDER_VERSION=1.2.3-test
```

### Building from Specific Tag

Checkout a tag and build:

```bash
git checkout v0.0.6
make build
# Version will be: 0.0.6
```

## CI/CD Workflows

### Workflow: Simple Release

**File**: `.github/workflows/simple-release.yml`

**Triggers**:
- Push of tags matching `v*.*.*`
- Manual dispatch with version input

**Outputs**:
- Provider binaries for multiple platforms
- GitHub release with checksums
- SDK tags created automatically

### Workflow: Publish SDKs

**File**: `.github/workflows/publish-sdks.yml`

**Triggers**:
- Release published (after simple-release completes)
- Manual dispatch with version input

**Outputs**:
- Generated SDK packages (Go, Node.js, Python)
- SDK release with installation instructions
- Packaged artifacts for distribution

### Manual Release Trigger

You can manually trigger a release from GitHub Actions UI:

1. Go to Actions → Simple Release → Run workflow
2. Enter version (e.g., `v1.2.3`)
3. Click "Run workflow"

## Version Injection Points

### Provider Binary

Version is injected via Go linker flags:

```bash
-ldflags "-X github.com/mapt-oss/pulumi-ibmcloud/provider/pkg/version.Version=1.2.3"
```

**File**: `provider/pkg/version/version.go`
```go
package version

// Version is initialized by the Go linker to contain the semver of this build.
var Version string
```

### Node.js SDK

Version is set in `package.json` during SDK generation:

**File**: `sdk/nodejs/package.json`
```json
{
  "name": "@pulumi/ibmcloud",
  "version": "1.2.3",
  "pulumi": {
    "version": "1.2.3"
  }
}
```

### Python SDK

Version is set in `pyproject.toml` during SDK generation:

**File**: `sdk/python/pyproject.toml`
```toml
[project]
  name = "pulumi_ibmcloud"
  version = "1.2.3"
```

### Go SDK

Version is referenced via the SDK tag:

```bash
go get github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud@sdk/v1.2.3
```

### Plugin Metadata

Version is included in `pulumi-plugin.json` files:

**Files**:
- `sdk/go/ibmcloud/pulumi-plugin.json`
- `sdk/python/pulumi_ibmcloud/pulumi-plugin.json`
- `sdk/nodejs/pulumi-plugin.json` (if generated)

## Troubleshooting

### Issue: Version shows "1.0.0-alpha.0+dev" instead of tag

**Cause**: Version detection script not found or not executable

**Solution**:
```bash
chmod +x scripts/get-version.sh
./scripts/get-version.sh  # Should output current version
```

### Issue: Wrong version in generated SDKs

**Cause**: SDKs were generated before tag was created

**Solution**:
```bash
# Ensure you're on the correct tag
git checkout v1.2.3

# Clean and regenerate
make clean
make generate_sdks PROVIDER_VERSION=1.2.3
```

### Issue: "PROVIDER_VERSION should not start with a 'v'" error

**Cause**: Version includes 'v' prefix

**Solution**:
```bash
# Correct
make build PROVIDER_VERSION=1.2.3

# Incorrect
make build PROVIDER_VERSION=v1.2.3
```

The version detection script automatically strips the 'v' prefix.

### Issue: Git tag exists but version not detected

**Cause**: Tag doesn't match expected pattern `v*.*.*`

**Solution**:
```bash
# List all tags
git tag -l

# Ensure tag follows pattern
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

### Issue: SDK generation fails with version mismatch

**Cause**: Provider schema was generated with different version

**Solution**:
```bash
# Regenerate everything with consistent version
make clean
make schema PROVIDER_VERSION=1.2.3
make generate_sdks PROVIDER_VERSION=1.2.3
```

## Pre-Release Versions

For alpha, beta, or release candidate versions:

```bash
# Alpha release
git tag -a v1.2.3-alpha.1 -m "Alpha release 1.2.3-alpha.1"
git push origin v1.2.3-alpha.1

# Beta release
git tag -a v1.2.3-beta.1 -m "Beta release 1.2.3-beta.1"
git push origin v1.2.3-beta.1

# Release candidate
git tag -a v1.2.3-rc.1 -m "Release candidate 1.2.3-rc.1"
git push origin v1.2.3-rc.1
```

Pre-release versions are automatically marked as "pre-release" in GitHub:

```yaml
prerelease: ${{ contains(steps.version.outputs.version, '-') }}
```

## Version Bumping Strategy

Follow semantic versioning guidelines:

### Patch Release (1.2.3 → 1.2.4)
- Bug fixes
- Documentation updates
- Performance improvements (without breaking changes)

### Minor Release (1.2.3 → 1.3.0)
- New features (backward compatible)
- New resources added
- Deprecations (but not removals)

### Major Release (1.2.3 → 2.0.0)
- Breaking API changes
- Resource removals
- Major architectural changes
- Incompatible provider changes

## Best Practices

1. **Always use annotated tags**:
   ```bash
   git tag -a v1.2.3 -m "Release v1.2.3"  # Good
   git tag v1.2.3                          # Avoid
   ```

2. **Update documentation before tagging**:
   - Update CHANGELOG.md
   - Update README.md if needed
   - Update version references in examples

3. **Test before releasing**:
   ```bash
   make build
   make test
   make test_provider
   ```

4. **Follow semantic versioning strictly**:
   - Breaking changes require major version bump
   - New features require minor version bump
   - Bug fixes require patch version bump

5. **Keep tags and releases synchronized**:
   - Don't delete tags after releasing
   - Don't force-push tags
   - If you need to fix a release, create a new patch version

## Version History

To see version history:

```bash
# List all version tags
git tag -l 'v*.*.*' --sort=-version:refname

# Show tag details
git show v1.2.3

# See commits between versions
git log v1.2.2..v1.2.3 --oneline
```

## Related Files

- `scripts/get-version.sh` - Version detection script
- `Makefile` - Build system with version injection
- `.github/workflows/simple-release.yml` - Provider release workflow
- `.github/workflows/publish-sdks.yml` - SDK publishing workflow
- `provider/pkg/version/version.go` - Version package
- `.ci-mgmt.yaml` - CI/CD configuration

## Support

For questions about versioning:

1. Check this documentation first
2. Review existing releases: https://github.com/mapt-oss/pulumi-ibmcloud/releases
3. Open an issue: https://github.com/mapt-oss/pulumi-ibmcloud/issues

---

*Last updated: December 1, 2025*
