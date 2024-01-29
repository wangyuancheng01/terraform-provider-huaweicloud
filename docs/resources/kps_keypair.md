---
subcategory: "Data Encryption Workshop (DEW)"
---

# huaweicloud_kps_keypair

Manages a keypair resource within HuaweiCloud.

By default, key pairs use the SSH-2 (RSA, 2048) algorithm for encryption and decryption.

Keys imported support the following cryptographic algorithms:

 * RSA-1024
 * RSA-2048
 * RSA-4096

## Example Usage

### Create a new keypair and export private key to current folder

```hcl
resource "huaweicloud_kps_keypair" "test-keypair" {
  name     = "my-keypair"
  key_file = "private_key.pem"
}
```

### Create a new keypair which scope is Tenant-level and the private key is managed by HuaweiCloud

```hcl
resource "huaweicloud_kms_key" "test" {
  key_alias = "kms_test"
}

resource "huaweicloud_kps_keypair" "test-keypair" {
  name            = "my-keypair"
  scope           = "account"
  encryption_type = "kms"
  kms_key_name    = huaweicloud_kms_key.test.key_alias
}
```

### Import an existing keypair

```hcl
resource "huaweicloud_kps_keypair" "test-keypair" {
  name       = "my-keypair"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAlJq5Pu+eizhou7nFFDxXofr2ySF8k/yuA9OnJdVF9Fbf85Z59CWNZBvcAT... root@terra-dev"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the keypair resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies a unique name for the keypair. The name can contain a maximum of 64
 characters, including letters, digits, underscores (_) and hyphens (-).
 Changing this parameter will create a new resource.

* `scope` - (Optional, String, ForceNew) Specifies the scope of key pair. The options are as follows:
  - **account**: Tenant-level, available to all users under the same account.
  - **user**: User-level, only available to that user.

 The default value is `user`.
 Changing this parameter will create a new resource.

* `encryption_type` - (Optional, String, ForceNew) Specifies encryption mode if manages the private key by HuaweiCloud.
 The options are as follows:
  - **default**: The default encryption mode. Applicable to sites where KMS is not deployed.
  - **kms**: KMS encryption mode.

 Changing this parameter will create a new resource.

* `kms_key_name` - (Optional, String, ForceNew) Specifies the KMS key name to encrypt private keys.
 It's mandatory when the `encryption_type` is `kms`. Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of key pair.

* `public_key` - (Optional, String, ForceNew) Specifies the imported OpenSSH-formatted public key.
 Changing this parameter will create a new resource.

* `key_file` - (Optional, String, ForceNew) Specifies the path of the created private key.
 The private key file (**.pem**) is created only after the resource is created.
 Changing this parameter will create a new resource.

  ~>**NOTE:** If the private key file already exists, it will be overwritten after a new keypair is created.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to name.

* `created_at` - The key pair creation time.

* `fingerprint` - Fingerprint information about an key pair.

* `is_managed` - Whether the private key is managed by HuaweiCloud.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

Keypairs can be imported using the `name`, e.g.

```
$ terraform import huaweicloud_kps_keypair.my-keypair test-keypair
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `encryption_type`,
and `kms_key_name`. It is generally recommended running `terraform plan` after importing a key pair.
You can then decide if changes should be applied to the key pair, or the resource definition
should be updated to align with the key pair. Also you can ignore changes as below.

```
resource "huaweicloud_kps_keypair" "test" {
    ...

  lifecycle {
    ignore_changes = [
      encryption_type, kms_key_name,
    ]
  }
}
```
