package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaApiManagementName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.ApiManagementServiceName,
	}
}

func SchemaApiManagementDataSourceName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validate.ApiManagementServiceName,
	}
}

// SchemaApiManagementChildName returns the Schema for the identifier
// used by resources within nested under the API Management Service resource
func SchemaApiManagementChildName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.ApiManagementChildName,
	}
}

// SchemaApiManagementChildDataSourceName returns the Schema for the identifier
// used by resources within nested under the API Management Service resource
func SchemaApiManagementChildDataSourceName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validate.ApiManagementChildName,
	}
}

func SchemaApiManagementUserName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.ApiManagementUserName,
	}
}

func SchemaApiManagementUserDataSourceName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validate.ApiManagementUserName,
	}
}

func SchemaApiManagementOperationRepresentation() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"content_type": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"form_parameter": SchemaApiManagementOperationParameterContract(),

				"sample": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"schema_id": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"type_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func ExpandApiManagementOperationRepresentation(input []interface{}) *[]apimanagement.RepresentationContract {
	if len(input) == 0 {
		return &[]apimanagement.RepresentationContract{}
	}

	outputs := make([]apimanagement.RepresentationContract, 0)

	for _, v := range input {
		vs := v.(map[string]interface{})

		contentType := vs["content_type"].(string)
		formParametersRaw := vs["form_parameter"].(*schema.Set).List()
		formParameters := ExpandApiManagementOperationParameterContract(formParametersRaw)
		sample := vs["sample"].(string)
		schemaId := vs["schema_id"].(string)
		typeName := vs["type_name"].(string)

		output := apimanagement.RepresentationContract{
			ContentType:    utils.String(contentType),
			FormParameters: formParameters,
			Sample:         utils.String(sample),
			SchemaID:       utils.String(schemaId),
			TypeName:       utils.String(typeName),
		}
		outputs = append(outputs, output)
	}

	return &outputs
}

func FlattenApiManagementOperationRepresentation(input *[]apimanagement.RepresentationContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	outputs := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		if v.ContentType != nil {
			output["content_type"] = *v.ContentType
		}

		output["form_parameter"] = FlattenApiManagementOperationParameterContract(v.FormParameters)

		if v.Sample != nil {
			output["sample"] = *v.Sample
		}

		if v.SchemaID != nil {
			output["schema_id"] = *v.SchemaID
		}

		if v.TypeName != nil {
			output["type_name"] = *v.TypeName
		}

		outputs = append(outputs, output)
	}

	return outputs
}

func SchemaApiManagementOperationParameterContract() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"description": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"type": {
					Type:     schema.TypeString,
					Required: true,
				},
				"default_value": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"required": {
					Type:     schema.TypeBool,
					Required: true,
				},
				"values": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Set: schema.HashString,
				},
			},
		},
	}
}

func ExpandApiManagementOperationParameterContract(input []interface{}) *[]apimanagement.ParameterContract {
	if len(input) == 0 {
		return &[]apimanagement.ParameterContract{}
	}

	outputs := make([]apimanagement.ParameterContract, 0)

	for _, v := range input {
		vs := v.(map[string]interface{})

		name := vs["name"].(string)
		description := vs["description"].(string)
		paramType := vs["type"].(string)
		defaultValue := vs["default_value"].(string)
		required := vs["required"].(bool)
		valuesRaw := vs["values"].(*schema.Set).List()

		output := apimanagement.ParameterContract{
			Name:         utils.String(name),
			Description:  utils.String(description),
			Type:         utils.String(paramType),
			Required:     utils.Bool(required),
			DefaultValue: utils.String(defaultValue),
			Values:       utils.ExpandStringArray(valuesRaw),
		}
		outputs = append(outputs, output)
	}

	return &outputs
}

func FlattenApiManagementOperationParameterContract(input *[]apimanagement.ParameterContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	outputs := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if v.Description != nil {
			output["description"] = *v.Description
		}

		if v.Type != nil {
			output["type"] = *v.Type
		}

		if v.Required != nil {
			output["required"] = *v.Required
		}

		if v.DefaultValue != nil {
			output["default_value"] = *v.DefaultValue
		}

		output["values"] = utils.FlattenStringArray(v.Values)

		outputs = append(outputs, output)
	}

	return outputs
}
