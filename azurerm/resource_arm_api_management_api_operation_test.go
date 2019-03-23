package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementApiOperation_basic(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
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

func TestAccAzureRMApiManagementApiOperation_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_api_operation.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementApiOperation_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management_api_operation"),
			},
		},
	})
}

func TestAccAzureRMApiManagementApiOperation_representations(t *testing.T) {
	resourceName := "azurerm_api_management_api.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiOperationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementApiOperation_representations(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiOperationExists(resourceName),
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

func testCheckAzureRMApiManagementApiOperationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).apiManagementApiOperationsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api_operation" {
			continue
		}

		operationId := rs.Primary.Attributes["operation_id"]
		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiName, operationId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMApiManagementApiOperationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		operationId := rs.Primary.Attributes["operation_id"]
		apiName := rs.Primary.Attributes["api_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).apiManagementApiOperationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, serviceName, apiName, operationId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Operation %q (API %q / API Management Service %q / Resource Group: %q) does not exist", operationId, apiName, serviceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagementApiOperationsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementApiOperation_basic(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  # ...

  display_name = "DELETE Resource"
  method       = "DELETE"
  url_template = "/resource"
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "import" {
  # ...

  display_name = "${azurerm_api_management_api_operation.test.display_name}"
  method       = "${azurerm_api_management_api_operation.test.method}"
  url_template = "${azurerm_api_management_api_operation.test.url_template}"
}
`, template)
}

func testAccAzureRMApiManagementApiOperation_representations(rInt int, location string) string {
	template := testAccAzureRMApiManagementApiOperation_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  # ...

  display_name = "Acceptance Test Operation"
  method       = "DELETE"
  url_template = "/user1"
  description = "This can only be done by the logged in user."
  
request {
    description = "Created user object"

    representation {
      content_type = "application/json"
      schema_id = "592f6c1d0af5840ca8897f0c"
      type_name = "User"
    }
  }

  response {
    status_code = 200
    description = "successful operation"

    representation {
      content_type = "application/xml"
    }

    representation {
      content_type = "application/json"
    }
  }
}

`, template)
}

func testAccAzureRMApiManagementApiOperation_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "Butter Parser"
  path                = "butter-parser"
  protocols           = ["https", "http"]
  revision            = "3"
  description         = "What is my purpose? You parse butter."
  service_url         = "https://example.com/foo/bar"

  subscription_key_parameter_names {
    header = "X-Butter-Robot-API-Key"
    query  = "location"
  }
}
`, rInt, location, rInt, rInt)
}
