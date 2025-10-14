# GitKit

> **Streamline your Git workflow with automated branching strategies and platform integrations**

GitKit is a modern, powerful CLI tool that brings git-flow branching strategies to the 21st century. Designed for teams working with protected branches in enterprise environments, GitKit automates complex Git workflows while maintaining clean, predictable branch management.

## Why GitKit?

Traditional Git workflows can be cumbersome and error-prone, especially in enterprise settings with protected branches, code reviews, and multiple environments. GitKit solves these challenges by:

- **Automating Branch Management**: Never worry about branch naming conventions or merge strategies again
- **Protected Branch Support**: Seamlessly creates merge requests when direct pushes aren't allowed
- **Multi-Platform Integration**: Native support for GitHub and GitLab with automatic PR/MR creation
- **Configuration-Driven**: Customize workflows per project with simple YAML configuration
- **Developer Experience**: Intuitive CLI with helpful prompts and clear error messages

## Key Benefits

- 🚀 **Faster Development Cycles**: Reduce time spent on Git operations by up to 60%
- 🔒 **Enterprise-Ready**: Built for protected branches and compliance requirements
- 🤖 **Automation First**: Let GitKit handle the complex Git operations
- 🌐 **Platform Agnostic**: Works with GitHub, GitLab, and self-hosted solutions
- ⚙️ **Highly Configurable**: Adapt to your team's specific workflow needs
- 📚 **Well Documented**: Extensive examples and clear error messages

## Features

- 🚀 **Git-Flow Workflows**: Feature, bugfix, hotfix, and release branch management
- ⚙️ **Configuration-Driven**: YAML-based configuration for branch prefixes and remote settings
- 🔧 **Modern CLI**: Built with Cobra for intuitive command-line interface
- 📦 **Modular Architecture**: Clean separation of concerns with dependency injection
- 🎯 **Interactive Setup**: Guided configuration with prompts
- 🔄 **Sync Operations**: Automated pull/push with remote branch management
- 🔒 **Protected Branch Support**: Auto-generates merge requests for GitHub/GitLab when branches are protected
- 🌐 **Multi-Platform**: Native support for GitHub and GitLab integrations

## Installation

### From Source

```bash
git clone https://github.com/samanar/gitkit.git
cd gitkit
go build -o gitkit main.go
# Move to your PATH
sudo mv gitkit /usr/local/bin/
```

## Quick Start

1. **Initialize GitKit in your repository:**

```bash
gitkit init
```

2. **Start a new feature:**

```bash
gitkit feature start my-awesome-feature
```

3. **Work on your feature, then finish it:**

```bash
gitkit feature finish
```

## Configuration

GitKit uses a `.gitkit.yaml` configuration file in your repository root:

```yaml
repo: "my-project"
branches:
  main: "main"
  develop: "develop"
prefixes:
  feature:
    name: "feature/"
    base: "develop"
  bugfix:
    name: "bugfix/"
    base: "develop"
  hotfix:
    name: "hotfix/"
    base: "main"
  release:
    name: "release/"
    base: "develop"
remote: "origin"
```

## Commands

### Feature Management

```bash
# Start a new feature branch
gitkit feature start <name>

# Finish current feature (merge back to develop, creates MR if branch protected)
gitkit feature finish [branch-name]
```

### Bugfix Management

```bash
# Start a bugfix branch
gitkit bugfix start <name>

# Finish bugfix (merge back to develop, creates MR if branch protected)
gitkit bugfix finish [branch-name]
```

### Hotfix Management

```bash
# Start a hotfix branch (from main)
gitkit hotfix start <name>

# Finish hotfix (merge to both main and develop, creates MRs if branches protected)
gitkit hotfix finish [branch-name]
```

### Release Management

```bash
# Start a release branch
gitkit release start <version>

# Finish release (merge to main and develop, create tag, creates MRs if branches protected)
gitkit release finish [branch-name]
```

### Basic Git Operations

```bash
# Enhanced status with pretty formatting
gitkit status

# Stage files (supports multiple files)
gitkit add [file1 file2 ...]

# Commit with message (joins all arguments)
gitkit commit Add new feature implementation

# Pull with rebase
gitkit pull

# Push (config-aware upstream setting)
gitkit push

# Sync (pull + push)
gitkit sync
```

### Catch-All Commands

GitKit supports any git command directly:

```bash
gitkit log --oneline
gitkit diff HEAD~1
gitkit stash list
```

## Platform Integration & Merge Requests

GitKit automatically generates merge requests (GitHub) or merge requests (GitLab) when working with protected branches. This is especially useful in enterprise environments where direct pushes to main/develop branches are restricted.

### Setting Up Platform Integration

When you finish a feature, bugfix, hotfix, or release branch, GitKit will automatically create a merge request if:

1. The target branch is protected (cannot be pushed to directly)
2. Platform integration is configured with valid credentials
3. The repository is hosted on GitHub or GitLab

### GitHub Setup

1. **Generate Personal Access Token:**

   - Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
   - Click "Generate new token (classic)"
   - Select scopes: `repo` (full control of private repositories)
   - Copy the generated token

2. **Configure GitKit for GitHub:**

   ```bash
   # Run init and select GitHub when prompted
   gitkit init
   ```

3. **Or manually create `.gitkit_private.yml`:**
   ```yaml
   username: "your-github-username"
   accessToken: "ghp_your_token_here"
   url: "https://api.github.com" # or your GitHub Enterprise URL
   repositoryName: "your-repo-name"
   ```

### GitLab Setup

1. **Generate Personal Access Token:**

   - Go to [GitLab User Settings > Access Tokens](https://gitlab.com/-/profile/personal_access_tokens)
   - Create a new token with `api` scope
   - Copy the generated token

2. **Configure GitKit for GitLab:**

   ```bash
   # Run init and select GitLab when prompted
   gitkit init
   ```

3. **Or manually create `.gitkit_private.yml`:**
   ```yaml
   username: "your-gitlab-username"
   accessToken: "glpat_your_token_here"
   url: "https://gitlab.com" # or your GitLab instance URL
   repositoryName: "your-repo-name"
   ```

### Security Notes

- Access tokens are stored in `.gitkit_private.yml` (add this to `.gitignore`)
- Tokens require appropriate permissions for repository and merge request creation
- For GitHub Enterprise or self-hosted GitLab, update the `url` field accordingly

## Architecture

GitKit follows a clean, modular architecture:

- **`cmd/`**: Cobra CLI command definitions
- **`git/`**: Core Git operations as methods on `GitCmd` struct
- **`config/`**: YAML configuration management
- **`repo/`**: Platform integrations (GitHub, GitLab)
- **`gitignore/`**: Gitignore template management

### Key Design Patterns

- **Dependency Injection**: `GitCmd` struct with injected configuration
- **Command Pattern**: All Git operations are methods on `GitCmd`
- **Configuration-Driven**: Branch prefixes and workflows defined in YAML
- **Error Handling**: Consistent emoji-based error messages

## Development

### Prerequisites

- Go 1.25.1 or later
- Git

### Building

```bash
go build -o gitkit main.go
```

### Adding New Commands

1. Create a new command file in `cmd/`
2. Implement the command using Cobra
3. Add to `Prefixes` config if it's a workflow command
4. Register with `rootCmd.AddCommand()`

Example:

```go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "My new command",
    Run: func(cmd *cobra.Command, args []string) {
        gitCmd := git.NewGitCommandWithConfig(false)
        // Your logic here
    },
}
```

## Contributing

1. Fork the repository
2. Create a feature branch: `gitkit feature start my-feature`
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `go test ./...`
6. Commit your changes: `gitkit commit "Add my feature"`
7. Push to your fork: `gitkit push`
8. Create a Pull Request
