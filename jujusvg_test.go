package jujusvg

import (
	"bytes"
	"strings"
	"testing"

	gc "gopkg.in/check.v1"
	"gopkg.in/juju/charm.v4"
)

func Test(t *testing.T) { gc.TestingT(t) }

type newSuite struct{}

var _ = gc.Suite(&newSuite{})

var bundle = `
services:
  mongodb:
    charm: "cs:precise/mongodb-21"
    num_units: 1
    annotations:
      "gui-x": "940.5"
      "gui-y": "388.7698359714502"
    constraints: "mem=2G cpu-cores=1"
  elasticsearch:
    charm: "cs:~charming-devs/precise/elasticsearch-2"
    num_units: 1
    annotations:
      "gui-x": "490.5"
      "gui-y": "369.7698359714502"
    constraints: "mem=2G cpu-cores=1"
  charmworld:
    charm: "cs:~juju-jitsu/precise/charmworld-58"
    num_units: 1
    expose: true
    annotations:
      "gui-x": "813.5"
      "gui-y": "112.23016402854975"
    options:
      charm_import_limit: -1
      source: "lp:~bac/charmworld/ingest-local-charms"
      revno: 511
relations:
  - - "charmworld:essearch"
    - "elasticsearch:essearch"
  - - "charmworld:database"
    - "mongodb:database"
series: precise
`

func iconURL(ref *charm.Reference) string {
	return "http://0.1.2.3/" + ref.Path() + ".svg"
}

func (s *newSuite) TestNewFromBundle(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	err = b.Verify(nil)
	c.Assert(err, gc.IsNil)

	cvs, err := NewFromBundle(b, iconURL)
	c.Assert(err, gc.IsNil)

	var buf bytes.Buffer
	cvs.Marshal(&buf)
	c.Logf("%s", buf.String())
	assertXMLEqual(c, buf.Bytes(), []byte(`
<?xml version="1.0"?>
<!-- Generated by SVGo -->
<svg width="546" height="372"
     xmlns="http://www.w3.org/2000/svg" 
     xmlns:xlink="http://www.w3.org/1999/xlink">
<defs>
<g id="healthCircle">
<circle cx="10" cy="10" r="10" style="stroke:#38B44A;fill:none;stroke-width:2px"/>
<circle cx="10" cy="10" r="5" style="fill:#38B44A"/>
</g>
</defs>
<g id="relations">
<line x1="371" y1="48" x2="48" y2="305" stroke="#38B44A" stroke-width="2px" stroke-dasharray="196.38, 20" />
<use x="199" y="166" xlink:href="#healthCircle" />
<line x1="371" y1="48" x2="498" y2="324" stroke="#38B44A" stroke-width="2px" stroke-dasharray="141.91, 20" />
<use x="424" y="176" xlink:href="#healthCircle" />
</g>
<g id="services">
<image x="323" y="0" width="96" height="96" xlink:href="http://0.1.2.3/~juju-jitsu/precise/charmworld-58.svg" />
<image x="0" y="257" width="96" height="96" xlink:href="http://0.1.2.3/~charming-devs/precise/elasticsearch-2.svg" />
<image x="450" y="276" width="96" height="96" xlink:href="http://0.1.2.3/precise/mongodb-21.svg" />
</g>
</svg>
`))
}

func (s *newSuite) TestWithBadBundle(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	b.Relations[0][0] = "evil-unknown-service"
	cvs, err := NewFromBundle(b, iconURL)
	c.Assert(err, gc.ErrorMatches, "cannot verify bundle: .*")
	c.Assert(cvs, gc.IsNil)
}

func (s *newSuite) TestWithBadPosition(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)

	b.Services["charmworld"].Annotations["gui-x"] = "bad"
	cvs, err := NewFromBundle(b, iconURL)
	c.Assert(err, gc.ErrorMatches, `service "charmworld" does not have a valid position`)
	c.Assert(cvs, gc.IsNil)

	b, err = charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)

	b.Services["charmworld"].Annotations["gui-y"] = "bad"
	cvs, err = NewFromBundle(b, iconURL)
	c.Assert(err, gc.ErrorMatches, `service "charmworld" does not have a valid position`)
	c.Assert(cvs, gc.IsNil)
}