---
subcategory: "Document Database Service (DDS)"
---

# huaweicloud_dds_instance

Manages dds instance resource within HuaweiCloud.

## Example Usage: Creating a Cluster Community Edition

```hcl
variable "dds_password" {}

resource "huaweicloud_dds_instance" "instance" {
  name = "dds-instance"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  availability_zone = "{{ availability_zone }}"
  vpc_id            = "{{ vpc_id }}"
  subnet_id         = "{{ subnet_network_id }}}"
  security_group_id = "{{ security_group_id }}"
  password          = var.dds_password
  mode              = "Sharding"
  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c3.medium.4.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c3.medium.4.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c3.large.2.config"
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }
}
```

## Example Usage: Creating a Replica Set Community Edition

```hcl
variable "dds_password" {}

resource "huaweicloud_dds_instance" "instance" {
  name = "dds-instance"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  availability_zone = "{{ availability_zone }}"
  vpc_id            = "{{ vpc_id }}"
  subnet_id         = "{{ subnet_network_id }}}"
  security_group_id = "{{ security_group_id }}"
  password          = var.dds_password
  mode              = "ReplicaSet"
  flavor {
    type      = "replica"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.c3.medium.4.repset"
  }
}
```

## Example Usage: Creating a Single Community Edition

```hcl
variable "dds_password" {}

resource "huaweicloud_dds_instance" "instance" {
  name = "dds-instance"
  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }
  availability_zone = "{{ availability_zone }}"
  vpc_id            = "{{ vpc_id }}"
  subnet_id         = "{{ subnet_network_id }}}"
  security_group_id = "{{ security_group_id }}"
  password          = var.dds_password
  mode              = "Single"
  
  flavor {
    type      = "single"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.s6.large.2.single"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region of the DDS instance. Changing this creates a new
  instance.

* `name` - (Required, String) Specifies the DB instance name. The DB instance name of the same type is unique in the
  same tenant.

* `datastore` - (Required, List, ForceNew) Specifies database information. The structure is described below. Changing
  this creates a new instance.

* `availability_zone` - (Required, String, ForceNew) Specifies the ID of the availability zone. Changing this creates a
  new instance.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this creates a new instance.

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet Network ID. Changing this creates a new instance.

* `security_group_id` - (Required, String) Specifies the security group ID of the DDS instance.

* `password` - (Optional, String) Specifies the Administrator password of the database instance.

* `disk_encryption_id` - (Optional, String, ForceNew) Specifies the disk encryption ID of the instance. Changing this
  creates a new instance.

* `mode` - (Required, String, ForceNew) Specifies the mode of the database instance. **Sharding**, **ReplicaSet**
  are supported. Changing this creates a new instance.

* `configuration` - (Optional, List, ForceNew) Specifies the configuration information.
  The structure is described below. Changing this creates a new instance.

* `flavor` - (Required, List, ForceNew) Specifies the flavors information. The structure is described below. Changing
  this creates a new instance.

* `port` - (Optional, Int) Specifies the database access port. The valid values are range from `2100` to `9500` and
  `27017`, `27018`, `27019`. Defaults to `8635`.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. The structure is described below.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the dds instance.

* `ssl` - (Optional, Bool) Specifies whether to enable or disable SSL. Defaults to true.

**NOTE:** The instance will be restarted in the background when switching SSL. Please operate with caution.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the instance.
  The valid values are as follows:
  + `prePaid`: indicates the yearly/monthly billing mode.
  + `postPaid`: indicates the pay-per-use billing mode.

  Default value is `postPaid`.
  Changing this creates a new instance.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the instance.
  Valid values are *month* and *year*. This parameter is mandatory if `charging_mode` is set to *prePaid*.
  Changing this creates a new instance.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the instance.
  If `period_unit` is set to *month*, the value ranges from 1 to 9.
  If `period_unit` is set to *year*, the value ranges from 1 to 3.
  This parameter is mandatory if `charging_mode` is set to *prePaid*.
  Changing this creates a new instance.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto-renew is enabled.
  Valid values are `true` and `false`, defaults to `false`.
  Changing this creates a new instance.

* `tags` - (Optional, Map) The key/value pairs to associate with the DDS instance.

The `datastore` block supports:

* `type` - (Required, String, ForceNew) Specifies the DB engine. **DDS-Community** is supported.

* `version` - (Required, String, ForceNew) Specifies the DB instance version. For the Community Edition, the valid
  values are `3.2`, `3.4`, `4.0`, `4.2`, or `4.4`.

* `storage_engine` - (Optional, String, ForceNew) Specifies the storage engine of the DB instance.
  If `version` is set to `3.2`, `3.4`, or `4.0`, the value is **wiredTiger**.
  If `version` is set to `4.2`, or `4.4`, the value is **rocksDB**.

The `configuration` block supports:

* `type` - (Required, String, ForceNew) Specifies the node type. Valid value:
  + For a Community Edition cluster instance, the value can be **mongos**, **shard** or **config**.
  + For a Community Edition replica set instance, the value is **replica**.
  + For a Community Edition single node instance, the value is **single**.
    Changing this creates a new instance.

* `id` - (Required, String, ForceNew) Specifies the ID of the template.
  Changing this creates a new instance.

The `flavor` block supports:

* `type` - (Required, String, ForceNew) Specifies the node type. Valid value:
  + For a Community Edition cluster instance, the value can be **mongos**, **shard**, or **config**.
  + For an Enhanced Edition cluster instance, the value is **shard**.
  + For a Community Edition replica set instance, the value is **replica**.
  + For a Community Edition single node instance, the value is **single**.

* `num` - (Required, Int) Specifies the node quantity. Valid value:
  + In a Community Edition cluster instance,the number of mongos ranges from 2 to 16.
  + In a Community Edition cluster instance,the number of shards ranges from 2 to 16.
  + In an Enhanced Edition cluster instance, the number of shards ranges from 2 to 12.
  + config: the value is 1.
  + replica: the value is 1.
  + single: The value is 1. This parameter can be updated when the value of `type` is mongos or shard.

* `storage` - (Optional, String, ForceNew) Specifies the disk type.
  Valid value: **ULTRAHIGH** which indicates the type SSD.

* `size` - (Optional, Int) Specifies the disk size. The value must be a multiple of 10. The unit is GB. This parameter
  is mandatory for nodes except mongos and invalid for mongos. This parameter can be updated when the value of `type` is
  shard, replica or single.

* `spec_code` - (Required, String) Specifies the resource specification code. In a cluster instance, multiple
  specifications need to be specified. All specifications must be of the same series, that is, general-purpose (s6),
  enhanced (c3), or enhanced II (c6). For example:
  + dds.mongodb.s6.large.4.mongos and dds.mongodb.s6.large.4.config have the same specifications.
  + dds.mongodb.s6.large.4.mongos and dds.mongodb.c3.large.4.config are not of the same specifications. This parameter
      can be updated when the value of `type` is mongos, shard, replica or single.

The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. The value cannot be empty. It must be a valid value in the
  "hh:mm-HH:MM" format. The current time is in the UTC format.
  + The HH value must be 1 greater than the hh value.
  + The values from mm and MM must be the same and must be set to any of the following 00, 15, 30, or 45.

* `keep_days` - (Required, Int) Specifies the number of days to retain the generated backup files. The value range is
  from 0 to 732.
  + If this parameter is set to 0, the automated backup policy is not set.
  + If this parameter is not transferred, the automated backup policy is enabled by default. Backup files are stored
      for seven days by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the the DB instance ID.
* `db_username` - Indicates the DB Administrator name.
* `status` - Indicates the the DB instance status.
* `port` - Indicates the database port number. The port range is 2100 to 9500.
* `nodes` - Indicates the instance nodes information. Structure is documented below.

The `nodes` block contains:

* `id` - Indicates the node ID.
* `name` - Indicates the node name.
* `role` - Indicates the node role.
* `type` - Indicates the node type.
* `private_ip` - Indicates the private IP address of a node. This parameter is valid only for mongos nodes, replica set
  instances, and single node instances.
* `public_ip` - Indicates the EIP that has been bound on a node. This parameter is valid only for mongos nodes of
  cluster instances, primary nodes and secondary nodes of replica set instances, and single node instances.
* `status` - Indicates the node status.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

DDS instance can be imported using the `id`, e.g.

```sh
terraform import huaweicloud_dds_instance.instance 9c6d6ff2cba3434293fd479571517e16in02
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `password`, `availability_zone`, `flavor`, configuration.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```
resource "huaweicloud_dds_instance" "instance" {
    ...

  lifecycle {
    ignore_changes = [
      password, availability_zone, flavor, configuration,
    ]
  }
}
```
