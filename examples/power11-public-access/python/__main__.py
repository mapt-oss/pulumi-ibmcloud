"""IBM Cloud Power11 Virtual Server with Public Access via Load Balancer"""

import pulumi
import pulumi_ibmcloud as ibmcloud

# Load configuration
config = pulumi.Config()
region = config.get("region") or "us-south"
zone = config.get("zone") or "us-south-1"
power_vs_zone = config.get("powerVsZone") or "us-south"
ssh_public_key = config.require("sshPublicKey")
instance_name = config.get("instanceName") or "power11-vm"

# Step 1: Create Resource Group
resource_group = ibmcloud.ResourceGroup(
    "power11-rg",
    name="power11-public-access-rg",
    tags=[
        "environment:production",
        "managed-by:pulumi",
        "power-architecture:power11",
    ],
)

# Step 2: Create Power VS Workspace for Power11
power_vs_workspace = ibmcloud.ResourceInstance(
    "power11-workspace",
    name="power11-workspace",
    service="power-iaas",
    plan="power-virtual-server-group",
    location=power_vs_zone,
    resource_group_id=resource_group.id,
    tags=["power11", "production"],
)

# Step 3: Create SSH Key
ssh_key = ibmcloud.PiKey(
    "power11-ssh-key",
    pi_key_name="power11-access-key",
    pi_ssh_key=ssh_public_key,
    pi_cloud_instance_id=power_vs_workspace.guid,
)

# Step 4: Create Private Network for Power VS
power_vs_network = ibmcloud.PiNetwork(
    "power11-network",
    pi_network_name="power11-private-net",
    pi_cloud_instance_id=power_vs_workspace.guid,
    pi_network_type="vlan",
    pi_cidr="192.168.50.0/24",
    pi_gateway="192.168.50.1",
    pi_dns_servers=["9.9.9.9", "1.1.1.1"],
)

# Step 5: Create Power11 Virtual Server Instance
# Power11 uses system type s1022 (Power10/11 architecture)
power11_instance = ibmcloud.PiInstance(
    "power11-instance",
    pi_instance_name=instance_name,
    pi_cloud_instance_id=power_vs_workspace.guid,
    # Power11 Specifications
    pi_memory=16.0,  # 16 GB RAM
    pi_processors=2.0,  # 2 cores
    pi_proc_type="dedicated",  # dedicated or shared
    pi_sys_type="s1022",  # Power10/11 system type
    # Operating System - RHEL 9 for Power11
    pi_image_id="rhel-9-2",  # Replace with actual image ID
    # Network Configuration
    pi_networks=[
        ibmcloud.PiInstancePiNetworkArgs(
            network_id=power_vs_network.network_id,
        )
    ],
    # SSH Access
    pi_key_pair_name=ssh_key.pi_key_name,
    # Storage - tier1 (NVMe) for best performance
    pi_storage_type="tier1",
    pi_storage_pool="Tier1-Flash",
    # Health
    pi_health_status="OK",
    opts=pulumi.ResourceOptions(
        custom_timeouts=pulumi.CustomTimeouts(
            create="30m",
            update="30m",
            delete="30m",
        )
    ),
)

# Step 6: Create VPC for Load Balancer
vpc = ibmcloud.IsVpc(
    "power11-vpc",
    name="power11-vpc",
    resource_group=resource_group.id,
    address_prefix_management="auto",
    tags=["power11", "loadbalancer"],
)

# Step 7: Create Subnet for Load Balancer
lb_subnet = ibmcloud.IsSubnet(
    "lb-subnet",
    name="power11-lb-subnet",
    vpc=vpc.id,
    zone=zone,
    ipv4_cidr_block="10.240.64.0/24",
    resource_group=resource_group.id,
)

# Step 8: Create Public Gateway
public_gateway = ibmcloud.IsPublicGateway(
    "public-gateway",
    name="power11-pgw",
    vpc=vpc.id,
    zone=zone,
    resource_group=resource_group.id,
)

# Attach Public Gateway to Subnet
pgw_attachment = ibmcloud.IsSubnetPublicGatewayAttachment(
    "pgw-attachment",
    subnet=lb_subnet.id,
    public_gateway=public_gateway.id,
)

# Step 9: Create Security Group
security_group = ibmcloud.IsSecurityGroup(
    "power11-sg",
    name="power11-security-group",
    vpc=vpc.id,
    resource_group=resource_group.id,
)

# Security Rules - HTTP
http_rule = ibmcloud.IsSecurityGroupRule(
    "allow-http",
    group=security_group.id,
    direction="inbound",
    remote="0.0.0.0/0",
    tcp=ibmcloud.IsSecurityGroupRuleTcpArgs(
        port_min=80,
        port_max=80,
    ),
)

# Security Rules - HTTPS
https_rule = ibmcloud.IsSecurityGroupRule(
    "allow-https",
    group=security_group.id,
    direction="inbound",
    remote="0.0.0.0/0",
    tcp=ibmcloud.IsSecurityGroupRuleTcpArgs(
        port_min=443,
        port_max=443,
    ),
)

# Security Rules - SSH (on custom port)
ssh_rule = ibmcloud.IsSecurityGroupRule(
    "allow-ssh",
    group=security_group.id,
    direction="inbound",
    remote="0.0.0.0/0",
    tcp=ibmcloud.IsSecurityGroupRuleTcpArgs(
        port_min=2222,
        port_max=2222,
    ),
)

# Security Rules - Allow all outbound
outbound_rule = ibmcloud.IsSecurityGroupRule(
    "allow-outbound",
    group=security_group.id,
    direction="outbound",
    remote="0.0.0.0/0",
)

# Step 10: Create Application Load Balancer
load_balancer = ibmcloud.IsLb(
    "power11-lb",
    name="power11-alb",
    subnets=[lb_subnet.id],
    type="public",
    resource_group=resource_group.id,
    security_groups=[security_group.id],
    tags=["power11", "public-access"],
    opts=pulumi.ResourceOptions(
        custom_timeouts=pulumi.CustomTimeouts(
            create="15m",
            update="15m",
            delete="15m",
        )
    ),
)

# Step 11: Create Backend Pool
backend_pool = ibmcloud.IsLbPool(
    "backend-pool",
    lb=load_balancer.id,
    name="power11-backend-pool",
    algorithm="round_robin",
    protocol="http",
    health_delay=10,
    health_retries=3,
    health_timeout=5,
    health_type="http",
    health_monitor_url="/",
    health_monitor_port=80,
)

# Step 12: Add Power11 Instance to Backend Pool
pool_member = ibmcloud.IsLbPoolMember(
    "power11-member",
    lb=load_balancer.id,
    pool=backend_pool.id,
    port=80,
    target_address=power11_instance.pi_networks[0].ip_address,
    weight=100,
)

# Step 13: Create HTTP Listener
http_listener = ibmcloud.IsLbListener(
    "http-listener",
    lb=load_balancer.id,
    default_pool=backend_pool.id,
    port=80,
    protocol="http",
    connection_limit=2000,
)

# Step 14: Create SSH Listener (TCP passthrough on port 2222)
ssh_listener = ibmcloud.IsLbListener(
    "ssh-listener",
    lb=load_balancer.id,
    default_pool=backend_pool.id,
    port=2222,
    protocol="tcp",
    connection_limit=100,
)

# Exports
pulumi.export("resourceGroupId", resource_group.id)
pulumi.export("resourceGroupName", resource_group.name)

pulumi.export("powerVsWorkspaceId", power_vs_workspace.id)
pulumi.export("powerVsWorkspaceGuid", power_vs_workspace.guid)

pulumi.export("powerVsNetworkId", power_vs_network.network_id)
pulumi.export("powerVsNetworkCidr", power_vs_network.pi_cidr)

pulumi.export("power11InstanceId", power11_instance.pi_instance_id)
pulumi.export("power11InstanceName", power11_instance.pi_instance_name)
pulumi.export("power11InstanceStatus", power11_instance.status)
pulumi.export("power11InstanceInternalIp", power11_instance.pi_networks[0].ip_address)
pulumi.export("power11InstanceMacAddress", power11_instance.pi_networks[0].mac_address)
pulumi.export("power11SystemType", power11_instance.pi_sys_type)
pulumi.export("power11Memory", power11_instance.pi_memory)
pulumi.export("power11Processors", power11_instance.pi_processors)

pulumi.export("vpcId", vpc.id)
pulumi.export("vpcCrn", vpc.crn)

pulumi.export("loadBalancerId", load_balancer.id)
pulumi.export("loadBalancerHostname", load_balancer.hostname)
pulumi.export("loadBalancerPublicIps", load_balancer.public_ips)

# Access Information
pulumi.export("httpUrl", load_balancer.hostname.apply(lambda h: f"http://{h}"))
pulumi.export("httpsUrl", load_balancer.hostname.apply(lambda h: f"https://{h}"))
pulumi.export("sshCommand", load_balancer.hostname.apply(lambda h: f"ssh -p 2222 root@{h}"))

# Summary
pulumi.export(
    "deploymentSummary",
    pulumi.Output.all(
        load_balancer.hostname,
        power11_instance.pi_networks[0].ip_address,
        power11_instance.status,
    ).apply(
        lambda args: {
            "message": "Power11 instance deployed successfully!",
            "publicAccess": {
                "http": f"http://{args[0]}",
                "https": f"https://{args[0]}",
                "ssh": f"ssh -p 2222 root@{args[0]}",
            },
            "instance": {
                "internalIp": args[1],
                "status": args[2],
                "systemType": "s1022 (Power10/11)",
            },
        }
    ),
)
