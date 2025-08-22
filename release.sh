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

# Update or create CHANGELOG.md
print_status "Updating CHANGELOG.md..."
update_changelog() {
    local TEMP_FILE=$(mktemp)
    local TODAY=$(date +%Y-%m-%d)
    
    # Check if CHANGELOG.md exists
    if [ ! -f CHANGELOG.md ]; then
        print_status "Creating new CHANGELOG.md..."
        cat > CHANGELOG.md << EOF
# Changelog

All notable changes to API Gladiatore will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

EOF
    fi
    
    # Prepare new version entry
    cat > "$TEMP_FILE" << EOF
# Changelog

All notable changes to API Gladiatore will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [$VERSION] - $TODAY

EOF

    # Analyze git changes to categorize updates
    print_status "Analyzing changes since last release..."
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_warning "Not a git repository. Initializing git..."
        git init
    fi
    
    # Get the last tag for comparison
    LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    
    if [ -z "$LAST_TAG" ]; then
        # First release - count all files in the repository
        echo "### Added" >> "$TEMP_FILE"
        echo "- Initial release of API Gladiatore platform" >> "$TEMP_FILE"
        echo "- JSON-to-Microservice generation engine" >> "$TEMP_FILE"
        echo "- Service management dashboard with React SPA" >> "$TEMP_FILE"
        echo "- Real-time service and endpoint status management" >> "$TEMP_FILE"
        echo "- MySQL database integration with GORM" >> "$TEMP_FILE"
        echo "- Tailwind CSS styled interface with modern UI/UX" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        
        # For first release, count all files in repository
        echo "### Summary" >> "$TEMP_FILE"
        local TOTAL_FILES=$(find . -type f -not -path "./.git/*" -not -path "./node_modules/*" -not -path "./frontend/node_modules/*" -not -path "./frontend/build/*" -not -path "./generated/*" 2>/dev/null | wc -l | tr -d ' ')
        local GO_FILES=$(find . -name "*.go" -not -path "./generated/*" 2>/dev/null | wc -l | tr -d ' ')
        local TS_FILES=$(find . -name "*.ts" -o -name "*.tsx" -not -path "./node_modules/*" -not -path "./frontend/node_modules/*" 2>/dev/null | wc -l | tr -d ' ')
        local JS_FILES=$(find . -name "*.js" -o -name "*.jsx" -not -path "./node_modules/*" -not -path "./frontend/node_modules/*" -not -path "./frontend/build/*" 2>/dev/null | wc -l | tr -d ' ')
        echo "- Total files: $TOTAL_FILES" >> "$TEMP_FILE"
        echo "- Go files: $GO_FILES" >> "$TEMP_FILE"
        echo "- TypeScript/React files: $TS_FILES" >> "$TEMP_FILE"
        echo "- JavaScript files: $JS_FILES" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    else
        # Analyze changes since last tag
        local ADDED_FILES=$(git diff --name-status "$LAST_TAG"..HEAD 2>/dev/null | grep "^A" | wc -l | tr -d ' ')
        local MODIFIED_FILES=$(git diff --name-status "$LAST_TAG"..HEAD 2>/dev/null | grep "^M" | wc -l | tr -d ' ')
        local DELETED_FILES=$(git diff --name-status "$LAST_TAG"..HEAD 2>/dev/null | grep "^D" | wc -l | tr -d ' ')
        
        # Check for specific changes
        local HAS_ADDED=false
        local HAS_CHANGED=false
        local HAS_FIXED=false
        local HAS_REMOVED=false
        
        # Analyze commit messages for categories
        if git log "$LAST_TAG"..HEAD --oneline 2>/dev/null | grep -qi "add\|feat\|feature\|new"; then
            HAS_ADDED=true
            echo "### Added" >> "$TEMP_FILE"
            git log "$LAST_TAG"..HEAD --oneline --grep="^feat\|^add\|^feature" --pretty=format:"- %s" 2>/dev/null | sed 's/^feat: //i; s/^add: //i; s/^feature: //i' >> "$TEMP_FILE"
            echo "" >> "$TEMP_FILE"
            echo "" >> "$TEMP_FILE"
        fi
        
        if git log "$LAST_TAG"..HEAD --oneline 2>/dev/null | grep -qi "change\|update\|improve\|enhance"; then
            HAS_CHANGED=true
            echo "### Changed" >> "$TEMP_FILE"
            git log "$LAST_TAG"..HEAD --oneline --grep="^change\|^update\|^improve" --pretty=format:"- %s" 2>/dev/null | sed 's/^change: //i; s/^update: //i; s/^improve: //i' >> "$TEMP_FILE"
            echo "" >> "$TEMP_FILE"
            echo "" >> "$TEMP_FILE"
        fi
        
        if git log "$LAST_TAG"..HEAD --oneline 2>/dev/null | grep -qi "fix\|bug\|patch\|correct"; then
            HAS_FIXED=true
            echo "### Fixed" >> "$TEMP_FILE"
            git log "$LAST_TAG"..HEAD --oneline --grep="^fix\|^bug" --pretty=format:"- %s" 2>/dev/null | sed 's/^fix: //i; s/^bug: //i' >> "$TEMP_FILE"
            echo "" >> "$TEMP_FILE"
            echo "" >> "$TEMP_FILE"
        fi
        
        if git log "$LAST_TAG"..HEAD --oneline 2>/dev/null | grep -qi "remove\|delete"; then
            HAS_REMOVED=true
            echo "### Removed" >> "$TEMP_FILE"
            git log "$LAST_TAG"..HEAD --oneline --grep="^remove\|^delete" --pretty=format:"- %s" 2>/dev/null | sed 's/^remove: //i; s/^delete: //i' >> "$TEMP_FILE"
            echo "" >> "$TEMP_FILE"
            echo "" >> "$TEMP_FILE"
        fi
        
        # If no categorized commits, show all recent commits
        if [ "$HAS_ADDED" = false ] && [ "$HAS_CHANGED" = false ] && [ "$HAS_FIXED" = false ] && [ "$HAS_REMOVED" = false ]; then
            # Check if there are any commits at all
            COMMIT_COUNT=$(git log "$LAST_TAG"..HEAD --oneline 2>/dev/null | wc -l | tr -d ' ')
            if [ "$COMMIT_COUNT" -gt 0 ]; then
                echo "### Changed" >> "$TEMP_FILE"
                git log "$LAST_TAG"..HEAD --oneline --pretty=format:"- %s" 2>/dev/null | head -10 >> "$TEMP_FILE"
                echo "" >> "$TEMP_FILE"
                echo "" >> "$TEMP_FILE"
            else
                echo "### Changed" >> "$TEMP_FILE"
                echo "- Minor updates and improvements" >> "$TEMP_FILE"
                echo "" >> "$TEMP_FILE"
            fi
        fi
        
        # Add summary statistics
        echo "### Summary" >> "$TEMP_FILE"
        echo "- Files added: $ADDED_FILES" >> "$TEMP_FILE"
        echo "- Files modified: $MODIFIED_FILES" >> "$TEMP_FILE"
        echo "- Files removed: $DELETED_FILES" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
    fi
    
    # Append existing changelog content (skip the header)
    if [ -f CHANGELOG.md ]; then
        # Skip the first few header lines and append the rest
        tail -n +7 CHANGELOG.md >> "$TEMP_FILE" 2>/dev/null || true
    fi
    
    # Replace the original file
    mv "$TEMP_FILE" CHANGELOG.md
    
    print_status "CHANGELOG.md updated successfully"
}

update_changelog

# Add all files including the updated CHANGELOG
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