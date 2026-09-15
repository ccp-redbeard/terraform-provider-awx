/*
Use this data source to look up a role on a credential, so it can be granted to a
user or team through role_entitlement.

# Example Usage

```hcl

	data "awx_credential" "mycred" {
	  credential_id = 42
	}

	data "awx_credential_role" "cred_use_role" {
	  name          = "Use"
	  credential_id = data.awx_credential.mycred.id
	}

```
*/
package awx

import (
	"context"
	"strconv"

	awx "github.com/denouche/goawx/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCredentialRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCredentialRoleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"credential_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceCredentialRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	credID := d.Get("credential_id").(int)
	credential, err := client.CredentialsService.GetCredentialsByID(credID, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Credential",
			"Fail to find the credential, got: %s",
			err.Error(),
		)
	}

	// A credential the token cannot read comes back without object roles.
	if credential.SummaryFields == nil || credential.SummaryFields.ObjectRoles == nil {
		return buildDiagnosticsMessage(
			"Get: Fail to fetch Credential Role",
			"Credential %d exposes no object roles; check the token has permission to read it",
			credID,
		)
	}

	roleslist := []*awx.ApplyRole{
		credential.SummaryFields.ObjectRoles.UseRole,
		credential.SummaryFields.ObjectRoles.AdminRole,
		credential.SummaryFields.ObjectRoles.ReadRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range roleslist {
			if v != nil && id == v.ID {
				d = setCredentialRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range roleslist {
			if v != nil && name == v.Name {
				d = setCredentialRoleData(d, v)
				return diags
			}
		}
	}

	return buildDiagnosticsMessage(
		"Failed to fetch credential role - Not Found",
		"The credential role was not found",
	)
}

func setCredentialRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	d.Set("name", r.Name)
	d.SetId(strconv.Itoa(r.ID))
	return d
}
