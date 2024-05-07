# Contributing Guidelines

The following is a set of guidelines for contributing to Inspector.

## Contributing

Please use GitHub issues for discussions, feature requests and bugs. When opening an issue please make sure to choose correct label.

### Reportin a bug

Please open an issue on GitHub and ensure the issue has not already been reported.

### Suggesting a new feature or improvement

Please create an issue on Github and choose the type 'Feature request'.

### Opening a pull request

- Fork the repo, create a branch, submit a PR when your changes are tested and ready for review

## Following style guides

### Git

- Keep a clean, concise and meaningful git commit history on your branch, rebasing locally and squashing before
  submitting a PR
- Follow the guidelines of writing a good commit message as described [here](https://chris.beams.io/posts/git-commit/)
  and summarized in the following points:

  - In the subject line, use the present tense ("Add feature" not "Added feature")
  - In the subject line, use the imperative mood ("Move cursor to..." not "Moves cursor to...")
  - Limit the subject line to 72 characters or less
  - Reference issues and pull requests liberally after the subject line
  - Add more detailed description in the body of the git message (`git commit -a` to give you more space and time in
    your text editor to write a good message instead of `git commit -am`)

### Go code higiene

- Run [gofumpt](https://github.com/mvdan/gofumpt) over your code to automatically resolve a lot of style issues.
- Run [staticcheck](https://staticcheck.dev) and `go vet` on your code too to catch linting issues.
- Use [gotestdox](https://github.com/bitfield/gotestdox) to help with writing meaningful test names that describe tested behavior.
