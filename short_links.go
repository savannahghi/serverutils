package base

import (
	"context"
	"fmt"

	"google.golang.org/api/firebasedynamiclinks/v1"
)

// FDLDomainEnvironmentVariableName is the name of the domain used for short
// links.
//
// e.g https://healthcloud.page.link or https://bwl.page.link
const FDLDomainEnvironmentVariableName = "FIREBASE_DYNAMIC_LINKS_DOMAIN"

// ServerPublicDomainEnvironmentVariableName is the name of the environment
// variable at which the server is deployed. It is used to generate long
// links for shortening
const ServerPublicDomainEnvironmentVariableName = "SERVER_PUBLIC_DOMAIN"

// ShortenLink shortens an FDL link
func ShortenLink(ctx context.Context, longLink string) (string, error) {
	fdlService, err := firebasedynamiclinks.NewService(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to initialize Firebase Dynamic Links service: %w", err)
	}

	dynamicLinkDomain, err := GetEnvVar(FDLDomainEnvironmentVariableName)
	if err != nil {
		return "", fmt.Errorf("environment variable %s missing", FDLDomainEnvironmentVariableName)
	}

	linkRequest := &firebasedynamiclinks.CreateShortDynamicLinkRequest{
		DynamicLinkInfo: &firebasedynamiclinks.DynamicLinkInfo{
			DomainUriPrefix: dynamicLinkDomain,
			Link:            longLink,
		},
	}
	linkReq := fdlService.ShortLinks.Create(linkRequest)
	linkResp, err := linkReq.Do()
	if err != nil {
		return "", fmt.Errorf("unable to shorten link: %w", err)
	}
	return linkResp.ShortLink, nil
}
