package outbound_attempt_limit

import (
	"fmt"
	"strings"

	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v191/platformclientv2"
)

func buildSdkOutboundAttemptLimitRecallEntryMap(recallEntries []interface{}) *map[string]platformclientv2.Recallentry {
	if len(recallEntries) == 0 {
		return nil
	}
	recallEntriesMap := map[string]platformclientv2.Recallentry{}
	if entriesMap, ok := recallEntries[0].(map[string]interface{}); ok {
		types := []string{"busy", "no_answer", "answering_machine", "fax"}
		for _, t := range types {
			entrySet := entriesMap[t].(*schema.Set).List()
			if len(entrySet) == 0 {
				continue
			}
			if entryMap, ok := entrySet[0].(map[string]interface{}); ok && len(entryMap) > 0 {
				recallEntriesMap[util.ToCamelCase(t)] = *buildSdkRecallEntry(entryMap)
			}
		}
	}
	return &recallEntriesMap
}

func buildSdkRecallEntry(entry map[string]interface{}) *platformclientv2.Recallentry {
	sdkRecallEntry := platformclientv2.Recallentry{}
	if nbrAttempts, ok := entry["nbr_attempts"].(int); ok {
		sdkRecallEntry.NbrAttempts = &nbrAttempts
	}
	if minsBetweenAttempts, ok := entry["minutes_between_attempts"].(int); ok {
		sdkRecallEntry.MinutesBetweenAttempts = &minsBetweenAttempts
	}
	return &sdkRecallEntry
}

func flattenSdkOutboundAttemptLimitRecallEntry(sdkRecallEntries *map[string]platformclientv2.Recallentry) []interface{} {
	recallEntries := make(map[string]interface{})
	for key, val := range *sdkRecallEntries {
		recallEntries[util.ToSnakeCase(key)] = flattenSdkRecallEntry(val)
	}
	return []interface{}{recallEntries}
}

func flattenSdkRecallEntry(sdkEntry platformclientv2.Recallentry) *schema.Set {
	var (
		entryMap = make(map[string]interface{})
		entrySet = schema.NewSet(schema.HashResource(recallSettings), []interface{}{})
	)
	entryMap["nbr_attempts"] = *sdkEntry.NbrAttempts
	entryMap["minutes_between_attempts"] = *sdkEntry.MinutesBetweenAttempts
	entrySet.Add(entryMap)
	return entrySet
}

func buildSdkAttemptLimits(d *schema.ResourceData) *platformclientv2.Attemptlimits {
	name := d.Get("name").(string)
	maxAttemptsPerContact := d.Get("max_attempts_per_contact").(int)
	maxAttemptsPerNumber := d.Get("max_attempts_per_number").(int)
	timeZoneId := d.Get("time_zone_id").(string)
	resetPeriod := d.Get("reset_period").(string)
	recallEntries := d.Get("recall_entries").([]interface{})

	sdkAttemptLimits := platformclientv2.Attemptlimits{}

	if name != "" {
		sdkAttemptLimits.Name = &name
	}
	if maxAttemptsPerContact != 0 {
		sdkAttemptLimits.MaxAttemptsPerContact = &maxAttemptsPerContact
	}
	if maxAttemptsPerNumber != 0 {
		sdkAttemptLimits.MaxAttemptsPerNumber = &maxAttemptsPerNumber
	}
	if timeZoneId != "" {
		sdkAttemptLimits.TimeZoneId = &timeZoneId
	}
	if resetPeriod != "" {
		sdkAttemptLimits.ResetPeriod = &resetPeriod
	}
	if len(recallEntries) > 0 {
		sdkAttemptLimits.RecallEntries = buildSdkOutboundAttemptLimitRecallEntryMap(recallEntries)
	}

	return &sdkAttemptLimits
}

func GenerateAttemptLimitResource(
	resourceLabel string,
	name string,
	maxAttemptsPerContact string,
	maxAttemptsPerNumber string,
	timeZoneId string,
	resetPeriod string,
	nestedBlocks ...string,
) string {
	if maxAttemptsPerContact != "" {
		maxAttemptsPerContact = fmt.Sprintf(`max_attempts_per_contact = %s`, maxAttemptsPerContact)
	}
	if maxAttemptsPerNumber != "" {
		maxAttemptsPerNumber = fmt.Sprintf(`max_attempts_per_number = %s`, maxAttemptsPerNumber)
	}
	if timeZoneId != "" {
		timeZoneId = fmt.Sprintf(`time_zone_id = "%s"`, timeZoneId)
	}
	if resetPeriod != "" {
		resetPeriod = fmt.Sprintf(`reset_period = "%s"`, resetPeriod)
	}
	return fmt.Sprintf(`
resource "genesyscloud_outbound_attempt_limit" "%s" {
	name = "%s"
	%s
	%s
	%s
	%s
	%s
}
	`, resourceLabel, name, maxAttemptsPerContact, maxAttemptsPerNumber, timeZoneId, resetPeriod, strings.Join(nestedBlocks, "\n"))
}

func GenerateOutboundAttemptLimitDataSource(dataSourceLabel string, attemptLimitName string, dependsOn string) string {
	return fmt.Sprintf(`
data "%s" "%s" {
	name = "%s"
	depends_on = [%s]
}
`, ResourceType, dataSourceLabel, attemptLimitName, dependsOn)
}
