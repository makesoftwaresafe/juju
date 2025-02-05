// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lifeflag

import (
	"github.com/juju/names/v4"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/common/lifeflag"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/core/watcher"
)

const (
	// ErrEntityNotFound is a convenience define of the
	// lifeflag.ErrEntityNotFound error. This define makes it so users are not
	// bound to the internal implementation details of this api client.
	ErrEntityNotFound = lifeflag.ErrEntityNotFound
)

// Client is the client used for connecting to the life flag facade.
type Client interface {
	Life(names.Tag) (life.Value, error)
	Watch(names.Tag) (watcher.NotifyWatcher, error)
}

// NewClient creates a new life flag client.
func NewClient(caller base.APICaller) Client {
	return lifeflag.NewClient(caller, "LifeFlag")
}
