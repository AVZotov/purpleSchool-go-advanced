#!/bin/bash

# Скрипт для создания релизов субпроекта order-api
# Запускается из purple_school/order-api/

if [ $# -eq 0 ]; then
    echo "Usage: ./release.sh <version>"
    echo "Example: ./release.sh 1.0.0"
    exit 1
fi

VERSION=$1

# Проверяем, что версия в правильном формате
if ! [[ $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must be in format X.Y.Z (e.g., 1.0.0)"
    exit 1
fi

# Находим корень Git репозитория
GIT_ROOT=$(git rev-parse --show-toplevel)
if [ $? -ne 0 ]; then
    echo "Error: Not in a Git repository"
    exit 1
fi

echo "Git root: $GIT_ROOT"
echo "Current directory: $(pwd)"

# Проверяем, что мы находимся в папке order-api
CURRENT_DIR=$(basename "$(pwd)")
if [ "$CURRENT_DIR" != "order-api" ]; then
    echo "Error: This script must be run from order-api directory"
    echo "Current directory: $(pwd)"
    exit 1
fi

# Проверяем наличие go.mod для подтверждения правильности расположения
if [ ! -f "go.mod" ]; then
    echo "Error: go.mod not found. Make sure you're in the order-api project root"
    exit 1
fi

# Определяем пути относительно текущей позиции
VERSION_FILE="version.go"
PROJECT_PATH="order-api"

# Переходим в корень Git репозитория для Git операций
cd "$GIT_ROOT" || exit

# Путь к файлу версии относительно Git root
VERSION_FILE_PATH="$PROJECT_PATH/$VERSION_FILE"

# Проверяем, что мы на ветке main
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "Error: You must be on main branch to create a release"
    echo "Current branch: $CURRENT_BRANCH"
    echo "Switch to main: git checkout main"
    exit 1
fi

# Проверяем, что нет незакоммиченных изменений
if ! git diff-index --quiet HEAD --; then
    echo "Error: You have uncommitted changes. Please commit or stash them first."
    git status --porcelain
    exit 1
fi

# Проверяем, что main актуальный
git fetch origin
LOCAL=$(git rev-parse main)
REMOTE=$(git rev-parse origin/main)
if [ "$LOCAL" != "$REMOTE" ]; then
    echo "Error: Your main branch is not up to date with origin/main"
    echo "Please run: git pull origin main"
    exit 1
fi

echo "Creating release $VERSION for order-api project"

# Обновляем version.go в корне order-api
cat > "$VERSION_FILE_PATH" << EOF
package main

const (
	Version   = "$VERSION"
	BuildDate = "$(date +%Y-%m-%d)"
	AppName   = "order-api"
)
EOF

echo "Updated $VERSION_FILE_PATH"

# Коммитим изменения
git add "$VERSION_FILE_PATH"
git commit -m "order-api: bump version to $VERSION"

# Создаем тег с префиксом проекта
TAG_NAME="order-api/v$VERSION"
git tag "$TAG_NAME"

echo "Created tag $TAG_NAME"
echo ""
echo "To push changes:"
echo "  git push origin main"
echo "  git push origin $TAG_NAME"
echo ""
echo "To build order-api:"
echo "  cd $PROJECT_PATH"
echo "  go build -o bin/order-api ./cmd/main.go"
echo ""
echo "Release $VERSION created successfully!"