# API Usage

### 1. Text

```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ],
    "cunet_id": "xxxxxxxxxx"
  }'
```

### 2. Image

```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: multipart/form-data" \
  -F "messages=[{\"role\":\"user\", \"content\":\"What is in this image?\"}]" \
  -F "cunet_id=xxxxxxxxxx" \
  -F "image=@/path/to/your/image.png"
```

### 3. File

```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: multipart/form-data" \
  -F "messages=[{\"role\":\"user\", \"content\":\"Summarize this document\"}]" \
  -F "cunet_id=xxxxxxxxxx" \
  -F "file=@/path/to/your/document.pdf"
```

**Note:** Ensure your file path is correct. The `@` symbol is required by `curl` to upload the file content.
