package jujuc_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/juju-core/cmd"
	"launchpad.net/juju-core/testing"
	"launchpad.net/juju-core/utils/set"
	"launchpad.net/juju-core/worker/uniter/jujuc"
)

type PortsSuite struct {
	ContextSuite
}

var _ = Suite(&PortsSuite{})

var portsTests = []struct {
	cmd    []string
	expect set.Strings
}{
	{[]string{"open-port", "80"}, set.NewStrings("80/tcp")},
	{[]string{"open-port", "99/tcp"}, set.NewStrings("80/tcp", "99/tcp")},
	{[]string{"close-port", "80/TCP"}, set.NewStrings("99/tcp")},
	{[]string{"open-port", "123/udp"}, set.NewStrings("99/tcp", "123/udp")},
	{[]string{"close-port", "9999/UDP"}, set.NewStrings("99/tcp", "123/udp")},
}

func (s *PortsSuite) TestOpenClose(c *C) {
	hctx := s.GetHookContext(c, -1, "")
	for _, t := range portsTests {
		com, err := jujuc.NewCommand(hctx, t.cmd[0])
		c.Assert(err, IsNil)
		ctx := testing.Context(c)
		code := cmd.Main(com, ctx, t.cmd[1:])
		c.Assert(code, Equals, 0)
		c.Assert(bufferString(ctx.Stdout), Equals, "")
		c.Assert(bufferString(ctx.Stderr), Equals, "")
		c.Assert(hctx.ports, DeepEquals, t.expect)
	}
}

var badPortsTests = []struct {
	args []string
	err  string
}{
	{nil, "no port specified"},
	{[]string{"0"}, `port must be in the range \[1, 65535\]; got "0"`},
	{[]string{"65536"}, `port must be in the range \[1, 65535\]; got "65536"`},
	{[]string{"two"}, `port must be in the range \[1, 65535\]; got "two"`},
	{[]string{"80/http"}, `protocol must be "tcp" or "udp"; got "http"`},
	{[]string{"blah/blah/blah"}, `expected <port>\[/<protocol>\]; got "blah/blah/blah"`},
	{[]string{"123", "haha"}, `unrecognized args: \["haha"\]`},
}

func (s *PortsSuite) TestBadArgs(c *C) {
	for _, name := range []string{"open-port", "close-port"} {
		for _, t := range badPortsTests {
			hctx := s.GetHookContext(c, -1, "")
			com, err := jujuc.NewCommand(hctx, name)
			c.Assert(err, IsNil)
			err = testing.InitCommand(com, t.args)
			c.Assert(err, ErrorMatches, t.err)
		}
	}
}

func (s *PortsSuite) TestHelp(c *C) {
	hctx := s.GetHookContext(c, -1, "")
	open, err := jujuc.NewCommand(hctx, "open-port")
	c.Assert(err, IsNil)
	flags := testing.NewFlagSet()
	c.Assert(string(open.Info().Help(flags)), Equals, `
usage: open-port <port>[/<protocol>]
purpose: register a port to open

The port will only be open while the service is exposed.
`[1:])

	close, err := jujuc.NewCommand(hctx, "close-port")
	c.Assert(err, IsNil)
	c.Assert(string(close.Info().Help(flags)), Equals, `
usage: close-port <port>[/<protocol>]
purpose: ensure a port is always closed
`[1:])
}

// Since the deprecation warning gets output during Run, we really need
// some valid commands to run
var portsFormatDeprectaionTests = []struct {
	cmd []string
}{
	{[]string{"open-port", "--format", "foo", "80"}},
	{[]string{"close-port", "--format", "foo", "80/TCP"}},
}

func (s *PortsSuite) TestOpenCloseDeprecation(c *C) {
	hctx := s.GetHookContext(c, -1, "")
	for _, t := range portsFormatDeprectaionTests {
		name := t.cmd[0]
		com, err := jujuc.NewCommand(hctx, name)
		c.Assert(err, IsNil)
		ctx := testing.Context(c)
		code := cmd.Main(com, ctx, t.cmd[1:])
		c.Assert(code, Equals, 0)
		c.Assert(testing.Stdout(ctx), Equals, "")
		c.Assert(testing.Stderr(ctx), Equals, "--format flag deprecated for command \""+name+"\"")
	}
}
