package sdkv2provider

import (
	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCloudflareRiskBehaviorsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		consts.AccountIDSchemaKey: {
			Description: consts.AccountIDSchemaDescription,
			Type:        schema.TypeString,
			Required:    true,
		},
		"behavior": {
			Description: "Zero Trust risk behaviors configured on this account",
			Type:        schema.TypeSet,
			Required:    true,
			Elem:        behaviorElem,
		},
	}
}

var behaviorElem = &schema.Resource{Schema: map[string]*schema.Schema{
	"enabled": {
		Description: "Whether this risk behavior type is enabled.",
		Type:        schema.TypeBool,
		Required:    true,
	},
	"risk_level": {
		Description:  "Flag controlling if this risk behavior type is enabled.",
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"low", "medium", "high"}, false),
	},
	"name": {
		Description: "Name of this risk behavior type",
		Type:        schema.TypeString,
		Required:    true,
	},
},
}
