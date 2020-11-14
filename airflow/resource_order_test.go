package airflow

import (
	"fmt"
	"testing"

	hc "github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAirflowOrderBasic(t *testing.T) {
	coffeeID := "1"
	quantity := "2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAirflowOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAirflowOrderConfigBasic(coffeeID, quantity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAirflowOrderExists("airflow_order.new"),
				),
			},
		},
	})
}

func testAccCheckAirflowOrderDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*hc.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "airflow_order" {
			continue
		}

		orderID := rs.Primary.ID

		err := c.DeleteOrder(orderID)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckAirflowOrderConfigBasic(coffeeID, quantity string) string {
	return fmt.Sprintf(`
	resource "airflow_order" "new" {
		items {
			coffee {
				id = %s
			}
    		quantity = %s
  		}
	}
	`, coffeeID, quantity)
}

func testAccCheckAirflowOrderExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OrderID set")
		}

		return nil
	}
}
