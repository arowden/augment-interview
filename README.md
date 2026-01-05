# Augment Fund Cap Table

A Go-based API for managing investment fund cap tables, tracking unit ownership and transfers between parties.

## Overview

This system provides a complete solution for:
- Creating and managing investment funds with fixed ownership units
- Tracking ownership (cap table) with full audit history
- Executing transfers between owners with validation and idempotency support
- React frontend for fund management

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    Frontend     │────▶│    Go API       │────▶│   PostgreSQL    │
│  (React/Vite)   │     │  (Chi Router)   │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                               │
                    ┌──────────┴──────────┐
                    ▼          ▼          ▼
              ┌─────────┐ ┌─────────┐ ┌─────────┐
              │  Fund   │ │Ownership│ │Transfer │
              │ Service │ │ Service │ │ Service │
              └─────────┘ └─────────┘ └─────────┘
```

### Domain Model

- **Fund**: An investment vehicle with a fixed number of ownership units
- **Cap Table**: The authoritative record of who owns what percentage of a fund
- **Transfer**: Movement of units from one owner to another

### Key Invariants

- Total units across all cap table entries must equal fund's total units
- Transfers must not result in negative ownership
- Owner cannot transfer to themselves
- Transfer units must be positive

## Database Schema

### Tables

#### `funds`
| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `name` | TEXT | Fund name (max 255 chars) |
| `total_units` | INTEGER | Total units issued (immutable) |
| `created_at` | TIMESTAMPTZ | Creation timestamp |

#### `cap_table_entries`
| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `fund_id` | UUID | Foreign key to funds |
| `owner_name` | TEXT | Owner identifier |
| `units` | INTEGER | Current units owned |
| `acquired_at` | TIMESTAMPTZ | Initial acquisition time |
| `updated_at` | TIMESTAMPTZ | Last modification time |
| `deleted_at` | TIMESTAMPTZ | Soft delete (audit trail) |

**Constraints**: `UNIQUE(fund_id, owner_name)`

#### `transfers`
| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `fund_id` | UUID | Foreign key to funds |
| `from_owner` | TEXT | Sender |
| `to_owner` | TEXT | Recipient |
| `units` | INTEGER | Units transferred |
| `idempotency_key` | UUID | Client deduplication key |
| `transferred_at` | TIMESTAMPTZ | Execution timestamp |

**Constraints**:
- `from_owner <> to_owner` (no self-transfers)
- Foreign keys to `cap_table_entries` for both owners
- Unique index on `idempotency_key` (when not null)

### Entity Relationship

```
┌─────────┐       ┌───────────────────┐       ┌───────────┐
│  funds  │──1:N──│ cap_table_entries │──N:1──│ transfers │
└─────────┘       └───────────────────┘       └───────────┘
                           │                        │
                           └────────────────────────┘
                            Both from_owner & to_owner
                            reference cap_table_entries
```

## API Endpoints

Base URL: `/api`

### Funds

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/funds` | List all funds (paginated) |
| `POST` | `/funds` | Create a new fund |
| `GET` | `/funds/{fundId}` | Get fund by ID |

### Cap Table

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/funds/{fundId}/cap-table` | Get ownership table |

### Transfers

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/funds/{fundId}/transfers` | List transfers (paginated) |
| `POST` | `/funds/{fundId}/transfers` | Execute a transfer |

### Pagination

All list endpoints support:
- `limit` (default: 100, max: 1000)
- `offset` (default: 0)

### Error Codes

| Code | Description |
|------|-------------|
| `INVALID_REQUEST` | Malformed request body |
| `INVALID_FUND` | Fund validation failed |
| `FUND_NOT_FOUND` | Fund does not exist |
| `OWNER_NOT_FOUND` | Owner not in cap table |
| `INSUFFICIENT_UNITS` | Sender lacks units |
| `SELF_TRANSFER` | Cannot transfer to self |
| `DUPLICATE_TRANSFER` | Idempotency key conflict |
| `INTERNAL_ERROR` | Server error |

### Idempotency

Transfer requests support an optional `idempotencyKey` (UUID):
- First request: Creates transfer, returns `201`
- Duplicate with same data: Returns original, `200`
- Duplicate with different data: Returns `409 Conflict`

## Project Structure

```
.
├── api/
│   └── openapi.yaml          # OpenAPI 3.0 specification
├── cmd/
│   └── server/
│       └── main.go           # Application entrypoint
├── frontend/                 # React frontend (Vite + Tailwind)
│   ├── src/
│   │   ├── components/       # React components
│   │   ├── hooks/            # React Query hooks
│   │   └── pages/            # Page components
│   └── e2e/                  # Playwright E2E tests
├── internal/
│   ├── config/               # Environment configuration
│   ├── fund/                 # Fund domain (entity, service, store)
│   ├── http/                 # HTTP handlers (generated + custom)
│   ├── ownership/            # Cap table domain
│   ├── postgres/             # Database pool, migrations
│   ├── transfer/             # Transfer domain
│   └── validation/           # Shared validation constraints
└── Makefile
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Node.js 18+ (for frontend)
- Docker (optional, for local development)

### Environment Variables

#### Required

| Variable | Description |
|----------|-------------|
| `DB_HOST` | PostgreSQL host |
| `DB_USER` | Database user |
| `DB_PASSWORD` | Database password |
| `DB_NAME` | Database name |

#### Optional

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_SSLMODE` | `require` | SSL mode |
| `DB_MAX_CONNS` | `25` | Max pool connections |
| `DB_MIN_CONNS` | `5` | Min pool connections |
| `SERVER_HOST` | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |

### Running Locally

1. **Start PostgreSQL** (using Docker):
   ```bash
   docker run -d --name augment-postgres \
     -e POSTGRES_USER=augment \
     -e POSTGRES_PASSWORD=secret \
     -e POSTGRES_DB=augment_fund \
     -p 5432:5432 \
     postgres:16
   ```

2. **Set environment variables**:
   ```bash
   export DB_HOST=localhost
   export DB_USER=augment
   export DB_PASSWORD=secret
   export DB_NAME=augment_fund
   export DB_SSLMODE=disable
   ```

3. **Run the server**:
   ```bash
   make run
   ```

4. **Start the frontend** (optional):
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

### API Examples

**Create a fund**:
```bash
curl -X POST http://localhost:8080/api/funds \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Growth Fund I",
    "totalUnits": 1000000,
    "initialOwner": "Founder LLC"
  }'
```

**Get cap table**:
```bash
curl http://localhost:8080/api/funds/{fundId}/cap-table
```

**Execute transfer**:
```bash
curl -X POST http://localhost:8080/api/funds/{fundId}/transfers \
  -H "Content-Type: application/json" \
  -d '{
    "fromOwner": "Founder LLC",
    "toOwner": "Investor A",
    "units": 100000,
    "idempotencyKey": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

## Development

### Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build server binary |
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage |
| `make lint` | Run golangci-lint |
| `make generate-api` | Regenerate OpenAPI code |
| `make run` | Build and run server |

### Running Tests

```bash
# Unit and integration tests (requires Docker for testcontainers)
make test

# With coverage report
make test-coverage
```

### Code Generation

The HTTP handlers are generated from the OpenAPI spec:

```bash
make generate-api
```

## AWS Deployment

The `deploy/terraform/` directory contains Terraform modules for deploying to AWS.

### Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                              VPC                                     │
│  ┌─────────────────────────┐    ┌─────────────────────────────────┐ │
│  │     Public Subnets      │    │       Private Subnets           │ │
│  │  ┌───────────────────┐  │    │  ┌───────────┐  ┌───────────┐  │ │
│  │  │        ALB        │  │    │  │ECS Fargate│  │    RDS    │  │ │
│  │  │   (HTTP/HTTPS)    │──┼────┼─▶│   Tasks   │──│PostgreSQL │  │ │
│  │  └───────────────────┘  │    │  └───────────┘  └───────────┘  │ │
│  └─────────────────────────┘    │        │                        │ │
│                                  │        ▼                        │ │
│  ┌─────────────────────────┐    │  ┌───────────┐                  │ │
│  │   S3 (Frontend SPA)     │    │  │    ECR    │ VPC Endpoints:   │ │
│  │   Static Website        │    │  │  (Images) │ - ECR API/Docker │ │
│  └─────────────────────────┘    │  └───────────┘ - S3, CloudWatch │ │
│                                  └─────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

### Infrastructure Components

| Component | Description | Cost (est.) |
|-----------|-------------|-------------|
| **VPC** | 2 AZ setup with public/private subnets | ~$0 |
| **VPC Endpoints** | ECR, S3, CloudWatch (no NAT Gateway) | ~$22/mo |
| **RDS PostgreSQL** | db.t3.micro, 20GB, encrypted, 7-day backups | ~$15/mo |
| **ECS Fargate** | 2 tasks (0.5 vCPU, 1GB each) | ~$37/mo |
| **ALB** | Application Load Balancer + data transfer | ~$20/mo |
| **ECR** | Container registry (storage-based) | ~$1/mo |
| **S3** | Frontend hosting (storage + requests) | ~$1/mo |
| **CloudWatch** | Logs, alarms, metrics | ~$1/mo |
| **Total** | | **~$97/mo** |

### Prerequisites

1. **AWS CLI** configured with credentials
2. **Terraform** 1.0+
3. **Docker** for building container images
4. **(Optional)** Domain name and ACM certificate for HTTPS

### Required AWS Permissions

The deploying IAM user/role needs these permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "vpc:*", "ec2:*", "rds:*", "ecs:*", "ecr:*",
        "elasticloadbalancing:*", "s3:*", "logs:*",
        "cloudwatch:*", "secretsmanager:*", "kms:*",
        "iam:CreateRole", "iam:DeleteRole", "iam:AttachRolePolicy",
        "iam:DetachRolePolicy", "iam:PutRolePolicy", "iam:GetRole",
        "iam:PassRole", "iam:CreatePolicy", "iam:DeletePolicy",
        "acm:*"
      ],
      "Resource": "*"
    }
  ]
}
```

### Deployment Steps

1. **Configure variables**:
   ```bash
   cd deploy/terraform
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your values
   ```

2. **Initialize Terraform**:
   ```bash
   terraform init
   ```

3. **Review plan**:
   ```bash
   terraform plan
   ```

4. **Deploy infrastructure**:
   ```bash
   terraform apply
   ```

5. **Build and push container image**:
   ```bash
   # Get ECR URL from outputs
   ECR_URL=$(terraform output -raw ecr_repository_url)

   # Login to ECR
   aws ecr get-login-password --region us-east-1 | \
     docker login --username AWS --password-stdin $ECR_URL

   # Build and push
   docker build -t $ECR_URL:latest .
   docker push $ECR_URL:latest
   ```

6. **Deploy frontend**:
   ```bash
   cd frontend
   npm run build
   aws s3 sync dist/ s3://$(terraform output -raw frontend_bucket_name)
   ```

### HTTPS Setup (Optional)

To enable HTTPS, you need an ACM certificate:

1. **Request certificate** in AWS Console (ACM) for your domain
2. **Validate** via DNS or email
3. **Update terraform.tfvars**:
   ```hcl
   domain_name = "api.yourdomain.com"
   ```
4. **Apply changes**:
   ```bash
   terraform apply
   ```
5. **Point DNS** to the ALB DNS name (from `terraform output alb_dns_name`)

### Outputs

After deployment, Terraform outputs:

| Output | Description |
|--------|-------------|
| `api_url` | API endpoint URL |
| `alb_dns_name` | ALB DNS name for DNS configuration |
| `ecr_repository_url` | ECR URL for Docker push |
| `frontend_website_url` | Frontend S3 website URL |
| `rds_endpoint` | RDS connection endpoint |

### Destroy

To tear down all resources:

```bash
terraform destroy
```

**Note**: RDS has `deletion_protection = true`. Disable it first:
```bash
terraform apply -var="enable_deletion_protection=false"
terraform destroy
```

## License

MIT
