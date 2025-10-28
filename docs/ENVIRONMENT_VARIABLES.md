## Environment Variables

The following environment variables can be used as fallbacks for function parameters:
You can also create a `.env` file in the project root directory with these variables.

The following environment variables can be used as fallbacks for function parameters:

**Required environment variables:**
- `CLOUDRU_KEY_ID`: Service account key ID (required)
- `CLOUDRU_KEY_SECRET`: Service account key secret (required)

To obtain access keys for authentication, please follow the instructions at:
https://cloud.ru/docs/console_api/ug/topics/quickstart

You will need a Key ID and Key Secret to use this service.

**Optional environment variables:**
- `CLOUDRU_REGISTRY_NAME`: Registry name
- `CLOUDRU_REGISTRY_DOMAIN`: Registry domain (defaults to 'cr.cloud.ru' if not set)
- `CLOUDRU_REPOSITORY_NAME`: Repository name (defaults to current directory name if not set)
- `CLOUDRU_PROJECT_ID`: Project ID for Container Apps (can be obtained from console.cloud.ru https://console.cloud.ru/**?projectId=YOUR_PROJECT_ID>)
- `CLOUDRU_CONTAINERAPP_NAME`: Container App name (optional)
- `CLOUDRU_DOCKERFILE`: Path to Dockerfile (defaults to 'Dockerfile' if not set)
- `CLOUDRU_DOCKERFILE_TARGET`: Target stage in a multi-stage Dockerfile (optional, defaults to '-' which means no target)
- `CLOUDRU_DOCKERFILE_FOLDER`: Dockerfile folder (build context, defaults to '.' which means current directory)
