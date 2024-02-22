package aom

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPromInstanceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	getPrometheusInstanceHttpUrl := "v1/{project_id}/aom/prometheus"
	getPrometheusInstanceHttpUrl = strings.ReplaceAll(getPrometheusInstanceHttpUrl, "{project_id}", client.ProjectID)
	getPrometheusInstanceHttpUrl += fmt.Sprintf("?prom_id=%s", state.Primary.ID)
	getPrometheusInstancePath := client.Endpoint + getPrometheusInstanceHttpUrl

	getPrometheusInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getPrometheusInstanceResp, err := client.Request("GET", getPrometheusInstancePath, &getPrometheusInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AOM prometheus instance: %s", err)
	}

	getPrometheusInstanceRespBody, err := utils.FlattenResponse(getPrometheusInstanceResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AOM prometheus instance: %s", err)
	}

	curJson, err := jmespath.Search("prometheus[0]", getPrometheusInstanceRespBody)
	if err != nil || curJson == nil {
		return nil, fmt.Errorf("error retrieving AOM prometheus instance: %s", err)
	}

	return curJson, nil
}

func TestAccPromInstance_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_prom_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPromInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesAOMPromInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "prom_name", rName),
					resource.TestCheckResourceAttr(resourceName, "prom_type", "VPC"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "prom_version", "1.5"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "remote_write_url"),
					resource.TestCheckResourceAttrSet(resourceName, "remote_read_url"),
					resource.TestCheckResourceAttrSet(resourceName, "prom_http_api_endpoint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func tesAOMPromInstance_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_prom_instance" "test" {
  prom_name             = "%s"
  prom_type             = "VPC"
  enterprise_project_id = "0"
  prom_version          = "1.5"
}
`, name)
}
