# Installation

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/your-org/artifactor
cd artifactor
make build
```

The binary is written to `bin/artifactor`.

## Configuration

Artifactor is configured with a YAML file. Point to it with the `CONFIG_PATH` environment variable:

```bash
CONFIG_PATH=/etc/artifactor/config.yml ./bin/artifactor
```

**Example config:**

```yaml
file_upload_limit: 100  # MB, maximum size of an uploaded artifact

mongo:
  connection_string: mongodb://localhost:27017
  database: artifactor
  token_collection: tokens
  product_collection: products

redis:
  addr: localhost:6379
  password: ""
  db: 0
```

| Field | Description |
|---|---|
| `file_upload_limit` | Max upload size in MB |
| `mongo.connection_string` | MongoDB connection URI |
| `mongo.database` | Database name |
| `mongo.token_collection` | Collection used to store tokens |
| `mongo.product_collection` | Collection used to store products |
| `redis.addr` | Redis address (`host:port`) |
| `redis.password` | Redis password (optional) |
| `redis.db` | Redis database index (optional) |

## Running

Start the server:

```bash
CONFIG_PATH=./config.yml ./bin/artifactor
```

By default the server listens on `0.0.0.0:8080`. Override with `SERVER_ADDR`:

```bash
SERVER_ADDR=0.0.0.0:9090 CONFIG_PATH=./config.yml ./bin/artifactor
```

**Starting dependencies with Docker Compose:**

A `docker-compose.yml` is included to spin up MongoDB and Redis locally:

```bash
docker compose up -d
```

## Authentication

Every request must include an `X-Api-Token` header with a valid token.

There are two token types:

- **Admin** — can register and manage tokens, create/delete products, and perform all operations.
- **Non-admin** — access is controlled per-product through token permissions (`upload`, `download`, `delete`, `maintainer`).

### Creating the first admin token

On first run, pass `--init-admin-token` to generate an initial admin token:

```bash
CONFIG_PATH=./config.yml ./bin/artifactor --init-admin-token
```

The token is printed to the log output. Remove the flag after the first use — it is a no-op if an admin token already exists.

Use this token in the `X-Api-Token` header for all subsequent admin operations, such as registering additional tokens via `PUT /api/register`.
