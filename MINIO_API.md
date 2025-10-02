# MinIO API Documentation

## Обзор

Этот модуль предоставляет два HTTP API эндпоинта для работы с MinIO:
1. **Upload API** - загрузка файлов в MinIO
2. **Download API** - скачивание файлов из MinIO

## Установка зависимостей

```bash
go get github.com/minio/minio-go/v7
```

Обновите `go.mod`:
```bash
go mod tidy
```

## Быстрый старт

### 1. Создание MinIO клиента

```go
import (
    "context"
    gotenberg "github.com/nativebpm/gotenberg-client"
)

config := gotenberg.MinioConfig{
    Endpoint:        "localhost:9000",
    AccessKeyID:     "minioadmin",
    SecretAccessKey: "minioadmin",
    BucketName:      "documents",
    UseSSL:          false,
}

ctx := context.Background()
minioClient, err := gotenberg.NewMinioClient(ctx, config)
if err != nil {
    log.Fatal(err)
}
```

### 2. Создание HTTP сервера

```go
api := gotenberg.NewMinioAPI(minioClient)

mux := http.NewServeMux()
api.RegisterRoutes(mux)

http.ListenAndServe(":8080", mux)
```

## API Эндпоинты

### 1. Upload File (Загрузка файла)

**Эндпоинт:** `POST /api/upload`

**Описание:** Загружает файл в MinIO storage.

**Content-Type:** `multipart/form-data`

**Параметры:**
- `file` (form-data, обязательный) - файл для загрузки
- `objectName` (query parameter, опциональный) - имя объекта в MinIO. Если не указано, используется оригинальное имя файла

**Пример запроса с curl:**

```bash
# Загрузка с оригинальным именем файла
curl -X POST http://localhost:8080/api/upload \
  -F "file=@/path/to/document.pdf"

# Загрузка с пользовательским именем
curl -X POST "http://localhost:8080/api/upload?objectName=my-document.pdf" \
  -F "file=@/path/to/document.pdf"
```

**Пример запроса с JavaScript:**

```javascript
const formData = new FormData();
formData.append('file', fileInput.files[0]);

fetch('http://localhost:8080/api/upload?objectName=custom-name.pdf', {
  method: 'POST',
  body: formData
})
  .then(response => response.json())
  .then(data => console.log(data));
```

**Успешный ответ (200 OK):**

```json
{
  "success": true,
  "object_name": "document.pdf",
  "size": 1024576,
  "etag": "d41d8cd98f00b204e9800998ecf8427e",
  "message": "File uploaded successfully"
}
```

**Ошибки:**
- `400 Bad Request` - неверный формат запроса или отсутствует файл
- `500 Internal Server Error` - ошибка при загрузке в MinIO

### 2. Download File (Скачивание файла)

**Эндпоинт:** `GET /api/download`

**Описание:** Скачивает файл из MinIO storage.

**Параметры:**
- `objectName` (query parameter, обязательный) - имя файла в MinIO

**Пример запроса с curl:**

```bash
# Скачивание файла
curl -X GET "http://localhost:8080/api/download?objectName=document.pdf" \
  -o downloaded-file.pdf

# Получение информации о файле (headers)
curl -I "http://localhost:8080/api/download?objectName=document.pdf"
```

**Пример запроса с JavaScript:**

```javascript
// Скачивание файла
fetch('http://localhost:8080/api/download?objectName=document.pdf')
  .then(response => response.blob())
  .then(blob => {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'document.pdf';
    a.click();
  });
```

**Успешный ответ (200 OK):**
- Headers:
  - `Content-Type`: MIME-тип файла
  - `Content-Length`: размер файла в байтах
  - `Content-Disposition`: attachment; filename="document.pdf"
  - `ETag`: хеш файла
- Body: бинарное содержимое файла

**Ошибки:**
- `400 Bad Request` - отсутствует параметр objectName
- `404 Not Found` - файл не найден в MinIO
- `500 Internal Server Error` - ошибка при скачивании из MinIO

## Дополнительные методы MinioClient

Помимо API эндпоинтов, `MinioClient` предоставляет дополнительные методы для работы с файлами:

```go
// Получение информации о файле
fileInfo, err := minioClient.GetFileInfo(ctx, "document.pdf")

// Удаление файла
err := minioClient.DeleteFile(ctx, "document.pdf")

// Получение списка файлов с префиксом
objectsCh := minioClient.ListFiles(ctx, "documents/")
for object := range objectsCh {
    if object.Err != nil {
        log.Println(object.Err)
        continue
    }
    log.Println(object.Key, object.Size)
}
```

## Запуск примера

```bash
cd examples
go run minio_api_server.go
```

Сервер запустится на порту 8080 (или порт из переменной окружения `PORT`).

## Конфигурация MinIO

Для работы с MinIO убедитесь, что:

1. MinIO сервер запущен и доступен
2. Указаны корректные credentials (AccessKeyID и SecretAccessKey)
3. Bucket существует или будет создан автоматически при первом подключении

### Docker пример для запуска MinIO:

```bash
docker run -p 9000:9000 -p 9001:9001 \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"
```

MinIO Console будет доступна по адресу: http://localhost:9001

## Безопасность

**Важно:** В production окружении:

1. Используйте HTTPS (установите `UseSSL: true`)
2. Используйте сильные пароли для MinIO credentials
3. Ограничьте размер загружаемых файлов
4. Добавьте аутентификацию для API эндпоинтов
5. Используйте CORS политики при необходимости
6. Валидируйте типы и размеры загружаемых файлов

## Лицензия

Используйте в соответствии с лицензией проекта.
