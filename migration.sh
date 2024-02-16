#!/bin/bash
# Исходный Nexus
SOURCE_NEXUS=""
SOURCE_REPO=""
SOURCE_USER=''
SOURCE_PASSWORD=''
# Целевой Nexus
TARGET_NEXUS=""
TARGET_REPO=""
TARGET_USER=''
TARGET_PASSWORD=''

# Получение начального значения continuationToken
firts_token=$(curl -s -k -u "$SOURCE_USER:$SOURCE_PASSWORD" "${SOURCE_NEXUS}/service/rest/v1/search/assets?repository=${SOURCE_REPO}")
initialContinuationToken="$firts_token"
ARTIFACTS_FILE="artifacts_list.txt"

# Функция для получения списка артефактов из исходного репозитория
fetch_artifacts() {
    local continuationToken="$initialContinuationToken" # Использование начального токена, если нужно
    > "$ARTIFACTS_FILE" # Очистка или создание временного файла

    while : ; do
        local response=$(curl -k -u "$SOURCE_USER:$SOURCE_PASSWORD" "${SOURCE_NEXUS}/service/rest/v1/search/assets?repository=${SOURCE_REPO}&continuationToken=${continuationToken}")
        local paths=$(echo "$response" | jq -r '.items[] | .path')

        # Запись артефактов в файл
        for path in $paths; do
            echo "$path" >> "$ARTIFACTS_FILE"
        done

        continuationToken=$(echo "$response" | jq -r '.continuationToken // empty')
        [[ -z "$continuationToken" ]] && break
    done
}

# Функция для загрузки и повторной загрузки артефактов
migrate_artifact() {
    local path="$1"
    local source_url="${SOURCE_NEXUS}/repository/${SOURCE_REPO}/${path}"
    local target_url="${TARGET_NEXUS}/repository/${TARGET_REPO}/${path}"

    # Скачивание артефакта из исходного Nexus
    curl -k -u "$SOURCE_USER:$SOURCE_PASSWORD" -o "temp_artifact" "$source_url"

    # Загрузка артефакта в целевой Nexus
    curl -k -u "$TARGET_USER:$TARGET_PASSWORD" --upload-file "temp_artifact" "$target_url"

    # Удаление временного файла
    rm "temp_artifact"
}

# Получение списка артефактов и запись в файл
fetch_artifacts

# Чтение списка артефактов из файла и их миграция
while IFS= read -r artifact; do
    echo "Мигрируем $artifact"
    migrate_artifact "$artifact"
done < "$ARTIFACTS_FILE"

# Удаление временного файла списка артефактов
rm "$ARTIFACTS_FILE"

echo "Миграция завершена."
