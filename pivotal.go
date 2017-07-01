package gong

import (
	"fmt"
	"gopkg.in/salsita/go-pivotaltracker.v1/v5/pivotal"
	"strconv"
	"regexp"
)

// PivotalClient : Struct implementing the generic Client interface
type PivotalClient struct {
	client *pivotal.Client
	config map[string]string
}

// NewPivotalClient : Returns a pointer to PivotalClient
func NewPivotalClient() *PivotalClient {
	return &PivotalClient{}
}

// GetName : Return the string name of the struct eg: pivotal
func (p *PivotalClient) GetName() string {
	return "pivotal"
}

// Browse : Browse to the URL of the issue related to the branch name
func (p *PivotalClient) Browse(branchName string) (string, error) {
	issueID := GetPivotalIssueID(branchName)
	domain := "https://www.pivotaltracker.com/story/show"
	url := fmt.Sprintf("%s/%s", domain, issueID)
	return url, nil
}

func GetPivotalIssueID(branchName string) string {
	re := regexp.MustCompile(`([\d]+)`)
	matches := re.FindAllString(branchName, -1)

	if len(matches) == 0 {
		return ""
	}

	return matches[0]
}

// Start : Start an issue
func (p *PivotalClient) Start(issueType string, issueID string) (string, error) {
	fmt.Println("calling pivotal.Start")

	projectIdInt, issueIdInt, err := p.GetProjectIdAndIssueId(issueID)

	if err != nil {
		fmt.Println(err)
	}

	request := &pivotal.StoryRequest{}
	request.State = "started"

	_, _, err = p.client.Stories.Update(projectIdInt, issueIdInt, request)

	if err != nil {
		fmt.Println(err)
	}

	branchName, err := p.GetBranchName(issueType, issueID)

	if err != nil {
		return "", err
	}

	fmt.Println("branch name:", branchName)
	return branchName, nil
}

func (p *PivotalClient) Comment(branchName, comment string) error {
	fmt.Println(branchName)
	issueID := GetPivotalIssueID(branchName)
	fmt.Println(issueID)

	projectIdInt, issueIdInt, err := p.GetProjectIdAndIssueId(issueID)

	if err != nil {
		fmt.Println(err)
	}

	commentObject := &pivotal.Comment{}
	commentObject.Text = comment

	_, _, err = p.client.Stories.AddComment(projectIdInt, issueIdInt, commentObject)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("added comment", comment)
	return nil
}

func (p *PivotalClient) PrepareCommitMessage(branchName, commitMessage string) string {
	issueID := GetPivotalIssueID(branchName)
	url, err := p.Browse(branchName)

	if err != nil {
		return commitMessage
	}

	patchedCommitMessage := fmt.Sprintf(`[%s](%s)`, issueID, url)

	return patchedCommitMessage
}

func (p *PivotalClient) GetBranchName(issueType string, issueID string) (string, error) {
	projectIdInt, issueIdInt, _ := p.GetProjectIdAndIssueId(issueID)

	story, _, err := p.client.Stories.Get(projectIdInt, issueIdInt)

	if err != nil {
		fmt.Println(err)
	}

	issueTitleSlug := SlugifyTitle(story.Name)

	return fmt.Sprintf("%s/%s-%s", issueType, issueID, issueTitleSlug), nil
}

func (p *PivotalClient) GetProjectIdAndIssueId(issueID string) (int, int, error) {
	fields, err := Load()
	if err != nil {
		return 0, 0, err
	}

	projectIdInt, err := strconv.Atoi(fields["projectId"])
	if err != nil {
		return 0, 0, err
	}

	fmt.Println("projectIdInt", projectIdInt)

	issueIdInt, err := strconv.Atoi(issueID)
	if err != nil {
		return 0, 0, err
	}

	fmt.Println("issueIdInt", issueIdInt)
	return projectIdInt, issueIdInt, nil
}

// FormatField : Returns a formatted field based on internal rules
func (p *PivotalClient) FormatField(fieldName string, value string) string {
	return value
}

// Authenticate : Authenticates using the fields passed in
func (p *PivotalClient) Authenticate(fields map[string]string) bool {
	pivotalClient := pivotal.NewClient(fields["apiToken"])
	_, _, err := pivotalClient.Me.Get()

	if err != nil {
		return false
	}

	p.client = pivotalClient
	p.config = fields

	return true
}

func (p *PivotalClient) GetAuthFields() map[string]bool {
	return map[string]bool{
		"projectId": false,
		"apiToken":  true,
	}
}
