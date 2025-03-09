package cli

import (
	"fmt"
	"io"
	"os"

	"gitea.darkeli.com/yezi/git-bump/git"
	vv "gitea.darkeli.com/yezi/git-bump/internal/version"
	"gitea.darkeli.com/yezi/git-bump/version"
	"github.com/Masterminds/semver/v3"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Major   bool `long:"major" description:"Bump up major version"`
	Minor   bool `long:"minor" description:"Bump up minor version"`
	Patch   bool `long:"patch" description:"Bump up patch version"`
	Version bool `long:"version" short:"v" description:"Show current version"`
}

type CLI struct {
	Options Options
	Stdout  io.Writer
	Stderr  io.Writer
	Git     *git.Repository
	Version *version.VersionManager
}

func NewCLI(opts Options) *CLI {
	return &CLI{
		Options: opts,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Version: version.NewVersionManager("v"),
	}
}

func (c *CLI) ParseArgs(args []string) ([]string, error) {
	parser := flags.NewParser(&c.Options, flags.Default)
	return parser.ParseArgs(args)
}

func (c *CLI) Run(args []string) error {
	if c.Options.Version {
		fmt.Fprintln(c.Stdout, vv.Version)
		return nil
	}

	path, err := git.FindGitRoot(".")
	if err != nil {
		return err
	}

	if len(args) > 0 {
		path = args[0]
	}

	repo, err := git.Open(path)
	if err != nil {
		return err
	}
	c.Git = repo

	current, hasTags, err := c.getCurrentVersion()
	if err != nil {
		return err
	}

	var tag string = c.Version.FormatVersion(current)
	if hasTags {
		tag, err = c.createNextVersion(current)
		if err != nil {
			return err
		}
	}

	if err := c.Git.CreateTag(tag); err != nil {
		return err
	}

	fmt.Println("Created tag:", tag)

	return nil
}

func (c *CLI) getCurrentVersion() (*semver.Version, bool, error) {
	tags, err := c.Git.GetTags()
	if err != nil {
		return nil, true, err
	}

	current, err := c.Version.FindCurrentVersion(tags)
	if err != nil {
		return nil, true, err
	}

	if current == nil {
		v, err := c.Version.PromptVersion()
		if err != nil {
			return nil, false, fmt.Errorf("failed to create new version: %w", err)
		}
		return v, false, nil
	}

	return current, true, nil
}

func (c *CLI) createNextVersion(current *semver.Version) (string, error) {
	var spec version.Spec
	switch {
	case c.Options.Major:
		spec = version.Major
	case c.Options.Minor:
		spec = version.Minor
	case c.Options.Patch:
		spec = version.Patch
	default:
		spec = version.Patch
	}

	next, err := c.Version.NextVersion(current, spec)
	if err != nil {
		return "", fmt.Errorf("failed to create next version: %w", err)
	}

	return c.Version.FormatVersion(next), nil
}
