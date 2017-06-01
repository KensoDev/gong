package gong

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestSlugger(t *testing.T) { TestingT(t) }

type SluggerSuite struct{}

var _ = Suite(&SluggerSuite{})

func (s *SluggerSuite) TestSlugCreationNormalText(c *C) {
	title := "This is a fake title"
	c.Assert(SlugifyTitle(title), Equals, "this-is-a-fake-title")
}

func (s *SluggerSuite) TestSlugCreationWithSpecialChars(c *C) {
	secondTitle := "This &&& *** is another title"
	c.Assert(SlugifyTitle(secondTitle), Equals, "this-is-another-title")
}

func (s *SluggerSuite) TestSlugCreationWithMultipleSpaces(c *C) {
	title := "This is a                    fake title"
	c.Assert(SlugifyTitle(title), Equals, "this-is-a-fake-title")
}
