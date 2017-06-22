package gong

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/andygrunwald/go-jira"
)

// JiraClient : Struct implementing the generic Client interface
type JiraClient struct {
	client *jira.Client
	config map[string]string
}

func (j *JiraClient) Create() (string, error) {
	domain, err := j.GetDomain()

	if err != nil {
		return "", err
	}

	return domain, nil
}

// NewJiraClient : Returns a pointer to JiraClient
func NewJiraClient() *JiraClient {
	return &JiraClient{}
}

// GetName : Return the string name of the struct eg: jira
func (j *JiraClient) GetName() string {
	return "jira"
}

// PrepareCommitMessage : Returns a string with the issue id in the link
func (j *JiraClient) PrepareCommitMessage(branchName, commitMessage string) string {
	issueID := GetIssueID(branchName)
	url, err := j.Browse(branchName)

	if err != nil {
		return commitMessage
	}

	patchedCommitMessage := fmt.Sprintf(`[%s](%s)`, issueID, url)

	return patchedCommitMessage
}

// GetIssueID : returns the issue id from a branch name
func GetIssueID(branchName string) string {
	re := regexp.MustCompile(`([A-Z]+-[\d]+)`)
	matches := re.FindAllString(branchName, -1)

	if len(matches) == 0 {
		return ""
	}

	return matches[0]
}

// Comment : Post a comment on a jira issue
func (j *JiraClient) Comment(branchName, comment string) error {
	fmt.Println(branchName)
	issueID := GetIssueID(branchName)
	fmt.Println(issueID)

	jiraComment := &jira.Comment{
		Body: comment,
	}
	_, _, err := j.client.Issue.AddComment(issueID, jiraComment)

	return err
}

func (j *JiraClient) GetDomain() (string, error) {

	domain, ok := j.config["domain"]

	if !ok {
		return "", errors.New("Could not locate domain in config")
	}

	return domain, nil
}

// Browse : Browse to the URL of the issue related to the branch name
func (j *JiraClient) Browse(branchName string) (string, error) {
	issueID := GetIssueID(branchName)

	domain, err := j.GetDomain()

	if err != nil {
		return "", err
	}

	if issueID == "" {
		return "", errors.New("Could not find issue id in the branch name")
	}

	url := fmt.Sprintf("%s/browse/%s", domain, issueID)

	return url, nil
}

// GetBranchName : Return the branch name from the issue id and issue type
func (j *JiraClient) GetBranchName(issueType string, issueID string) (string, error) {
	issue, _, err := j.client.Issue.Get(issueID, nil)

	if err != nil {
		return "", err
	}

	issueTitleSlug := SlugifyTitle(issue.Fields.Summary)
	return fmt.Sprintf("%s/%s-%s", issueType, issueID, issueTitleSlug), nil
}

func indexOf(status string, data []string) int {
	for k, v := range data {
		if status == v {
			return k
		}
	}
	return -1
}

// Start : Start an issue
func (j *JiraClient) Start(issueType string, issueID string) (string, error) {
	allowed := []string{"Ready", "Start"}

	fmt.Println(issueID)

	transitions, response, err := j.client.Issue.GetTransitions(issueID)

	if err != nil {
		fmt.Println(err)
		fmt.Println(response.Body)
		return "", err
	}

	nextTransition := transitions[0]

	if indexOf(nextTransition.Name, allowed) > -1 {
		_, err := j.client.Issue.DoTransition(issueID, nextTransition.ID)

		if err != nil {
			return "", err
		}

		_, _ = j.Start(issueType, issueID)
	}

	branchName, err := j.GetBranchName(issueType, issueID)

	if err != nil {
		return "", err
	}

	return branchName, nil
}

// FormatField : Returns a formatted field based on internal rules
func (j *JiraClient) FormatField(fieldName string, value string) string {
	if fieldName == "domain" {
		return fmt.Sprintf("https://%s", value)
	}

	return value
}

// GetAuthFields : Get a map of auth fields
func (j *JiraClient) GetAuthFields() map[string]bool {
	return map[string]bool{
		"username":       true,
		"domain":         true,
		"password":       true,
		"project_prefix": false,
	}
}

// Authenticate : Authenticates using the fields passed in
func (j *JiraClient) Authenticate(fields map[string]string) bool {
	jiraClient, err := jira.NewClient(nil, fields["domain"])

	if err != nil {
		return false
	}

	res, err := jiraClient.Authentication.AcquireSessionCookie(fields["username"], fields["password"])

	if err != nil || res == false {
		return false
	}

	j.client = jiraClient
	j.config = fields

	return true
}
