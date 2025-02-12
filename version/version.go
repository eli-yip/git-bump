package version

import (
	"errors"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/manifoldco/promptui"
)

type Spec int

const (
	Major Spec = iota
	Minor
	Patch
)

func (s Spec) String() string {
	switch s {
	case Major:
		return "major"
	case Minor:
		return "minor"
	case Patch:
		return "patch"
	default:
		return "unknown"
	}
}

type VersionManager struct{ Prefix string }

func NewVersionManager(prefix string) *VersionManager { return &VersionManager{Prefix: prefix} }

func (vm *VersionManager) FindCurrentVersion(tags []string) (*semver.Version, error) {
	vs := make([]*semver.Version, 0)
	for _, tag := range tags {
		v, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}
		vs = append(vs, v)
	}

	sort.Sort(semver.Collection(vs))

	if len(vs) == 0 {
		return nil, nil
	}

	return vs[len(vs)-1], nil
}

func (vm *VersionManager) NextVersion(current *semver.Version, spec Spec) (*semver.Version, error) {
	if current == nil {
		return nil, errors.New("current version is nil")
	}

	var next *semver.Version
	switch spec {
	case Major:
		v := current.IncMajor()
		next = &v
	case Minor:
		v := current.IncMinor()
		next = &v
	case Patch:
		v := current.IncPatch()
		next = &v
	default:
		return nil, errors.New("invalid version spec")
	}

	return next, nil
}

func (vm *VersionManager) FormatVersion(v *semver.Version) string {
	version := v.String()
	if strings.HasPrefix(v.Original(), vm.Prefix) {
		version = vm.Prefix + version
	}
	return version
}

func (vm *VersionManager) PromptVersion() (*semver.Version, error) {
	validate := func(input string) error {
		_, err := semver.NewVersion(input)
		return err
	}

	prompt := promptui.Prompt{
		Label:    "New version",
		Validate: validate,
	}

	v, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return semver.NewVersion(v)
}
