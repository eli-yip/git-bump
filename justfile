current_branch := `git branch --show-current`

lint:
  golangci-lint run -v --timeout 5m

update:
  go get -u ./...
  go mod tidy

commit:
  git add -A
  git commit -v

push:
  git push -u origin {{current_branch}}
  git push --tags

fpush:
  git push -u origin {{current_branch}} --force-with-lease --tags

conclude:
  git diff --stat @{0.day.ago.midnight} | sort -k3nr

dtag +tags:
  #!/usr/bin/env bash
  for tag in {{tags}}; do
    git tag -d "${tag}" && \
    git push origin --delete "${tag}"
  done

@ltag:
  git tag --list --sort -v:refname -n

lol:
  git lol
