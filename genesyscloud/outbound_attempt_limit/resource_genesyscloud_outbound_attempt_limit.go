package outbound_attempt_limit

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/consistency_checker"

	resourceExporter "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_exporter"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v191/platformclientv2"
)

func getAllAttemptLimits(ctx context.Context, clientConfig *platformclientv2.Configuration) (resourceExporter.ResourceIDMetaMap, diag.Diagnostics) {
	resources := make(resourceExporter.ResourceIDMetaMap)
	proxy := getOutboundAttemptLimitProxy(clientConfig)

	attemptLimits, resp, err := proxy.getAllOutboundAttemptLimits(ctx)
	if err != nil {
		return nil, util.BuildAPIDiagnosticError(ResourceType, fmt.Sprintf("Failed to get outbound attempt limits error: %s", err), resp)
	}

	for _, attemptLimit := range *attemptLimits {
		resources[*attemptLimit.Id] = &resourceExporter.ResourceMeta{BlockLabel: *attemptLimit.Name}
	}
	return resources, nil
}

func createOutboundAttemptLimit(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	proxy := getOutboundAttemptLimitProxy(sdkConfig)

	sdkAttemptLimits := buildSdkAttemptLimits(d)

	log.Printf("Creating Outbound Attempt Limit %s", *sdkAttemptLimits.Name)
	attemptLimit, resp, err := proxy.createOutboundAttemptLimit(ctx, sdkAttemptLimits)
	if err != nil {
		return util.BuildAPIDiagnosticError(ResourceType, fmt.Sprintf("Failed to create Outbound Attempt Limit %s error: %s", *sdkAttemptLimits.Name, err), resp)
	}

	d.SetId(*attemptLimit.Id)
	log.Printf("Created Outbound Attempt Limit %s %s", *attemptLimit.Name, *attemptLimit.Id)
	return readOutboundAttemptLimit(ctx, d, meta)
}

func updateOutboundAttemptLimit(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	proxy := getOutboundAttemptLimitProxy(sdkConfig)

	sdkAttemptLimits := buildSdkAttemptLimits(d)

	log.Printf("Updating Outbound Attempt Limit %s", *sdkAttemptLimits.Name)
	diagErr := util.RetryWhen(util.IsVersionMismatch, func() (*platformclientv2.APIResponse, diag.Diagnostics) {
		_, resp, err := proxy.updateOutboundAttemptLimit(ctx, d.Id(), sdkAttemptLimits)
		if err != nil {
			return resp, util.BuildAPIDiagnosticError(ResourceType, fmt.Sprintf("Failed to update outbound attempt limit %s error: %s", d.Id(), err), resp)
		}
		return nil, nil
	})
	if diagErr != nil {
		return diagErr
	}

	log.Printf("Updated Outbound Attempt Limit %s", *sdkAttemptLimits.Name)
	return readOutboundAttemptLimit(ctx, d, meta)
}

func readOutboundAttemptLimit(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	proxy := getOutboundAttemptLimitProxy(sdkConfig)
	cc := consistency_checker.NewConsistencyCheck(ctx, d, meta, ResourceOutboundAttemptLimit(), constants.ConsistencyChecks(), ResourceType)

	log.Printf("Reading Outbound Attempt Limit %s", d.Id())

	return util.WithRetriesForRead(ctx, d, func() *retry.RetryError {
		sdkAttemptLimits, resp, getErr := proxy.getOutboundAttemptLimitById(ctx, d.Id())
		if getErr != nil {
			if util.IsStatus404(resp) {
				return retry.RetryableError(util.BuildWithRetriesApiDiagnosticError(ResourceType, fmt.Sprintf("failed to read Outbound Attempt Limit %s | error: %s", d.Id(), getErr), resp))
			}
			return retry.NonRetryableError(util.BuildWithRetriesApiDiagnosticError(ResourceType, fmt.Sprintf("failed to read Outbound Attempt Limit %s | error: %s", d.Id(), getErr), resp))
		}

		if sdkAttemptLimits.Name != nil {
			_ = d.Set("name", *sdkAttemptLimits.Name)
		}
		if sdkAttemptLimits.MaxAttemptsPerContact != nil {
			_ = d.Set("max_attempts_per_contact", *sdkAttemptLimits.MaxAttemptsPerContact)
		}
		if sdkAttemptLimits.MaxAttemptsPerNumber != nil {
			_ = d.Set("max_attempts_per_number", *sdkAttemptLimits.MaxAttemptsPerNumber)
		}
		if sdkAttemptLimits.TimeZoneId != nil {
			_ = d.Set("time_zone_id", *sdkAttemptLimits.TimeZoneId)
		}
		if sdkAttemptLimits.ResetPeriod != nil {
			_ = d.Set("reset_period", *sdkAttemptLimits.ResetPeriod)
		}

		if sdkAttemptLimits.RecallEntries != nil && len(*sdkAttemptLimits.RecallEntries) > 0 {
			_ = d.Set("recall_entries", flattenSdkOutboundAttemptLimitRecallEntry(sdkAttemptLimits.RecallEntries))
		}

		log.Printf("Read Outbound Attempt Limit %s %s", d.Id(), *sdkAttemptLimits.Name)
		return cc.CheckState(d)
	})
}

func deleteOutboundAttemptLimit(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	proxy := getOutboundAttemptLimitProxy(sdkConfig)

	diagErr := util.RetryWhen(util.IsStatus400, func() (*platformclientv2.APIResponse, diag.Diagnostics) {
		log.Printf("Deleting Outbound Attempt Limit")
		resp, err := proxy.deleteOutboundAttemptLimit(ctx, d.Id())
		if err != nil {
			return resp, util.BuildAPIDiagnosticError(ResourceType, fmt.Sprintf("Failed to delete outbound attempt limit %s error: %s", d.Id(), err), resp)
		}
		return resp, nil
	})
	if diagErr != nil {
		return diagErr
	}

	return util.WithRetries(ctx, 30*time.Second, func() *retry.RetryError {
		_, resp, err := proxy.getOutboundAttemptLimitById(ctx, d.Id())
		if err != nil {
			if util.IsStatus404(resp) {
				log.Printf("Deleted Outbound Attempt Limit %s", d.Id())
				return nil
			}
			return retry.NonRetryableError(util.BuildWithRetriesApiDiagnosticError(ResourceType, fmt.Sprintf("error deleting Outbound Attempt Limit %s | error: %s", d.Id(), err), resp))
		}
		return retry.RetryableError(util.BuildWithRetriesApiDiagnosticError(ResourceType, fmt.Sprintf("Outbound Attempt Limit %s still exists", d.Id()), resp))
	})
}
