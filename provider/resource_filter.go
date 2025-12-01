package ibmcloud

import (
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
)

// FilterUnusedResources removes resources you don't need to reduce SDK size
// This significantly reduces compilation memory usage
func FilterUnusedResources(info *tfbridge.ProviderInfo) {
	// Example: Keep only VPC, Compute, and Storage resources
	// Uncomment categories you NEED, comment out what you DON'T need

	keepPrefixes := map[string]bool{
		// Core Infrastructure (highly recommended)
		"ibm_is_":       true, // VPC Infrastructure (VPC, Subnets, VMs, etc.)
		"ibm_resource_": true, // Resource Groups, Instances

		// Compute
		"ibm_compute_": true, // Virtual Servers (Classic)
		"ibm_pi_":      true, // Power Systems

		// Storage
		"ibm_cos_":     true,  // Cloud Object Storage
		"ibm_storage_": false, // Classic Storage (disable if not needed)

		// Kubernetes & Containers
		"ibm_container_": true,  // Kubernetes Service
		"ibm_cr_":        false, // Container Registry (disable if not needed)
		"ibm_ob_":        false, // OpenShift on IBM Cloud (disable if not needed)

		// Database
		"ibm_database_": false, // Cloud Databases (disable if not needed)
		"ibm_cd_":       false, // Continuous Delivery (disable if not needed)

		// Networking
		"ibm_dns_": true,  // DNS Services
		"ibm_tg_":  false, // Transit Gateway (disable if not needed)
		"ibm_dl_":  false, // Direct Link (disable if not needed)

		// Security & IAM
		"ibm_iam_": true,  // Identity & Access Management
		"ibm_kms_": false, // Key Protect (disable if not needed)
		"ibm_sm_":  false, // Secrets Manager (disable if not needed)

		// AI & Watson
		"ibm_watson_": false, // Watson services (disable if not needed)

		// VMware
		"ibm_vmaas_": false, // VMware (disable if not needed)

		// Satellite
		"ibm_satellite_": false, // Satellite (disable if not needed)

		// Event Streams & Messaging
		"ibm_event_": false, // Event Streams/Kafka (disable if not needed)
		"ibm_mq_":    false, // MQ (disable if not needed)

		// Backup & Recovery
		"ibm_backup_recovery_": false, // Backup Recovery (this has HUGE types!)
	}

	// Filter resources
	for key := range info.Resources {
		keep := false
		for prefix, enabled := range keepPrefixes {
			if enabled && len(key) >= len(prefix) && key[:len(prefix)] == prefix {
				keep = true
				break
			}
		}
		if !keep {
			delete(info.Resources, key)
		}
	}

	// Filter data sources
	for key := range info.DataSources {
		keep := false
		for prefix, enabled := range keepPrefixes {
			if enabled && len(key) >= len(prefix) && key[:len(prefix)] == prefix {
				keep = true
				break
			}
		}
		if !keep {
			delete(info.DataSources, key)
		}
	}
}
