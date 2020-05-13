// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func setTestCheckMonitoringSloId(res string, sloId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		updateId, err := getTestResourceMonitoringSloId(res, s)
		if err != nil {
			return err
		}
		*sloId = updateId
		return nil
	}
}

func testCheckMonitoringSloIdAfterUpdate(res string, sloId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		updateId, err := getTestResourceMonitoringSloId(res, s)
		if err != nil {
			return err
		}

		if sloId == nil {
			return fmt.Errorf("unexpected error, slo ID was not set")
		}

		if *sloId != updateId {
			return fmt.Errorf("unexpected mismatch in slo ID after update, resource was recreated. Initial %q, Updated %q",
				*sloId, updateId)
		}
		return nil
	}
}

func getTestResourceMonitoringSloId(res string, s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[res]
	if !ok {
		return "", fmt.Errorf("not found: %s", res)
	}

	if rs.Primary.ID == "" {
		return "", fmt.Errorf("no ID is set for %s", res)
	}

	if v, ok := rs.Primary.Attributes["slo_id"]; ok {
		return v, nil
	}

	return "", fmt.Errorf("slo_id not set on resource %s", res)
}

func TestAccMonitoringSlo_basic(t *testing.T) {
	t.Parallel()

	var generatedId string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringSloDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSlo_basic(),
				Check:  setTestCheckMonitoringSloId("google_monitoring_slo.primary", &generatedId),
			},
			{
				ResourceName:      "google_monitoring_slo.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSlo_basicUpdate(),
				Check:  testCheckMonitoringSloIdAfterUpdate("google_monitoring_slo.primary", &generatedId),
			},
			{
				ResourceName:      "google_monitoring_slo.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}

func TestAccMonitoringSlo_requestBased(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"random_suffix": randString(t, 10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringSloDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSlo_requestBasedDistribution(context),
			},
			{
				ResourceName:      "google_monitoring_slo.request_based_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSlo_requestBasedGoodBadRatio(context),
			},
			{
				ResourceName:      "google_monitoring_slo.request_based_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSlo_requestBasedGoodTotalRatio(context),
			},
			{
				ResourceName:      "google_monitoring_slo.request_based_slo",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}

func testAccMonitoringSlo_basic() string {
	return `
data "google_monitoring_app_engine_service" "ae" {
  module_id = "default"
}

resource "google_monitoring_slo" "primary" {
  service = data.google_monitoring_app_engine_service.ae.service_id

  goal = 0.9
  rolling_period_days = 1

  basic_sli {
    latency {
      threshold = "1s"
    }
  }
}
`
}

func testAccMonitoringSlo_basicUpdate() string {
	return `
data "google_monitoring_app_engine_service" "ae" {
  module_id = "default"
}

resource "google_monitoring_slo" "primary" {
  service = data.google_monitoring_app_engine_service.ae.service_id

  goal = 0.8
  display_name = "Terraform Test updated SLO"
  calendar_period = "WEEK"

  basic_sli {
    latency {
      threshold = "2s"
    }
  }
}
`
}

func testAccMonitoringSlo_requestBasedDistribution(context map[string]interface{}) string {
	return Nprintf(`
resource "google_monitoring_custom_service" "srv" {
  service_id = "tf-test-custom-srv%{random_suffix}"
  display_name = "My Custom Service"
}

resource "google_monitoring_slo" "request_based_slo" {
  service = google_monitoring_custom_service.srv.service_id
  slo_id = "tf-test-consumed-api-slo%{random_suffix}"
  display_name = "Terraform Test SLO with request based SLI"

  goal = 0.9
  rolling_period_days = 30

  request_based_sli {
    distribution_cut {
      distribution_filter = join(" AND ", [
        "metric.type=\"serviceruntime.googleapis.com/api/request_latencies\"",
        "resource.type=\"consumed_api\"",
        "resource.label.\"project_id\"=\"%{project}\"",
      ])

      range {
        max = 10
      }
    }
  }
}
`, context)
}

func testAccMonitoringSlo_requestBasedGoodTotalRatio(context map[string]interface{}) string {
	return Nprintf(`
resource "google_monitoring_custom_service" "srv" {
  service_id = "tf-test-custom-srv%{random_suffix}"
  display_name = "My Custom Service"
}

resource "google_monitoring_slo" "request_based_slo" {
  service = google_monitoring_custom_service.srv.service_id
  slo_id = "tf-test-consumed-api-slo%{random_suffix}"
  display_name = "Terraform Test SLO with request based SLI (good total ratio)"

  goal = 0.9
  rolling_period_days = 30

  request_based_sli {
    good_total_ratio {
      good_service_filter = join(" AND ", [
        "metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
        "resource.type=\"consumed_api\"",
        "resource.label.\"project_id\"=\"%{project}\"",
        "metric.label.\"response_code\"=\"200\"",
      ])
      total_service_filter = join(" AND ", [
        "metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
        "resource.type=\"consumed_api\"",
        "resource.label.\"project_id\"=\"%{project}\"",
      ])
    }
  }
}
`, context)
}

func testAccMonitoringSlo_requestBasedGoodBadRatio(context map[string]interface{}) string {
	return Nprintf(`
resource "google_monitoring_custom_service" "srv" {
  service_id = "tf-test-custom-srv%{random_suffix}"
  display_name = "My Custom Service"
}

resource "google_monitoring_slo" "request_based_slo" {
  service = google_monitoring_custom_service.srv.service_id
  slo_id = "tf-test-consumed-api-slo%{random_suffix}"
  display_name = "Terraform Test SLO with request based SLI (good total ratio)"

  goal = 0.9
  rolling_period_days = 30

  request_based_sli {
    good_total_ratio {
      good_service_filter = join(" AND ", [
        "metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
        "resource.type=\"consumed_api\"",
        "resource.label.\"project_id\"=\"%{project}\"",
        "metric.label.\"response_code\"=\"200\"",
      ])
      bad_service_filter = join(" AND ", [
        "metric.type=\"serviceruntime.googleapis.com/api/request_count\"",
        "resource.type=\"consumed_api\"",
        "resource.label.\"project_id\"=\"%{project}\"",
				"metric.label.\"response_code\"=\"400\"",
      ])
    }
  }
}
`, context)
}
