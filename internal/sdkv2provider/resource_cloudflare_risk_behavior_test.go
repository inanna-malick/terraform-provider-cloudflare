package sdkv2provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudflareRiskBehaviors_Basic(t *testing.T) {
	t.Parallel()
	// accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := generateRandomResourceName()
	name := fmt.Sprintf("cloudflare_risk_behavior.%s", rnd)
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckEmail(t)
			testAccPreCheckApiKey(t)
			testAccPreCheckAccount(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareRiskBehaviors(rnd, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, consts.AccountIDSchemaKey, accountID),
					resource.TestCheckResourceAttr(name, "behavior.#", "2"),
				),
			},
		},
	})
}

// NOTE: seems to require all valid names to be set (eg missing high_dlp causes an error)
func testAccCloudflareRiskBehaviors(name, accountId string) string {
	return fmt.Sprintf(`
	resource cloudflare_risk_behavior %s {
		account_id = "%s"
		behavior {
			name = "imp_travel"
			enabled = true
			risk_level = "medium"
		}
		behavior {
			name = "high_dlp"
			enabled = true
			risk_level = "medium"
		}
	}`, name, accountId)
}
