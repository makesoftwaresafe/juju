// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package resources

import (
	"testing"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

//go:generate go run go.uber.org/mock/mockgen -package mocks -destination mocks/charmhub.go github.com/juju/juju/apiserver/facades/client/resources CharmHub
//go:generate go run go.uber.org/mock/mockgen -package mocks -destination mocks/logger.go github.com/juju/juju/apiserver/facades/client/resources Logger
