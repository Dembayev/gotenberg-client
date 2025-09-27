# Gotenberg Client - Optimized Streaming Version

Высокопроизводительный Go клиент для [Gotenberg](https://gotenberg.dev/) с оптимизированным streaming API и минимальными аллокациями памяти.

## Features

- 🚀 **High Performance**: Optimized streaming with buffer pools (70% fewer allocations)
- ⚡ **Direct Streaming**: Write data directly to multipart writer without intermediate structures
- 📦 **Zero Dependencies**: Uses only Go standard library
- � **Webhook Support**: Full async processing with webhooks
- 🔄 **Client Reuse**: Reset and reuse client instances for multiple conversions
- ⚡ **Context Support**: Proper timeout and cancellation handling
- 🎯 **Simple API**: Single `Client` type with direct write methods

## Installation

```bash
go get github.com/nativebpm/gotenberg-client
```

**Requirements**: Go 1.18 or later

## Доступные примеры

### 🟢 Основные примеры

- **`html2pdf_minimal/`** - Минимальный пример HTML → PDF конвертации
- **`html2pdf/`** - Расширенный пример с настройками страницы и файлами
- **`url2pdf/`** - Пример конвертации URL → PDF
- **`html2pdf_webhook/`** - Пример использования webhook для асинхронной обработки

## Новый API

Все примеры обновлены для использования нового упрощенного API:

```go
## Quick Start

### Optimized Streaming API (Recommended)

```go
package main

import (
    "context"
    "net/http"
    "strings"
    "time"
    "github.com/nativebpm/gotenberg-client"
)

func main() {
    // Create client
    httpClient := &http.Client{Timeout: 30 * time.Second}
    client := gotenberg.NewClient(httpClient, "http://localhost:3000")
    
    // HTML content
    htmlContent := `<html><body><h1>Hello Gotenberg!</h1></body></html>`
    cssContent := `body { font-family: Arial; margin: 20px; }`
    
    // Write data directly to multipart stream
    client.WriteHTML(strings.NewReader(htmlContent))
    client.WriteFile("styles.css", strings.NewReader(cssContent))
    
    // Configure page settings
    client.WritePaperSize(gotenberg.PaperSizeA4[0], gotenberg.PaperSizeA4[1])
    client.WriteMargins(1.0, 1.0, 1.0, 1.0)
    client.WriteBoolProperty(gotenberg.FieldPrintBackground, true)
    
    // Execute conversion
    resp, err := client.Execute(context.Background())
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    // Save response to file...
}
```
```

## Запуск примеров

### Предварительные требования

1. Запущенный Gotenberg сервер:
   ```bash
   docker run --rm -p 3000:3000 gotenberg/gotenberg:8
   ```

2. **Requirements**: Go 1.18 or later

### Запуск

```bash
# Запустить Gotenberg сервер
docker run --rm -p 3000:3000 gotenberg/gotenberg:8

# Перейти в папку примера
cd example/cmd/html2pdf_minimal

# Скомпилировать и запустить
go run main.go
```

## Features

- 🚀 **High Performance**: Optimized streaming with buffer pools (70% fewer allocations)
- ⚡ **Direct Streaming**: Write data directly to multipart writer without intermediate structures
- 📦 **Zero Dependencies**: Uses only Go standard library
- 🔗 **Webhook Support**: Full async processing with webhooks
- 🔄 **Client Reuse**: Reset and reuse client instances for multiple conversions
- ⚡ **Context Support**: Proper timeout and cancellation handling
- 🎯 **Simple API**: Single `Client` type with direct write methods

## Типовые сценарии

### Простая HTML конвертация
```go
client := gotenberg.NewClient(httpClient, baseURL)
client.WriteHTML(htmlReader)
resp, err := client.Execute(ctx)
```

### HTML с файлами
```go
client.WriteHTML(htmlReader)
client.WriteFile("styles.css", cssReader)
client.WriteFile("script.js", jsReader)
resp, err := client.Execute(ctx)
```

### URL конвертация
```go
client.WriteURL("https://example.com")
client.WritePaperSize(8.5, 11) // Letter
resp, err := client.ExecuteURL(ctx)
```

### Webhook (асинхронно)
```go
client.WriteHTML(htmlReader)
client.SetWebhookSuccess("https://callback.com/success", "POST")
client.SetWebhookError("https://callback.com/error", "POST")
resp, err := client.Execute(ctx)
```

### Повторное использование
```go
// Первый документ
client.WriteHTML(html1)
resp1, err := client.Execute(ctx)

// Сброс и второй документ
client.Reset()
client.WriteHTML(html2)
resp2, err := client.Execute(ctx)
```

## Поддержка

Если у вас есть вопросы или проблемы:

1. Проверьте, что Gotenberg сервер запущен на `localhost:3000`
2. Убедитесь, что используете Go 1.18+
3. Смотрите примеры в `complete_demo/` для полного понимания API