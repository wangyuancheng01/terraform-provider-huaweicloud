package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const DhcpLeaseType = "51"

// API: VPC PUT /v1/{project_id}/ports/{port_id}
// API: VPC POST /v1/{project_id}/ports
// API: VPC GET /v1/{project_id}/ports/{port_id}
// API: VPC DELETE /v1/{project_id}/ports/{port_id}
func ResourceNetworkInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkInterfaceCreate,
		ReadContext:   resourceNetworkInterfaceRead,
		UpdateContext: resourceNetworkInterfaceUpdate,
		DeleteContext: resourceNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"fixed_ip_v4": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"allowed_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dhcp_lease_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// Computed
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_security_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_efi": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipv6_bandwidth_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceNetworkInterfaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC network client: %s", err)
	}

	opts := ports.CreateOpts{
		NetworkId:      d.Get("subnet_id").(string),
		Name:           d.Get("name").(string),
		SecurityGroups: utils.ExpandToStringList(d.Get("security_group_ids").([]interface{})),
	}

	// build fixed_ips
	fixedIps := buildFixedIps(d)
	if fixedIps != nil {
		opts.FixedIps = fixedIps
	}
	// build allowed_address_pairs
	allowedAddrPairs := buildAllowedAddrPairs(d)
	if allowedAddrPairs != nil {
		opts.AllowedAddressPairs = allowedAddrPairs
	}
	// build extra_dhcp_opt
	extraDhcpOpts := buildExtraDhcpOpts(d)
	if extraDhcpOpts != nil {
		opts.ExtraDhcpOpts = extraDhcpOpts
	}

	networkinterface, err := ports.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating VPC network interface: %s", err)
	}
	d.SetId(networkinterface.ID)

	return resourceNetworkInterfaceRead(ctx, d, meta)
}

func buildFixedIps(d *schema.ResourceData) []ports.FixedIp {
	if fixedIpv4, ok := d.GetOk("fixed_ip_v4"); ok {
		fixedIp := ports.FixedIp{
			IpAddress: fixedIpv4.(string),
		}
		fixIpsOpt := make([]ports.FixedIp, 1)
		fixIpsOpt[0] = fixedIp
		return fixIpsOpt
	}
	return nil
}

func buildSecurityGroups(d *schema.ResourceData) []string {
	if securityGroup, ok := d.GetOk("security_group_ids"); ok {
		s := utils.ExpandToStringList(securityGroup.([]interface{}))
		if len(s) > 0 {
			return s
		}
	}
	return nil
}

func buildAllowedAddrPairs(d *schema.ResourceData) []ports.AddressPair {
	if allowedAddress, ok := d.GetOk("allowed_addresses"); ok {
		allowedAddressOpts := make([]ports.AddressPair, len(allowedAddress.([]interface{})))
		for i, allowedAdd := range allowedAddress.([]interface{}) {
			allowedAddressOpts[i] = ports.AddressPair{
				IpAddress: allowedAdd.(string),
			}
		}
		return allowedAddressOpts
	}
	return nil
}

func buildExtraDhcpOpts(d *schema.ResourceData) []ports.ExtraDhcpOpt {
	if extraDhcpOpts, ok := d.GetOk("dhcp_lease_time"); ok && extraDhcpOpts != "" {
		extraDhcpOpt := make([]ports.ExtraDhcpOpt, 1)
		extraDhcpOpt[0] = ports.ExtraDhcpOpt{
			OptName:  DhcpLeaseType,
			OptValue: extraDhcpOpts.(string),
		}
		return extraDhcpOpt
	}
	return nil
}

func resourceNetworkInterfaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC network client: %s", err)
	}
	networkinterface, err := ports.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPC network interface")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", networkinterface.Name),
		d.Set("status", networkinterface.Status),
		d.Set("device_id", networkinterface.DeviceId),
		d.Set("device_owner", networkinterface.DeviceOwner),
		d.Set("mac_address", networkinterface.MacAddress),
		d.Set("subnet_id", networkinterface.NetworkId),
		d.Set("security_group_ids", networkinterface.SecurityGroups),
		d.Set("fixed_ip_v4", flattenFixedIps(networkinterface.FixedIps)),
		d.Set("allowed_addresses", flattenAllowedAddr(networkinterface.AllowedAddressPairs)),
		d.Set("dhcp_lease_time", flattenExtraDhcpOpts(networkinterface.ExtraDhcpOpts)),
		d.Set("instance_id", networkinterface.InstanceId),
		d.Set("instance_type", networkinterface.InstanceType),
		d.Set("availability_zone", networkinterface.ZoneId),
		d.Set("port_security_enabled", networkinterface.PortSecurityEnabled),
		d.Set("enable_efi", networkinterface.EnableEfi),
		d.Set("ipv6_bandwidth_id", networkinterface.Ipv6BandwidthId),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFixedIps(fixedIps []ports.FixedIp) string {
	if len(fixedIps) > 0 {
		fixedIpv4 := fixedIps[0]
		return fixedIpv4.IpAddress
	}
	return ""
}
func flattenAllowedAddr(addressPairs []ports.AddressPair) []interface{} {
	if addressPairs == nil {
		return nil
	}
	addressPairsRes := make([]interface{}, len(addressPairs))
	for i, addressPair := range addressPairs {
		addressPairsRes[i] = addressPair.IpAddress
	}
	return addressPairsRes
}
func flattenExtraDhcpOpts(extraDhcpOpts []ports.ExtraDhcpOpt) string {
	if len(extraDhcpOpts) > 0 {
		extraDhcpLeaseTime := extraDhcpOpts[0]
		return extraDhcpLeaseTime.OptValue
	}
	return ""
}
func resourceNetworkInterfaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC network client: %s", err)
	}

	opts := ports.UpdateOpts{
		Name: d.Get("name").(string),
	}

	if d.HasChange("security_group_ids") {
		secuirtyGroups := buildSecurityGroups(d)
		if len(secuirtyGroups) > 0 {
			opts.SecurityGroups = secuirtyGroups
		}
	}
	if d.HasChange("allowed_addresses") {
		allowedAddrPairs := buildAllowedAddrPairs(d)
		opts.AllowedAddressPairs = allowedAddrPairs
		if allowedAddrPairs == nil {
			opts.AllowedAddressPairs = make([]ports.AddressPair, 0)
		}
	}
	if d.HasChange("dhcp_lease_time") {
		extraDhcpOpts := buildExtraDhcpOpts(d)
		if extraDhcpOpts != nil {
			opts.ExtraDhcpOpts = extraDhcpOpts
		}
	}

	_, err = ports.Update(client, d.Id(), opts)
	if err != nil {
		return diag.Errorf("error updating VPC network interface: %s", err)
	}

	return resourceNetworkInterfaceRead(ctx, d, meta)
}

func resourceNetworkInterfaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC network client: %s", err)
	}
	err = ports.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting VPC network interface: %s", err)
	}

	return nil
}
