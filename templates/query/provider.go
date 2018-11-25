/*
* @Author: Bartuccio Antoine
* @Date:   2018-09-16 18:00:04
* @Last Modified by:   Bartuccio Antoine
* @Last Modified time: 2018-11-24 23:12:22
 */

package query

import (
	"fmt"

	"gitlab.com/GT-RIMi/RIMo/media"
)

var providers = make(map[string]Provider)

// Provider interface for providers
type Provider interface {
	IsCompatible(mediaType media.MediaType) bool
	Search(query QueryStruct) QueryResponse
	Discover(query QueryDiscoverStruct) QueryResponse
	Init()
}

// RegisterProvider register a new provider to call
func RegisterProvider(provider Provider, providerName string) error {

	if provider == nil {
		return fmt.Errorf("no provider found for provider name %s", providerName)
	}
	if _, exists := providers[providerName]; exists {
		return fmt.Errorf("provider with name %s already exists", providerName)
	}
	provider.Init()
	providers[providerName] = provider
	return nil
}
