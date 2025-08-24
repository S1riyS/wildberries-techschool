# Задача L2.5
## 📝 Задание
Что выведет программа?

Объяснить вывод программы.

```go
package main

type customError struct {
  msg string
}

func (e *customError) Error() string {
  return e.msg
}

func test() *customError {
  // ... do something
  return nil
}

func main() {
  var err error
  err = test()
  if err != nil {
    println("error")
    return
  }
  println("ok")
}
```

## ✅ Решение
### Вывод
```
error
```

### Объяснение
Как было описано в задании [L2.3](../3):
> При сравнении непустого интерфеса с `nil` проверяется, является ли сам интерфейс `nil` (и тип, и значение равны `nil`).

В данном случае переменная `error` имеет значение `nil` и тип `*customError`. Поэтому условие `if err != nil` выполняется и в stdout мы видим сообщение "error".


---
Решение приведено в файле в файле [`main.go`](./main.go).
