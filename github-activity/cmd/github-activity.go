package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode"

	"github.com/urfave/cli/v2"
)

// All possible events the API can provide us
type EventType string

const (
	COMMIT_COMMENT_EVENT              EventType = "CommitCommentEvent"
	CREATE_EVENT                      EventType = "CreateEvent"
	DELETE_EVENT                      EventType = "DeleteEvent"
	FORK_EVENT                        EventType = "ForkEvent"
	GOLLUM_EVENT                      EventType = "GollumEvent"
	ISSUE_COMMENT_EVENT               EventType = "IssueCommentEvent"
	ISSUES_EVENT                      EventType = "IssuesEvent"
	MEMBER_EVENT                      EventType = "MemberEvent"
	PUBLIC_EVENT                      EventType = "PublicEvent"
	PULL_REQUEST_EVENT                EventType = "PullRequestEvent"
	PULL_REQUEST_REVIEW_EVENT         EventType = "PullRequestReviewEvent"
	PULL_REQUEST_REVIEW_COMMENT_EVENT EventType = "PullRequestReviewCommentEvent"
	PULL_REQUEST_REVIEW_THREAD_EVENT  EventType = "PullRequestReviewThreadEvent"
	PUSH_EVENT                        EventType = "PushEvent"
	RELEASE_EVENT                     EventType = "ReleaseEvent"
	SPONSORSHIP_EVENT                 EventType = "SponsorshipEvent"
	WATCH_EVENT                       EventType = "WatchEvent"
)

const (
	GITHUB_API = "https://api.github.com/users/%s/events"
	TRUNC_LEN  = 32
)

// Represents a Github Event
//
// This doesn't actually hold all of the data that is received, only the data
// relevant to us
type Event struct {
	Type    EventType `json:"type"`
	Repo    Repo      `json:"repo"`
	Payload Payload   `json:"payload"`
}

// Represents the repository in which the event happened
type Repo struct {
	Name string `json:"name"`
}

// Represents the actual important data in an event
//
// This varies a lot (*A LOT* a lot) from event to event, so make sure to check
// the EventType value in the Event to make sure you're accessing the correct
// payload data
type Payload struct {
	Action      Action      `json:"action"`
	Ref         string      `json:"ref"`
	Assignee    User        `json:"assignee"`
	Member      User        `json:"member"`
	Reason      string      `json:"reason"`
	Release     Release     `json:"release"`
	Size        json.Number `json:"size"`
	Forkee      Repo        `json:"forkee"`
	RefType     RefType     `json:"ref_type"`
	Comment     Comment     `json:"comment"`
	Review      Review      `json:"review"`
	PullRequest PullRequest `json:"pull_request"`
	Issue       Issue       `json:"issue"`
	Pages       []Page      `json:"pages"`
	Commits     []Commit    `json:"commits"`
}

type Action string

const (
	ACT_ADDED                  string = "addded"
	ACT_ASSIGNED               Action = "assigned"
	ACT_CLOSED                 Action = "closed"
	ACT_CREATED                Action = "created"
	ACT_DELETED                Action = "deleted"
	ACT_DEQUEUED               Action = "dequeued"
	ACT_EDITED                 Action = "edited"
	ACT_LABELED                Action = "labeled"
	ACT_OPENED                 Action = "opened"
	ACT_PUBLISHED              Action = "published"
	ACT_REOPENED               Action = "reopened"
	ACT_RESOLVED               Action = "resolved"
	ACT_REVIEW_REQUESTED       Action = "review_requested"
	ACT_REVIEW_REQUEST_REMOVED Action = "review_request_removed"
	ACT_SYNCHRONIZE            Action = "synchronize"
	ACT_UNASSIGNED             Action = "unassigned"
	ACT_UNLABELED              Action = "unlabeled"
	ACT_UNRESOLVED             Action = "unresolved"
)

func (a Action) String() string {
	switch a {
	case ACT_REVIEW_REQUESTED:
		return "requested review"
	case ACT_REVIEW_REQUEST_REMOVED:
		return "removed review request"
	case ACT_SYNCHRONIZE:
		return "synchronized"
	default:
		return string(a)
	}
}

type RefType string

const (
	RT_BRANCH RefType = "branch"
	RT_TAG    RefType = "tag"
	RT_REPO   RefType = "repository"
)

// A comment on a issue/PR
type Comment struct {
	Body string `json:"body"`
}

// A review of a commit
type Review struct {
	State string `json:"state"`
}

// A PR
type PullRequest struct {
	Number json.Number `json:"number"`
	Title  string      `json:"title"`
	State  string      `json:"state"`
	User   User        `json:"login"`
}

func (pr PullRequest) String() string {
	return fmt.Sprintf("PR #%s: '%s'", pr.Number, truncate(pr.Title, TRUNC_LEN))
}

// An issue
type Issue struct {
	Number json.Number `json:"number"`
	Title  string      `json:"title"`
	State  string      `json:"state"`
	User   User        `json:"login"`
}

func (i Issue) String() string {
	return fmt.Sprintf("Issue #%s: '%s'", i.Number, truncate(i.Title, TRUNC_LEN))
}

// A commit
type Commit struct {
	Message  string `json:"message"`
	SHA      string `json:"string"`
	Author   Author `json:"author"`
	Distinct bool   `json:"distinct"`
}

func (c Commit) String() string {
	return fmt.Sprintf("Commit %s (%s): '%s'", truncate(c.SHA, 8), c.Author, truncate(c.Message, TRUNC_LEN))
}

// The author of a commit
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (a Author) String() string {
	return fmt.Sprintf("%s <%s>", a.Name, a.Email)
}

// Represents a Github user
type User struct {
	Login string `json:"login"`
}

// Represents a wiki Page
type Page struct {
	PageName string `json:"page_name"`
	Title    string `json:"title"`
	Action   Action `json:"action"`
	HTMLUrl  string `json:"html_url"`
}

type Release struct {
	Title string `json:"title"`
}

func (p Page) String() string {
	if p.Action == ACT_CREATED {
		return fmt.Sprintf("Created page '%s' @ '%s'", p.PageName, p.HTMLUrl)
	}

	return fmt.Sprintf("Edited page '%s' @ '%s'", p.PageName, p.HTMLUrl)
}

func fetchFromAPI(user string) ([]Event, error) {
	// Get formatted API URL
	url := fmt.Sprintf(GITHUB_API, user)

	// Make request
	res, err := http.Get(url)
	if err != nil {
		return []Event{}, err
	}

	// Check if user was not found to give the user a friendlier message
	if res.StatusCode == http.StatusNotFound {
		return []Event{}, fmt.Errorf("No such user '%s'!", user)
	}

	// Check if the response went OK otherwise
	if res.StatusCode != http.StatusOK {
		return []Event{}, fmt.Errorf("Received bad status code: %v\n", res.Status)
	}

	// Parse response into struct
	var data []Event
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return []Event{}, err
	}

	return data, nil
}

func truncate(text string, maxLength int) string {
	lastSpaceIx := maxLength
	length := 0
	for i, r := range text {
		if r == '\n' {
			return text[:i] + "..."
		}

		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}

		length++
		if length > maxLength {
			return text[:lastSpaceIx] + "..."
		}
	}
	// If here, string is shorter or equal to maxLen
	return text
}

type EventPrinter func(Event)

var eventPrinters map[EventType]EventPrinter = map[EventType]EventPrinter{
	COMMIT_COMMENT_EVENT: func(e Event) {
		fmt.Printf(
			"commented on a commit on repo %s:\n  '%s'",
			e.Repo.Name, truncate(e.Payload.Comment.Body, TRUNC_LEN),
		)
	},
	CREATE_EVENT: func(e Event) {
		fmt.Printf("created a new %s ", e.Payload.RefType)

		if e.Payload.RefType != RT_REPO {
			fmt.Printf("'%s' on ", e.Payload.Ref)
		}

		fmt.Printf("'%s'", e.Repo.Name)
	},
	DELETE_EVENT: func(e Event) {
		fmt.Printf("deleted the '%s' %s", e.Payload.Ref, e.Payload.RefType)
	},
	FORK_EVENT: func(e Event) {
		fmt.Printf("forked the '%s' repo", e.Payload.Forkee.Name)
	},
	GOLLUM_EVENT: func(e Event) {
		if len(e.Payload.Pages) == 1 {
			fmt.Printf("interacted with the following page:\n")
		} else {
			fmt.Printf("interacted with the following pages:\n")
		}

		for idx, page := range e.Payload.Pages {
			fmt.Print("• ", page)

			if idx == len(e.Payload.Pages)-1 {
				fmt.Println(".")
			} else {
				fmt.Println(";")
			}

		}
	},
	ISSUE_COMMENT_EVENT: func(e Event) {
		fmt.Printf(
			"%s comment under %s on repo '%s':\n  '%s'",
			e.Payload.Action, e.Payload.Issue, e.Repo.Name,
			truncate(e.Payload.Comment.Body, TRUNC_LEN),
		)
	},
	ISSUES_EVENT: func(e Event) {
		fmt.Print(e.Payload.Action, " ")

		switch e.Payload.Action {
		case ACT_LABELED, ACT_UNLABELED:
			fmt.Print("label from ")
		case ACT_ASSIGNED, ACT_UNASSIGNED:
			fmt.Printf("user %s from ", e.Payload.Assignee.Login)
		}

		fmt.Print(e.Payload.Issue)
	},
	MEMBER_EVENT: func(e Event) {
		if e.Payload.Action == ACT_ASSIGNED {
			fmt.Printf("Assigned user '%s' to repo '%s'", e.Payload.Member.Login, e.Repo.Name)
			return
		}

		fmt.Printf("Edited user '%s''s permissions on repo '%s'", e.Payload.Member.Login, e.Repo.Name)
	},
	PUBLIC_EVENT: func(e Event) {
		fmt.Printf("Made repo '%s' public", e.Repo.Name)
	},
	PULL_REQUEST_EVENT: func(e Event) {
		fmt.Printf("%s %s in repo '%s'", e.Payload.Action, e.Payload.PullRequest, e.Repo.Name)

		if e.Payload.Action == ACT_DEQUEUED {
			fmt.Printf("\n  reason: '%s'", truncate(e.Payload.Reason, TRUNC_LEN))
		}
	},
	PULL_REQUEST_REVIEW_EVENT: func(e Event) {
		fmt.Printf("%s review under %s", e.Payload.Action, e.Payload.PullRequest)
	},
	PULL_REQUEST_REVIEW_COMMENT_EVENT: func(e Event) {
		fmt.Printf(
			"%s review comment under %s:\n  '%s'",
			e.Payload.Action, e.Payload.PullRequest,
			truncate(e.Payload.Comment.Body, TRUNC_LEN),
		)
	},
	PULL_REQUEST_REVIEW_THREAD_EVENT: func(e Event) {
		fmt.Printf("marked review thread under %s as %s", e.Payload.PullRequest, e.Payload.Action)
	},
	PUSH_EVENT: func(e Event) {
		if e.Payload.Size == "1" {
			fmt.Printf("pushed 1 commit to repo '%s':\n", e.Repo.Name)
		} else {
			fmt.Printf("pushed %s commits to repo '%s':\n", e.Payload.Size, e.Repo.Name)
		}

		for idx, commit := range e.Payload.Commits {
			if !commit.Distinct {
				continue
			}

			if idx < len(e.Payload.Commits)-1 {
				fmt.Printf("├─ %s", commit)
				fmt.Println()
			} else {
				fmt.Printf("└─ %s", commit)
			}
		}
	},
	RELEASE_EVENT: func(e Event) {
		fmt.Printf("%s release '%s'", e.Payload.Action, truncate(e.Payload.Release.Title, TRUNC_LEN))
	},
	SPONSORSHIP_EVENT: func(e Event) {
		fmt.Printf("something related to the sponsors of repo '%s'", e.Repo.Name)
	},
	WATCH_EVENT: func(e Event) {
		fmt.Printf("starred repo '%s'", e.Repo.Name)
	},
}

func parseUserData(username string, events []Event) error {
	fmt.Printf("%s's GitHub activity:\n", username)

	for idx, event := range events {
		fmt.Print("• ")
		eventPrinters[event.Type](event)

		if idx == len(events)-1 {
			fmt.Println(".")
		} else {
			fmt.Print(";\n\n")
		}
	}

	return nil
}

func Handle(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf(
			"github-activity accepts only one argument (a Github user), but received %v arguments!\n",
			ctx.NArg(),
		)
	}

	username := ctx.Args().First()

	res, err := fetchFromAPI(username)
	if err != nil {
		return err
	}

	if err := parseUserData(username, res); err != nil {
		return err
	}

	return nil
}
