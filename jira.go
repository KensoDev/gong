package gong

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/manifoldco/promptui"
)

// JiraClient : Struct implementing the generic Client interface
type JiraClient struct {
	client *jira.Client
	config map[string]string
}

// Create will create a new issue interactively and return the branch name
func (j *JiraClient) Create(projectKey string) (string, error) {
	if projectKey == "" {
		return "", errors.New("project key is required")
	}

	// Fetch create metadata for the project
	metaProject, err := j.getProjectMeta(projectKey)
	if err != nil {
		return "", fmt.Errorf("failed to get project metadata: %w", err)
	}

	// Prompt for issue type
	issueType, err := j.promptIssueType(metaProject)
	if err != nil {
		return "", err
	}

	// Gather required fields
	fields, err := j.promptRequiredFields(issueType, projectKey)
	if err != nil {
		return "", err
	}

	// Create the issue
	issue := &jira.Issue{
		Fields: fields,
	}

	fmt.Println("\nCreating issue...")
	createdIssue, resp, err := j.client.Issue.Create(issue)
	if err != nil {
		// Print response body for debugging
		if resp != nil && resp.Body != nil {
			bodyBytes := make([]byte, 2048)
			n, _ := resp.Body.Read(bodyBytes)
			if n > 0 {
				fmt.Printf("\nJIRA Error Response:\n%s\n", string(bodyBytes[:n]))
			}
		}
		return "", fmt.Errorf("failed to create issue: %w", err)
	}

	fmt.Printf("✓ Created issue: %s\n", createdIssue.Key)

	// Determine branch type from issue type
	branchType := j.determineBranchType(issueType.Name)

	// Get branch name
	branchName, err := j.GetBranchName(branchType, createdIssue.Key)
	if err != nil {
		return "", err
	}

	// Create and checkout the branch
	err = j.createGitBranch(branchName)
	if err != nil {
		return "", err
	}

	fmt.Printf("✓ Created and checked out branch: %s\n", branchName)

	return branchName, nil
}

// getProjectMeta fetches the create metadata for a specific project
func (j *JiraClient) getProjectMeta(projectKey string) (*jira.MetaProject, error) {
	options := &jira.GetQueryOptions{
		Expand: "projects.issuetypes.fields",
	}

	meta, _, err := j.client.Issue.GetCreateMetaWithOptions(options)
	if err != nil {
		return nil, err
	}

	// Get the specific project
	project := meta.GetProjectWithKey(projectKey)
	if project == nil {
		return nil, fmt.Errorf("project %s not found or you don't have permission to create issues in it", projectKey)
	}

	return project, nil
}

// promptIssueType prompts the user to select an issue type
func (j *JiraClient) promptIssueType(project *jira.MetaProject) (*jira.MetaIssueType, error) {
	if len(project.IssueTypes) == 0 {
		return nil, errors.New("no issue types available for this project")
	}

	// Build issue type list
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Name | cyan }}",
		Inactive: "  {{ .Name }}",
		Selected: "✓ {{ .Name | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select issue type",
		Items:     project.IssueTypes,
		Templates: templates,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return project.IssueTypes[idx], nil
}

// promptRequiredFields prompts the user for all required fields
func (j *JiraClient) promptRequiredFields(issueType *jira.MetaIssueType, projectKey string) (*jira.IssueFields, error) {
	fields := &jira.IssueFields{
		Type: jira.IssueType{
			ID:   issueType.Id,
			Name: issueType.Name,
		},
		Project: jira.Project{
			Key: projectKey,
		},
	}

	// Get mandatory field keys to know which are required
	mandatoryFields, err := issueType.GetMandatoryFields()
	if err != nil {
		return nil, fmt.Errorf("failed to get mandatory fields: %w", err)
	}

	// Create a map of fieldID -> isRequired for easy lookup
	requiredMap := make(map[string]bool)
	for _, fieldID := range mandatoryFields {
		requiredMap[fieldID] = true
	}

	// Always prompt for these common fields in order
	commonFields := []struct {
		id    string
		label string
	}{
		{"summary", "Summary"},
		{"description", "Description"},
	}

	for _, field := range commonFields {
		// Skip if already set
		if field.id == "project" || field.id == "issuetype" {
			continue
		}

		isRequired := requiredMap[field.id]
		label := field.label
		if !isRequired {
			label = label + " (optional, press Enter to skip)"
		}

		// Prompt based on field ID
		switch field.id {
		case "summary":
			summary, err := j.promptString(label, "", isRequired)
			if err != nil {
				return nil, err
			}
			fields.Summary = summary

		case "description":
			description, err := j.promptMultiline(label, isRequired)
			if err != nil {
				return nil, err
			}
			if description != "" {
				fields.Description = description
			}
		}
	}

	// Prompt for any other mandatory fields not in common list
	for fieldName, fieldID := range mandatoryFields {
		// Skip if already handled
		if fieldID == "project" || fieldID == "issuetype" || fieldID == "summary" || fieldID == "description" {
			continue
		}

		// For other required fields, provide basic string input
		value, err := j.promptString(fieldName, "", true)
		if err != nil {
			return nil, err
		}
		// Store in Unknowns map for custom fields
		if fields.Unknowns == nil {
			fields.Unknowns = make(map[string]interface{})
		}
		fields.Unknowns[fieldID] = value
	}

	return fields, nil
}

// promptString prompts the user for a string input
func (j *JiraClient) promptString(label string, defaultValue string, required bool) (string, error) {
	validate := func(input string) error {
		if required && strings.TrimSpace(input) == "" {
			return errors.New("this field is required")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Default:  defaultValue,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// promptMultiline prompts the user for multi-line input using their default editor
func (j *JiraClient) promptMultiline(label string, required bool) (string, error) {
	fmt.Printf("%s (press Enter to open editor, or Ctrl+C to skip)\n", label)

	// Wait for user to press Enter
	var dummy string
	fmt.Scanln(&dummy)

	// Get editor from environment
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default to vim
	}

	// Create temporary file
	tmpfile, err := os.CreateTemp("", "gong-description-*.txt")
	if err != nil {
		return "", err
	}
	tmpfilePath := tmpfile.Name()
	defer os.Remove(tmpfilePath)

	// Write instructions to the file
	instructions := "# Enter your description below. Lines starting with # will be ignored.\n# Save and close the editor to continue.\n\n"
	tmpfile.WriteString(instructions)
	tmpfile.Close()

	// Open editor
	cmd := exec.Command(editor, tmpfilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to open editor: %w", err)
	}

	// Read the file
	content, err := os.ReadFile(tmpfilePath)
	if err != nil {
		return "", err
	}

	// Remove comment lines and trim
	lines := strings.Split(string(content), "\n")
	var result []string
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") {
			result = append(result, line)
		}
	}

	description := strings.TrimSpace(strings.Join(result, "\n"))

	if required && description == "" {
		return "", errors.New("description is required")
	}

	return description, nil
}

// determineBranchType maps JIRA issue types to branch types
func (j *JiraClient) determineBranchType(issueTypeName string) string {
	issueTypeLower := strings.ToLower(issueTypeName)

	// Check environment variable for default
	if defaultType := os.Getenv("GONG_DEFAULT_BRANCH_TYPE"); defaultType != "" {
		return defaultType
	}

	// Map common JIRA issue types to branch types
	switch issueTypeLower {
	case "bug", "defect":
		return "bugfix"
	case "story", "task", "feature":
		return "feature"
	case "epic":
		return "epic"
	case "improvement", "enhancement":
		return "enhancement"
	case "sub-task", "subtask":
		return "task"
	default:
		return "feature"
	}
}

// createGitBranch creates and checks out a new git branch
func (j *JiraClient) createGitBranch(branchName string) error {
	// Check if git is available
	_, err := exec.LookPath("git")
	if err != nil {
		return errors.New("git is not installed or not in PATH")
	}

	// Create and checkout the branch
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create git branch: %w", err)
	}

	return nil
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
	issueID := GetIssueID(branchName)

	jiraComment := &jira.Comment{
		Body: comment,
	}
	_, _, err := j.client.Issue.AddComment(issueID, jiraComment)

	return err
}

// GetDomain : Get the domain from the config
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
	allowed := strings.Split(j.config["transitions"], ",")

	transitions, response, err := j.client.Issue.GetTransitions(issueID)

	for _, transition := range transitions {
		if indexOf(transition.Name, allowed) > -1 {
			_, err := j.client.Issue.DoTransition(issueID, transition.ID)

			if err != nil {
				fmt.Println(err)
				fmt.Println(response.Body)
				return "", err
			}
		}
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
		"username":       false,
		"domain":         false,
		"password":       true,
		"project_prefix": false,
		"transitions":    false,
	}
}

// Authenticate : Authenticates using the fields passed in
func (j *JiraClient) Authenticate(fields map[string]string) bool {
	tp := jira.BasicAuthTransport{
		Username: fields["username"],
		Password: fields["password"],
	}

	jiraClient, err := jira.NewClient(tp.Client(), fields["domain"])

	if err != nil {
		return false
	}

	j.client = jiraClient
	j.config = fields

	return true
}
