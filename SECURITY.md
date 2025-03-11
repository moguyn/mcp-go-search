# Security Policy and Best Practices

This document outlines security considerations and best practices for the Bocha AI Search MCP Server.

## API Key Security

The Bocha AI Search MCP Server requires an API key to authenticate with the Bocha AI Search API. This API key should be treated as a sensitive credential.

### Best Practices for API Key Management

1. **Environment Variables**: Always use environment variables to provide the API key to the application rather than hardcoding it in configuration files or source code.

   ```bash
   export BOCHA_API_KEY="your-api-key-here"
   ```

2. **Configuration Files**: If you must use a configuration file, ensure it:
   - Is not committed to version control (add it to `.gitignore`)
   - Has restricted file permissions (`chmod 600 config.yaml`)
   - Is stored in a secure location

3. **Secrets Management**: In production environments, consider using a secrets management solution like:
   - HashiCorp Vault
   - AWS Secrets Manager
   - Google Secret Manager
   - Azure Key Vault

## Secure Deployment

### Network Security

1. **Firewall Rules**: Restrict network access to the server to only trusted clients.

2. **TLS**: If exposing the server over a network, ensure all communication is encrypted using TLS 1.2 or higher.

3. **Reverse Proxy**: Consider using a reverse proxy like Nginx or Caddy to handle TLS termination and provide additional security features.

### Container Security

If deploying in containers:

1. **Minimal Base Images**: Use minimal base images like `alpine` or `distroless`.

2. **Non-root User**: Run the container as a non-root user.

   ```dockerfile
   USER nobody
   ```

3. **Read-only Filesystem**: Mount the filesystem as read-only where possible.

   ```dockerfile
   VOLUME ["/tmp"]
   ```

## Input Validation and Rate Limiting

The server implements several security measures:

1. **Input Validation**: All user inputs are validated before processing.

2. **Query Sanitization**: Search queries are sanitized to prevent potential injection attacks.

3. **Rate Limiting**: The server implements rate limiting to prevent abuse.

4. **Timeout Handling**: Requests have a timeout to prevent resource exhaustion.

## Error Handling

The server implements secure error handling to prevent information leakage:

1. **Sanitized Error Messages**: Error messages are sanitized to remove sensitive information like API keys and URLs.

2. **Appropriate Error Responses**: The server returns appropriate error responses without exposing internal details.

## Dependency Management

1. **Regular Updates**: Keep dependencies updated to patch security vulnerabilities.

   ```bash
   go get -u
   ```

2. **Vulnerability Scanning**: Regularly scan dependencies for known vulnerabilities using tools like:
   - `govulncheck`
   - GitHub Dependabot
   - Snyk

## Security Monitoring and Logging

1. **Logging**: The server logs important events, but sanitizes sensitive information.

2. **Monitoring**: Consider implementing monitoring to detect unusual patterns that might indicate abuse.

## Reporting Security Issues

If you discover a security vulnerability, please report it by sending an email to [security@example.com](mailto:security@example.com). Please do not disclose security vulnerabilities publicly until they have been handled by the security team.

## Security Updates

Security updates will be released as needed. Users are encouraged to stay updated with the latest version of the software. 