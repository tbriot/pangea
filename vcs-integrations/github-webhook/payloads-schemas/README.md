# GitHub Webhook Schemas

This directory contains JSON schemas for validating GitHub webhook payloads.

## Source

GitHub provides official schemas for webhook payloads in the following repository:
- **Repository**: [github/rest-api-description](https://github.com/github/rest-api-description)

## Schema Location

The OpenAPI 3.0.x specification is available at:
- **Path**: `/descriptions/api.github.com/api.github.com.json`

Webhook schemas are defined in the `x-webhooks` section of the OpenAPI specification.

## ⚠️ Important Notice

> **Warning**: Schemas in an OpenAPI 3.0.x file are **not** compliant with the JSON schema specification. Manual adjustments are necessary to make them valid JSON schemas.

## Schema Conversion

To convert the OpenAPI schema to a compliant JSON schema, apply the following transformations:

### 1. Add Schema Specification
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#"
}
```

### 2. Update Schema References
- Move references from `.components.schema` to `.definitions`

### 3. Handle Nullable Types
Replace this pattern:
```json
{
  "type": "string",
  "nullable": "true"
}
```

With:
```json
{
  "type": ["string", "null"]
}
```

### 4. Remove Example Entries
- Remove all `example` entries from the schema definitions

## Usage

These converted schemas can be used to:
- Validate incoming webhook payloads
- Documentation generation
- API testing and mocking
                
                