Этот пакет реализует некоторые алгоритмы кодирования информации с помощью штрих-кодов. 
Пакет разледен на две части: линейные штрих-коды (https://godoc.org/github.com/6d7367/gobarcode/linear) и штрих-коды, используемые для почтовых отправлений (https://godoc.org/github.com/6d7367/gobarcode/post)

Небольшой пример:

```go
msg := "1234567890"
f, _ := os.Create("codabar.png")
b := linear.NewCodabar(msg)
img := b.GetImage()
png.Encode(f, img)
f.Close()
```

![codabar.png](codabar.png)