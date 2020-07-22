module {{ .RepoProvider }}/{{ .RepoOwner }}/{{ .RepoName }}/backend

go 1.13

require (
    github.com/lyft/clutch/backend v0.0.0
)

replace github.com/lyft/clutch/backend => ../../../../github.com/lyft/clutch-preview/backend
