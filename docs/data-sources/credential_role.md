---
layout: "awx"
page_title: "AWX: awx_credential_role"
sidebar_current: "docs-awx-datasource-credential_role"
description: |-
  Use this data source to look up a role on a credential, so it can be granted to a user or team through role_entitlement.
---

# awx_credential_role

Use this data source to look up a role on a credential, so it can be granted to a
user or team through role_entitlement.

## Example Usage

```hcl
data "awx_credential" "mycred" {
  credential_id = 42
}

data "awx_credential_role" "cred_use_role" {
  name          = "Use"
  credential_id = data.awx_credential.mycred.id
}
```

## Argument Reference

The following arguments are supported:

* `credential_id` - (Required) ID of the credential to reference for credential roles
* `id` - (Optional)
* `name` - (Optional) One of `Use`, `Admin` or `Read`
