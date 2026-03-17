## 2024-05-18 - [Fix Missing Authentication on Grafana Endpoint]
**Vulnerability:** Grafana was configured with anonymous authentication enabled (`GF_AUTH_ANONYMOUS_ENABLED=true`) providing unauthenticated "Viewer" access. This publicly exposed sensitive network monitoring data, host active state, and network topology.
**Learning:** Development-friendly settings (like anonymous access) were accidentally left enabled in a production configuration (`docker-compose.yml`), despite the `README.md` highlighting secure deployment and password protection via `.env`.
**Prevention:** Ensure any infrastructure-as-code or Docker Compose configurations strictly disable unauthenticated access unless explicitly required and verified. Do not use development "convenience" settings in default production templates.
## 2026-03-17 - [Nginx Missing Security Headers & Weak TLS]
**Vulnerability:** The Nginx reverse proxy configuration lacked basic security headers (HSTS, X-Content-Type-Options, etc.) and did not enforce modern TLS versions or strong ciphers, leaving the application vulnerable to clickjacking, MIME sniffing, and downgrade attacks.
**Learning:** Default Nginx configurations are optimized for compatibility, not security. When deploying Nginx as a reverse proxy, security hardening must be explicitly configured in the `http` or `server` blocks.
**Prevention:** Always include a standard block of security headers and TLS hardening directives in any new Nginx configuration file.
