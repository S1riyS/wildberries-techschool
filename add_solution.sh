#!/bin/bash

# Проверяем, что переданы два аргумента
if [ "$#" -ne 2 ]; then
    echo "Использование: $0 <L{0..6}> <N>"
    exit 1
fi

level=$1
number=$2

# Проверяем формат первого аргумента
if [[ ! "$level" =~ ^L[0-6]$ ]]; then
    echo "Первый аргумент должен быть в формате L{0..6}"
    exit 1
fi

# Проверяем, что второй аргумент - число
if [[ ! "$number" =~ ^[0-9]+$ ]]; then
    echo "Второй аргумент должен быть числом"
    exit 1
fi

# Создаем директорию
dir_path="./${level}/${number}"
mkdir -p "$dir_path"

# Создаем файл main.go
cat > "${dir_path}/main.go" <<EOF
package main

func main() {

}
EOF

# Создаем файл README.md
cat > "${dir_path}/README.md" <<EOF
# Задача ${level}.${number}
## 📝 Задание

## ✅ Решение
Решение приведено в файле в файле [\`main.go\`](./main.go).
EOF

echo "Создана структура в ${dir_path}"