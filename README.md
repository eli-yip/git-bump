# Git-Bump

> `git-bump` is a fork of [`git-bump`](https://github.com/babarot/git-bump). Compared to the original version, it removes unnecessary features and updates dependencies.

This is a minimalist tool designed for automatically incrementing version numbers.

## Usage

```bash
# git-bump will automatically increment the Patch version by default.
git-bump
# Use git-bump --patch/--minor/--major to specify which version to increment.
git-bump --patch
git-bump --minor
git-bump --major
# If no tags exist in the repository, the tool will prompt you to input an initial version number.
git-bump
```
