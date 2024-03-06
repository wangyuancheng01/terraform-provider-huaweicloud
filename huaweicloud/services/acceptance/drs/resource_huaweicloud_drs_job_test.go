package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/drs/v3/jobs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getDrsJobResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DrsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DRS client, err: %s", err)
	}
	detailResp, err := jobs.Get(client, jobs.QueryJobReq{Jobs: []string{state.Primary.ID}})
	if err != nil {
		return nil, err
	}
	status := detailResp.Results[0].Status
	if status == "DELETED" {
		return nil, golangsdk.ErrDefault404{}
	}
	return detailResp, nil
}

func TestAccResourceDrsJob_basic(t *testing.T) {
	var obj jobs.BatchCreateJobReq
	resourceName := "huaweicloud_drs_job.test"
	name := acceptance.RandomAccResourceName()
	dbName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	pwd := "TestDrs@123"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDrsJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJob_migrate_mysql(name, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "migration"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readnoly", "true"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.instance_id",
						"huaweicloud_rds_instance.test2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.region",
						"huaweicloud_rds_instance.test2", "region"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", name),
				),
			},
			{
				Config: testAccDrsJob_migrate_mysql(updateName, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "type", "migration"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readnoly", "true"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", updateName),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.instance_id",
						"huaweicloud_rds_instance.test2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.region",
						"huaweicloud_rds_instance.test2", "region"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", updateName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"source_db.0.password", "destination_db.0.password",
					"expired_days", "migrate_definer", "force_destroy"},
			},
		},
	})
}

func testAccDrsJob_mysql(index int, name, pwd, ip string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_instance" "test%d" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.ingress,
    huaweicloud_networking_secgroup_rule.egress,
  ]
  name                = "%s%d"
  flavor              = "rds.mysql.x1.large.2.ha"
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  fixed_ip            = "%s"
  ha_replication_mode = "semisync"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[3],
  ]

  db {
    password = "%s"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, index, name, index, ip, pwd)
}

func testAccDrsJob_migrate_mysql(name, dbName, pwd string) string {
	netConfig := common.TestBaseNetwork(name)
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")

	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 3306
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_networking_secgroup_rule" "egress" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_availability_zones" "test" {}

%s
%s

resource "huaweicloud_drs_job" "test" {
  name           = "%s"
  type           = "migration"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "eip"
  migration_type = "FULL_INCR_TRANS"
  description    = "%s"
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test1.fixed_ip
    port        = 3306
    user        = "root"
    password    = "%s"
  }


  destination_db {
    region      = huaweicloud_rds_instance.test2.region
    ip          = huaweicloud_rds_instance.test2.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "%s"
    instance_id = huaweicloud_rds_instance.test2.id
    subnet_id   = huaweicloud_rds_instance.test2.subnet_id
  }

  tags = {
    key = "%s"
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, netConfig, sourceDb, destDb, name, name, pwd, pwd, name)
}
