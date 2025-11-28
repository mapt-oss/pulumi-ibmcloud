import * as pulumi from "@pulumi/pulumi";
import * as ibmcloud from "@pulumi/ibmcloud";

// Load configuration
const config = new pulumi.Config();
const region = config.get("region") || "us-south";
const zone = config.get("zone") || "us-south-1";
const powerVsZone = config.get("powerVsZone") || "us-south"; // Power11 zones: us-south, us-east, eu-de, etc.
const sshPublicKey = config.require("sshPublicKey");
const instanceName = config.get("instanceName") || "power11-vm";

// Step 1: Create Resource Group
const resourceGroup = new ibmcloud.ResourceGroup("power11-rg", {
    name: "power11-public-access-rg",
    tags: [
        "environment:production",
        "managed-by:pulumi",
        "power-architecture:power11",
    ],
});

// Step 2: Create Power VS Workspace for Power11
const powerVsWorkspace = new ibmcloud.ResourceInstance("power11-workspace", {
    name: "power11-workspace",
    service: "power-iaas",
    plan: "power-virtual-server-group",
    location: powerVsZone,
    resourceGroupId: resourceGroup.id,
    tags: [
        "power11",
        "production",
    ],
});

// Step 3: Create SSH Key
const sshKey = new ibmcloud.PiKey("power11-ssh-key", {
    piKeyName: "power11-access-key",
    piSshKey: sshPublicKey,
    piCloudInstanceId: powerVsWorkspace.guid,
});

// Step 4: Create Private Network for Power VS
const powerVsNetwork = new ibmcloud.PiNetwork("power11-network", {
    piNetworkName: "power11-private-net",
    piCloudInstanceId: powerVsWorkspace.guid,
    piNetworkType: "vlan",
    piCidr: "192.168.50.0/24",
    piGateway: "192.168.50.1",
    piDnsServers: [
        "9.9.9.9",    // Quad9 DNS
        "1.1.1.1",    // Cloudflare DNS
    ],
});

// Step 5: Create Power11 Virtual Server Instance
// Power11 uses system type s1022 (Power10/11 architecture)
const power11Instance = new ibmcloud.PiInstance("power11-instance", {
    piInstanceName: instanceName,
    piCloudInstanceId: powerVsWorkspace.guid,

    // Power11 Specifications
    piMemory: 16,              // 16 GB RAM (adjust as needed)
    piProcessors: 2,           // 2 cores (adjust as needed)
    piProcType: "dedicated",   // Use "dedicated" for production, "shared" for dev/test
    piSysType: "s1022",        // Power10/11 system type (alternatives: e1080)

    // Operating System - Use RHEL 9 or AIX 7.3 for Power11
    piImageId: "rhel-9-2",     // Replace with actual image ID from your region

    // Network Configuration
    piNetworks: [{
        networkId: powerVsNetwork.networkId,
    }],

    // SSH Access
    piKeyPairName: sshKey.piKeyName,

    // Storage - tier1 (NVMe) recommended for Power11 performance
    piStorageType: "tier1",
    piStoragePool: "Tier1-Flash",

    // Health
    piHealthStatus: "OK",
}, {
    // Power11 instances can take 15-20 minutes to provision
    customTimeouts: {
        create: "30m",
        update: "30m",
        delete: "30m",
    },
});

// Step 6: Create VPC for Load Balancer
const vpc = new ibmcloud.IsVpc("power11-vpc", {
    name: "power11-vpc",
    resourceGroup: resourceGroup.id,
    addressPrefixManagement: "auto",
    tags: [
        "power11",
        "loadbalancer",
    ],
});

// Step 7: Create Subnet for Load Balancer
const lbSubnet = new ibmcloud.IsSubnet("lb-subnet", {
    name: "power11-lb-subnet",
    vpc: vpc.id,
    zone: zone,
    ipv4CidrBlock: "10.240.64.0/24",
    resourceGroup: resourceGroup.id,
    publicGateway: undefined, // Will attach separately
});

// Step 8: Create Public Gateway for Internet Access
const publicGateway = new ibmcloud.IsPublicGateway("public-gateway", {
    name: "power11-pgw",
    vpc: vpc.id,
    zone: zone,
    resourceGroup: resourceGroup.id,
});

// Attach Public Gateway to Subnet
const pgwAttachment = new ibmcloud.IsSubnetPublicGatewayAttachment("pgw-attachment", {
    subnet: lbSubnet.id,
    publicGateway: publicGateway.id,
});

// Step 9: Create Security Group
const securityGroup = new ibmcloud.IsSecurityGroup("power11-sg", {
    name: "power11-security-group",
    vpc: vpc.id,
    resourceGroup: resourceGroup.id,
});

// Security Rules - HTTP
const httpRule = new ibmcloud.IsSecurityGroupRule("allow-http", {
    group: securityGroup.id,
    direction: "inbound",
    remote: "0.0.0.0/0",
    tcp: {
        portMin: 80,
        portMax: 80,
    },
});

// Security Rules - HTTPS
const httpsRule = new ibmcloud.IsSecurityGroupRule("allow-https", {
    group: securityGroup.id,
    direction: "inbound",
    remote: "0.0.0.0/0",
    tcp: {
        portMin: 443,
        portMax: 443,
    },
});

// Security Rules - SSH
const sshRule = new ibmcloud.IsSecurityGroupRule("allow-ssh", {
    group: securityGroup.id,
    direction: "inbound",
    remote: "0.0.0.0/0",
    tcp: {
        portMin: 22,
        portMax: 22,
    },
});

// Security Rules - Allow all outbound
const outboundRule = new ibmcloud.IsSecurityGroupRule("allow-outbound", {
    group: securityGroup.id,
    direction: "outbound",
    remote: "0.0.0.0/0",
});

// Step 10: Create Application Load Balancer
const loadBalancer = new ibmcloud.IsLb("power11-lb", {
    name: "power11-alb",
    subnets: [lbSubnet.id],
    type: "public",
    resourceGroup: resourceGroup.id,
    securityGroups: [securityGroup.id],
    tags: [
        "power11",
        "public-access",
    ],
}, {
    customTimeouts: {
        create: "15m",
        update: "15m",
        delete: "15m",
    },
});

// Step 11: Create Backend Pool
const backendPool = new ibmcloud.IsLbPool("backend-pool", {
    lb: loadBalancer.id,
    name: "power11-backend-pool",
    algorithm: "round_robin",
    protocol: "http",
    healthDelay: 10,
    healthRetries: 3,
    healthTimeout: 5,
    healthType: "http",
    healthMonitorUrl: "/",
    healthMonitorPort: 80,
});

// Step 12: Add Power11 Instance to Backend Pool
const poolMember = new ibmcloud.IsLbPoolMember("power11-member", {
    lb: loadBalancer.id,
    pool: backendPool.id,
    port: 80,
    targetAddress: power11Instance.piNetworks[0].ipAddress,
    weight: 100,
});

// Step 13: Create HTTP Listener
const httpListener = new ibmcloud.IsLbListener("http-listener", {
    lb: loadBalancer.id,
    defaultPool: backendPool.id,
    port: 80,
    protocol: "http",
    connectionLimit: 2000,
});

// Step 14: Create HTTPS Listener (Optional - requires certificate)
// Uncomment and configure if you have a certificate in IBM Certificate Manager
/*
const httpsListener = new ibmcloud.IsLbListener("https-listener", {
    lb: loadBalancer.id,
    defaultPool: backendPool.id,
    port: 443,
    protocol: "https",
    connectionLimit: 2000,
    certificateInstance: "crn:v1:bluemix:public:cloudcerts:...", // Your cert CRN
});
*/

// Step 15: Create SSH Listener for Direct SSH Access via Load Balancer
const sshListener = new ibmcloud.IsLbListener("ssh-listener", {
    lb: loadBalancer.id,
    defaultPool: backendPool.id,
    port: 2222,  // Use non-standard port to avoid conflicts
    protocol: "tcp",
    connectionLimit: 100,
});

// Exports
export const resourceGroupId = resourceGroup.id;
export const resourceGroupName = resourceGroup.name;

export const powerVsWorkspaceId = powerVsWorkspace.id;
export const powerVsWorkspaceGuid = powerVsWorkspace.guid;

export const powerVsNetworkId = powerVsNetwork.networkId;
export const powerVsNetworkCidr = powerVsNetwork.piCidr;

export const power11InstanceId = power11Instance.piInstanceId;
export const power11InstanceName = power11Instance.piInstanceName;
export const power11InstanceStatus = power11Instance.status;
export const power11InstanceInternalIp = power11Instance.piNetworks[0].ipAddress;
export const power11InstanceMacAddress = power11Instance.piNetworks[0].macAddress;
export const power11SystemType = power11Instance.piSysType;
export const power11Memory = power11Instance.piMemory;
export const power11Processors = power11Instance.piProcessors;

export const vpcId = vpc.id;
export const vpcCrn = vpc.crn;

export const loadBalancerId = loadBalancer.id;
export const loadBalancerHostname = loadBalancer.hostname;
export const loadBalancerPublicIps = loadBalancer.publicIps;
export const loadBalancerPrivateIps = loadBalancer.privateIps;

// Access Information
export const httpUrl = pulumi.interpolate`http://${loadBalancer.hostname}`;
export const httpsUrl = pulumi.interpolate`https://${loadBalancer.hostname}`;
export const sshCommand = pulumi.interpolate`ssh -p 2222 root@${loadBalancer.hostname}`;

// Summary
export const deploymentSummary = pulumi.all([
    loadBalancer.hostname,
    power11Instance.piNetworks[0].ipAddress,
    power11Instance.status,
]).apply(([hostname, ip, status]) => ({
    message: "Power11 instance deployed successfully!",
    publicAccess: {
        http: `http://${hostname}`,
        https: `https://${hostname}`,
        ssh: `ssh -p 2222 root@${hostname}`,
    },
    instance: {
        internalIp: ip,
        status: status,
        systemType: "s1022 (Power10/11)",
    },
}));
