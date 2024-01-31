package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/myklst/terraform-provider-st-cdnetworks/cdnetworks/utils"
	"github.com/myklst/terraform-provider-st-cdnetworks/cdnetworksapi"
)

const (
	API_VERSION = "1.0.0"
)

var DomainSchema = schema.Schema{
	Description: "This resource provides the configuration of acceleration domain",
	Attributes: map[string]schema.Attribute{
		"domain_id": &schema.StringAttribute{
			Description: "Id of acceleration domain, generated by cdnetworks.",
			Computed:    true,
		},
		"domain": &schema.StringAttribute{
			Description: "CDN accelerated domain name.",
			Required:    true,
		},
		"cname": &schema.StringAttribute{
			Description: "Cname",
			Computed:    true,
		},
		"comment": &schema.StringAttribute{
			Description: "Remarks, up to 1000 characters",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
		},
		"status": &schema.StringAttribute{
			Description: "The deployment status of the accelerate domain name. Deployed indicates that the accelerated domain name configuration is complete. InProgress indicates that the deployment task of the accelerated domain name configuration is still in progress, and may be in queue, deployed, or failed.",
			Computed:    true,
		},
		"contract_id": &schema.StringAttribute{
			Description: "The id of contract",
			Computed:    true,
		},
		"item_id": &schema.StringAttribute{
			Description: "The id of item",
			Computed:    true,
		},
		"enabled": &schema.BoolAttribute{
			Description: "Speed up the activation of the domain name. This is false when the accelerated domain name service is disabled; true when the accelerated domain name service is enabled.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(true),
		},
		"header_of_clientip": &schema.StringAttribute{
			Description: "Pass the response header of client IP. The optional values are Cdn-Src-Ip and X-Forwarded-For. The default value is Cdn-Src-Ip.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("Cdn-Src-Ip"),
		},
		"cdn_service_status": &schema.StringAttribute{
			Description: "Accelerate the CDN service status of the domain name, true means to enable the CDN acceleration service; false means to cancel the CDN acceleration service.",
			Computed:    true,
		},
		"service_type": &schema.StringAttribute{
			Description: "Accelerated domain name service types, including the following: 1028 : Content Acceleration; 1115 : Dynamic Web Acceleration; 1369 : Media Acceleration - RTMP 1391 : Download Acceleration 1348 : Media Acceleration Live Broadcast 1551 : Flood Shield",
			Computed:    true,
		},
		"origin_config": &schema.SingleNestedAttribute{
			Description: "Back-to-origin policy setting, which is used to set the origin site information and the back-to-origin policy of the accelerated domain name",
			Attributes: map[string]schema.Attribute{
				"origin_ips": schema.ListAttribute{
					ElementType: types.StringType,
					Description: `Origin site address, which can be an IP or a domain name.
						1. Only one domain name can be entered. IP and domain names cannot be entered at the same time.
						2. Maximum character limit is 500.`,
					Required: true,
				},
				"default_origin_host_header": schema.StringAttribute{
					Description: `The Origin HOST for changing the HOST field in the return source HTTP request header.
						Note: It should be domain or IP format. For domain name format, each segement separated by a dot, does not exceed 62 characters, the total length should not exceed 128 characters.`,
					Optional: true,
					Computed: true,
				},
			},
			Required: true,
		},
	},
}

var originConfigModelAttributeTypes = map[string]attr.Type{
	"origin_ips":                 types.ListType{}.WithElementType(types.StringType),
	"default_origin_host_header": types.StringType,
}

type DomainResourceModel struct {
	DomainId         types.String `tfsdk:"domain_id"`
	Domain           types.String `tfsdk:"domain"`
	Cname            types.String `tfsdk:"cname"`
	Comment          types.String `tfsdk:"comment"`
	Status           types.String `tfsdk:"status"`
	ContractId       types.String `tfsdk:"contract_id"`
	ItemId           types.String `tfsdk:"item_id"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	HeaderOfClientIp types.String `tfsdk:"header_of_clientip"`
	CdnServiceStatus types.String `tfsdk:"cdn_service_status"`
	ServiceType      types.String `tfsdk:"service_type"`
	OriginConfig     types.Object `tfsdk:"origin_config"`
}

func (model *DomainResourceModel) BuildApiOriginConfig() *cdnetworksapi.OriginConfig {
	config := &cdnetworksapi.OriginConfig{}
	for k, v := range model.OriginConfig.Attributes() {
		if k == "origin_ips" {
			list := make([]string, 0)
			v.(types.List).ElementsAs(nil, &list, false)
			s := strings.Join(list, utils.Separator)
			config.OriginIps = &s
		} else if k == "default_origin_host_header" {
			config.DefaultOriginHostHeader = v.(types.String).ValueStringPointer()
		}
	}

	return config
}

func (model *DomainResourceModel) UpdateDomainFromApiConfig(ctx context.Context, config *cdnetworksapi.QueryCdnDomainResponse) {
	model.DomainId = types.StringPointerValue(config.DomainId)
	model.Domain = types.StringPointerValue(config.DomainName)
	model.Comment = types.StringPointerValue(config.Comment)
	model.Cname = types.StringPointerValue(config.Cname)
	model.Status = types.StringPointerValue(config.Status)
	model.ServiceType = types.StringPointerValue(config.ServiceType)
	model.ContractId = types.StringPointerValue(config.ContractId)
	model.ItemId = types.StringPointerValue(config.ItemId)
	model.CdnServiceStatus = types.StringPointerValue(config.CdnServiceStatus)
	model.HeaderOfClientIp = types.StringPointerValue(config.HeaderOfClientIp)
	model.Enabled = types.BoolPointerValue(config.Enabled)
	if config.OriginConfig != nil {
		defaultOriginHeader := model.Domain
		if config.OriginConfig.DefaultOriginHostHeader != nil {
			defaultOriginHeader = types.StringPointerValue(config.OriginConfig.DefaultOriginHostHeader)
		}
		iplistModel, _ := types.ListValueFrom(ctx, types.StringType, strings.Split(*config.OriginConfig.OriginIps, utils.Separator))
		model.OriginConfig = types.ObjectValueMust(originConfigModelAttributeTypes, map[string]attr.Value{
			"origin_ips":                 iplistModel,
			"default_origin_host_header": defaultOriginHeader,
		})
	}
}

func (model *DomainResourceModel) CopyComputedFields(src *cdnetworksapi.QueryCdnDomainResponse) {
	if src == nil {
		return
	}
	model.Cname = types.StringPointerValue(src.Cname)
	model.ContractId = types.StringPointerValue(src.ContractId)
	model.ItemId = types.StringPointerValue(src.ItemId)
	model.Status = types.StringPointerValue(src.Status)
	model.ServiceType = types.StringPointerValue(src.ServiceType)
	model.CdnServiceStatus = types.StringPointerValue(src.CdnServiceStatus)
}

func (model *DomainResourceModel) Check() error {
	iplist := model.OriginConfig.Attributes()["origin_ips"].(types.List)
	if len(iplist.Elements()) > 15 {
		return fmt.Errorf("The number of IPs cannot exceed 15.")
	}
	return nil
}

func (model *DomainResourceModel) Fill() {
	defaulOriginHostHeader := model.OriginConfig.Attributes()["default_origin_host_header"]
	if defaulOriginHostHeader.IsUnknown() {
		ips := model.OriginConfig.Attributes()["origin_ips"]
		model.OriginConfig = types.ObjectValueMust(originConfigModelAttributeTypes, map[string]attr.Value{
			"origin_ips":                 ips,
			"default_origin_host_header": model.Domain,
		})
	}
}