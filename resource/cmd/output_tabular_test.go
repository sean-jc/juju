// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package cmd

import (
	"time"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	charmresource "gopkg.in/juju/charm.v6-unstable/resource"

	"github.com/juju/juju/resource"
)

var _ = gc.Suite(&CharmTabularSuite{})

type CharmTabularSuite struct {
	testing.IsolationSuite
}

func (s *CharmTabularSuite) TestFormatCharmTabularOkay(c *gc.C) {
	res := charmRes(c, "spam", ".tgz", "...", "")
	formatted := []FormattedCharmResource{FormatCharmResource(res)}

	data, err := FormatCharmTabular(formatted)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(string(data), gc.Equals, `
RESOURCE REVISION
spam     1
`[1:])
}

func (s *CharmTabularSuite) TestFormatCharmTabularMinimal(c *gc.C) {
	res := charmRes(c, "spam", "", "", "")
	formatted := []FormattedCharmResource{FormatCharmResource(res)}

	data, err := FormatCharmTabular(formatted)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(string(data), gc.Equals, `
RESOURCE REVISION
spam     1
`[1:])
}

func (s *CharmTabularSuite) TestFormatCharmTabularUpload(c *gc.C) {
	res := charmRes(c, "spam", "", "", "")
	res.Origin = charmresource.OriginUpload
	formatted := []FormattedCharmResource{FormatCharmResource(res)}

	data, err := FormatCharmTabular(formatted)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(string(data), gc.Equals, `
RESOURCE REVISION
spam     1
`[1:])
}

func (s *CharmTabularSuite) TestFormatCharmTabularMulti(c *gc.C) {
	formatted := []FormattedCharmResource{
		FormatCharmResource(charmRes(c, "spam", ".tgz", "spamspamspamspam", "")),
		FormatCharmResource(charmRes(c, "eggs", "", "...", "")),
		FormatCharmResource(charmRes(c, "somethingbig", ".zip", "", "")),
		FormatCharmResource(charmRes(c, "song", ".mp3", "your favorite", "")),
		FormatCharmResource(charmRes(c, "avatar", ".png", "your picture", "")),
	}
	formatted[1].Revision = 2

	data, err := FormatCharmTabular(formatted)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(string(data), gc.Equals, `
RESOURCE     REVISION
spam         1
eggs         2
somethingbig 1
song         1
avatar       1
`[1:])
}

func (s *CharmTabularSuite) TestFormatCharmTabularBadValue(c *gc.C) {
	bogus := "should have been something else"
	_, err := FormatCharmTabular(bogus)

	c.Check(err, gc.ErrorMatches, `expected value of type .*`)
}

var _ = gc.Suite(&SvcTabularSuite{})

type SvcTabularSuite struct {
	testing.IsolationSuite
}

func (s *SvcTabularSuite) TestFormatOkay(c *gc.C) {
	res := resource.Resource{

		Resource: charmresource.Resource{
			Meta: charmresource.Meta{
				Name:    "openjdk",
				Comment: "the java runtime",
			},
			Origin:   charmresource.OriginStore,
			Revision: 7,
		},
		Timestamp: time.Now(),
	}

	formatted := []FormattedSvcResource{FormatSvcResource(res)}

	data, err := FormatSvcTabular(formatted)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(string(data), gc.Equals, `
RESOURCE SUPPLIED BY REVISION COMMENT
openjdk  charmstore  7        the java runtime
`[1:])
}

func (s *SvcTabularSuite) TestFormatCharmTabularMulti(c *gc.C) {
	res := []resource.Resource{
		{
			Resource: charmresource.Resource{
				Meta: charmresource.Meta{
					Name:    "openjdk",
					Comment: "the java runtime",
				},
				Origin:   charmresource.OriginStore,
				Revision: 7,
			},
		},
		{
			Resource: charmresource.Resource{
				Meta: charmresource.Meta{
					Name:    "website",
					Comment: "your website data",
				},
				Origin: charmresource.OriginUpload,
			},
		},
		{
			Resource: charmresource.Resource{
				Meta: charmresource.Meta{
					Name:    "openjdk2",
					Comment: "another java runtime",
				},
				Origin:   charmresource.OriginStore,
				Revision: 8,
			},
			Timestamp: time.Date(2012, 12, 12, 12, 12, 12, 0, time.UTC),
		},
		{
			Resource: charmresource.Resource{
				Meta: charmresource.Meta{
					Name:    "website2",
					Comment: "your website data",
				},
				Origin: charmresource.OriginUpload,
			},
			Username:  "Bill User",
			Timestamp: time.Date(2012, 12, 12, 12, 12, 12, 0, time.UTC),
		},
	}

	formatted := make([]FormattedSvcResource, len(res))
	for i := range res {
		formatted[i] = FormatSvcResource(res[i])
	}

	data, err := FormatSvcTabular(formatted)
	c.Assert(err, jc.ErrorIsNil)

	// Notes: sorted by name, then by revision, newest first.
	c.Check(string(data), gc.Equals, `
RESOURCE SUPPLIED BY REVISION         COMMENT
openjdk  charmstore  7                the java runtime
website  upload      -                your website data
openjdk2 charmstore  8                another java runtime
website2 Bill User   2012-12-12T12:12 your website data
`[1:])
}

func (s *SvcTabularSuite) TestFormatCharmTabularBadValue(c *gc.C) {
	bogus := "should have been something else"
	_, err := FormatSvcTabular(bogus)

	c.Check(err, gc.ErrorMatches, `expected value of type .*`)
}
