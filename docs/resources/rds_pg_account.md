---
subcategory: "Relational Database Service (RDS)"
---

# huaweicloud_rds_pg_account

Manages RDS PostgreSQL account resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "account_password" {}

resource "huaweicloud_rds_pg_account" "test" {
  instance_id = var.instance_id
  name        = "test_account_name"
  password    = var.account_password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RDS PostgreSQL instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the username of the DB account. The username contains 1 to 63
  characters, including letters, digits, and underscores (_). It cannot start with pg or a digit and must be different
  from system usernames. System users include **rdsAdmin**, **rdsMetric**, **rdsBackup**, **rdsRepl**, **rdsProxy**,
  and **rdsDdm**.

  Changing this parameter will create a new resource.

* `password` - (Required, String) Specifies the password of the DB account. The value must be 8 to 32 characters long
  and contain at least three types of the following characters: uppercase letters, lowercase letters, digits, and special
  characters (~!@#%^*-_=+?,). The value cannot contain the username or the username spelled backwards.

* `description` - (Optional, String) Specifies the remarks of the DB account. The parameter must be 1 to 512 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of account which is formatted `<instance_id>/<name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS PostgreSQL account can be imported using the `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_pg_account.test <instance_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password`. It is generally recommended
running `terraform plan` after importing the RDS PostgreSQL account. You can then decide if changes should be applied to
the RDS PostgreSQL account, or the resource definition should be updated to align with the RDS PostgreSQL account. Also
you can ignore changes as below.

```hcl
resource "huaweicloud_rds_pg_account" "account_1" {
    ...

  lifecycle {
    ignore_changes = [
      password
    ]
  }
}
```
