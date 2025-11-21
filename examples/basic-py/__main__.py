# Copyright 2024, Pulumi Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
IBM Cloud Resource Group and Cloud Object Storage example
Creates a resource group and a Cloud Object Storage instance
"""

import pulumi
import pulumi_ibmcloud as ibmcloud

# Create a resource group
resource_group = ibmcloud.ResourceGroup(
    "example-rg",
    name=f"pulumi-example-rg-{pulumi.get_stack()}",
    tags=["pulumi", "example", "python"]
)

# Create a Cloud Object Storage instance
cos_instance = ibmcloud.ResourceInstance(
    "example-cos",
    name=f"pulumi-example-cos-{pulumi.get_stack()}",
    service="cloud-object-storage",
    plan="standard",
    location="global",
    resource_group_id=resource_group.id,
    tags=["pulumi", "example", "storage"]
)

# Export the resource IDs
pulumi.export("resource_group_id", resource_group.id)
pulumi.export("resource_group_name", resource_group.name)
pulumi.export("cos_instance_id", cos_instance.id)
pulumi.export("cos_instance_crn", cos_instance.crn)
