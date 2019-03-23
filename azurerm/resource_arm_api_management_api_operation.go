package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO: a Data Source, but without Request/Response things?

func resourceArmApiManagementApiOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementApiOperationCreateUpdate,
		Read:   resourceArmApiManagementApiOperationRead,
		Update: resourceArmApiManagementApiOperationCreateUpdate,
		Delete: resourceArmApiManagementApiOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"operation_id": azure.SchemaApiManagementChildName(),

			"api_name": azure.SchemaApiManagementChildName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"resource_group_name": resourceGroupNameSchema(),

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"method": {
				Type:     schema.TypeString,
				Required: true,
			},

			"url_template": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"request": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"header": azure.SchemaApiManagementOperationParameterContract(),

						"query_parameter": azure.SchemaApiManagementOperationParameterContract(),

						"representation": azure.SchemaApiManagementOperationRepresentation(),
					},
				},
			},

			"response": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"status_code": {
							Type:     schema.TypeSet,
							Optional: true,
						},

						"header": azure.SchemaApiManagementOperationParameterContract(),

						"representation": azure.SchemaApiManagementOperationRepresentation(),
					},
				},
			},

			"template_parameter": azure.SchemaApiManagementOperationParameterContract(),
		},
	}
}

func resourceArmApiManagementApiOperationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementApiOperationsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiId := d.Get("api_name").(string)
	operationId := d.Get("operation_id").(string)

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	method := d.Get("method").(string)
	urlTemplate := d.Get("url_template").(string)

	requestContractRaw := d.Get("request").([]interface{})
	requestContract := expandApiManagementOperationRequestContract(requestContractRaw)

	responseContractsRaw := d.Get("response").(*schema.Set).List()
	responseContracts := expandApiManagementOperationResponseContract(responseContractsRaw)

	templateParametersRaw := d.Get("template_parameter").(*schema.Set).List()
	templateParameters := azure.ExpandApiManagementOperationParameterContract(templateParametersRaw)

	parameters := apimanagement.OperationContract{
		OperationContractProperties: &apimanagement.OperationContractProperties{
			Description:        utils.String(description),
			DisplayName:        utils.String(displayName),
			Method:             utils.String(method),
			Request:            requestContract,
			Responses:          responseContracts,
			TemplateParameters: templateParameters,
			URLTemplate:        utils.String(urlTemplate),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiId, operationId, parameters, ""); err != nil {
		return fmt.Errorf("Error creating/updating API Operation %q (API %q / API Management Service %q / Resource Group %q): %+v", operationId, apiId, serviceName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiId, operationId)
	if err != nil {
		return fmt.Errorf("Error retrieving API Operation %q (API %q / API Management Service %q / Resource Group %q): %+v", operationId, apiId, serviceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmApiManagementApiOperationRead(d, meta)
}

func resourceArmApiManagementApiOperationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementApiOperationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiId := id.Path["apis"]
	operationId := id.Path["operations"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiId, operationId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] API Operation %q (API %q / API Management Service %q / Resource Group %q) was not found - removing from state!", operationId, apiId, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving API Operation %q (API %q / API Management Service %q / Resource Group %q): %+v", operationId, apiId, serviceName, resourceGroup, err)
	}

	d.Set("operation_id", operationId)
	d.Set("api_name", apiId)
	d.Set("api_management_name", serviceName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.OperationContractProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("method", props.Method)
		d.Set("url_template", props.URLTemplate)

		flattenedRequest := flattenApiManagementOperationRequestContract(props.Request)
		if err := d.Set("request", flattenedRequest); err != nil {
			return fmt.Errorf("Error flattening `request`: %+v", err)
		}

		flattenedResponse := flattenApiManagementOperationResponseContract(props.Responses)
		if err := d.Set("response", flattenedResponse); err != nil {
			return fmt.Errorf("Error flattening `response`: %+v", err)
		}

		flattenedTemplateParams := azure.FlattenApiManagementOperationParameterContract(props.TemplateParameters)
		if err := d.Set("template_parameter", flattenedTemplateParams); err != nil {
			return fmt.Errorf("Error flattening `template_parameter`: %+v", err)
		}
	}

	return nil
}

func resourceArmApiManagementApiOperationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementApiOperationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiId := id.Path["apis"]
	operationId := id.Path["operations"]

	resp, err := client.Delete(ctx, resourceGroup, serviceName, apiId, operationId, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting API Operation %q (API %q / API Management Service %q / Resource Group %q): %+v", operationId, apiId, serviceName, resourceGroup, err)
		}
	}

	return nil
}

func expandApiManagementOperationRequestContract(input []interface{}) *apimanagement.RequestContract {
	if len(input) == 0 {
		return nil
	}

	vs := input[0].(map[string]interface{})
	description := vs["description"].(string)

	headersRaw := vs["header"].(*schema.Set).List()
	headers := azure.ExpandApiManagementOperationParameterContract(headersRaw)

	queryParametersRaw := vs["query_parameter"].(*schema.Set).List()
	queryParameters := azure.ExpandApiManagementOperationParameterContract(queryParametersRaw)

	representationsRaw := vs["representation"].(*schema.Set).List()
	representations := azure.ExpandApiManagementOperationRepresentation(representationsRaw)

	return &apimanagement.RequestContract{
		Description:     utils.String(description),
		Headers:         headers,
		QueryParameters: queryParameters,
		Representations: representations,
	}
}

func flattenApiManagementOperationRequestContract(input *apimanagement.RequestContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	if input.Description != nil {
		output["description"] = *input.Description
	}

	output["headers"] = azure.FlattenApiManagementOperationParameterContract(input.Headers)
	output["query_parameter"] = azure.FlattenApiManagementOperationParameterContract(input.QueryParameters)
	output["representation"] = azure.FlattenApiManagementOperationRepresentation(input.Representations)

	return []interface{}{output}
}

func expandApiManagementOperationResponseContract(input []interface{}) *[]apimanagement.ResponseContract {
	if len(input) == 0 {
		return &[]apimanagement.ResponseContract{}
	}

	outputs := make([]apimanagement.ResponseContract, 0)

	for _, v := range input {
		vs := v.(map[string]interface{})

		description := vs["description"].(string)
		statusCode := vs["status_code"].(int)

		headersRaw := vs["header"].(*schema.Set).List()
		headers := azure.ExpandApiManagementOperationParameterContract(headersRaw)

		representationsRaw := vs["representation"].(*schema.Set).List()
		representations := azure.ExpandApiManagementOperationRepresentation(representationsRaw)

		output := apimanagement.ResponseContract{
			Description:     utils.String(description),
			Headers:         headers,
			Representations: representations,
		}

		if statusCode != 0 {
			output.StatusCode = utils.Int32(int32(statusCode))
		}

		outputs = append(outputs, output)
	}

	return &outputs
}

func flattenApiManagementOperationResponseContract(input *[]apimanagement.ResponseContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	outputs := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		if v.Description != nil {
			output["description"] = *v.Description
		}

		if v.StatusCode != nil {
			output["status_code"] = int(*v.StatusCode)
		}

		output["header"] = azure.FlattenApiManagementOperationParameterContract(v.Headers)
		output["representation"] = azure.FlattenApiManagementOperationRepresentation(v.Representations)

		outputs = append(outputs, output)
	}

	return outputs
}
