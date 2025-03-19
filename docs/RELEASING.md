# Creating Releases

This document outlines the process for creating new releases of the NetGeX project.

## Prerequisites

Make sure you have the following tools installed:

- [GitHub CLI (gh)](https://cli.github.com/) - For creating GitHub releases
- [yq](https://github.com/mikefarah/yq) - For parsing YAML files
- [envsubst](https://www.gnu.org/software/gettext/manual/html_node/envsubst-Invocation.html) - Part of the gettext package

You'll also need appropriate permissions to create releases in the GitHub repository.

## Release Process Overview

1. Prepare a new release
2. Update the CHANGELOG.md with relevant changes
3. Verify the release configuration
4. Commit and tag the changes
5. Create the GitHub release

## Step-by-Step Guide

### 1. Prepare a new release

Run the prepare-release task with the new version:

```bash
task prepare-release -- 1.2.3
```

This will:
- Update the VERSION file
- Update the CHANGELOG.md with a new version entry
- Create a .release.yml file from the template

### 2. Update the CHANGELOG

Edit the CHANGELOG.md file to add details about the changes in this release:

```markdown
## [1.2.3] - 2024-03-19

### Added
- New feature X
- New feature Y

### Changed
- Improved performance of Z
- Updated dependency A to version 2.0.0

### Fixed
- Bug in component B
- Issue with configuration C
```

### 3. Verify the release configuration

Review and update the .release.yml file if needed:

```yaml
# Version should match the content of the VERSION file
version: "1.2.3" 

# Release title (${VERSION} will be replaced with the version)
title: "NetGeX v${VERSION}"

# Release notes
notes: |
  ## What's Changed

  This release adds several new features and fixes important bugs.

  ## Full Changelog
  <!-- The task will automatically include this -->
```

### 4. Commit and tag the changes

```bash
git add VERSION CHANGELOG.md .release.yml
git commit -m "chore: prepare release v1.2.3"
git tag v1.2.3
git push && git push --tags
```

### 5. Create the GitHub release

Run the release task:

```bash
task release
```

This will:
- Build release artifacts for multiple platforms
- Generate checksums for the artifacts
- Create a GitHub release with the artifacts
- Include the changelog and release notes in the release description

## After Release

After creating a release, you should:

1. Verify the release on GitHub
2. Update the VERSION file and CHANGELOG.md in the main branch to start the next development cycle
3. Share the release announcement with the community

## Troubleshooting

### Common Issues

**Issue**: The release task fails with "yq not found"
**Solution**: Install yq using your package manager or from the [releases page](https://github.com/mikefarah/yq/releases)

**Issue**: GitHub CLI authentication error
**Solution**: Run `gh auth login` to authenticate with GitHub

**Issue**: Release already exists
**Solution**: Delete the existing tag and release with `git tag -d v1.2.3 && git push origin :refs/tags/v1.2.3 && gh release delete v1.2.3` 