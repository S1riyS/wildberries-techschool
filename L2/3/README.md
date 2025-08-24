# Задача L2.3
## 📝 Задание
Что выведет программа?

Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
  "fmt"
  "os"
)

func Foo() error {
  var err *os.PathError = nil
  return err
}

func main() {
  err := Foo()
  fmt.Println(err)
  fmt.Println(err == nil)
}
```

## ✅ Решение
### Вывод программы
```
<nil>
false
```

### Объяснение
**Непустой интерфейс**:
```go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}

type itab struct {
    inter  *interfacetype
    _type  *_type
    link   *itab
    bad    int32
    unused int32
    fun    [1]uintptr // variable sized
}
```
Таким образом интерфейс содержит в себе указатель на сами данные и информацию о самом интерфейсе.

**Пустой интерфейс**:
```go
type eface struct {
	_type *_type
	data  unsafe.Pointer
}
```
В пустом интерфейсе отсутствует таблица интерфейса (`itab`), так как любой тип удовлетворяет пустому интерфейсу. Поэтому хранит в себе только указатель на тип и данные.

Таким образом, при сравнении непустого интерфеса с `nil` проверяется, является ли сам интерфейс `nil` (и тип, и значение равны `nil`)

---
Код задания находится в файле [`main.go`](./main.go).
