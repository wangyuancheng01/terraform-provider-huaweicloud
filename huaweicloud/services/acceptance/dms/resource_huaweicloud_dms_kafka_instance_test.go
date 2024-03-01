package dms

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getKafkaInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud DMS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID).Extract()
}

func TestAccKafkaInstance_basic(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := rName + "update"
	resourceName := "huaweicloud_dms_kafka_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	// DMS instances use the tenant-level shared lock, the instances cannot be created or modified in parallel.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "security_protocol", "SASL_PLAINTEXT"),
					resource.TestCheckResourceAttr(resourceName, "enabled_mechanisms.0", "SCRAM-SHA-512"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestMatchResourceAttr(resourceName, "cross_vpc_accesses.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
			{
				Config: testAccKafkaInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "kafka test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "enable_auto_topic", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"manager_password",
					"used_storage_space",
					"cross_vpc_accesses",
				},
			},
		},
	})
}

func TestAccKafkaInstance_prePaid(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := rName + "update"
	resourceName := "huaweicloud_dms_kafka_instance.test"
	baseNetwork := common.TestBaseNetwork(rName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	// DMS instances use the tenant-level shared lock, the instances cannot be created or modified in parallel.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_newFormat_prePaid(baseNetwork, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
				),
			},
			{
				Config: testAccKafkaInstance_newFormat_prePaid_update(baseNetwork, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "kafka test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"manager_password",
					"used_storage_space",
					"cross_vpc_accesses",
					"auto_renew",
					"period",
					"period_unit",
				},
			},
		},
	})
}

func TestAccKafkaInstance_withEpsId(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccKafkaInstance_compatible(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_compatible(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrPair(resourceName, "storage_space", "data.huaweicloud_dms_product.test", "storage"),
				),
			},
		},
	})
}

func TestAccKafkaInstance_newFormat(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_newFormat(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "security_protocol", "SASL_PLAINTEXT"),
					resource.TestCheckResourceAttr(resourceName, "enabled_mechanisms.0", "SCRAM-SHA-512"),
					resource.TestCheckResourceAttrPair(resourceName, "broker_num",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.properties.0.min_broker"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "storage_spec_code",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.ios.0.storage_spec_code"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "3"),

					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.1.advertised_ip", "www.terraform-test.com"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.2.advertised_ip", "192.168.0.53"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "min.insync.replicas"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "2"),
				),
			},
			{
				Config: testAccKafkaInstance_newFormatUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "4"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "storage_spec_code",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.ios.0.storage_spec_code"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "600"),

					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.0.advertised_ip", "192.168.0.61"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.1.advertised_ip", "test.terraform.com"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.2.advertised_ip", "192.168.0.62"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.3.advertised_ip", "192.168.0.63"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "auto.create.groups.enable"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "false"),
				),
			},
		},
	})
}

func testAccKafkaInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_product" "test" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.7"
  bandwidth         = "100MB"
  storage_spec_code = "dms.physical.storage.ultra"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name               = "%s"
  description        = "kafka test"
  access_user        = "user"
  password           = "Kafkatest@123"
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code

  manager_user       = "kafka-user"
  manager_password   = "Kafkatest@123"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccKafkaInstance_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_product" "test" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.7"
  bandwidth         = "300MB"
  storage_spec_code = "dms.physical.storage.ultra"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name               = "%s"
  description        = "kafka test update"
  access_user        = "user"
  password           = "Kafkatest@1234"
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code

  manager_user       = "kafka-user"
  manager_password   = "Kafkatest@123"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]
  enable_auto_topic  = true

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}
`, common.TestBaseNetwork(rName), updateName)
}

func testAccKafkaInstance_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_product" "test" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.7"
  bandwidth         = "300MB"
  storage_spec_code = "dms.physical.storage.ultra"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name                  = "%s"
  description           = "kafka test"
  access_user           = "user"
  password              = "Kafkatest@123"
  vpc_id                = huaweicloud_vpc.test.id
  network_id            = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zones    = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
  product_id            = data.huaweicloud_dms_product.test.id
  engine_version        = data.huaweicloud_dms_product.test.version
  storage_spec_code     = data.huaweicloud_dms_product.test.storage_spec_code

  manager_user          = "kafka-user"
  manager_password      = "Kafkatest@123"
  enterprise_project_id = "%s"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccKafkaInstance_compatible(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_product" "test" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.7"
  bandwidth         = "300MB"
  storage_spec_code = "dms.physical.storage.ultra"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name        = "%s"
  description = "kafka test"

  availability_zones    = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code
  storage_space     = data.huaweicloud_dms_product.test.storage
  # use deprecated argument "bandwidth"
  bandwidth         = data.huaweicloud_dms_product.test.bandwidth

  access_user      = "user"
  password         = "Kafkatest@123"
  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccKafkaInstance_newFormat(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = local.flavor.id
  storage_spec_code  = local.flavor.ios[0].storage_spec_code
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2]
  ]
  engine_version = "2.7"
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3

  access_user        = "user"
  password           = "Kafkatest@123"
  manager_user       = "kafka-user"
  manager_password   = "Kafkatest@123"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]

  cross_vpc_accesses {
    advertised_ip = ""
  }
  cross_vpc_accesses {
    advertised_ip = "www.terraform-test.com"
  }
  cross_vpc_accesses {
    advertised_ip = "192.168.0.53"
  }

  parameters {
    name  = "min.insync.replicas"
    value = "2"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccKafkaInstance_newFormatUpdate(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.4u8g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = local.flavor.id
  storage_spec_code  = local.flavor.ios[0].storage_spec_code
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2]
  ]
  engine_version = "2.7"
  storage_space  = 600
  broker_num     = 4

  access_user        = "user"
  password           = "Kafkatest@123"
  manager_user       = "kafka-user"
  manager_password   = "Kafkatest@123"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]

  cross_vpc_accesses {
    advertised_ip = "192.168.0.61"
  }
  cross_vpc_accesses {
    advertised_ip = "test.terraform.com"
  }
  cross_vpc_accesses {
    advertised_ip = "192.168.0.62"
  }
  cross_vpc_accesses {
    advertised_ip = "192.168.0.63"
  }

  parameters {
    name  = "auto.create.groups.enable"
    value = "false"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccKafkaInstance_newFormat_prePaid(baseNetwork, rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  engine_version = "2.7"
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3

  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, baseNetwork, rName)
}

func testAccKafkaInstance_newFormat_prePaid_update(baseNetwork, updateName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  description       = "kafka test update"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  engine_version = "2.7"
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3

  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}`, baseNetwork, updateName)
}
