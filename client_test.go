package gong

import (
	"fmt"
	. "gopkg.in/check.v1"
	"testing"
)

func TestClient(t *testing.T) { TestingT(t) }

type ClientSuite struct{}

var _ = Suite(&ClientSuite{})

type FakeClient struct{}

func (f *FakeClient) FormatField(fieldName string, value string) string {
	return ""
}

func (f *FakeClient) GetAuthFields() map[string]bool {
	return map[string]bool{}
}

func (f *FakeClient) GetName() string {
	return "fakeclient"
}

func (f *FakeClient) Authenticate(field map[string]string) bool {
	return false
}

func (f *FakeClient) Start(issueType, issueID string) (string, error) {
	return fmt.Sprintf("%s/%s", issueType, issueID), nil
}

func (f *FakeClient) Browse(branchName string) (string, error) {
	return "https://www.fake.com/FAKE-1111", nil
}

func (f *FakeClient) Comment(branchName, comment string) error {
	return nil
}

func (f *FakeClient) PrepareCommitMessage(branchName, commitMessage string) string {
	return "Fake commit message"
}

func (s *ClientSuite) TestClientStartIssue(c *C) {
	client := &FakeClient{}
	branchName, _ := Start(client, "feature", "111")
	c.Assert(branchName, Equals, "feature/111")
}

func (s *ClientSuite) TestBrowse(c *C) {
	client := &FakeClient{}
	branchName := "feature/FAKE-1111-some-issue-title"
	url, _ := Browse(client, branchName)
	c.Assert(url, Equals, "https://www.fake.com/FAKE-1111")
}

func (s *ClientSuite) TestComment(c *C) {
	client := &FakeClient{}
	branchName := "feature/FAKE-1111-some-issue-title"
	comment := "This is a sample comment"
	err := Comment(client, branchName, comment)
	c.Assert(err, Equals, nil)
}

func (s *ClientSuite) TestPrepareCommitMessage(c *C) {
	client := &FakeClient{}
	branchName := "feature/FAKE-1111-some-issue-title"
	commitMessage := "This is a sample comment"
	newMessage := PrepareCommitMessage(client, branchName, commitMessage)
	c.Assert(newMessage, Equals, "Fake commit message")
}
