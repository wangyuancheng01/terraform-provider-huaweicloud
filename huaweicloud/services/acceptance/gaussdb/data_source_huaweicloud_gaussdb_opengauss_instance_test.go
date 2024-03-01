package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccOpenGaussInstanceDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_gaussdb_opengauss_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstanceDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "volume.0.size", "40"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstanceDataSource_haModeCentralized(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_gaussdb_opengauss_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstanceDataSource_haModeCentralized(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstanceDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "volume.0.size", "40"),
				),
			},
		},
	})
}

func testAccCheckOpenGaussInstanceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GaussDB opengauss instance data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GaussDB opengauss data source ID not set ")
		}

		return nil
	}
}

func testAccOpenGaussInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  name        = "%[2]s"
  password    = "Test@12345678"
  flavor      = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  availability_zone = "cn-north-4a,cn-north-4a,cn-north-4a"
  security_group_id = huaweicloud_networking_secgroup.test.id

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }
  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  sharding_num = 1
  coordinator_num = 2
}

data "huaweicloud_gaussdb_opengauss_instance" "test" {
  name = huaweicloud_gaussdb_opengauss_instance.test.name
  depends_on = [
    huaweicloud_gaussdb_opengauss_instance.test,
  ]
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccOpenGaussInstanceDataSource_haModeCentralized(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  name        = "%[2]s"
  password    = "Test@12345678"
  flavor      = "gaussdb.opengauss.ee.m6.2xlarge.x868.ha"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  availability_zone = "cn-north-4a,cn-north-4a,cn-north-4a"
  security_group_id = huaweicloud_networking_secgroup.test.id

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
  }
  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  replica_num = 3
}

data "huaweicloud_gaussdb_opengauss_instance" "test" {
  name = huaweicloud_gaussdb_opengauss_instance.test.name
  depends_on = [
    huaweicloud_gaussdb_opengauss_instance.test,
  ]
}
`, common.TestBaseNetwork(rName), rName)
}
