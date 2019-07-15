package mongodb

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProvider schema.Provider
var testAccProviders map[string]terraform.ResourceProvider

func init() {
	testAccProvider := Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"mongodb": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err=%s", err)
	}
}

// testAccPreCheck Acceptance Tests pre-checks
func testAccPreCheck(t *testing.T) {
	if _, err := getHostPort(); err != nil {
		t.Fatalf("Error reading host port from environment: %v", err)
	}

	// ensure a Docker daemon is available
}

// getHostPort loads the host port (SSH) to use for testing (usually bound to a Docker container)
func getHostPort() (int, error) {
	return strconv.Atoi(os.Getenv("TF_VAR_host_port"))
}

func TestMongoDBProviderSchema_unit(t *testing.T) {
	// ensure we don't accidentally define duplicate keys in our provider schema
	_, err := getMergedProviderSchema()
	if err != nil {
		t.Fatalf("Duplicate keys detected in provider schema: %v", err)
	}
}
