// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package testing contains the testing infrastructure to mock out the LXD API.
// Run 'go generate' to regenerate the mock interfaces.
package testing

//go:generate go run go.uber.org/mock/mockgen -package testing -destination lxd_mock.go -write_package_comment=false github.com/canonical/lxd/client Operation,RemoteOperation,Server,ImageServer,InstanceServer
