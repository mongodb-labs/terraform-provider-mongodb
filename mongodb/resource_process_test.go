package mongodb

// https://www.terraform.io/docs/extend/best-practices/testing.html
import (
	"fmt"
	"strconv"
	"testing"

	_ "github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/types"
)

func TestAccDeployProcess(t *testing.T) {
	var process types.ProcessConfig
	port, _ := getHostPort()
	name := "mongodb_process.standalone"

	// https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDestroyProcess,
		Steps: []resource.TestStep{
			{
				Config: testAccProcessResource(port),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(name, "host.%", "1"),
					resource.TestCheckResourceAttr(name, "host.0.hostname", "127.0.0.1"),
					resource.TestCheckResourceAttr(name, "host.0.port", strconv.Itoa(port)),
					resource.TestCheckResourceAttr(name, "host.0.user", "root"),
					resource.TestCheckResourceAttr(name, "mongod.%", "1"),
					resource.TestCheckResourceAttr(name, "mongod.0.bindip", "0.0.0.0"),
					resource.TestCheckResourceAttr(name, "mongod.0.dbpath", "/var/lib/mongodb/data"),
					resource.TestCheckResourceAttr(name, "mongod.0.logpath", "/var/log/mongodb/mongod.log"),
					resource.TestCheckResourceAttr(name, "mongod.0.port", "27017"),
					testAccCheckProcessExists(name, port, &process),
				),
			},
		},
	})
}

func testAccCheckProcessExists(name string, port int, res *types.ProcessConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID not set for: %s ", name)
		}

		_ = testAccProvider.Meta().(*ProviderConfig)

		// connect and verify MDB is up and running

		db := &types.ProcessConfig{}

		*res = *db
		return nil
	}
}

func testAccDestroyProcess(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		// ignore everything but mongodb_process resource types
		if rs.Type != "mongodb_process" {
			continue
		}

		// check that the process is no longer accessible
		return nil
	}

	return nil
}

func testAccProcessResource(port int) string {
	return fmt.Sprintf(`
	  resource "mongodb_process" "standalone" {
		"host" {
		  user     = "root"
		  hostname = "127.0.0.1"
		  port     = "%d"
		}

		"mongod" {
		  binary = "http://downloads.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1804-4.0.10.tgz"
		  bindip = "0.0.0.0"
		  workdir = "/opt/mongodb"
		}
	  }
	  `, port)
}

func TestMongoDB_unit(*testing.T) {
	// https://www.terraform.io/docs/extend/testing/unit-testing.html
}
