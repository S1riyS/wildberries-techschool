# Задача L2.11

## 📝 Задание
Текст можно найти в [этом](./docs/TASK.md) файле

## 🎯 Решение
### Ручное тестирование
```bash
task run
```

### Качество кода
```bash
task install-deps # install golangci-lint в ./bin
task test # Run tests
task lint # Run linter
task vet # Run go vet
```

> [!NOTE]
> Полный список команд доступен при помощи команды `task`

> [!TIP]
> Если у вас не установлен `Task`, то его можно установить через команду:
> ```bash
> go install github.com/go-task/task/v3/cmd/task@latest
> ```
> Полная инструкция по установке есть [тут](https://taskfile.dev/docs/installation)