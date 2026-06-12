package outbound_attempt_limit

import (
	"context"
	"fmt"
	"log"

	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"

	"github.com/mypurecloud/platform-client-sdk-go/v191/platformclientv2"
)

/*
The genesyscloud_outbound_attempt_limit_proxy.go file contains the proxy structures and methods that interact
with the Genesys Cloud SDK. We use composition here for each function on the proxy so individual functions can be stubbed
out during testing.
*/

// internalProxy holds a proxy instance that can be used throughout the package
var internalProxy *outboundAttemptLimitProxy

// Type definitions for each func on our proxy so we can easily mock them out later
type createOutboundAttemptLimitFunc func(ctx context.Context, p *outboundAttemptLimitProxy, attemptLimits *platformclientv2.Attemptlimits) (*platformclientv2.Attemptlimits, *platformclientv2.APIResponse, error)
type getAllOutboundAttemptLimitsFunc func(ctx context.Context, p *outboundAttemptLimitProxy) (*[]platformclientv2.Attemptlimits, *platformclientv2.APIResponse, error)
type getOutboundAttemptLimitIdByNameFunc func(ctx context.Context, p *outboundAttemptLimitProxy, name string) (id string, retryable bool, response *platformclientv2.APIResponse, err error)
type getOutboundAttemptLimitByIdFunc func(ctx context.Context, p *outboundAttemptLimitProxy, id string) (attemptLimits *platformclientv2.Attemptlimits, response *platformclientv2.APIResponse, err error)
type updateOutboundAttemptLimitFunc func(ctx context.Context, p *outboundAttemptLimitProxy, id string, attemptLimits *platformclientv2.Attemptlimits) (*platformclientv2.Attemptlimits, *platformclientv2.APIResponse, error)
type deleteOutboundAttemptLimitFunc func(ctx context.Context, p *outboundAttemptLimitProxy, id string) (response *platformclientv2.APIResponse, err error)

// outboundAttemptLimitProxy contains all of the methods that call genesys cloud APIs.
type outboundAttemptLimitProxy struct {
	clientConfig                        *platformclientv2.Configuration
	outboundApi                         *platformclientv2.OutboundApi
	createOutboundAttemptLimitAttr      createOutboundAttemptLimitFunc
	getAllOutboundAttemptLimitsAttr     getAllOutboundAttemptLimitsFunc
	getOutboundAttemptLimitIdByNameAttr getOutboundAttemptLimitIdByNameFunc
	getOutboundAttemptLimitByIdAttr     getOutboundAttemptLimitByIdFunc
	updateOutboundAttemptLimitAttr      updateOutboundAttemptLimitFunc
	deleteOutboundAttemptLimitAttr      deleteOutboundAttemptLimitFunc
}

// newOutboundAttemptLimitProxy initializes the outbound attempt limit proxy with all of the data needed to communicate with Genesys Cloud
func newOutboundAttemptLimitProxy(clientConfig *platformclientv2.Configuration) *outboundAttemptLimitProxy {
	api := platformclientv2.NewOutboundApiWithConfig(clientConfig)
	return &outboundAttemptLimitProxy{
		clientConfig:                        clientConfig,
		outboundApi:                         api,
		createOutboundAttemptLimitAttr:      createOutboundAttemptLimitFn,
		getAllOutboundAttemptLimitsAttr:     getAllOutboundAttemptLimitsFn,
		getOutboundAttemptLimitIdByNameAttr: getOutboundAttemptLimitIdByNameFn,
		getOutboundAttemptLimitByIdAttr:     getOutboundAttemptLimitByIdFn,
		updateOutboundAttemptLimitAttr:      updateOutboundAttemptLimitFn,
		deleteOutboundAttemptLimitAttr:      deleteOutboundAttemptLimitFn,
	}
}

// getOutboundAttemptLimitProxy acts as a singleton for the internalProxy. It also ensures
// that we can still proxy our tests by directly setting internalProxy package variable
func getOutboundAttemptLimitProxy(clientConfig *platformclientv2.Configuration) *outboundAttemptLimitProxy {
	if internalProxy == nil {
		internalProxy = newOutboundAttemptLimitProxy(clientConfig)
	}
	return internalProxy
}

// createOutboundAttemptLimit creates a Genesys Cloud outbound attempt limit
func (p *outboundAttemptLimitProxy) createOutboundAttemptLimit(ctx context.Context, attemptLimits *platformclientv2.Attemptlimits) (*platformclientv2.Attemptlimits, *platformclientv2.APIResponse, error) {
	return p.createOutboundAttemptLimitAttr(ctx, p, attemptLimits)
}

// getAllOutboundAttemptLimits retrieves all Genesys Cloud outbound attempt limits
func (p *outboundAttemptLimitProxy) getAllOutboundAttemptLimits(ctx context.Context) (*[]platformclientv2.Attemptlimits, *platformclientv2.APIResponse, error) {
	return p.getAllOutboundAttemptLimitsAttr(ctx, p)
}

// getOutboundAttemptLimitIdByName returns a single Genesys Cloud outbound attempt limit by name
func (p *outboundAttemptLimitProxy) getOutboundAttemptLimitIdByName(ctx context.Context, name string) (id string, retryable bool, response *platformclientv2.APIResponse, err error) {
	return p.getOutboundAttemptLimitIdByNameAttr(ctx, p, name)
}

// getOutboundAttemptLimitById returns a single Genesys Cloud outbound attempt limit by Id
func (p *outboundAttemptLimitProxy) getOutboundAttemptLimitById(ctx context.Context, id string) (attemptLimits *platformclientv2.Attemptlimits, response *platformclientv2.APIResponse, err error) {
	return p.getOutboundAttemptLimitByIdAttr(ctx, p, id)
}

// updateOutboundAttemptLimit updates a Genesys Cloud outbound attempt limit
func (p *outboundAttemptLimitProxy) updateOutboundAttemptLimit(ctx context.Context, id string, attemptLimits *platformclientv2.Attemptlimits) (*platformclientv2.Attemptlimits, *platformclientv2.APIResponse, error) {
	return p.updateOutboundAttemptLimitAttr(ctx, p, id, attemptLimits)
}

// deleteOutboundAttemptLimit deletes a Genesys Cloud outbound attempt limit by Id
func (p *outboundAttemptLimitProxy) deleteOutboundAttemptLimit(ctx context.Context, id string) (response *platformclientv2.APIResponse, err error) {
	return p.deleteOutboundAttemptLimitAttr(ctx, p, id)
}

// createOutboundAttemptLimitFn is an implementation function for creating a Genesys Cloud outbound attempt limit
func createOutboundAttemptLimitFn(ctx context.Context, p *outboundAttemptLimitProxy, attemptLimits *platformclientv2.Attemptlimits) (*platformclientv2.Attemptlimits, *platformclientv2.APIResponse, error) {
	ctx = provider.EnsureResourceContext(ctx, ResourceType)

	attemptLimit, resp, err := p.outboundApi.PostOutboundAttemptlimits(*attemptLimits)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create outbound attempt limit: %s", err)
	}
	return attemptLimit, resp, nil
}

// getAllOutboundAttemptLimitsFn is the implementation for retrieving all outbound attempt limits in Genesys Cloud
func getAllOutboundAttemptLimitsFn(ctx context.Context, p *outboundAttemptLimitProxy) (*[]platformclientv2.Attemptlimits, *platformclientv2.APIResponse, error) {
	ctx = provider.EnsureResourceContext(ctx, ResourceType)

	var allAttemptLimits []platformclientv2.Attemptlimits
	const pageSize = 100

	attemptLimits, resp, err := p.outboundApi.GetOutboundAttemptlimits(pageSize, 1, true, "", "", "", "")
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get outbound attempt limits: %v", err)
	}
	if attemptLimits.Entities == nil || len(*attemptLimits.Entities) == 0 {
		return &allAttemptLimits, resp, nil
	}

	allAttemptLimits = append(allAttemptLimits, *attemptLimits.Entities...)

	for pageNum := 2; pageNum <= *attemptLimits.PageCount; pageNum++ {
		attemptLimits, resp, err := p.outboundApi.GetOutboundAttemptlimits(pageSize, pageNum, true, "", "", "", "")
		if err != nil {
			return nil, resp, fmt.Errorf("failed to get outbound attempt limits: %v", err)
		}
		if attemptLimits.Entities == nil || len(*attemptLimits.Entities) == 0 {
			break
		}
		allAttemptLimits = append(allAttemptLimits, *attemptLimits.Entities...)
	}

	return &allAttemptLimits, resp, nil
}

// getOutboundAttemptLimitIdByNameFn is an implementation of the function to get a Genesys Cloud outbound attempt limit by name
func getOutboundAttemptLimitIdByNameFn(ctx context.Context, p *outboundAttemptLimitProxy, name string) (id string, retryable bool, response *platformclientv2.APIResponse, err error) {
	ctx = provider.EnsureResourceContext(ctx, ResourceType)

	attemptLimits, resp, err := p.outboundApi.GetOutboundAttemptlimits(100, 1, true, "", name, "", "")
	if err != nil {
		return "", false, resp, fmt.Errorf("error searching outbound attempt limit %s: %v", name, err)
	}

	if attemptLimits.Entities == nil || len(*attemptLimits.Entities) == 0 {
		return "", true, resp, fmt.Errorf("no outbound attempt limit found with name %s", name)
	}

	for _, attemptLimit := range *attemptLimits.Entities {
		if *attemptLimit.Name == name {
			log.Printf("Retrieved the outbound attempt limit id %s by name %s", *attemptLimit.Id, name)
			return *attemptLimit.Id, false, resp, nil
		}
	}

	return "", true, resp, fmt.Errorf("unable to find outbound attempt limit with name %s", name)
}

// getOutboundAttemptLimitByIdFn is an implementation of the function to get a Genesys Cloud outbound attempt limit by Id
func getOutboundAttemptLimitByIdFn(ctx context.Context, p *outboundAttemptLimitProxy, id string) (attemptLimits *platformclientv2.Attemptlimits, response *platformclientv2.APIResponse, err error) {
	ctx = provider.EnsureResourceContext(ctx, ResourceType)

	attemptLimit, resp, err := p.outboundApi.GetOutboundAttemptlimit(id)
	if err != nil {
		return nil, resp, err
	}
	return attemptLimit, resp, nil
}

// updateOutboundAttemptLimitFn is an implementation of the function to update a Genesys Cloud outbound attempt limit
func updateOutboundAttemptLimitFn(ctx context.Context, p *outboundAttemptLimitProxy, id string, attemptLimits *platformclientv2.Attemptlimits) (*platformclientv2.Attemptlimits, *platformclientv2.APIResponse, error) {
	ctx = provider.EnsureResourceContext(ctx, ResourceType)

	current, resp, err := getOutboundAttemptLimitByIdFn(ctx, p, id)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to read outbound attempt limit %s: %s", id, err)
	}
	attemptLimits.Version = current.Version

	attemptLimit, resp, err := p.outboundApi.PutOutboundAttemptlimit(id, *attemptLimits)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update outbound attempt limit %s: %s", id, err)
	}
	return attemptLimit, resp, nil
}

// deleteOutboundAttemptLimitFn is an implementation function for deleting a Genesys Cloud outbound attempt limit
func deleteOutboundAttemptLimitFn(ctx context.Context, p *outboundAttemptLimitProxy, id string) (response *platformclientv2.APIResponse, err error) {
	ctx = provider.EnsureResourceContext(ctx, ResourceType)

	return p.outboundApi.DeleteOutboundAttemptlimit(id)
}
