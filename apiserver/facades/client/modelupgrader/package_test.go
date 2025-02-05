// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelupgrader

import (
	stdtesting "testing"

	"github.com/juju/version/v2"

	"github.com/juju/juju/testing"
	coretools "github.com/juju/juju/tools"
)

//go:generate go run go.uber.org/mock/mockgen -package mocks -destination mocks/state_mock.go github.com/juju/juju/apiserver/facades/client/modelupgrader StatePool,State,Model
//go:generate go run go.uber.org/mock/mockgen -package mocks -destination mocks/agents_mock.go github.com/juju/juju/apiserver/common ToolsFinder
//go:generate go run go.uber.org/mock/mockgen -package mocks -destination mocks/environs_mock.go github.com/juju/juju/environs BootstrapEnviron
//go:generate go run go.uber.org/mock/mockgen -package mocks -destination mocks/common_mock.go github.com/juju/juju/apiserver/common BlockCheckerInterface

func TestAll(t *stdtesting.T) {
	testing.MgoTestPackage(t)
}

func (m *ModelUpgraderAPI) FindTools(
	st State, model Model,
	majorVersion, minorVersion int, agentVersion version.Number, osType, arch, agentStream string,
) (coretools.Versions, error) {
	return m.findTools(st, model, majorVersion, minorVersion, agentVersion, osType, arch, agentStream)
}

func (m *ModelUpgraderAPI) DecideVersion(
	targetVersion, agentVersion version.Number, agentStream string, st State, model Model,
) (_ version.Number, err error) {
	return m.decideVersion(targetVersion, agentVersion, agentStream, st, model)
}
