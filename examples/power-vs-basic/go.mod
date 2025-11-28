module power-vs-basic

go 1.23

require (
	github.com/mapt-oss/pulumi-ibmcloud/sdk v0.0.0
	github.com/pulumi/pulumi/sdk/v3 v3.207.0
)

// Use local SDK during development
replace github.com/mapt-oss/pulumi-ibmcloud/sdk => ../../sdk
