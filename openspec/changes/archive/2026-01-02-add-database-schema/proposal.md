# Change: Add Database Schema

## Why
A well-designed PostgreSQL schema is the foundation for data integrity. The schema enforces constraints, enables efficient queries, and supports the domain model's invariants at the database level.

## What Changes
- Add database migration for funds, cap_table_entries, and transfers tables
- Add constraints for data integrity
- Add indexes for query performance
- Add connection pooling and migration infrastructure

## Impact
- Affected specs: database-schema (new)
- Affected code: `/migrations/`, `internal/postgres/`
