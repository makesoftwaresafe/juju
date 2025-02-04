// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter_test

import (
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api/agent/uniter"
	"github.com/juju/juju/api/base/testing"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/rpc/params"
	coretesting "github.com/juju/juju/testing"
)

type stateSuite struct {
	coretesting.BaseSuite
}

var _ = gc.Suite(&stateSuite{})

func (s *stateSuite) TestProviderType(c *gc.C) {
	apiCaller := testing.APICallerFunc(func(objType string, version int, id, request string, arg, result interface{}) error {
		c.Assert(objType, gc.Equals, "Uniter")
		c.Assert(request, gc.Equals, "ProviderType")
		c.Assert(arg, gc.IsNil)
		c.Assert(result, gc.FitsTypeOf, &params.StringResult{})
		*(result.(*params.StringResult)) = params.StringResult{
			Result: "somecloud",
		}
		return nil
	})
	client := uniter.NewState(apiCaller, names.NewUnitTag("mysql/0"))

	providerType, err := client.ProviderType()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(providerType, gc.Equals, "somecloud")
}

func (s *stateSuite) TestAllMachinePorts(c *gc.C) {
	apiCaller := testing.APICallerFunc(func(objType string, version int, id, request string, arg, result interface{}) error {
		c.Assert(objType, gc.Equals, "Uniter")
		c.Assert(request, gc.Equals, "AllMachinePorts")
		c.Assert(arg, gc.DeepEquals, params.Entities{Entities: []params.Entity{{Tag: "machine-666"}}})
		c.Assert(result, gc.FitsTypeOf, &params.MachinePortsResults{})
		*(result.(*params.MachinePortsResults)) = params.MachinePortsResults{
			Results: []params.MachinePortsResult{{
				Ports: []params.MachinePortRange{{
					UnitTag:     "unit-mysql-0",
					RelationTag: "",
					PortRange:   params.PortRange{100, 200, "tcp"},
				}, {
					UnitTag:     "unit-mysql-1",
					RelationTag: "",
					PortRange:   params.PortRange{10, 20, "udp"},
				}},
			}},
		}
		return nil
	})
	caller := testing.BestVersionCaller{apiCaller, 1}
	client := uniter.NewState(caller, names.NewUnitTag("mysql/0"))

	portsMap, err := client.AllMachinePorts(names.NewMachineTag("666"))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(portsMap, jc.DeepEquals, map[network.PortRange]params.RelationUnit{
		{100, 200, "tcp"}: {Unit: "unit-mysql-0"},
		{10, 20, "udp"}:   {Unit: "unit-mysql-1"},
	})
}

func (s *stateSuite) TestOpenedMachinePortRangesByEndpoint(c *gc.C) {
	apiCaller := testing.APICallerFunc(func(objType string, version int, id, request string, arg, result interface{}) error {
		c.Assert(objType, gc.Equals, "Uniter")
		c.Assert(request, gc.Equals, "OpenedMachinePortRangesByEndpoint")
		c.Assert(arg, gc.DeepEquals, params.Entities{Entities: []params.Entity{{Tag: "machine-42"}}})
		c.Assert(result, gc.FitsTypeOf, &params.OpenMachinePortRangesByEndpointResults{})
		*(result.(*params.OpenMachinePortRangesByEndpointResults)) = params.OpenMachinePortRangesByEndpointResults{
			Results: []params.OpenMachinePortRangesByEndpointResult{
				{
					UnitPortRanges: map[string][]params.OpenUnitPortRangesByEndpoint{
						"unit-mysql-0": {
							{
								Endpoint:   "",
								PortRanges: []params.PortRange{{100, 200, "tcp"}},
							},
							{
								Endpoint:   "server",
								PortRanges: []params.PortRange{{3306, 3306, "tcp"}},
							},
						},
						"unit-wordpress-0": {
							{
								Endpoint:   "monitoring-port",
								PortRanges: []params.PortRange{{1337, 1337, "udp"}},
							},
						},
					},
				},
			},
		}
		return nil
	})
	caller := testing.BestVersionCaller{apiCaller, 17}
	client := uniter.NewState(caller, names.NewUnitTag("mysql/0"))

	portRangesMap, err := client.OpenedMachinePortRangesByEndpoint(names.NewMachineTag("42"))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(portRangesMap, jc.DeepEquals, map[names.UnitTag]network.GroupedPortRanges{
		names.NewUnitTag("mysql/0"): {
			"":       []network.PortRange{network.MustParsePortRange("100-200/tcp")},
			"server": []network.PortRange{network.MustParsePortRange("3306/tcp")},
		},
		names.NewUnitTag("wordpress/0"): {
			"monitoring-port": []network.PortRange{network.MustParsePortRange("1337/udp")},
		},
	})
}
