// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package repository

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	"github.com/juju/charm/v8"
	charmresource "github.com/juju/charm/v8/resource"
	"github.com/juju/collections/set"
	"github.com/juju/errors"
	"gopkg.in/macaroon.v2"

	"github.com/juju/juju/charmhub"
	"github.com/juju/juju/charmhub/transport"
	corecharm "github.com/juju/juju/core/charm"
	coreseries "github.com/juju/juju/core/series"
)

// CharmHubClient describes the API exposed by the charmhub client.
type CharmHubClient interface {
	DownloadAndRead(ctx context.Context, resourceURL *url.URL, archivePath string, options ...charmhub.DownloadOption) (*charm.CharmArchive, error)
	Refresh(ctx context.Context, config charmhub.RefreshConfig) ([]transport.RefreshResponse, error)
}

// CharmHubRepository provides an API for charm-related operations using charmhub.
type CharmHubRepository struct {
	logger Logger
	client CharmHubClient
}

// NewCharmHubRepository returns a new repository instance using the provided
// charmhub client.
func NewCharmHubRepository(logger Logger, chClient CharmHubClient) *CharmHubRepository {
	return &CharmHubRepository{
		logger: logger,
		client: chClient,
	}
}

// ResolveWithPreferredChannel defines a way using the given charm URL and
// charm origin (platform and channel) to locate a matching charm against the
// Charmhub API.
//
// There are a few things to note in the attempt to resolve the charm and it's
// supporting series.
//
//  1. The algorithm for this is terrible. For charmstore lookups, only one
//     request is required, unfortunately for Charmhub the worst case for this
//     will be 2.
//     Most of the initial requests from the client will hit this first time
//     around (think `juju deploy foo`) without a series (client can then
//     determine what to call the real request with) will be default of 2
//     requests.
//  2. Attempting to find the default series will require 2 requests so that
//     we can find the correct charm ID ensure that the default series exists
//     along with the revision.
//  3. In theory we could just return most of this information without the
//     re-request, but we end up with missing data and potential incorrect
//     charm downloads later.
//
// When charmstore goes, we could potentially rework how the client requests
// the store.
func (c *CharmHubRepository) ResolveWithPreferredChannel(charmURL *charm.URL, argOrigin corecharm.Origin, macaroons macaroon.Slice) (*charm.URL, corecharm.Origin, []string, error) {
	c.logger.Tracef("Resolving CharmHub charm %q with origin %v", charmURL, argOrigin)

	requestedOrigin, err := c.validateAndFixOrigin(argOrigin)
	if err != nil {
		return nil, corecharm.Origin{}, nil, err
	}

	// First attempt to find the charm based on the only input provided.
	res, err := c.refreshOne(charmURL, requestedOrigin, macaroons)
	if err != nil {
		return nil, corecharm.Origin{}, nil, errors.Annotatef(err, "resolving with preferred channel")
	}

	// resolvableBases holds a slice of supported bases from the subsequent
	// refresh API call. The bases can inform the consumer of the API about what
	// they can also install *IF* the retry resolution uses a base that doesn't
	// match their requirements. This can happen in the client if the series
	// selection also wants to consider model-config default-series after the
	// call.
	var (
		effectiveChannel  string
		resolvableBases   []corecharm.Platform
		chSuggestedOrigin = requestedOrigin
	)
	switch {
	case res.Error != nil:
		retryResult, err := c.retryResolveWithPreferredChannel(charmURL, requestedOrigin, macaroons, res.Error)
		if err != nil {
			return nil, corecharm.Origin{}, nil, errors.Trace(err)
		}

		res = retryResult.refreshResponse
		resolvableBases = retryResult.bases
		chSuggestedOrigin = retryResult.origin

		// Fill these back on the origin, so that we can fix the issue of
		// bundles passing back "all" on the response type.
		// Note: we can be sure these have at least one, because of the
		// validation logic in retry method.
		requestedOrigin.Platform.OS = resolvableBases[0].OS
		requestedOrigin.Platform.Channel = resolvableBases[0].Channel

		effectiveChannel = res.EffectiveChannel
	case requestedOrigin.Revision != nil && *requestedOrigin.Revision != -1:
		if len(res.Entity.Bases) > 0 {
			for _, v := range res.Entity.Bases {
				resolvableBases = append(resolvableBases, corecharm.NormalisePlatformSeries(corecharm.Platform{
					Architecture: v.Architecture,
					OS:           v.Name,
					Channel:      v.Channel,
				}))
			}
		}
		// Entities installed by revision do not have an effective channel in the data.
		// A charm with revision X can be in multiple different channels.  The user
		// specified channel for future refresh calls, is located in the origin.
		if res.Entity.Type == transport.CharmType && requestedOrigin.Channel != nil {
			effectiveChannel = requestedOrigin.Channel.String()
		} else if res.Entity.Type == transport.BundleType {
			// This is a hack. A bundle does not require a channel moving forward as refresh
			// by bundle is not implemented, however the following code expects a channel.
			effectiveChannel = "stable"
		}
	default:
		effectiveChannel = res.EffectiveChannel
	}

	// Use the channel that was actually picked by the API. This should
	// account for the closed tracks in a given channel.
	channel, err := charm.ParseChannelNormalize(effectiveChannel)
	if err != nil {
		return nil, corecharm.Origin{}, nil, errors.Annotatef(err, "invalid channel for %q", charmURL)
	}

	// Ensure we send the updated charmURL back, with all the correct segments.
	revision := res.Entity.Revision
	resCurl := charmURL.
		WithArchitecture(chSuggestedOrigin.Platform.Architecture).
		WithRevision(revision)
	// TODO(wallyworld) - does charm url still need a series?
	if chSuggestedOrigin.Platform.Channel != "" {
		series, err := coreseries.VersionSeries(chSuggestedOrigin.Platform.Channel)
		if err != nil {
			return nil, corecharm.Origin{}, nil, errors.Trace(err)
		}
		resCurl = resCurl.WithSeries(series)
	}

	// Create a resolved origin.  Keep the original values for ID and Hash, if
	// any were passed in.  ResolveWithPreferredChannel is called for both
	// charms to be deployed, and charms which are being upgraded.
	// Only charms being upgraded will have an ID and Hash. Those values should
	// only ever be updated in DownloadURL.
	resOrigin := corecharm.Origin{
		Source:   requestedOrigin.Source,
		ID:       requestedOrigin.ID,
		Hash:     requestedOrigin.Hash,
		Type:     string(res.Entity.Type),
		Channel:  &channel,
		Revision: &revision,
		Platform: chSuggestedOrigin.Platform,
	}

	outputOrigin, err := sanitizeCharmOrigin(resOrigin, requestedOrigin)
	if err != nil {
		return nil, corecharm.Origin{}, nil, errors.Trace(err)
	}
	c.logger.Tracef("Resolved CharmHub charm %q with origin %v", resCurl, outputOrigin)

	// If the callee of the API defines a series and that series is pick and
	// identified as being selected (think `juju deploy --series`) then we will
	// never have to retry. The API will never give us back any other supported
	// series, so we can just pass back what the callee requested.
	// This is the happy path for resolving a charm.
	//
	// Unfortunately, most deployments will not pass a series flag, so we will
	// have to ask the API to give us back a potential base. The supported
	// bases can be passed back as a slice of supported series. The callee can
	// then determine which base they want to use and deploy that accordingly,
	// without another API request.
	// TODO(juju3) - we should use supported channels not series
	var series string
	if outputOrigin.Platform.Channel != "" {
		series, err = coreseries.VersionSeries(outputOrigin.Platform.Channel)
		if err != nil {
			return nil, corecharm.Origin{}, nil, errors.Trace(err)
		}
	}
	supportedSeries := []string{
		series,
	}
	if len(resolvableBases) > 0 {
		supportedSeries = make([]string, len(resolvableBases))
		for k, base := range resolvableBases {
			series, err = coreseries.VersionSeries(base.Channel)
			if err != nil {
				return nil, corecharm.Origin{}, nil, errors.Trace(err)
			}
			supportedSeries[k] = series
		}
	}

	return resCurl, outputOrigin, supportedSeries, nil
}

// validateAndFixOrigin, validate the origin and maybe fix as follows:
//
//	Platform must have an architecture.
//	Platform can have both an empty Channel AND os.
//	If channel defined and OS an empty string, add data.
//	Platform must have channel if os defined.
func (c *CharmHubRepository) validateAndFixOrigin(origin corecharm.Origin) (corecharm.Origin, error) {
	p := origin.Platform

	if p.Architecture == "" {
		return corecharm.Origin{}, errors.BadRequestf("origin.Platform requires an Architecture")
	}

	if p.OS != "" && p.Channel == "" {
		return corecharm.Origin{}, errors.BadRequestf("origin.Platform requires a Channel, if OS set")
	}

	if p.Channel != "" && p.OS == "" {
		series, err := coreseries.VersionSeries(origin.Platform.Channel)
		if err != nil {
			return corecharm.Origin{}, errors.NewBadRequest(err, "")
		}
		os, err := coreseries.GetOSFromSeries(series)
		if err != nil {
			return corecharm.Origin{}, errors.NewBadRequest(err, "")
		}
		p.OS = strings.ToLower(os.String())
		c.logger.Tracef("origin.platform is missing os value, though channel provided, setting to %q", p.OS)
		origin.Platform = p
	}

	return origin, nil
}

type retryResolveResult struct {
	refreshResponse transport.RefreshResponse
	origin          corecharm.Origin
	bases           []corecharm.Platform
}

// retryResolveWithPreferredChannel will attempt to inspect the transport
// APIError and determine if a retry is possible with the information gathered
// from the error.
func (c *CharmHubRepository) retryResolveWithPreferredChannel(charmURL *charm.URL, origin corecharm.Origin, macaroons macaroon.Slice, resErr *transport.APIError) (*retryResolveResult, error) {
	var (
		err   error
		bases []corecharm.Platform
	)
	switch resErr.Code {
	case transport.ErrorCodeInvalidCharmPlatform, transport.ErrorCodeInvalidCharmBase:
		c.logger.Tracef("Invalid charm base %q %v - Default Base: %v", charmURL, origin, resErr.Extra.DefaultBases)

		if bases, err = c.selectNextBases(resErr.Extra.DefaultBases, origin); err != nil {
			return nil, errors.Annotatef(err, "selecting next bases")
		}

	case transport.ErrorCodeRevisionNotFound:
		c.logger.Tracef("Revision not found %q %v - Releases: %v", charmURL, origin, resErr.Extra.Releases)

		return nil, errors.Annotatef(c.handleRevisionNotFound(resErr.Extra.Releases, origin), "selecting releases")

	default:
		return nil, errors.Errorf("resolving error: %s", resErr.Message)
	}

	if len(bases) == 0 {
		return nil, errors.Wrap(resErr, errors.Errorf("no releases found for channel %q", origin.Channel.String()))
	}
	base := bases[0]

	p := origin.Platform
	p.OS = base.OS
	p.Channel = base.Channel

	origin.Platform = corecharm.NormalisePlatformSeries(p)

	if origin.Platform.Channel == "" {
		return nil, errors.NotValidf("channel for %s", charmURL.Name)
	}

	c.logger.Tracef("Refresh again with %q %v", charmURL, origin)
	res, err := c.refreshOne(charmURL, origin, macaroons)
	if err != nil {
		return nil, errors.Annotatef(err, "retrying")
	}
	if resErr := res.Error; resErr != nil {
		return nil, errors.Errorf("resolving retry error: %s", resErr.Message)
	}
	return &retryResolveResult{
		refreshResponse: res,
		origin:          origin,
		bases:           bases,
	}, nil
}

// DownloadCharm retrieves specified charm from the store and saves its
// contents to the specified path.
func (c *CharmHubRepository) DownloadCharm(charmURL *charm.URL, requestedOrigin corecharm.Origin, macaroons macaroon.Slice, archivePath string) (corecharm.CharmArchive, corecharm.Origin, error) {
	c.logger.Tracef("DownloadCharm %q, origin: %q", charmURL, requestedOrigin)

	// Resolve charm URL to a link to the charm blob and keep track of the
	// actual resolved origin which may be different from the requested one.
	resURL, actualOrigin, err := c.GetDownloadURL(charmURL, requestedOrigin, macaroons)
	if err != nil {
		return nil, corecharm.Origin{}, errors.Trace(err)
	}

	// TODO(achilleasa): pass macaroons to client when charmhub rolls out support for private charms.
	charmArchive, err := c.client.DownloadAndRead(context.TODO(), resURL, archivePath)
	if err != nil {
		return nil, corecharm.Origin{}, errors.Trace(err)
	}

	return charmArchive, actualOrigin, nil
}

// GetDownloadURL returns the url from which to download the CharmHub charm/bundle
// defined by the provided curl and charm origin.  An updated charm origin is
// also returned with the ID and hash for the charm to be downloaded.  If the
// provided charm origin has no ID, it is assumed that the charm is being
// installed, not refreshed.
func (c *CharmHubRepository) GetDownloadURL(charmURL *charm.URL, requestedOrigin corecharm.Origin, macaroons macaroon.Slice) (*url.URL, corecharm.Origin, error) {
	c.logger.Tracef("GetDownloadURL %q, origin: %q", charmURL, requestedOrigin)

	refreshRes, err := c.refreshOne(charmURL, requestedOrigin, macaroons)
	if err != nil {
		return nil, corecharm.Origin{}, errors.Trace(err)
	}
	if refreshRes.Error != nil {
		return nil, corecharm.Origin{}, errors.Errorf("%s: %s", refreshRes.Error.Code, refreshRes.Error.Message)
	}

	resOrigin := requestedOrigin

	// We've called Refresh with the install action.  Now update the
	// charm ID and Hash values saved.  This is the only place where
	// they should be saved.
	resOrigin.ID = refreshRes.Entity.ID
	resOrigin.Hash = refreshRes.Entity.Download.HashSHA256

	durl, err := url.Parse(refreshRes.Entity.Download.URL)
	if err != nil {
		return nil, corecharm.Origin{}, errors.Trace(err)
	}
	outputOrigin, err := sanitizeCharmOrigin(resOrigin, requestedOrigin)
	return durl, outputOrigin, errors.Trace(err)
}

// ListResources returns the resources for a given charm and origin.
func (c *CharmHubRepository) ListResources(charmURL *charm.URL, origin corecharm.Origin, macaroons macaroon.Slice) ([]charmresource.Resource, error) {
	c.logger.Tracef("ListResources %q", charmURL)

	resCurl, resOrigin, _, err := c.ResolveWithPreferredChannel(charmURL, origin, macaroons)
	if isErrSelection(err) {
		var channel string
		if origin.Channel != nil {
			channel = origin.Channel.String()
		}
		return nil, errors.Errorf("unable to locate charm %q with matching channel %q", charmURL.Name, channel)
	} else if err != nil {
		return nil, errors.Trace(err)
	}

	// If a revision is included with an install action, no resources will be
	// returned. Resources are dependent on a channel, a specific revision can
	// be in multiple channels.  refreshOne gives priority to a revision if
	// specified.  ListResources is used by the "charm-resources" cli cmd,
	// therefore specific charm revisions are less important.
	resOrigin.Revision = nil
	resp, err := c.refreshOne(resCurl, resOrigin, macaroons)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results := make([]charmresource.Resource, len(resp.Entity.Resources))
	for i, resource := range resp.Entity.Resources {
		results[i], err = resourceFromRevision(resource)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return results, nil
}

// GetEssentialMetadata resolves each provided MetadataRequest and returns back
// a slice with the results. The results include the minimum set of metadata
// that is required for deploying each charm.
func (c *CharmHubRepository) GetEssentialMetadata(reqs ...corecharm.MetadataRequest) ([]corecharm.EssentialMetadata, error) {
	// Group reqs in batches based on the provided macaroons
	var urlToReqIdx = make(map[*charm.URL]int)
	var reqGroups = make(map[string][]corecharm.MetadataRequest)
	for reqIdx, req := range reqs {
		urlToReqIdx[req.CharmURL] = reqIdx
		if len(req.Macaroons) == 0 {
			reqGroups[""] = append(reqGroups[""], req)
			continue
		}

		// Calculate the concatenated signatures for all specified
		// macaroons, convert them to a hex string and use that as
		// the map key; this allows us to group together requests that
		// reference the same set of macaroons.
		var macSigs []byte
		for _, mac := range req.Macaroons {
			macSigs = append(macSigs, mac.Signature()...)
		}
		macSig := hex.EncodeToString(macSigs)
		reqGroups[macSig] = append(reqGroups[macSig], req)
	}

	// Make a batch request for each group
	var res = make([]corecharm.EssentialMetadata, len(reqs))
	for _, reqGroup := range reqGroups {
		resGroup, err := c.getEssentialMetadataForBatch(reqGroup)
		if err != nil {
			return nil, errors.Trace(err)
		}

		for groupResIdx, groupRes := range resGroup {
			reqURL := reqGroup[groupResIdx].CharmURL
			res[urlToReqIdx[reqURL]] = groupRes
		}
	}

	return res, nil
}

func (c *CharmHubRepository) getEssentialMetadataForBatch(reqs []corecharm.MetadataRequest) ([]corecharm.EssentialMetadata, error) {
	if len(reqs) == 0 {
		return nil, nil
	}

	// TODO(achilleasa): this method is invoked with a batch of requests
	// that share the same set of macaroons. Once charmhub adds support
	// for macaroons they should be extracted from the first request entry
	// and passed to the charmhub client.
	//   macaroons := reqs[0].Macaroons

	resolvedOrigins := make([]corecharm.Origin, len(reqs))
	refreshCfgs := make([]charmhub.RefreshConfig, len(reqs))
	for reqIdx, req := range reqs {
		// TODO(achilleasa): We should add support for resolving origin
		// batches and move this outside the loop.
		_, resolvedOrigin, _, err := c.ResolveWithPreferredChannel(req.CharmURL, req.Origin, req.Macaroons)
		if err != nil {
			return nil, errors.Annotatef(err, "resolving origin for %q", req.CharmURL)
		}

		refreshCfg, err := refreshConfig(req.CharmURL, resolvedOrigin)
		if err != nil {
			return nil, errors.Trace(err)
		}

		resolvedOrigins[reqIdx] = resolvedOrigin
		refreshCfgs[reqIdx] = refreshCfg
	}

	refreshResults, err := c.client.Refresh(context.TODO(), charmhub.RefreshMany(refreshCfgs...))
	if err != nil {
		return nil, errors.Trace(err)
	}

	var metaRes = make([]corecharm.EssentialMetadata, len(reqs))
	for resIdx, refreshResult := range refreshResults {
		if len(refreshResult.Entity.MetadataYAML) == 0 {
			return nil, errors.Errorf("charmhub refresh response for %q does not include the contents of metadata.yaml", reqs[resIdx].CharmURL)
		}
		chMeta, err := charm.ReadMeta(strings.NewReader(refreshResult.Entity.MetadataYAML))
		if err != nil {
			return nil, errors.Annotatef(err, "parsing metadata.yaml for %q", reqs[resIdx].CharmURL)
		}

		if len(refreshResult.Entity.ConfigYAML) == 0 {
			return nil, errors.Errorf("charmhub refresh response for %q does not include the contents of config.yaml", reqs[resIdx].CharmURL)
		}
		chConfig, err := charm.ReadConfig(strings.NewReader(refreshResult.Entity.ConfigYAML))
		if err != nil {
			return nil, errors.Annotatef(err, "parsing config.yaml for %q", reqs[resIdx].CharmURL)
		}

		chManifest := new(charm.Manifest)
		for _, base := range refreshResult.Entity.Bases {
			baseCh, err := charm.ParseChannelNormalize(base.Channel)
			if err != nil {
				return nil, errors.Annotatef(err, "parsing base channel for %q", reqs[resIdx].CharmURL)
			}

			chManifest.Bases = append(chManifest.Bases, charm.Base{
				Name:          base.Name,
				Channel:       baseCh,
				Architectures: []string{base.Architecture},
			})
		}

		resolvedOrigins[resIdx].ID = refreshResult.ID
		metaRes[resIdx].ResolvedOrigin = resolvedOrigins[resIdx]
		metaRes[resIdx].Meta = chMeta
		metaRes[resIdx].Config = chConfig
		metaRes[resIdx].Manifest = chManifest
	}

	return metaRes, nil
}

func (c *CharmHubRepository) refreshOne(charmURL *charm.URL, origin corecharm.Origin, _ macaroon.Slice) (transport.RefreshResponse, error) {
	// TODO(achilleasa): pass macaroons to client when charmhub rolls out support for private charms.
	cfg, err := refreshConfig(charmURL, origin)
	if err != nil {
		return transport.RefreshResponse{}, errors.Trace(err)
	}
	c.logger.Tracef("Locate charm using: %v", cfg)
	result, err := c.client.Refresh(context.TODO(), cfg)
	if err != nil {
		return transport.RefreshResponse{}, errors.Trace(err)
	}
	if len(result) != 1 {
		return transport.RefreshResponse{}, errors.Errorf("more than 1 result found")
	}

	return result[0], nil
}

func (c *CharmHubRepository) selectNextBases(bases []transport.Base, origin corecharm.Origin) ([]corecharm.Platform, error) {
	if len(bases) == 0 {
		return nil, errors.Errorf("no bases available")
	}
	// We've got a invalid charm platform error, the error should contain
	// a valid platform to query again to get the right information. If
	// the platform is empty, consider it a failure.
	var compatible []transport.Base
	for _, base := range bases {
		if base.Architecture != origin.Platform.Architecture {
			continue
		}
		compatible = append(compatible, base)
	}
	if len(compatible) == 0 {
		return nil, errors.NotFoundf("bases matching architecture %q", origin.Platform.Architecture)
	}

	// Serialize all the platforms into core entities.
	results := make([]corecharm.Platform, len(compatible))
	for k, base := range compatible {
		track, err := corecharm.ChannelTrack(base.Channel)
		if err != nil {
			return nil, errors.Annotate(err, "base")
		}
		platform := corecharm.NormalisePlatformSeries(corecharm.Platform{
			Architecture: base.Architecture,
			OS:           base.Name,
			Channel:      track,
		})
		results[k] = platform
	}

	return results, nil
}

func (c *CharmHubRepository) handleRevisionNotFound(releases []transport.Release, origin corecharm.Origin) error {
	if len(releases) == 0 {
		return errors.Errorf("no releases available")
	}
	// If the user passed in a branch, but not enough information about the
	// arch and channel, then we can help by giving a better error message.
	if origin.Channel != nil && origin.Channel.Branch != "" {
		return errors.Errorf("ambiguous arch and series with channel %q, specify both arch and series along with channel", origin.Channel.String())
	}
	// Help the user out by creating a list of channel/base suggestions to try.
	suggestions := c.composeSuggestions(releases, origin)
	var s string
	if len(suggestions) > 0 {
		s = fmt.Sprintf("\navailable releases are:\n  %v", strings.Join(suggestions, "\n  "))
	}
	// If the origin's channel is nil, one wasn't specified by the user,
	// so we requested "stable", which indicates the charm's default channel.
	// However, at the time we're writing this message, we do not know what
	// the charm's default channel is.
	var channelString string
	if origin.Channel != nil {
		channelString = fmt.Sprintf("for channel %q", origin.Channel.String())
	} else {
		channelString = "in the charm's default channel"
	}

	return errSelection{
		err: errors.Errorf(
			"charm or bundle not found %s, base %q%s",
			channelString, origin.Platform.String(), s),
	}
}

type errSelection struct {
	err error
}

func (e errSelection) Error() string {
	return e.err.Error()
}

func isErrSelection(err error) bool {
	_, ok := errors.Cause(err).(errSelection)
	return ok
}

// Method describes the method for requesting the charm using the RefreshAPI.
type Method string

const (
	// MethodRevision utilizes an install action by the revision only. A
	// channel must be in the origin, however it's not used in this request,
	// but saved in the origin for future use.
	MethodRevision Method = "revision"
	// MethodChannel utilizes an install action by the channel only.
	MethodChannel Method = "channel"
	// MethodID utilizes an refresh action by the id, revision and
	// channel (falls back to latest/stable if channel is not found).
	MethodID Method = "id"
)

// refreshConfig creates a RefreshConfig for the given input.
// If the origin.ID is not set, a install refresh config is returned. For
// install. Channel and Revision are mutually exclusive in the api, only
// one will be used.
//
// If the origin.ID is set, a refresh config is returned.
//
// NOTE: There is one idiosyncrasy of this method.  The charm URL and and
// origin have a revision number in them when called by FindDownloadURL
// to install a charm. Potentially causing an unexpected install by revision.
// This is okay as all of the data is ready and correct in the origin.
func refreshConfig(charmURL *charm.URL, origin corecharm.Origin) (charmhub.RefreshConfig, error) {
	// Work out the correct install method.
	rev := -1
	var method Method
	if origin.Revision != nil && *origin.Revision >= 0 {
		rev = *origin.Revision
	}
	if origin.ID == "" && rev != -1 {
		method = MethodRevision
	}

	var (
		channel         string
		nonEmptyChannel = origin.Channel != nil && !origin.Channel.Empty()
	)

	// Select the appropriate channel based on the supplied origin.
	// We need to ensure that we always, always normalize the incoming channel
	// before we hit the refresh API.
	if nonEmptyChannel {
		channel = origin.Channel.Normalize().String()
	} else if method != MethodRevision {
		channel = corecharm.DefaultChannel.Normalize().String()
	}

	if method != MethodRevision && channel != "" {
		method = MethodChannel
	}

	// Bundles can not use method IDs, which in turn forces a refresh.
	if !transport.BundleType.Matches(origin.Type) && origin.ID != "" {
		method = MethodID
	}

	var (
		cfg charmhub.RefreshConfig
		err error

		base = charmhub.RefreshBase{
			Architecture: origin.Platform.Architecture,
			Name:         origin.Platform.OS,
			Channel:      corecharm.ComputeBaseChannel(origin.Platform).Channel,
		}
	)
	switch method {
	case MethodChannel:
		// Install from just the name and the channel. If there is no origin ID,
		// we haven't downloaded this charm before.
		// Try channel first.
		cfg, err = charmhub.InstallOneFromChannel(charmURL.Name, channel, base)
	case MethodRevision:
		// If there is a revision, install it using that. If there is no origin
		// ID, we haven't downloaded this charm before.
		cfg, err = charmhub.InstallOneFromRevision(charmURL.Name, rev)
	case MethodID:
		// This must be a charm upgrade if we have an ID.  Use the refresh
		// action for metric keeping on the CharmHub side.
		cfg, err = charmhub.RefreshOne(origin.InstanceKey, origin.ID, rev, channel, base)
	default:
		return nil, errors.NotValidf("origin %v", origin)
	}
	return cfg, err
}

func (c *CharmHubRepository) composeSuggestions(releases []transport.Release, origin corecharm.Origin) []string {
	charmRisks := set.NewStrings()
	for _, v := range charm.Risks {
		charmRisks.Add(string(v))
	}
	channelSeries := make(map[string][]string)
	for _, release := range releases {
		base := corecharm.NormalisePlatformSeries(corecharm.Platform{
			Architecture: release.Base.Architecture,
			OS:           release.Base.Name,
			Channel:      release.Base.Channel,
		})
		arch := base.Architecture
		track, err := corecharm.ChannelTrack(base.Channel)
		if err != nil {
			c.logger.Errorf("invalid base channel %v: %s", base.Channel, err)
			continue
		}
		if track == "all" {
			track = origin.Platform.Channel
		}
		series, err := coreseries.VersionSeries(track)
		if err != nil {
			c.logger.Errorf("converting version to series: %s", err)
			continue
		}
		if arch == "all" {
			arch = origin.Platform.Architecture
		}
		if arch != origin.Platform.Architecture {
			continue
		}
		// Now that we have default tracks other than latest:
		// If a channel is risk only, add latest as the track
		// to be more clear for the user facing error message.
		// At this point, we do not know the default channel,
		// or if the charm has one, therefore risk only output
		// is ambiguous.
		charmChannel := release.Channel
		if charmRisks.Contains(charmChannel) {
			charmChannel = "latest/" + charmChannel
		}
		channelSeries[charmChannel] = append(channelSeries[charmChannel], series)
	}

	var suggestions []string
	for channel, values := range channelSeries {
		suggestions = append(suggestions, fmt.Sprintf("channel %q: available series are: %s", channel, strings.Join(values, ", ")))
	}
	return suggestions
}

// TODO (stickupkid) - Find a common place for this as it's duplicated from
// apiserver/client/resources
func resourceFromRevision(rev transport.ResourceRevision) (charmresource.Resource, error) {
	resType, err := charmresource.ParseType(rev.Type)
	if err != nil {
		return charmresource.Resource{}, errors.Trace(err)
	}
	fp, err := charmresource.ParseFingerprint(rev.Download.HashSHA384)
	if err != nil {
		return charmresource.Resource{}, errors.Trace(err)
	}
	r := charmresource.Resource{
		Meta: charmresource.Meta{
			Name:        rev.Name,
			Type:        resType,
			Path:        rev.Filename,
			Description: rev.Description,
		},
		Origin:      charmresource.OriginStore,
		Revision:    rev.Revision,
		Fingerprint: fp,
		Size:        int64(rev.Download.Size),
	}
	return r, nil
}
