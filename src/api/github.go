package api

import (
	"context"
	"fmt"
	Err "github.com/42milez/NexusModsWatcher/src/error"
	"github.com/42milez/NexusModsWatcher/src/util"
	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
	"os"
	"time"
)

const accessTokenName = "GITHUB_TOKEN"
const (
	repo        = "NexusModsWatcher"
	repoOwner   = "42milez"
	baseBranch  = "main"
	authorName  = "Akihiro TAKASE"
	authorEmail = "42milez@gmail.com"
)

var GitHub *gitHubClient

type gitHubClient struct {
	accessToken string
	client      *github.Client
	ctx         context.Context
}

func (p *gitHubClient) CreatePullRequest(sub string, desc string, branch string, files []string, msg string) (string, error) {
	if sub == "" {
		return "", Err.NoSubjectProvided
	}
	if desc == "" {
		return "", Err.NoDescriptionProvided
	}
	if branch == "" {
		return "", Err.NoBranchProvided
	}
	if len(files) == 0 {
		return "", Err.NoFileProvided
	}
	if msg == "" {
		return "", Err.NoCommitMessageProvided
	}

	ref, cfErr := p.createRef(branch)
	if cfErr != nil {
		return "", Err.CreateRefFailed
	}

	tree, ctErr := p.createTree(ref, files)
	if ctErr != nil {
		return "", Err.CreateTreeFailed
	}

	if pcErr := p.pushCommit(ref, tree, msg); pcErr != nil {
		return "", Err.PushCommitFailed
	}

	newPR := &github.NewPullRequest{
		Title:               github.String(sub),
		Head:                github.String(branch),
		Base:                github.String(baseBranch),
		Body:                github.String(desc),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := p.client.PullRequests.Create(p.ctx, repoOwner, repo, newPR)
	if err != nil {
		return "", Err.CreatePullRequestFailed
	}

	return pr.GetHTMLURL(), nil
}

func (p *gitHubClient) createRef(branch string) (*github.Reference, error) {
	if branch == baseBranch {
		return nil, Err.InvalidBranchName
	}

	baseRef, _, grErr := p.client.Git.GetRef(p.ctx, repoOwner, repo, "refs/heads/"+baseBranch)
	if grErr != nil {
		return nil, Err.GetRefFailed
	}

	config := &github.Reference{
		Ref: github.String("refs/heads/" + branch),
		Object: &github.GitObject{
			SHA: baseRef.Object.SHA,
		},
	}
	newRef, _, crErr := p.client.Git.CreateRef(p.ctx, repoOwner, repo, config)
	if crErr != nil {
		return nil, Err.CreateRefFailed
	}

	return newRef, nil
}

func (p *gitHubClient) createTree(ref *github.Reference, files []string) (*github.Tree, error) {
	readFile := func(f string) ([]byte, error) {
		if len(files) == 0 {
			return nil, Err.NoFileProvided
		}
		return os.ReadFile(f)
	}

	var entries []*github.TreeEntry

	for _, f := range files {
		content, err := readFile(f)
		if err != nil {
			return nil, Err.ReadFileFailed
		}
		entries = append(entries, &github.TreeEntry{
			Path:    github.String(f),
			Type:    github.String("blob"),
			Content: github.String(string(content)),
			Mode:    github.String("100644"),
		})
	}

	tree, _, err := p.client.Git.CreateTree(p.ctx, repoOwner, repo, *ref.Object.SHA, entries)
	if err != nil {
		return nil, Err.CreateTreeFailed
	}

	return tree, nil
}

func (p *gitHubClient) pushCommit(ref *github.Reference, tree *github.Tree, msg string) error {
	parent, _, gcErr := p.client.Repositories.GetCommit(p.ctx, repoOwner, repo, *ref.Object.SHA)
	if gcErr != nil {
		return Err.PushCommitFailed
	}
	parent.Commit.SHA = parent.SHA

	date := time.Now()
	author := &github.CommitAuthor{
		Date:  &date,
		Name:  github.String(authorName),
		Email: github.String(authorEmail),
	}
	commit := &github.Commit{
		Author:  author,
		Message: github.String(msg),
		Tree:    tree,
		Parents: []*github.Commit{parent.Commit},
	}
	newCommit, _, ccErr := p.client.Git.CreateCommit(p.ctx, repoOwner, repo, commit)
	if ccErr != nil {
		return Err.CreateCommitFailed
	}

	ref.Object.SHA = newCommit.SHA
	_, _, err := p.client.Git.UpdateRef(p.ctx, repoOwner, repo, ref, false)
	if err != nil {
		return Err.UpdateRefFailed
	}

	return nil
}

func init() {
	GitHub = &gitHubClient{}

	var err error

	if GitHub.accessToken, err = getSecret("auth", accessTokenName); err != nil {
		util.Exit(fmt.Errorf("%s (%s)", Err.GetSecretFailed, accessTokenName))
	}
	if GitHub.accessToken == "" {
		util.Exit(Err.InvalidAccessToken)
	}

	// https://github.com/google/go-github#authentication
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GitHub.accessToken},
	)
	GitHub.client = github.NewClient(oauth2.NewClient(ctx, ts))

	GitHub.ctx = context.Background()
}
