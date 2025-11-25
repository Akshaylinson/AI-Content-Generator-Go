# AI Content Generator

A simple local AI content generation system with Go backend, local model support, and web interface.

## Features

- **Go Backend**: Single-file REST API server
- **Local Model Support**: Automatically detects and uses GGUF models
- **Web Interface**: Built-in frontend for content generation
- **SQLite Database**: Local job storage and history
- **Fallback Mode**: Works without models using templates

## Quick Start

1. **Run the application**:
   ```bash
   go mod tidy
   go run main.go
   ```

2. **Open browser**: http://localhost:8080

3. **Generate content**: Enter a topic and click "Generate Content"

## Adding Local Models

1. Download a GGUF model file (e.g., from Hugging Face)
2. Place it in the `models/` directory
3. Restart the application

**Recommended models**:
- Mistral-7B-Instruct (Q4_K_M.gguf)
- Llama-3-8B-Instruct (Q4_K_M.gguf)

## Project Structure

```
ai-content-automator/
├── main.go          # Complete application
├── go.mod           # Go dependencies
├── models/          # Place GGUF models here
└── content.db       # SQLite database (auto-created)
```

## API Endpoints

- `GET /` - Web interface
- `POST /api/jobs` - Create content generation job
- `GET /api/jobs` - List all jobs
- `GET /api/jobs/{id}` - Get specific job

## Usage Example

**Create job via API**:
```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"topic": "Artificial Intelligence"}'
```

**Get jobs**:
```bash
curl http://localhost:8080/api/jobs
```

## Requirements

- Go 1.21+
- No external dependencies for basic operation
- Optional: GGUF model files for enhanced generation
