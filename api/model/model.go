package model

import "time"

//GetURL ...
type GetURL struct {
	URL string `json:"url" valid:"required"`
}

//TrackIssueInfo ...
type TrackIssueInfo struct {
	IssueURL  string `json:"url"`
	State     string `json:"state"`
	OwnerName string `json:"owner_name"`
	RepoName  string `json:"repo_name"`
	CreatedAt string `json:"created_at"`
}

//OpenIssues ... this model is for get rquire format of data to display in front end in single call
type OpenIssues struct {
	OneDay         int `json:"one_day"`
	Upto7Day       int `json:"seven_day"`
	OlderThen7Day  int `json:"older_then_7day"`
	TotalOpenIssue int `json:"open_issues"`
	Status         int `json:"status"`
}

//IssuesResp ...
type IssuesResp struct {
	URL           string `json:"url"`
	RepositoryURL string `json:"repository_url"`
	LabelsURL     string `json:"labels_url"`
	CommentsURL   string `json:"comments_url"`
	EventsURL     string `json:"events_url"`
	HTMLURL       string `json:"html_url"`
	ID            int    `json:"id"`
	NodeID        string `json:"node_id"`
	Number        int    `json:"number"`
	Title         string `json:"title"`
	User          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"user"`
	Labels []struct {
		ID      int    `json:"id"`
		NodeID  string `json:"node_id"`
		URL     string `json:"url"`
		Name    string `json:"name"`
		Color   string `json:"color"`
		Default bool   `json:"default"`
	} `json:"labels"`
	State             string        `json:"state"`
	Locked            bool          `json:"locked"`
	Assignee          interface{}   `json:"assignee"`
	Assignees         []interface{} `json:"assignees"`
	Milestone         interface{}   `json:"milestone"`
	Comments          int           `json:"comments"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	ClosedAt          interface{}   `json:"closed_at"`
	AuthorAssociation string        `json:"author_association"`
	Body              string        `json:"body"`
}
