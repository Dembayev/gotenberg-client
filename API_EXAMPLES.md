# API Usage Examples

## cURL Examples

### Upload File

```bash
# Upload file with original name
curl -X POST http://localhost:8080/api/upload \
  -F "file=@/path/to/document.pdf"

# Upload file with custom name
curl -X POST "http://localhost:8080/api/upload?objectName=my-custom-name.pdf" \
  -F "file=@/path/to/document.pdf"
```

### Download File

```bash
# Download file
curl -X GET "http://localhost:8080/api/download?objectName=document.pdf" \
  -o downloaded.pdf

# Get file headers only
curl -I "http://localhost:8080/api/download?objectName=document.pdf"
```

## JavaScript/TypeScript Examples

### Using Fetch API

```javascript
// Upload file
async function uploadFile(file, customName = null) {
  const formData = new FormData();
  formData.append('file', file);
  
  const url = customName 
    ? `http://localhost:8080/api/upload?objectName=${encodeURIComponent(customName)}`
    : 'http://localhost:8080/api/upload';
  
  try {
    const response = await fetch(url, {
      method: 'POST',
      body: formData
    });
    
    const result = await response.json();
    console.log('Upload result:', result);
    return result;
  } catch (error) {
    console.error('Upload error:', error);
    throw error;
  }
}

// Download file
async function downloadFile(objectName, saveAs) {
  try {
    const response = await fetch(
      `http://localhost:8080/api/download?objectName=${encodeURIComponent(objectName)}`
    );
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = saveAs || objectName;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error('Download error:', error);
    throw error;
  }
}

// Usage example
const fileInput = document.getElementById('fileInput');
fileInput.addEventListener('change', async (event) => {
  const file = event.target.files[0];
  if (file) {
    await uploadFile(file, 'custom-name.pdf');
  }
});

// Download example
downloadFile('document.pdf', 'my-document.pdf');
```

### Using Axios

```javascript
import axios from 'axios';

// Upload file
async function uploadFile(file, customName = null) {
  const formData = new FormData();
  formData.append('file', file);
  
  const params = customName ? { objectName: customName } : {};
  
  try {
    const response = await axios.post(
      'http://localhost:8080/api/upload',
      formData,
      {
        params,
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      }
    );
    
    console.log('Upload result:', response.data);
    return response.data;
  } catch (error) {
    console.error('Upload error:', error.response?.data || error);
    throw error;
  }
}

// Download file
async function downloadFile(objectName, saveAs) {
  try {
    const response = await axios.get(
      'http://localhost:8080/api/download',
      {
        params: { objectName },
        responseType: 'blob'
      }
    );
    
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', saveAs || objectName);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  } catch (error) {
    console.error('Download error:', error.response?.data || error);
    throw error;
  }
}
```

## Python Examples

### Using requests library

```python
import requests

# Upload file
def upload_file(file_path, custom_name=None):
    url = 'http://localhost:8080/api/upload'
    
    params = {}
    if custom_name:
        params['objectName'] = custom_name
    
    with open(file_path, 'rb') as file:
        files = {'file': file}
        response = requests.post(url, files=files, params=params)
    
    if response.status_code == 200:
        result = response.json()
        print(f"Upload successful: {result}")
        return result
    else:
        print(f"Upload failed: {response.status_code}")
        print(response.text)
        return None

# Download file
def download_file(object_name, save_path):
    url = 'http://localhost:8080/api/download'
    params = {'objectName': object_name}
    
    response = requests.get(url, params=params, stream=True)
    
    if response.status_code == 200:
        with open(save_path, 'wb') as file:
            for chunk in response.iter_content(chunk_size=8192):
                file.write(chunk)
        print(f"File downloaded to: {save_path}")
        return True
    else:
        print(f"Download failed: {response.status_code}")
        print(response.text)
        return False

# Usage
if __name__ == '__main__':
    # Upload
    upload_file('/path/to/document.pdf', 'my-document.pdf')
    
    # Download
    download_file('my-document.pdf', '/path/to/save/document.pdf')
```

## Go Examples

### Using standard library

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
)

// UploadFile uploads a file to the API
func UploadFile(filePath string, customName string) error {
    // Open file
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    // Create multipart form
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    // Add file part
    part, err := writer.CreateFormFile("file", filepath.Base(filePath))
    if err != nil {
        return fmt.Errorf("failed to create form file: %w", err)
    }

    _, err = io.Copy(part, file)
    if err != nil {
        return fmt.Errorf("failed to copy file: %w", err)
    }

    writer.Close()

    // Create request
    url := "http://localhost:8080/api/upload"
    if customName != "" {
        url = fmt.Sprintf("%s?objectName=%s", url, customName)
    }

    req, err := http.NewRequest("POST", url, body)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    // Send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("upload failed with status: %d", resp.StatusCode)
    }

    fmt.Println("Upload successful")
    return nil
}

// DownloadFile downloads a file from the API
func DownloadFile(objectName, savePath string) error {
    // Create request
    url := fmt.Sprintf("http://localhost:8080/api/download?objectName=%s", objectName)
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("failed to download: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("download failed with status: %d", resp.StatusCode)
    }

    // Create file
    out, err := os.Create(savePath)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer out.Close()

    // Copy content
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return fmt.Errorf("failed to write file: %w", err)
    }

    fmt.Printf("File downloaded to: %s\n", savePath)
    return nil
}

func main() {
    // Upload
    err := UploadFile("/path/to/document.pdf", "my-document.pdf")
    if err != nil {
        fmt.Printf("Upload error: %v\n", err)
    }

    // Download
    err = DownloadFile("my-document.pdf", "/path/to/save/document.pdf")
    if err != nil {
        fmt.Printf("Download error: %v\n", err)
    }
}
```

## React Component Example

```typescript
import React, { useState } from 'react';
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080';

export const FileUploadDownload: React.FC = () => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [uploadStatus, setUploadStatus] = useState<string>('');
  const [objectName, setObjectName] = useState<string>('');

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files.length > 0) {
      setSelectedFile(event.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setUploadStatus('Please select a file');
      return;
    }

    const formData = new FormData();
    formData.append('file', selectedFile);

    try {
      setUploadStatus('Uploading...');
      const response = await axios.post(
        `${API_BASE_URL}/api/upload`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        }
      );

      setUploadStatus(`Upload successful! Object name: ${response.data.object_name}`);
      setObjectName(response.data.object_name);
    } catch (error) {
      setUploadStatus(`Upload failed: ${error}`);
    }
  };

  const handleDownload = async () => {
    if (!objectName) {
      alert('Please enter an object name');
      return;
    }

    try {
      const response = await axios.get(
        `${API_BASE_URL}/api/download`,
        {
          params: { objectName },
          responseType: 'blob',
        }
      );

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', objectName);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (error) {
      alert(`Download failed: ${error}`);
    }
  };

  return (
    <div style={{ padding: '20px' }}>
      <h2>MinIO File Upload/Download</h2>
      
      <div style={{ marginBottom: '20px' }}>
        <h3>Upload File</h3>
        <input type="file" onChange={handleFileSelect} />
        <button onClick={handleUpload} disabled={!selectedFile}>
          Upload
        </button>
        {uploadStatus && <p>{uploadStatus}</p>}
      </div>

      <div>
        <h3>Download File</h3>
        <input
          type="text"
          placeholder="Object name"
          value={objectName}
          onChange={(e) => setObjectName(e.target.value)}
          style={{ marginRight: '10px' }}
        />
        <button onClick={handleDownload} disabled={!objectName}>
          Download
        </button>
      </div>
    </div>
  );
};
```

## PHP Example

```php
<?php

function uploadFile($filePath, $customName = null) {
    $url = 'http://localhost:8080/api/upload';
    
    if ($customName) {
        $url .= '?objectName=' . urlencode($customName);
    }
    
    $file = new CURLFile($filePath);
    $postData = array('file' => $file);
    
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_POST, 1);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $postData);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    if ($httpCode === 200) {
        echo "Upload successful\n";
        return json_decode($response, true);
    } else {
        echo "Upload failed with status: $httpCode\n";
        return null;
    }
}

function downloadFile($objectName, $savePath) {
    $url = 'http://localhost:8080/api/download?objectName=' . urlencode($objectName);
    
    $fp = fopen($savePath, 'w+');
    
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_FILE, $fp);
    curl_setopt($ch, CURLOPT_FOLLOWLOCATION, true);
    
    curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    fclose($fp);
    
    if ($httpCode === 200) {
        echo "File downloaded to: $savePath\n";
        return true;
    } else {
        echo "Download failed with status: $httpCode\n";
        unlink($savePath);
        return false;
    }
}

// Usage
uploadFile('/path/to/document.pdf', 'my-document.pdf');
downloadFile('my-document.pdf', '/path/to/save/document.pdf');
?>
```
