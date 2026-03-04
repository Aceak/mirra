#!/bin/bash

set -e

# 版本文件路径
VERSION_FILE="internal/version/version.go"
CHANGELOG_FILE="CHANGELOG.md"

# 默认操作：仅打 tag（不修改版本号）
BUMP_TYPE=""
TAG_ONLY=false

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --major)
            BUMP_TYPE="major"
            shift
            ;;
        --minor)
            BUMP_TYPE="minor"
            shift
            ;;
        --patch)
            BUMP_TYPE="patch"
            shift
            ;;
        --tag-only)
            TAG_ONLY=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --major      Bump major version and create tag (x.0.0)"
            echo "  --minor      Bump minor version and create tag (0.x.0)"
            echo "  --patch      Bump patch version and create tag (0.0.x)"
            echo "  --tag-only   Only create tag with current version (default behavior)"
            echo "  --help       Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                    # Tag current version (e.g., v0.0.0)"
            echo "  $0 --major            # Bump to 1.0.0 and tag"
            echo "  $0 --minor            # Bump to 0.1.0 and tag"
            echo "  $0 --patch            # Bump to 0.0.1 and tag"
            echo "  $0 --tag-only         # Same as default, only tag current version"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# 检查版本文件是否存在
if [[ ! -f "$VERSION_FILE" ]]; then
    echo "Error: Version file not found: $VERSION_FILE"
    exit 1
fi

# 提取当前版本号
CURRENT_VERSION=$(grep 'const BaseVersion' "$VERSION_FILE" | cut -d'"' -f2)
echo "Current version: ${CURRENT_VERSION}"

# 确定新版本号
if [[ -n "$BUMP_TYPE" ]]; then
    # 解析版本号
    IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

    # 验证版本号格式
    if ! [[ "$MAJOR" =~ ^[0-9]+$ ]] || ! [[ "$MINOR" =~ ^[0-9]+$ ]] || ! [[ "$PATCH" =~ ^[0-9]+$ ]]; then
        echo "Error: Invalid version format: ${CURRENT_VERSION}"
        exit 1
    fi

    # 提升版本号
    case $BUMP_TYPE in
        major)
            MAJOR=$((MAJOR + 1))
            MINOR=0
            PATCH=0
            ;;
        minor)
            MINOR=$((MINOR + 1))
            PATCH=0
            ;;
        patch)
            PATCH=$((PATCH + 1))
            ;;
    esac

    NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
    echo "Bumping version: ${CURRENT_VERSION} -> ${NEW_VERSION}"
else
    # 不打版本，使用当前版本
    NEW_VERSION="$CURRENT_VERSION"
    if [[ "$TAG_ONLY" == true ]]; then
        echo "Tag-only mode: will create tag with current version"
    fi
fi

NEW_TAG="v${NEW_VERSION}"

# 检查 tag 是否已存在
if git rev-parse "${NEW_TAG}" >/dev/null 2>&1; then
    echo "Error: Tag ${NEW_TAG} already exists"
    exit 1
fi

# 如果需要提升版本号，更新版本文件
if [[ -n "$BUMP_TYPE" ]]; then
    # 更新版本文件
    sed -i "s/const BaseVersion = \"${CURRENT_VERSION}\"/const BaseVersion = \"${NEW_VERSION}\"/" "$VERSION_FILE"

    # 检查 git 状态
    if ! git diff --quiet "$VERSION_FILE"; then
        echo "Changes made to version file:"
        git diff "$VERSION_FILE"

        # 提交更改
        git add "$VERSION_FILE"
        git commit -m "chore: bump version to ${NEW_TAG}"
        echo "Version file updated and committed"
    else
        echo "Error: Failed to update version file"
        exit 1
    fi
fi

# 打 tag
git tag -a "$NEW_TAG" -m "Release ${NEW_TAG}"
echo "Tag ${NEW_TAG} created"

# 从 CHANGELOG 提取 release notes
RELEASE_NOTES=""
if [[ -f "$CHANGELOG_FILE" ]]; then
    # 提取当前版本的 changelog 内容
    RELEASE_NOTES=$(awk -v tag="${NEW_VERSION}" '
        /^## \[/ {
            if (p) exit
            if ($0 ~ "\\[" tag "\\]") { p=1; next }
        }
        p { print }
    ' "$CHANGELOG_FILE")

    if [[ -n "$RELEASE_NOTES" ]]; then
        echo "Found release notes in CHANGELOG:"
        echo "---"
        echo "$RELEASE_NOTES"
        echo "---"
    else
        echo "No release notes found for ${NEW_VERSION} in CHANGELOG"
        RELEASE_NOTES="Release ${NEW_TAG}"
    fi
fi

echo ""
echo "Done! To push the commit and tag, run:"
echo "  git push origin main && git push origin ${NEW_TAG}"
echo ""
if [[ -n "$BUMP_TYPE" ]]; then
    echo "Note: Version was bumped from ${CURRENT_VERSION} to ${NEW_VERSION}"
fi
