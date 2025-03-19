#!/bin/bash

set -e

# Check if version is provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 1.0.0"
  exit 1
fi

VERSION=$1

# Validate version format (semver)
if ! [[ $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9\.]+)?(\+[a-zA-Z0-9\.]+)?$ ]]; then
  echo "Error: Version must follow semantic versioning (x.y.z)"
  exit 1
fi

# Get current version
CURRENT_VERSION=$(cat VERSION)
echo "Current version: $CURRENT_VERSION"
echo "New version: $VERSION"

# Check if version already exists in CHANGELOG.md
if grep -q "\[$VERSION\]" CHANGELOG.md; then
  echo "Error: Version $VERSION already exists in CHANGELOG.md"
  exit 1
fi

# Update VERSION file
echo "$VERSION" > VERSION
echo "✅ Updated VERSION file"

# Update CHANGELOG.md
DATE=$(date +"%Y-%m-%d")
sed -i.bak "s/## \[Unreleased\]/## \[Unreleased\]\n\n### Added\n- None\n\n## \[$VERSION\] - $DATE/" CHANGELOG.md
rm CHANGELOG.md.bak
echo "✅ Updated CHANGELOG.md file"

# Create a release config from template
cp .release.yml.template .release.yml
sed -i.bak "s/version: \".*\"/version: \"$VERSION\"/" .release.yml
rm .release.yml.bak
echo "✅ Created .release.yml file"

# Instructions for next steps
echo
echo "🎉 Release preparation complete!"
echo
echo "Next steps:"
echo "1. Update the CHANGELOG.md with the changes for this release"
echo "2. Verify the .release.yml file"
echo "3. Commit the changes: git add VERSION CHANGELOG.md .release.yml"
echo "4. Create a commit: git commit -m \"chore: prepare release v$VERSION\""
echo "5. Create a tag: git tag v$VERSION"
echo "6. Push changes: git push && git push --tags"
echo "7. Run 'task release' to create the GitHub release" 