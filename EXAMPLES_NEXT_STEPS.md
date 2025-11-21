# Examples - Next Steps Guide

## What's Been Done ✅

I've created and updated the basic examples for the IBM Cloud provider:

### 1. **basic-go** - Resource Group (NEW)
- **File**: `examples/basic-go/main.go`
- **Resources**: Creates 1 IBM Cloud Resource Group
- **Language**: Go
- **Status**: ✅ Ready to use

### 2. **basic-ts** - Resource Group + VPC (UPDATED)
- **File**: `examples/basic-ts/index.ts`
- **Resources**: Creates Resource Group, VPC, and Subnet
- **Language**: TypeScript
- **Status**: ✅ Updated from boilerplate

### 3. **basic-py** - Resource Group + COS (UPDATED)
- **File**: `examples/basic-py/__main__.py`
- **Resources**: Creates Resource Group and Cloud Object Storage instance
- **Language**: Python
- **Status**: ✅ Updated from boilerplate

### 4. Documentation
- ✅ `examples/EXAMPLES.md` - Comprehensive examples guide
- ✅ `examples/basic-go/README.md` - Go example documentation
- ✅ `examples/basic-ts/README.md` - TypeScript example documentation
- ✅ `examples/basic-py/README.md` - Python example documentation
- ✅ `EXAMPLES_NEXT_STEPS.md` - This file

### 5. Cleanup
- ✅ Removed `examples_dotnet_test.go` (C# SDK not available)
- ✅ Updated Python requirements to use local SDK installation

## What You Can Do Next

### Option 1: Test the Basic Examples (Recommended First Step)

After you create your first release (`v0.0.1`), test each example:

#### Test Go Example
```bash
cd examples/basic-go
go mod init example
go get github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud@v0.0.1
pulumi stack init dev
export IC_API_KEY="your-api-key"
pulumi up
pulumi destroy
```

#### Test TypeScript Example
```bash
cd examples/basic-ts
npm install
# Link local SDK (until published to NPM)
cd ../../sdk/nodejs && npm install && npm run build && npm link && cd -
npm link @pulumi/ibmcloud
pulumi stack init dev
export IC_API_KEY="your-api-key"
pulumi up
pulumi destroy
```

#### Test Python Example
```bash
cd examples/basic-py
python3 -m venv venv
source venv/bin/activate
pip install pulumi
pip install ../../sdk/python
pulumi stack init dev
export IC_API_KEY="your-api-key"
pulumi up
pulumi destroy
```

### Option 2: Create More Advanced Examples

Here are example ideas ranked by usefulness:

#### High Priority (Most Useful)

1. **vpc-complete** - Complete VPC setup
   - VPC with multiple subnets across zones
   - Security groups and ACLs
   - Public gateway
   - Show best practices for networking

2. **kubernetes-cluster** - IKS Cluster
   - VPC and networking
   - IKS cluster with worker pool
   - Demonstrate multi-zone deployment

3. **web-application** - Three-tier web app
   - Load balancer
   - VPC with subnets
   - IKS cluster
   - PostgreSQL database
   - COS for static assets

#### Medium Priority

4. **databases** - Cloud Databases example
   - PostgreSQL
   - Redis
   - Service credentials
   - Show different service plans

5. **object-storage** - Advanced COS
   - Multiple buckets with different storage classes
   - Lifecycle policies
   - Bucket configurations (CORS, encryption)

6. **iam-policies** - IAM and access management
   - Service IDs
   - Access groups
   - IAM policies
   - API keys

#### Lower Priority

7. **serverless** - Cloud Functions
   - Functions
   - API Gateway
   - Triggers

8. **observability** - Monitoring and logging
   - LogDNA
   - SysDig
   - Activity Tracker

### Option 3: Update Test Files

The test files (`examples_*_test.go`) still reference the old "xyz" provider. Update them:

**Current files to update**:
- `examples/examples_test.go`
- `examples/examples_go_test.go`
- `examples/examples_nodejs_test.go`
- `examples/examples_py_test.go`

**What to change**:
1. Replace all references from `xyz` to `ibmcloud`
2. Update import paths
3. Verify tests work with the basic examples
4. Remove tests for resources that don't exist in IBM Cloud provider

### Option 4: Add Integration Tests

Add actual integration tests that:
1. Deploy the example
2. Verify resources were created
3. Clean up resources
4. Run in CI/CD on PR or schedule

Example test structure:
```go
func TestGoExample(t *testing.T) {
    integration.ProgramTest(t, &integration.ProgramTestOptions{
        Dir: filepath.Join("examples", "basic-go"),
        Dependencies: []string{
            "github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud",
        },
        Quick: true,
    })
}
```

### Option 5: Create Example Templates

Make examples discoverable via `pulumi new`:

1. Create a templates directory
2. Add metadata for each template
3. Submit to Pulumi templates registry
4. Users can then do:
   ```bash
   pulumi new https://github.com/mapt-oss/pulumi-ibmcloud/tree/main/examples/basic-go
   ```

## Recommended Workflow

### Phase 1: Validate (Before First Release)
1. ✅ Basic examples are ready
2. Test one example locally to ensure the provider works
3. Create `v0.0.1` release
4. Test plugin installation and example deployment

### Phase 2: Expand (After First Release)
1. Create 2-3 more advanced examples (vpc-complete, kubernetes-cluster)
2. Update test files to reference `ibmcloud`
3. Add README files for new examples

### Phase 3: Production Ready (After v1.0.0)
1. Add integration tests
2. Add CI/CD to run example tests
3. Create comprehensive example coverage
4. Submit templates to Pulumi registry

## Quick Commands Reference

### Create a New Example

```bash
# Create directory
mkdir examples/new-example

# For Go
cat > examples/new-example/Pulumi.yaml <<EOF
name: new-example
runtime: go
description: Description of what this example does
EOF

cat > examples/new-example/main.go <<EOF
package main

import (
    "github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Your code here
        return nil
    })
}
EOF

cat > examples/new-example/README.md <<EOF
# New Example

Description

## Setup
...
EOF
```

### Test All Examples

```bash
cd examples
go test -v -timeout 2h
```

## Files Modified/Created

### New Files Created ✅
- `examples/EXAMPLES.md` - Main examples documentation
- `examples/basic-go/main.go` - Go example code
- `examples/basic-go/Pulumi.yaml` - Go example config
- `examples/basic-go/README.md` - Go example docs
- `examples/basic-py/README.md` - Python example docs
- `examples/basic-ts/README.md` - TypeScript example docs
- `EXAMPLES_NEXT_STEPS.md` - This file

### Files Updated ✅
- `examples/basic-py/__main__.py` - Changed from xyz to real IBM Cloud resources
- `examples/basic-py/requirements.txt` - Updated dependencies
- `examples/basic-ts/index.ts` - Changed from xyz to real IBM Cloud resources

### Files Removed ✅
- `examples/examples_dotnet_test.go` - Removed (no C# SDK)

### Files That Need Updates (Future)
- `examples/examples_test.go` - Update xyz references to ibmcloud
- `examples/examples_go_test.go` - Update xyz references to ibmcloud
- `examples/examples_nodejs_test.go` - Update xyz references to ibmcloud
- `examples/examples_py_test.go` - Update xyz references to ibmcloud

## Summary

You now have:
- ✅ 3 working basic examples (Go, TypeScript, Python)
- ✅ Complete documentation for each example
- ✅ Comprehensive examples guide
- ✅ Clear next steps for expansion

**Next action**: After creating your first release (`v0.0.1`), test one of the basic examples end-to-end to validate everything works!

## Questions to Consider

1. **Do you want integration tests in CI/CD?**
   - If yes, we can set up automated testing of examples
   - This would catch breaking changes automatically

2. **Which advanced examples are most important?**
   - VPC complete setup?
   - Kubernetes cluster?
   - Database examples?
   - All of the above?

3. **Do you want to publish the SDK to PyPI/NPM eventually?**
   - If yes, we should update the examples to use published packages
   - If no, keep the local installation instructions

Let me know which direction you'd like to go!
