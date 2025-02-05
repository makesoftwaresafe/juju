// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/juju/charm/v8"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/version/v2"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"
)

type strategySuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&strategySuite{})

func (s strategySuite) TestValidate(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	curl := charm.MustParseURL("cs:redis-0")

	mockStore := NewMockStore(ctrl)
	mockStore.EXPECT().Validate(curl).Return(nil)

	strategy := &Strategy{
		charmURL: curl,
		store:    mockStore,
		logger:   &fakeLogger{},
	}
	err := strategy.Validate()
	c.Assert(err, jc.ErrorIsNil)
}

func (s strategySuite) TestValidateWithError(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	curl := charm.MustParseURL("cs:redis-0")

	mockStore := NewMockStore(ctrl)
	mockStore.EXPECT().Validate(curl).Return(errors.New("boom"))

	strategy := &Strategy{
		charmURL: curl,
		store:    mockStore,
		logger:   &fakeLogger{},
	}
	err := strategy.Validate()
	c.Assert(err, gc.ErrorMatches, "boom")
}

func (s strategySuite) TestDownloadResult(c *gc.C) {
	file, err := ioutil.TempFile("", "foo")
	c.Assert(err, jc.ErrorIsNil)

	_, _ = fmt.Fprintln(file, "meshuggah")
	err = file.Sync()
	c.Assert(err, jc.ErrorIsNil)

	strategy := &Strategy{logger: &fakeLogger{}}
	result, err := strategy.downloadResult(file.Name(), AlwaysMatchChecksum)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.SHA256, gc.Equals, "4e97ed7423be2ea12939e8fdd592cfb3dcd4d0097d7d193ef998ab6b4db70461")
	c.Assert(result.Size, gc.Equals, int64(10))
}

func (s strategySuite) TestDownloadResultWithOpenError(c *gc.C) {
	strategy := &Strategy{}
	_, err := strategy.downloadResult("foo-123", AlwaysMatchChecksum)
	c.Assert(err, gc.ErrorMatches, "cannot read downloaded charm: open foo-123: no such file or directory")
}

func (s strategySuite) TestRunWithCharmAlreadyUploaded(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	curl := charm.MustParseURL("cs:redis-0")

	mockStore := NewMockStore(ctrl)
	mockStore.EXPECT().DownloadOrigin(curl, gomock.AssignableToTypeOf(Origin{})).DoAndReturn(
		func(curl *charm.URL, origin Origin) (Origin, error) {
			return origin, nil
		},
	)
	mockVersionValidator := NewMockJujuVersionValidator(ctrl)

	mockStateCharm := NewMockStateCharm(ctrl)
	mockStateCharm.EXPECT().IsUploaded().Return(true)

	mockState := NewMockState(ctrl)
	mockState.EXPECT().PrepareCharmUpload(curl).Return(mockStateCharm, nil)

	strategy := &Strategy{
		charmURL: curl,
		store:    mockStore,
		logger:   &fakeLogger{},
	}
	_, alreadyExists, obtainedOrigin, err := strategy.Run(mockState, mockVersionValidator, Origin{Source: CharmHub})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(alreadyExists, jc.IsTrue)
	c.Assert(obtainedOrigin, jc.DeepEquals, Origin{
		Source: CharmHub,
		Platform: Platform{
			Architecture: "amd64",
		},
	})
}

func (s strategySuite) TestRunWithPrepareUploadError(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	curl := charm.MustParseURL("cs:redis-0")

	mockStore := NewMockStore(ctrl)
	mockVersionValidator := NewMockJujuVersionValidator(ctrl)

	mockState := NewMockState(ctrl)
	mockState.EXPECT().PrepareCharmUpload(curl).Return(nil, errors.New("boom"))

	strategy := &Strategy{
		charmURL: curl,
		store:    mockStore,
		logger:   &fakeLogger{},
	}
	_, alreadyExists, _, err := strategy.Run(mockState, mockVersionValidator, Origin{})
	c.Assert(err, gc.ErrorMatches, "boom")
	c.Assert(alreadyExists, jc.IsFalse)
}

func (s strategySuite) TestRun(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	curl := charm.MustParseURL("cs:redis-0")
	meta := &charm.Meta{
		MinJujuVersion: version.Number{Major: 2},
	}

	mockVersionValidator := NewMockJujuVersionValidator(ctrl)
	mockVersionValidator.EXPECT().Validate(meta).Return(nil)

	mockStateCharm := NewMockStateCharm(ctrl)
	mockStateCharm.EXPECT().IsUploaded().Return(false)

	mockStoreCharm := NewMockStoreCharm(ctrl)
	mockStoreCharm.EXPECT().Meta().Return(meta)

	// We're replicating a charm without a LXD profile here and ensuring it
	// correctly handles nil.
	mockStoreCharm.EXPECT().LXDProfile().Return(nil)

	mockState := NewMockState(ctrl)
	mockState.EXPECT().PrepareCharmUpload(curl).Return(mockStateCharm, nil)

	mockStore := NewMockStore(ctrl)
	mockStore.EXPECT().Download(curl, gomock.Any(), gomock.AssignableToTypeOf(Origin{})).DoAndReturn(mustWriteToTempFile(c, mockStoreCharm))

	strategy := &Strategy{
		charmURL: curl,
		store:    mockStore,
		logger:   &fakeLogger{},
	}
	_, alreadyExists, _, err := strategy.Run(mockState, mockVersionValidator, Origin{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(alreadyExists, jc.IsFalse)
}

func (s strategySuite) TestRunWithPlatform(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	curl := charm.MustParseURL("cs:redis-0")
	meta := &charm.Meta{
		MinJujuVersion: version.Number{Major: 2},
	}

	mockVersionValidator := NewMockJujuVersionValidator(ctrl)
	mockVersionValidator.EXPECT().Validate(meta).Return(nil)

	mockStateCharm := NewMockStateCharm(ctrl)
	mockStateCharm.EXPECT().IsUploaded().Return(false)

	mockStoreCharm := NewMockStoreCharm(ctrl)
	mockStoreCharm.EXPECT().Meta().Return(meta)

	// We're replicating a charm without a LXD profile here and ensuring it
	// correctly handles nil.
	mockStoreCharm.EXPECT().LXDProfile().Return(nil)

	mockState := NewMockState(ctrl)
	mockState.EXPECT().PrepareCharmUpload(curl).Return(mockStateCharm, nil)

	mockStore := NewMockStore(ctrl)
	mockStore.EXPECT().Download(curl, gomock.Any(), gomock.AssignableToTypeOf(Origin{})).DoAndReturn(mustWriteToTempFile(c, mockStoreCharm))

	strategy := &Strategy{
		charmURL: curl,
		store:    mockStore,
		logger:   &fakeLogger{},
	}
	_, alreadyExists, origin, err := strategy.Run(mockState, mockVersionValidator, Origin{
		Platform: Platform{
			Channel: "20.04",
			OS:      "Ubuntu",
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(alreadyExists, jc.IsFalse)
	c.Assert(origin, gc.DeepEquals, Origin{
		Platform: Platform{
			Architecture: "amd64",
			Channel:      "20.04",
			OS:           "ubuntu", // notice lower case
		},
	})
}

func (s strategySuite) TestRunWithInvalidLXDProfile(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	curl := charm.MustParseURL("cs:redis-0")
	meta := &charm.Meta{
		MinJujuVersion: version.Number{Major: 2},
	}

	mockVersionValidator := NewMockJujuVersionValidator(ctrl)
	mockVersionValidator.EXPECT().Validate(meta).Return(nil)

	mockStateCharm := NewMockStateCharm(ctrl)
	mockStateCharm.EXPECT().IsUploaded().Return(false)

	mockStoreCharm := NewMockStoreCharm(ctrl)
	mockStoreCharm.EXPECT().Meta().Return(meta)

	// Handle a failure from LXDProfiles
	lxdProfile := &charm.LXDProfile{
		Config: map[string]string{
			"boot": "",
		},
	}

	mockStoreCharm.EXPECT().LXDProfile().Return(lxdProfile)

	mockState := NewMockState(ctrl)
	mockState.EXPECT().PrepareCharmUpload(curl).Return(mockStateCharm, nil)

	mockStore := NewMockStore(ctrl)
	mockStore.EXPECT().Download(curl, gomock.Any(), gomock.AssignableToTypeOf(Origin{})).DoAndReturn(mustWriteToTempFile(c, mockStoreCharm))

	strategy := &Strategy{
		charmURL: curl,
		store:    mockStore,
		logger:   &fakeLogger{},
	}
	_, alreadyExists, _, err := strategy.Run(mockState, mockVersionValidator, Origin{})
	c.Assert(err, gc.ErrorMatches, `cannot add charm: invalid lxd-profile.yaml: contains config value "boot"`)
	c.Assert(alreadyExists, jc.IsFalse)
}

func (s strategySuite) TestFinishAfterRun(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	curl := charm.MustParseURL("cs:redis-0")
	meta := &charm.Meta{
		MinJujuVersion: version.Number{Major: 2},
	}

	mockVersionValidator := NewMockJujuVersionValidator(ctrl)
	mockVersionValidator.EXPECT().Validate(meta).Return(nil)

	mockStateCharm := NewMockStateCharm(ctrl)
	mockStateCharm.EXPECT().IsUploaded().Return(false)

	mockStoreCharm := NewMockStoreCharm(ctrl)
	mockStoreCharm.EXPECT().Meta().Return(meta)
	mockStoreCharm.EXPECT().LXDProfile().Return(nil)

	mockState := NewMockState(ctrl)
	mockState.EXPECT().PrepareCharmUpload(curl).Return(mockStateCharm, nil)

	var tmpFile string

	mockStore := NewMockStore(ctrl)
	mockStore.EXPECT().Download(curl, gomock.Any(), gomock.AssignableToTypeOf(Origin{})).DoAndReturn(
		func(curl *charm.URL, file string, origin Origin) (StoreCharm, ChecksumCheckFn, Origin, error) {
			tmpFile = file
			return mustWriteToTempFile(c, mockStoreCharm)(curl, file, origin)
		},
	)

	strategy := &Strategy{
		charmURL: curl,
		store:    mockStore,
		logger:   &fakeLogger{},
	}
	_, alreadyExists, _, err := strategy.Run(mockState, mockVersionValidator, Origin{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(alreadyExists, jc.IsFalse)

	err = strategy.Finish()
	c.Assert(err, jc.ErrorIsNil)

	_, err = os.Stat(tmpFile)
	c.Assert(os.IsNotExist(err), jc.IsTrue)
}

func mustWriteToTempFile(c *gc.C, mockCharm *MockStoreCharm) func(*charm.URL, string, Origin) (StoreCharm, ChecksumCheckFn, Origin, error) {
	return func(curl *charm.URL, file string, origin Origin) (StoreCharm, ChecksumCheckFn, Origin, error) {
		err := ioutil.WriteFile(file, []byte("meshuggah"), 0644)
		c.Assert(err, jc.ErrorIsNil)

		return mockCharm, AlwaysMatchChecksum, origin, nil
	}
}

type fakeLogger struct {
}

func (l *fakeLogger) Errorf(_ string, _ ...interface{})   {}
func (l *fakeLogger) Debugf(_ string, _ ...interface{})   {}
func (l *fakeLogger) Tracef(_ string, _ ...interface{})   {}
func (l *fakeLogger) Warningf(_ string, _ ...interface{}) {}
func (l *fakeLogger) Child(string) Logger {
	return &fakeLogger{}
}
