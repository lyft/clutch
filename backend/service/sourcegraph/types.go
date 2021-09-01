package sourcegraph

// Compare commits api
type compareCommitsResponse struct {
	Data Repository `json:"data"`
}

type Repository struct {
	Repository Comparision `json:"repository"`
}

type Comparision struct {
	Comparision Commits `json:"comparison"`
}

type Commits struct {
	Commits CommitNodes `json:"commits"`
}

type CommitNodes struct {
	Nodes []Commit `json:"nodes"`
}

type Commit struct {
	Message string       `json:"message"`
	Oid     string       `json:"oid"`
	Author  CommitAuthor `json:"author"`
}

type CommitAuthor struct {
	Person Person `json:"person"`
}

type Person struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}
