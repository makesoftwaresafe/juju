// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package context

import (
	stdcontext "context"

	"github.com/juju/errors"
	"github.com/juju/testing"
	gc "gopkg.in/check.v1"
)

type CloudCallContextSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&CloudCallContextSuite{})

func (s *CloudCallContextSuite) TestCloudCallContext(c *gc.C) {
	stdctx := stdcontext.TODO()
	ctx := NewCloudCallContext(stdctx)
	c.Assert(ctx, gc.NotNil)
	c.Assert(ctx.Context, gc.Equals, stdctx)

	err := ctx.InvalidateCredential("call")
	c.Assert(errors.Is(err, errors.NotImplemented), gc.Equals, true)
}
