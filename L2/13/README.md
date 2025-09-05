# Задача L2.13

## 📝 Задание
Текст можно найти в [этом](./docs/TASK.md) файле

## 🎯 Решение
Ключевой особенностью моего решения является сортировка и слияние отрезков, заданных при помощи флага `--fields`. 
Это позволяет проверить, входит ли индекс в список отрезков, за `O(log(n))` (при помощи бинарного поиска).

Основная задача:
- ✅ Чтение из файла и STDIN
- ✅ Все флаги
- ✅ Требования к реализации (линтеры, тесты и т.п.)

### Сборка и запуск
```bash
task build
./build/cut --help
```
Для проверки работоспособности команды можно использовать заранее подготовленные файлы из следующих директорий:
- [`./examples`](./examples) 
- [`./integration/testdata`](./integration/testdata)

Пример использования:
```bash
./build/cut -d "|" -f 10-100,2,1-3 -s examples/file.txt
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