# Задача L2.10

## 📝 Задание
Текст можно найти в [этом](./docs/TASK.md) файле

## 🎯 Решение
Мое решение использует [внешнюю сортировку](https://en.wikipedia.org/wiki/External_sorting). Это позволяет обрабатывать данные, объем которых превышает объем оперативной памяти.

- 🚧 Чтение из файла
- ✅ Обязательные флаги
- ✅ Дополнительные флаги
- 🚧 Требования к реализации (линтеры, тесты и т.п.)

### Сборка и запуск
```bash
task build
./build/sort --help
```

### Качество кода
```bash
task install-deps # install golangci-lint to ./bin
task lint # Run linter
task vet # Run go vet
task test # Run tests
```

> [!NOTE]
> Полный список команд доступен при помощи команды `task`

> [!TIP]
> Если у вас не установлен `Task`, то его можно установить через команду:
> ```bash
> go install github.com/go-task/task/v3/cmd/task@latest
> ```
> Полная инструкция по установке есть [тут](https://taskfile.dev/docs/installation)