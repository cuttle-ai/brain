// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package appctx declares the app context to be passed between apps. It will considerably reduce the arguments in the function signature
package appctx

import "github.com/cuttle-ai/brain/log"

//Discovery interface must be implemented to facilitate the discovery API calls
type Discovery interface {
	//DiscoveryAddress returns the address of the discovery service
	DiscoveryAddress() string
	//DiscoveryToken returns the token to access the discovery service
	DiscoveryToken() string
}

//AppContext should be implemented to facilitate the appcontext passing
type AppContext interface {
	//Logger returns the logger
	Logger() log.Log
	//AccessToken returns the token for the authentication
	AccessToken() string
	Discovery
}

type appCtxImpl struct {
	accesToken       string
	discoveryToken   string
	discoveryAddress string
	l                log.Log
}

//Logger returns the logger
func (a *appCtxImpl) Logger() log.Log {
	return a.l
}

//AccessToken returns the access token of the app context
func (a *appCtxImpl) AccessToken() string {
	return a.accesToken
}

//DiscoveryAddress returns the discovery service address of the app context
func (a *appCtxImpl) DiscoveryAddress() string {
	return a.discoveryAddress
}

//DiscoveryToken returns the discovery service access token of the app context
func (a *appCtxImpl) DiscoveryToken() string {
	return a.discoveryToken
}

//NewAppCtx returns a app context
func NewAppCtx(accessToken, discoverToken, discoveryAddress string) AppContext {
	return &appCtxImpl{
		accesToken:       accessToken,
		discoveryAddress: discoveryAddress,
		discoveryToken:   discoverToken,
		l:                log.NewLogger(),
	}
}

//WithAccessToken will return a app context with all the proerties from base app context except access token
func WithAccessToken(ctx AppContext, accessToken string) AppContext {
	return &appCtxImpl{
		accesToken:       accessToken,
		discoveryAddress: ctx.DiscoveryAddress(),
		discoveryToken:   ctx.DiscoveryToken(),
		l:                ctx.Logger(),
	}
}
