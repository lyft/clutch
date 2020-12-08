package github

import (
	"github.com/shurcooL/githubv4"
)

// For general docs on using the GraphQL API, see:
//  - https://developer.github.com/v4/
//  - https://github.com/shurcooL/githubv4#simple-query

type getFileQuery struct {
	Repository struct {
		// Get more information about desired ref and last modified ref for the file.
		Ref struct {
			Commit struct {
				ID      githubv4.ID
				OID     githubv4.GitObjectID
				History struct {
					Nodes []struct {
						CommittedDate githubv4.DateTime
						OID           githubv4.GitObjectID
					}
				} `graphql:"history(path:$path,first:1)"`
			} `graphql:"... on Commit"`
		} `graphql:"ref: object(expression:$ref)"`

		// Fetch requested blob.
		Object struct {
			Blob struct {
				ID          githubv4.ID
				IsBinary    githubv4.Boolean
				IsTruncated githubv4.Boolean
				OID         githubv4.GitObjectID
				Text        githubv4.String
			} `graphql:"... on Blob"`
		} `graphql:"object(expression:$refPath)"`
	} `graphql:"repository(owner:$owner,name:$name)"`
}

type getDefaultBranchQuery struct {
	Repository struct {
		DefaultBranchRef struct {
			Name string
		}
	} `graphql:"repository(owner:$owner,name:$name)"`
}
