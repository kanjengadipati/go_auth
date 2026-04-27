# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-04-27

### Added
- **Health Checks**: Implemented liveness (`/health/live`) and readiness (`/health/ready`) probes for orchestration and monitoring.
- **Integration Tests**: Added a full HTTP integration test covering the `/auth/login` flow using `httptest` and `sqlmock`.
- **Project Rebranding**: Officially rebranded the project from "GoKit Auth" to **Pleco**.

### Changed
- **API Routes**: Renamed `POST /auth/admin/audit-logs/investigate` to `POST /auth/admin/audit-logs/investigations` to match plural naming conventions.
- **Project Structure**: Updated documentation and configuration to reflect the Pleco rebranding.
- **Postman Collection**: Updated investigation endpoints and rebranded environment variables.

### Fixed
- **Code Formatting**: Fixed formatting issues across multiple files to satisfy CI checks.
- **API Documentation**: Updated OpenAPI specifications and Postman collection with corrected routes.

### Removed
- **Unused Dependencies**: Pruned heavy and irrelevant dependencies:
  - Removed `go.mongodb.org/mongo-driver/v2` and `github.com/quic-go/quic-go` (downgraded Gin to v1.10.0).
  - Removed `gorm.io/driver/sqlite` and replaced it with `sqlmock` for testing.
