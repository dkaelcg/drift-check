# drift-check

> CLI tool that detects configuration drift between live cloud resources and their Terraform state definitions.

---

## Installation

```bash
go install github.com/yourusername/drift-check@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/drift-check.git
cd drift-check
go build -o drift-check .
```

---

## Usage

Point `drift-check` at your Terraform state file and let it compare live cloud resources against the expected configuration.

```bash
# Check drift using a local state file
drift-check --state terraform.tfstate --provider aws

# Check drift against a remote S3 backend
drift-check --state s3://my-bucket/terraform.tfstate --region us-east-1

# Output results as JSON
drift-check --state terraform.tfstate --output json
```

### Example Output

```
[DRIFT DETECTED] aws_instance.web
  - instance_type: expected "t3.micro", got "t3.small"
  - tags.Env:       expected "production", got "staging"

[OK] aws_s3_bucket.assets
[OK] aws_security_group.default

Summary: 1 resource drifted, 2 resources in sync.
```

---

## Supported Providers

- AWS
- GCP *(coming soon)*
- Azure *(coming soon)*

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)