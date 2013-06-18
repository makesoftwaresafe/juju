// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package openstack_test

import (
	"flag"
	. "launchpad.net/gocheck"
	"launchpad.net/goose/identity"
	"launchpad.net/goose/nova"
	"launchpad.net/juju-core/environs/openstack"
	"launchpad.net/juju-core/instance"
	"testing"
)

var live = flag.Bool("live", false, "Include live OpenStack tests")

func Test(t *testing.T) {
	if *live {
		cred, err := identity.CompleteCredentialsFromEnv()
		if err != nil {
			t.Fatalf("Error setting up test suite: %s", err.Error())
		}
		registerLiveTests(cred)
	}
	registerLocalTests()
	TestingT(t)
}

// localTests contains tests which do not require a live service or test double to run.
type localTests struct{}

var _ = Suite(&localTests{})

// ported from lp:juju/juju/providers/openstack/tests/test_machine.py
var addressTests = []struct {
	summary  string
	private  []nova.IPAddress
	public   []nova.IPAddress
	networks []string
	expected string
	failure  error
}{
	{
		summary:  "missing",
		expected: "",
		failure:  instance.ErrNoDNSName,
	},
	{
		summary:  "empty",
		private:  []nova.IPAddress{},
		networks: []string{"private"},
		expected: "",
		failure:  instance.ErrNoDNSName,
	},
	{
		summary:  "private only",
		private:  []nova.IPAddress{{4, "127.0.0.4"}},
		networks: []string{"private"},
		expected: "127.0.0.4",
		failure:  nil,
	},
	{
		summary:  "private plus (HP cloud)",
		private:  []nova.IPAddress{{4, "127.0.0.4"}, {4, "8.8.4.4"}},
		networks: []string{"private"},
		expected: "8.8.4.4",
		failure:  nil,
	},
	{
		summary:  "public only",
		public:   []nova.IPAddress{{4, "8.8.8.8"}},
		networks: []string{"", "public"},
		expected: "8.8.8.8",
		failure:  nil,
	},
	{
		summary:  "public and private",
		private:  []nova.IPAddress{{4, "127.0.0.4"}},
		public:   []nova.IPAddress{{4, "8.8.4.4"}},
		networks: []string{"private", "public"},
		expected: "8.8.4.4",
		failure:  nil,
	},
	{
		summary:  "public private plus",
		private:  []nova.IPAddress{{4, "127.0.0.4"}, {4, "8.8.4.4"}},
		public:   []nova.IPAddress{{4, "8.8.8.8"}},
		networks: []string{"private", "public"},
		expected: "8.8.8.8",
		failure:  nil,
	},
	{
		summary:  "custom only",
		private:  []nova.IPAddress{{4, "127.0.0.2"}},
		networks: []string{"special"},
		expected: "127.0.0.2",
		failure:  nil,
	},
	{
		summary:  "custom and public",
		private:  []nova.IPAddress{{4, "127.0.0.2"}},
		public:   []nova.IPAddress{{4, "8.8.8.8"}},
		networks: []string{"special", "public"},
		expected: "8.8.8.8",
		failure:  nil,
	},
	{
		summary:  "non-IPv4",
		private:  []nova.IPAddress{{6, "::dead:beef:f00d"}},
		networks: []string{"private"},
		expected: "",
		failure:  instance.ErrNoDNSName,
	},
}

func (t *localTests) TestGetServerAddresses(c *C) {
	for i, t := range addressTests {
		c.Logf("#%d. %s -> %s (%v)", i, t.summary, t.expected, t.failure)
		addresses := make(map[string][]nova.IPAddress)
		if t.private != nil {
			if len(t.networks) < 1 {
				addresses["private"] = t.private
			} else {
				addresses[t.networks[0]] = t.private
			}
		}
		if t.public != nil {
			if len(t.networks) < 2 {
				addresses["public"] = t.public
			} else {
				addresses[t.networks[1]] = t.public
			}
		}
		addr, err := openstack.InstanceAddress(addresses)
		c.Assert(err, Equals, t.failure)
		c.Assert(addr, Equals, t.expected)
	}
}
