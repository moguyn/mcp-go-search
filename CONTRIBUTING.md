# Contributing to Bocha AI Search MCP Server

Thank you for considering contributing to the Bocha AI Search MCP Server! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to abide by our code of conduct:

- Be respectful and inclusive
- Be patient and welcoming
- Be thoughtful
- Be collaborative
- When disagreeing, try to understand why

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue tracker to see if the problem has already been reported. When you are creating a bug report, please include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples to demonstrate the steps**
- **Describe the behavior you observed after following the steps**
- **Explain which behavior you expected to see instead and why**
- **Include screenshots or animated GIFs if possible**
- **Include details about your configuration and environment**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

- **Use a clear and descriptive title**
- **Provide a step-by-step description of the suggested enhancement**
- **Provide specific examples to demonstrate the steps**
- **Describe the current behavior and explain which behavior you expected to see instead**
- **Explain why this enhancement would be useful**

### Pull Requests

- Fill in the required template
- Do not include issue numbers in the PR title
- Include screenshots and animated GIFs in your pull request whenever possible
- Follow the Go style guide
- Include tests for new features or bug fixes
- Document new code
- End all files with a newline

## Development Process

### Setting Up the Development Environment

1. Fork the repository
2. Clone your fork locally:
   ```bash
   git clone https://github.com/yourusername/mcp-go-search.git
   cd mcp-go-search
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Create a branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

### Testing

Before submitting a pull request, make sure all tests pass:

```bash
make test
```

### Linting

Ensure your code follows our style guidelines:

```bash
make lint
```

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

## Style Guides

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Document all exported functions, types, and constants

### Documentation Style Guide

- Use Markdown for documentation
- Reference function and variable names using backticks: `functionName`
- Use code blocks for examples
- Keep documentation up to date with code changes

## Additional Notes

### Issue and Pull Request Labels

This project uses the following labels to track issues and pull requests:

- `bug`: Indicates a confirmed bug or problem
- `documentation`: Indicates a need for improvements or additions to documentation
- `enhancement`: Indicates new feature requests or improvements
- `good first issue`: Indicates issues which are good for newcomers
- `help wanted`: Indicates issues where help is particularly desired

## Releasing

The project maintainers are responsible for releasing new versions. The process typically involves:

1. Updating the version number
2. Creating release notes
3. Creating a new GitHub release
4. Building and publishing the release artifacts

## Questions?

If you have any questions or need help, please open an issue or contact the project maintainers.

Thank you for contributing to the Bocha AI Search MCP Server! 