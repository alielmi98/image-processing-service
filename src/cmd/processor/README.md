# Image Processor Service

This service consumes image processing messages from a RabbitMQ queue and processes them accordingly.

## Prerequisites

- Go 1.16 or higher
- RabbitMQ server running (or Docker)

## Configuration

Edit the `configs/config.yaml` file to match your RabbitMQ server settings.

## Building and Running

1. Build the application:
   ```bash
   cd src/cmd/processor
   go build -o image-processor
   ```

2. Run the processor:
   ```bash
   ./image-processor
   ```

## Environment Variables

You can override configuration using environment variables:

```bash
export RABBITMQ_URL=amqp://user:pass@localhost:5672/
export RABBITMQ_QUEUE_NAME=my_image_queue
```

## Message Format

The processor expects messages in the following JSON format:

```json
{
  "job_id": 123,
  "image_id": 456,
  "user_id": 789,
  "processing_type": "resize",
  "parameters": {
    "width": 800,
    "height": 600,
    "maintain_ratio": true,
    "quality": 90,
    "format": "jpg"
  },
  "source_path": "/path/to/source/image.jpg",
  "destination_dir": "/path/to/output/",
  "priority": 5,
  "timestamp": "2023-01-01T12:00:00Z",
  "retry_count": 0,
  "max_retries": 3
}
```

## Supported Processing Types

- `resize`: Resize the image
- `crop`: Crop the image
- `rotate`: Rotate the image
- `filter`: Apply image filters
- `watermark`: Add watermark
- `compress`: Compress the image
- `format`: Convert image format

## Logging

Logs are written to stdout. You can redirect them to a file:

```bash
./image-processor > processor.log 2>&1
```
