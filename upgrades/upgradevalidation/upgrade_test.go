// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgradevalidation_test

import (
	"github.com/juju/names/v4"
	"github.com/juju/replicaset/v2"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/version/v2"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	environscloudspec "github.com/juju/juju/environs/cloudspec"
	"github.com/juju/juju/provider/lxd"
	"github.com/juju/juju/state"
	coretesting "github.com/juju/juju/testing"
	"github.com/juju/juju/upgrades/upgradevalidation"
	"github.com/juju/juju/upgrades/upgradevalidation/mocks"
)

func (s *upgradeValidationSuite) TestValidatorsForControllerUpgradeJuju3(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	s.PatchValue(&upgradevalidation.MinMajorUpgradeVersions, map[int]version.Number{
		3: version.MustParse("2.9.1"),
	})

	ctrlModelTag := names.NewModelTag("deadpork-0bad-400d-8000-4b1d0d06f00d")
	model1ModelTag := coretesting.ModelTag
	statePool := mocks.NewMockStatePool(ctrl)

	ctrlState := mocks.NewMockState(ctrl)
	ctrlModel := mocks.NewMockModel(ctrl)

	state1 := mocks.NewMockState(ctrl)
	model1 := mocks.NewMockModel(ctrl)

	server := mocks.NewMockServer(ctrl)
	serverFactory := mocks.NewMockServerFactory(ctrl)
	s.PatchValue(&upgradevalidation.NewServerFactory,
		func(_ lxd.NewHTTPClientFunc) lxd.ServerFactory {
			return serverFactory
		},
	)
	cloudSpec := environscloudspec.CloudSpec{Type: "lxd"}

	gomock.InOrder(
		// 1. Check controller model.
		// - check agent version;
		ctrlModel.EXPECT().AgentVersion().Return(version.MustParse("2.9.1"), nil),
		// - check mongo status;
		ctrlState.EXPECT().MongoCurrentStatus().Return(&replicaset.Status{
			Members: []replicaset.MemberStatus{
				{
					Id:      1,
					Address: "1.1.1.1",
					State:   replicaset.PrimaryState,
				},
				{
					Id:      2,
					Address: "2.2.2.2",
					State:   replicaset.SecondaryState,
				},
				{
					Id:      3,
					Address: "3.3.3.3",
					State:   replicaset.SecondaryState,
				},
			},
		}, nil),
		// - check mongo version;
		statePool.EXPECT().MongoVersion().Return("4.4", nil),
		// - check if the model has win machines;
		ctrlState.EXPECT().MachineCountForSeries(
			"win10", "win2008r2", "win2012", "win2012hv", "win2012hvr2", "win2012r2",
			"win2016", "win2016hv", "win2019", "win7", "win8", "win81",
		).Return(nil, nil),
		ctrlState.EXPECT().MachineCountForSeries(
			"artful",
			"bionic",
			"cosmic",
			"disco",
			"eoan",
			"groovy",
			"hirsute",
			"impish",
			"kinetic",
			"lunar",
			"mantic",
			"precise",
			"quantal",
			"raring",
			"saucy",
			"trusty",
			"utopic",
			"vivid",
			"wily",
			"xenial",
			"yakkety",
			"zesty",
		).Return(nil, nil),
		// - check LXD version.
		serverFactory.EXPECT().RemoteServer(cloudSpec).Return(server, nil),
		server.EXPECT().ServerVersion().Return("5.2"),

		// 2. Check hosted models.
		// - check agent version;
		model1.EXPECT().AgentVersion().Return(version.MustParse("2.9.1"), nil),
		//  - check if model migration is ongoing;
		model1.EXPECT().MigrationMode().Return(state.MigrationModeNone),
		// - check if the model has win machines;
		state1.EXPECT().MachineCountForSeries(
			"win10", "win2008r2", "win2012", "win2012hv", "win2012hvr2", "win2012r2",
			"win2016", "win2016hv", "win2019", "win7", "win8", "win81",
		).Return(nil, nil),
		state1.EXPECT().MachineCountForSeries(
			"artful",
			"bionic",
			"cosmic",
			"disco",
			"eoan",
			"groovy",
			"hirsute",
			"impish",
			"kinetic",
			"lunar",
			"mantic",
			"precise",
			"quantal",
			"raring",
			"saucy",
			"trusty",
			"utopic",
			"vivid",
			"wily",
			"xenial",
			"yakkety",
			"zesty",
		).Return(nil, nil),
		// - check LXD version.
		serverFactory.EXPECT().RemoteServer(cloudSpec).Return(server, nil),
		server.EXPECT().ServerVersion().Return("5.2"),
	)

	targetVersion := version.MustParse("3.0.0")
	validators := upgradevalidation.ValidatorsForControllerUpgrade(true, targetVersion, cloudSpec)
	checker := upgradevalidation.NewModelUpgradeCheck(ctrlModelTag.Id(), statePool, ctrlState, ctrlModel, validators...)
	blockers, err := checker.Validate()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(blockers, gc.IsNil)

	validators = upgradevalidation.ValidatorsForControllerUpgrade(false, targetVersion, cloudSpec)
	checker = upgradevalidation.NewModelUpgradeCheck(model1ModelTag.Id(), statePool, state1, model1, validators...)
	blockers, err = checker.Validate()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(blockers, gc.IsNil)
}

func (s *upgradeValidationSuite) TestValidatorsForControllerUpgradeJuju2(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	s.PatchValue(&upgradevalidation.MinMajorUpgradeVersions, map[int]version.Number{
		3: version.MustParse("2.9.1"),
	})

	ctrlModelTag := names.NewModelTag("deadpork-0bad-400d-8000-4b1d0d06f00d")
	model1ModelTag := coretesting.ModelTag
	statePool := mocks.NewMockStatePool(ctrl)

	ctrlState := mocks.NewMockState(ctrl)
	ctrlModel := mocks.NewMockModel(ctrl)

	state1 := mocks.NewMockState(ctrl)
	model1 := mocks.NewMockModel(ctrl)

	gomock.InOrder(
		// 1. Check controller model.
		// - check agent version;
		ctrlModel.EXPECT().AgentVersion().Return(version.MustParse("2.9.1"), nil),
		// - check mongo status;
		ctrlState.EXPECT().MongoCurrentStatus().Return(&replicaset.Status{
			Members: []replicaset.MemberStatus{
				{
					Id:      1,
					Address: "1.1.1.1",
					State:   replicaset.PrimaryState,
				},
				{
					Id:      2,
					Address: "2.2.2.2",
					State:   replicaset.SecondaryState,
				},
				{
					Id:      3,
					Address: "3.3.3.3",
					State:   replicaset.SecondaryState,
				},
			},
		}, nil),

		// 2. Check hosted models.
		// - check agent version;
		model1.EXPECT().AgentVersion().Return(version.MustParse("2.9.1"), nil),
		//  - check if model migration is ongoing;
		model1.EXPECT().MigrationMode().Return(state.MigrationModeNone),
	)

	targetVersion := version.MustParse("2.9.99")
	validators := upgradevalidation.ValidatorsForControllerUpgrade(true, targetVersion, environscloudspec.CloudSpec{Type: "lxd"})
	checker := upgradevalidation.NewModelUpgradeCheck(ctrlModelTag.Id(), statePool, ctrlState, ctrlModel, validators...)
	blockers, err := checker.Validate()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(blockers, gc.IsNil)

	validators = upgradevalidation.ValidatorsForControllerUpgrade(false, targetVersion, environscloudspec.CloudSpec{Type: "lxd"})
	checker = upgradevalidation.NewModelUpgradeCheck(model1ModelTag.Id(), statePool, state1, model1, validators...)
	blockers, err = checker.Validate()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(blockers, gc.IsNil)
}

func (s *upgradeValidationSuite) TestValidatorsForModelUpgradeJuju3(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	s.PatchValue(&upgradevalidation.MinMajorUpgradeVersions, map[int]version.Number{
		3: version.MustParse("2.9.1"),
	})

	modelTag := coretesting.ModelTag
	statePool := mocks.NewMockStatePool(ctrl)
	state := mocks.NewMockState(ctrl)
	model := mocks.NewMockModel(ctrl)

	server := mocks.NewMockServer(ctrl)
	serverFactory := mocks.NewMockServerFactory(ctrl)
	s.PatchValue(&upgradevalidation.NewServerFactory,
		func(_ lxd.NewHTTPClientFunc) lxd.ServerFactory {
			return serverFactory
		},
	)
	cloudSpec := environscloudspec.CloudSpec{Type: "lxd"}

	gomock.InOrder(
		// - check no upgrade series in process.
		state.EXPECT().HasUpgradeSeriesLocks().Return(false, nil),
		// - check if the model has win machines.
		state.EXPECT().MachineCountForSeries(
			"win10", "win2008r2", "win2012", "win2012hv", "win2012hvr2", "win2012r2",
			"win2016", "win2016hv", "win2019", "win7", "win8", "win81",
		).Return(nil, nil),
		state.EXPECT().MachineCountForSeries(
			"artful",
			"bionic",
			"cosmic",
			"disco",
			"eoan",
			"groovy",
			"hirsute",
			"impish",
			"kinetic",
			"lunar",
			"mantic",
			"precise",
			"quantal",
			"raring",
			"saucy",
			"trusty",
			"utopic",
			"vivid",
			"wily",
			"xenial",
			"yakkety",
			"zesty",
		).Return(nil, nil),
		// - check LXD version.
		serverFactory.EXPECT().RemoteServer(cloudSpec).Return(server, nil),
		server.EXPECT().ServerVersion().Return("5.2"),
	)

	targetVersion := version.MustParse("3.0.0")
	validators := upgradevalidation.ValidatorsForModelUpgrade(false, targetVersion, cloudSpec)
	checker := upgradevalidation.NewModelUpgradeCheck(modelTag.Id(), statePool, state, model, validators...)
	blockers, err := checker.Validate()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(blockers, gc.IsNil)
}

func (s *upgradeValidationSuite) TestValidatorsForModelUpgradeJuju2(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	s.PatchValue(&upgradevalidation.MinMajorUpgradeVersions, map[int]version.Number{
		3: version.MustParse("2.9.1"),
	})

	modelTag := coretesting.ModelTag
	statePool := mocks.NewMockStatePool(ctrl)
	state := mocks.NewMockState(ctrl)
	model := mocks.NewMockModel(ctrl)

	gomock.InOrder(
		// - check no upgrade series in process.
		state.EXPECT().HasUpgradeSeriesLocks().Return(false, nil),
	)

	targetVersion := version.MustParse("2.9.99")
	validators := upgradevalidation.ValidatorsForModelUpgrade(false, targetVersion, environscloudspec.CloudSpec{Type: "lxd"})
	checker := upgradevalidation.NewModelUpgradeCheck(modelTag.Id(), statePool, state, model, validators...)
	blockers, err := checker.Validate()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(blockers, gc.IsNil)
}
