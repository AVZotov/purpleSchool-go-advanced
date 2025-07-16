#!/bin/bash

# Глобальные переменные
PROJECT_NAME="order-api"

# Скрипт для создания релизов субпроекта $PROJECT_NAME
# Запускается из purple_school/$PROJECT_NAME/

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
if ! GIT_ROOT=$(git rev-parse --show-toplevel); then
    echo "Error: Not in a Git repository"
    exit 1
fi

echo "Git root: $GIT_ROOT"
echo "Current directory: $(pwd)"

# Проверяем, что мы находимся в папке проекта
CURRENT_DIR=$(basename "$(pwd)")
if [ "$CURRENT_DIR" != "$PROJECT_NAME" ]; then
    echo "Error: This script must be run from $PROJECT_NAME directory"
    echo "Current directory: $(pwd)"
    exit 1
fi

# Проверяем наличие go.mod для подтверждения правильности расположения
if [ ! -f "go.mod" ]; then
    echo "Error: go.mod not found. Make sure you're in the $PROJECT_NAME project root"
    exit 1
fi

# Определяем пути относительно текущей позиции
VERSION_FILE="version.go"

# Переходим в корень Git репозитория для Git операций
if ! cd "$GIT_ROOT"; then
    echo "Error: Failed to change to Git root directory"
    exit 1
fi

# Путь к файлу версии относительно Git root
VERSION_FILE_PATH="$PROJECT_NAME/$VERSION_FILE"

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
if ! git fetch origin; then
    echo "Error: Failed to fetch from origin"
    exit 1
fi

LOCAL=$(git rev-parse main)
REMOTE=$(git rev-parse origin/main)
if [ "$LOCAL" != "$REMOTE" ]; then
    echo "Error: Your main branch is not up to date with origin/main"
    echo "Please run: git pull origin main"
    exit 1
fi

echo "Creating release $VERSION for $PROJECT_NAME project"

# Обновляем version.go в корне проекта
cat > "$VERSION_FILE_PATH" << EOF
package main

const (
	Version   = "$VERSION"
	BuildDate = "$(date +%Y-%m-%d)"
	AppName   = "$PROJECT_NAME"
)
EOF

echo "Updated $VERSION_FILE_PATH"

# Коммитим изменения
if ! git add "$VERSION_FILE_PATH"; then
    echo "Error: Failed to add $VERSION_FILE_PATH to git"
    exit 1
fi

if ! git commit -m "$PROJECT_NAME: bump version to $VERSION"; then
    echo "Error: Failed to commit changes"
    exit 1
fi

# Создаем тег с префиксом проекта
TAG_NAME="$PROJECT_NAME/v$VERSION"
if ! git tag "$TAG_NAME"; then
    echo "Error: Failed to create tag $TAG_NAME"
    exit 1
fi

echo "Created tag $TAG_NAME"
echo ""
echo "To push changes:"
echo "  git push origin main"
echo "  git push origin $TAG_NAME"
echo ""
echo "To build $PROJECT_NAME:"
echo "  cd $PROJECT_NAME"
echo "  go build -o bin/$PROJECT_NAME ./cmd/main.go"
echo ""
echo "Release $VERSION created successfully!"