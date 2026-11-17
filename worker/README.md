# Cloudflare Worker

This app runs on Cloudflare Workers and D1. The previous deployment is not changed automatically.

## Local development

Requires Node.js 24 and pnpm 12.16.1.

```sh
cd worker
pnpm install --frozen-lockfile
pnpm exec wrangler d1 migrations apply risuhunnik --local
pnpm dev
```

Run `pnpm check` for TypeScript and Biome checks, `pnpm format` to format, and `pnpm test` for Worker integration tests.

## Deploy a new empty site

1. Authenticate with `pnpm exec wrangler login` and create a D1 database with `pnpm exec wrangler d1 create risuhunnik`.
2. Replace the placeholder `database_id` in `wrangler.jsonc` with the ID returned by Cloudflare. Choose a separate database for staging and production.
3. Run `pnpm exec wrangler d1 migrations apply risuhunnik --remote`, then `pnpm deploy`.
4. Test the `workers.dev` URL before connecting a custom domain.

The migrations include the original seed data. Do not apply them blindly to a database containing existing data. For an existing production site, first back up the SQLite database and plan and test an import into a separate D1 database, preserving IDs, comments, stars, authors and tokens. Freeze writes during the final import and verify the data before switching traffic. No production data migration or domain cutover is performed by this project.

The Worker currently does not implement the previous server's rate limiter. Configure appropriate Cloudflare protection before exposing a production deployment.
