package outbound_attempt_limit

import (
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	resourceExporter "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_exporter"
	registrar "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_register"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

/*
resource_genesyscloud_outbound_attempt_limit_schema.go holds four functions within it:

1.  The registration code that registers the Datasource, Resource and Exporter for the package.
2.  The resource schema definitions for the outbound_attempt_limit resource.
3.  The datasource schema definitions for the outbound_attempt_limit datasource.
4.  The resource exporter configuration for the outbound_attempt_limit exporter.
*/

const ResourceType = "genesyscloud_outbound_attempt_limit"

var (
	recallSettings = &schema.Resource{
		Schema: map[string]*schema.Schema{
			`nbr_attempts`: {
				Description: `Number of recall attempts. Must be less than max_attempts_per_contact.`,
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
			},
			`minutes_between_attempts`: {
				Description:  `Number of minutes between attempts. Must be greater than or equal to 5.`,
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(5),
			},
		},
	}

	attemptLimitRecallSettingsResource = &schema.Resource{
		Schema: map[string]*schema.Schema{
			`answering_machine`: {
				Optional: true,
				MaxItems: 1,
				Type:     schema.TypeSet,
				Elem:     recallSettings,
			},
			`busy`: {
				Optional: true,
				MaxItems: 1,
				Type:     schema.TypeSet,
				Elem:     recallSettings,
			},
			`fax`: {
				Optional: true,
				MaxItems: 1,
				Type:     schema.TypeSet,
				Elem:     recallSettings,
			},
			`no_answer`: {
				Optional: true,
				MaxItems: 1,
				Type:     schema.TypeSet,
				Elem:     recallSettings,
			},
		},
	}
)

// SetRegistrar registers all of the resources, datasources and exporters in the package
func SetRegistrar(regInstance registrar.Registrar) {
	regInstance.RegisterDataSource(ResourceType, DataSourceOutboundAttemptLimit())
	regInstance.RegisterResource(ResourceType, ResourceOutboundAttemptLimit())
	regInstance.RegisterExporter(ResourceType, OutboundAttemptLimitExporter())
}

// ResourceOutboundAttemptLimit registers the genesyscloud_outbound_attempt_limit resource with Terraform
func ResourceOutboundAttemptLimit() *schema.Resource {
	return &schema.Resource{
		Description: `Genesys Cloud Outbound Attempt Limit`,

		CreateContext: provider.CreateWithPooledClient(createOutboundAttemptLimit),
		ReadContext:   provider.ReadWithPooledClient(readOutboundAttemptLimit),
		UpdateContext: provider.UpdateWithPooledClient(updateOutboundAttemptLimit),
		DeleteContext: provider.DeleteWithPooledClient(deleteOutboundAttemptLimit),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			`name`: {
				Description: `The name for the attempt limit.`,
				Optional:    true,
				Type:        schema.TypeString,
			},
			`max_attempts_per_contact`: {
				Description: `The maximum number of times a contact can be called within the resetPeriod. Required if maxAttemptsPerNumber is not defined.`,
				Optional:    true,
				Type:        schema.TypeInt,
			},
			`max_attempts_per_number`: {
				Description: `The maximum number of times a phone number can be called within the resetPeriod. Required if maxAttemptsPerContact is not defined.`,
				Optional:    true,
				Type:        schema.TypeInt,
			},
			`time_zone_id`: {
				Description: `If the resetPeriod is TODAY, this specifies the timezone in which TODAY occurs. Required if the resetPeriod is TODAY.`,
				Optional:    true,
				Type:        schema.TypeString,
			},
			`reset_period`: {
				Description:  `After how long the number of attempts will be set back to 0.`,
				Optional:     true,
				Type:         schema.TypeString,
				Default:      `NEVER`,
				ValidateFunc: validation.StringInSlice([]string{`NEVER`, `TODAY`}, true),
			},
			`recall_entries`: {
				Description: `Configuration for recall attempts.`,
				Optional:    true,
				MaxItems:    1,
				Type:        schema.TypeList,
				Elem:        attemptLimitRecallSettingsResource,
			},
		},
	}
}

// OutboundAttemptLimitExporter returns the resourceExporter object used to hold the genesyscloud_outbound_attempt_limit exporter's config
func OutboundAttemptLimitExporter() *resourceExporter.ResourceExporter {
	return &resourceExporter.ResourceExporter{
		GetResourcesFunc: provider.GetAllWithPooledClient(getAllAttemptLimits),
	}
}

// DataSourceOutboundAttemptLimit registers the genesyscloud_outbound_attempt_limit data source
func DataSourceOutboundAttemptLimit() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud Outbound Attempt Limits. Select an attempt limit by name.",
		ReadContext: provider.ReadWithPooledClient(dataSourceOutboundAttemptLimitRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Attempt Limit name.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}
