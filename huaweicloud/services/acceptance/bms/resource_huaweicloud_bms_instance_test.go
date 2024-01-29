package bms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/bms/v1/baremetalservers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccBmsInstance_basic(t *testing.T) {
	var instance baremetalservers.CloudServer

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_bms_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckUserId(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccBmsInstance_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", rName)),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.tag1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.tag2", "value2"),
				),
			},
		},
	})
}

func testAccCheckBmsInstanceDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	bmsClient, err := cfg.BmsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating bms client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_bms_instance" {
			continue
		}

		server, err := baremetalservers.Get(bmsClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckBmsInstanceExists(n string, instance *baremetalservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		bmsClient, err := cfg.BmsV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating bms client: %s", err)
		}

		found, err := baremetalservers.Get(bmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccBmsInstance_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_bms_flavors" "test" {
  availability_zone = try(element(data.huaweicloud_availability_zones.test.names, 0), "")
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%s"
}`, common.TestBaseNetwork(rName), rName)
}

func testAccBmsInstance_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[2]s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = data.huaweicloud_bms_flavors.test.flavors[0].id
  key_pair          = huaweicloud_kps_keypair.test.name
  image_id          = "519ea918-1fea-4ebc-911a-593739b1a3bc" # CentOS 7.4 64bit for BareMetal

  name                  = "%[2]s"
  user_id               = "%[3]s"
  enterprise_project_id = "%[4]s"

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
    key = "value"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "false"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HW_USER_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccBmsInstance_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = data.huaweicloud_bms_flavors.test.flavors[0].id
  key_pair          = huaweicloud_kps_keypair.test.name
  image_id          = "519ea918-1fea-4ebc-911a-593739b1a3bc" # CentOS 7.4 64bit for BareMetal

  name                  = "%[2]s_update"
  user_id               = "%[3]s"
  enterprise_project_id = "%[4]s"

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  tags = {
    tag1 = "value1"
    tag2 = "value2"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HW_USER_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
