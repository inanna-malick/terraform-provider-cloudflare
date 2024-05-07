package sdkv2provider

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cloudflare/cloudflare-go"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func resourceCloudflareRiskBehavior() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceCloudflareRiskBehaviorsSchema(),
		CreateContext: resourceCloudflareRiskBehaviorUpdate,
		ReadContext:   resourceCloudflareRiskBehaviorRead,
		UpdateContext: resourceCloudflareRiskBehaviorUpdate,
		DeleteContext: resourceCloudflareRiskBehaviorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudflareRiskBehaviorImport,
		},
		Description: heredoc.Doc(`
			Used to manage configuration for risk behaviors.
		`),
	}
}

func resourceCloudflareRiskBehaviorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	accountID := d.Get(consts.AccountIDSchemaKey).(string)

	behaviors, err := client.Behaviors(ctx, accountID)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to get risk behaviors"))
	}

	d.SetId(accountID)

	behaviorsSet := []interface{}{}

	for k, b := range behaviors.Behaviors {
		behavior :=
			map[string]interface{}{
				"name":       k,
				"enabled":    b.Enabled,
				"risk_level": fmt.Sprint(b.RiskLevel),
			}

		behaviorsSet = append(behaviorsSet, behavior)
	}

	d.SetId(accountID)
	d.Set("behavior", schema.NewSet(HashByMapKey("name"), behaviorsSet))
	return nil
}

func resourceCloudflareRiskBehaviorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	accountID := d.Get(consts.AccountIDSchemaKey).(string)

	behaviors := d.Get("behavior").(*schema.Set).List()
	tflog.Debug(ctx, fmt.Sprintf("Setting zero trust risk behaviors to: %v", behaviors))

	behaviorsMap := map[string]cloudflare.Behavior{}
	for _, b := range behaviors {
		b := b.(map[string]interface{})

		riskLevel, err := cloudflare.RiskLevelFromString(b["risk_level"].(string))
		if err != nil {
			return diag.FromErr(errors.Wrap(err, "failed to get risk behaviors"))
		}

		enabled := b["enabled"].(bool)

		behavior := cloudflare.Behavior{
			Enabled:   &enabled,
			RiskLevel: *riskLevel,
		}

		behaviorsMap[b["name"].(string)] = behavior
	}

	_, err := client.UpdateBehaviors(ctx, accountID, cloudflare.Behaviors{Behaviors: behaviorsMap})
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to update Risk Behavior values"))
	}

	return resourceCloudflareRiskBehaviorRead(ctx, d, meta)
}

func resourceCloudflareRiskBehaviorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	accountID := d.Get(consts.AccountIDSchemaKey).(string)

	tflog.Debug(ctx, "Resetting all zero trust risk behaviors to enabled: false, risk_level: low")

	behaviors, err := client.Behaviors(ctx, accountID)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to get risk behaviors"))
	}

	// set all risk behavior values to false/low before running update
	for _, behavior := range behaviors.Behaviors {
		behavior.Enabled = cloudflare.BoolPtr(false)
		behavior.RiskLevel = cloudflare.Low
	}

	_, err = client.UpdateBehaviors(ctx, accountID, behaviors)
	if err != nil {
		return diag.FromErr(errors.Wrap(err, "failed to reset Regional Tiered Cache value"))
	}

	return nil
}

func resourceCloudflareRiskBehaviorImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set(consts.ZoneIDSchemaKey, d.Id())

	resourceCloudflareRiskBehaviorRead(ctx, d, meta)
	return []*schema.ResourceData{d}, nil
}
