#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://github.com/adipras/api-gladiatore"
BRANCH="main"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if git is installed
if ! command -v git &> /dev/null; then
    print_error "Git is not installed"
    exit 1
fi

# Initialize git if not already initialized
if [ ! -d .git ]; then
    print_status "Initializing git repository..."
    git init
    git remote add origin $REPO_URL
fi

# Check if remote origin exists
if ! git remote | grep -q "origin"; then
    print_status "Adding remote origin..."
    git remote add origin $REPO_URL
fi

# Get version from user or use default
if [ -z "$1" ]; then
    print_warning "No version specified. Using 'v0.1.0' as default"
    VERSION="v0.1.0"
else
    VERSION=$1
fi

# Get commit message
if [ -z "$2" ]; then
    COMMIT_MSG="Release $VERSION"
else
    COMMIT_MSG="$2"
fi

print_status "Starting release process for version $VERSION"

# Add all files
print_status "Adding files to git..."
git add .

# Check if there are changes to commit
if git diff --staged --quiet; then
    print_warning "No changes to commit"
else
    # Commit changes
    print_status "Committing changes..."
    git commit -m "$COMMIT_MSG"
fi

# Create or update branch
print_status "Checking branch..."
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "$BRANCH" ]; then
    if git show-ref --verify --quiet refs/heads/$BRANCH; then
        print_status "Switching to $BRANCH branch..."
        git checkout $BRANCH
    else
        print_status "Creating $BRANCH branch..."
        git checkout -b $BRANCH
    fi
fi

# Push to remote
print_status "Pushing to GitHub..."
if git push -u origin $BRANCH; then
    print_status "Successfully pushed to $BRANCH branch"
else
    print_error "Failed to push. You may need to set up authentication or pull first."
    print_status "Trying to pull and merge..."
    git pull origin $BRANCH --allow-unrelated-histories
    git push -u origin $BRANCH
fi

# Create a tag for release
print_status "Creating tag $VERSION..."
git tag -a $VERSION -m "Release $VERSION"

# Push tags
print_status "Pushing tags..."
if git push origin --tags; then
    print_status "Successfully pushed tags"
else
    print_error "Failed to push tags"
    exit 1
fi

# Create GitHub release using gh CLI if available
if command -v gh &> /dev/null; then
    print_status "Creating GitHub release..."
    
    # Check if authenticated
    if gh auth status &> /dev/null; then
        # Generate release notes
        RELEASE_NOTES="# Release $VERSION

## What's New
- API Gladiatore microservice generation platform
- JSON configuration validation
- Golang service code generation
- MySQL database integration
- Web interface for configuration input

## Features
- Generate complete Golang microservices from JSON config
- Automatic CRUD endpoint creation
- Database schema generation
- Docker support
- API documentation

## Installation
\`\`\`bash
git clone $REPO_URL
cd api-gladiatore
go mod download
go run cmd/server/main.go
\`\`\`

## Usage
Visit http://localhost:8080 and input your service configuration in JSON format.
"
        
        echo "$RELEASE_NOTES" | gh release create $VERSION \
            --title "Release $VERSION" \
            --notes-file - \
            --target $BRANCH
        
        print_status "GitHub release created successfully!"
    else
        print_warning "GitHub CLI not authenticated. Run 'gh auth login' to enable release creation"
        print_status "You can manually create a release at: $REPO_URL/releases/new"
    fi
else
    print_warning "GitHub CLI (gh) not installed. Install it to automatically create releases"
    print_status "You can manually create a release at: $REPO_URL/releases/new"
fi

print_status "✅ Release process completed!"
print_status "Repository: $REPO_URL"
print_status "Version: $VERSION"
print_status "Branch: $BRANCH"

# Show recent commits
echo ""
print_status "Recent commits:"
git log --oneline -5