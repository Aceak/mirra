#!/bin/bash

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 版本文件路径
VERSION_FILE="internal/version/version.go"

# 解析参数
BUMP_TYPE="patch"
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
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --major    Bump major version (x.0.0)"
            echo "  --minor    Bump minor version (0.x.0)"
            echo "  --help     Show this help message"
            echo ""
            echo "Default: Bump patch version (0.0.x)"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# 检查版本文件是否存在
if [[ ! -f "$VERSION_FILE" ]]; then
    echo -e "${RED}Error: Version file not found: $VERSION_FILE${NC}"
    exit 1
fi

# 提取当前版本号
CURRENT_VERSION=$(grep 'const BaseVersion' "$VERSION_FILE" | cut -d'"' -f2)
echo -e "${YELLOW}Current version: ${CURRENT_VERSION}${NC}"

# 解析版本号
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# 验证版本号格式
if ! [[ "$MAJOR" =~ ^[0-9]+$ ]] || ! [[ "$MINOR" =~ ^[0-9]+$ ]] || ! [[ "$PATCH" =~ ^[0-9]+$ ]]; then
    echo -e "${RED}Error: Invalid version format: ${CURRENT_VERSION}${NC}"
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
NEW_TAG="${NEW_VERSION}"

echo -e "${GREEN}New version: ${NEW_TAG}${NC}"

# 更新版本文件
sed -i "s/const BaseVersion = \"${CURRENT_VERSION}\"/const BaseVersion = \"${NEW_VERSION}\"/" "$VERSION_FILE"

# 检查 git 状态
if ! git diff --quiet "$VERSION_FILE"; then
    echo -e "${YELLOW}Changes made to version file:${NC}"
    git diff "$VERSION_FILE"

    # 提交更改
    git add "$VERSION_FILE"
    git commit -m "chore: bump version to ${NEW_TAG}"
    echo -e "${GREEN}Version file updated and committed${NC}"

    # 打 tag
    git tag -a "$NEW_TAG" -m "Release ${NEW_TAG}"
    echo -e "${GREEN}Tag ${NEW_TAG} created${NC}"

    echo ""
    echo -e "${GREEN}Done! To push the commit and tag, run:${NC}"
    echo "  git push origin main && git push origin ${NEW_TAG}"
else
    echo -e "${RED}Error: Failed to update version file${NC}"
    exit 1
fi
