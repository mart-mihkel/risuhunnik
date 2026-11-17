# Show available commands.
help:
    @just --list

# Start the Worker locally.
dev: migrate
    pnpm --prefix worker dev

# Apply D1 migrations locally.
migrate:
    pnpm --prefix worker exec wrangler d1 migrations apply risuhunnik --local

# Check types, formatting, and lint.
check:
    pnpm --prefix worker check

# Format Worker files.
format:
    pnpm --prefix worker format

# Run Worker tests.
test:
    pnpm --prefix worker test

# Deploy the Worker to Cloudflare.
deploy:
    pnpm --prefix worker deploy
