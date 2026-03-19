## 2024-05-18 - [Fix Missing Authentication on Grafana Endpoint]
**Vulnerability:** Grafana was configured with anonymous authentication enabled (`GF_AUTH_ANONYMOUS_ENABLED=true`) providing unauthenticated "Viewer" access. This publicly exposed sensitive network monitoring data, host active state, and network topology.
**Learning:** Development-friendly settings (like anonymous access) were accidentally left enabled in a production configuration (`docker-compose.yml`), despite the `README.md` highlighting secure deployment and password protection via `.env`.
**Prevention:** Ensure any infrastructure-as-code or Docker Compose configurations strictly disable unauthenticated access unless explicitly required and verified. Do not use development "convenience" settings in default production templates.
## 2026-03-17 - [Nginx Missing Security Headers & Weak TLS]
**Vulnerability:** The Nginx reverse proxy configuration lacked basic security headers (HSTS, X-Content-Type-Options, etc.) and did not enforce modern TLS versions or strong ciphers, leaving the application vulnerable to clickjacking, MIME sniffing, and downgrade attacks.
**Learning:** Default Nginx configurations are optimized for compatibility, not security. When deploying Nginx as a reverse proxy, security hardening must be explicitly configured in the `http` or `server` blocks.
**Prevention:** Always include a standard block of security headers and TLS hardening directives in any new Nginx configuration file.
## 2024-05-19 - [Missing Input Length Limit on IP Subnet Parsing]
**Vulnerability:** The `GenerateTargets` function blindly expanded CIDR blocks provided in configuration into a slice of strings containing every valid IP address. A user could configure a block like `10.0.0.0/8`, causing the application to loop 16.7 million times and append to a slice, resulting in immediate Out of Memory (OOM) crashes and Denial of Service (DoS) of the monitoring application.
**Learning:** Functions that translate compact notations (like CIDR subnets) into expanded memory representations must enforce upper bounds to prevent malicious or accidental resource exhaustion.
**Prevention:** Always validate the size of the generated collection against a safe upper limit (e.g., maximum `/16` subnet size which yields 65,536 addresses) before entering the allocation loop.
