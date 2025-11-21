/**
 * IBM Cloud Resource Group and VPC example
 * Creates a resource group and a VPC with basic networking
 */

import * as pulumi from "@pulumi/pulumi";
import * as ibmcloud from "@pulumi/ibmcloud";

const stack = pulumi.getStack();

// Create a resource group
const resourceGroup = new ibmcloud.ResourceGroup("example-rg", {
    name: `pulumi-example-rg-${stack}`,
    tags: ["pulumi", "example", "typescript"],
});

// Create a VPC
const vpc = new ibmcloud.IsVpc("example-vpc", {
    name: `pulumi-example-vpc-${stack}`,
    resourceGroup: resourceGroup.id,
    tags: ["pulumi", "example", "networking"],
});

// Create a subnet in the VPC
const subnet = new ibmcloud.IsSubnet("example-subnet", {
    name: `pulumi-example-subnet-${stack}`,
    vpc: vpc.id,
    zone: "us-south-1",
    ipv4CidrBlock: "10.240.0.0/24",
    resourceGroup: resourceGroup.id,
});

// Exports
export const resourceGroupId = resourceGroup.id;
export const resourceGroupName = resourceGroup.name;
export const vpcId = vpc.id;
export const vpcName = vpc.name;
export const subnetId = subnet.id;
export const subnetCidr = subnet.ipv4CidrBlock;
