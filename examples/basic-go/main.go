// IBM Cloud Resource Group example
// Creates a simple resource group to demonstrate the provider

package main

import (
	"fmt"

	"github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a resource group
		rg, err := ibmcloud.NewResourceGroup(ctx, "example-rg", &ibmcloud.ResourceGroupArgs{
			Name: pulumi.String(fmt.Sprintf("pulumi-example-rg-%s", ctx.Stack())),
			Tags: pulumi.StringArray{
				pulumi.String("pulumi"),
				pulumi.String("example"),
				pulumi.String("go"),
			},
		})
		if err != nil {
			return err
		}

		// Export the resource group ID and name
		ctx.Export("resourceGroupId", rg.ID())
		ctx.Export("resourceGroupName", rg.Name)

		return nil
	})
}
