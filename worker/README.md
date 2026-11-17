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

The migrations define a fresh database and include the original seed data. The migration history was reset, so use a new D1 database rather than applying it to a database with the old migration history. If existing data matters, back it up and test an import into a separate D1 database, preserving IDs, comments, stars, and authors. Freeze writes during the final import and verify the data before switching traffic. No production data migration or domain cutover is performed by this project.

## Admin access

Generate a random 256-bit secret with `openssl rand -hex 32` and store it securely in a password manager. Set it as the `ADMIN_SECRET` Cloudflare Worker secret using `pnpm exec wrangler secret put ADMIN_SECRET` and enter it at the prompt. For local development, set `ADMIN_SECRET=<secret>` in `worker/.dev.vars`; this file is ignored by version control. Use a different secret for each environment.

Open `/admin` and sign in with the secret. There is no public link to this page. Admin sessions expire after 12 hours and can be revoked by signing out, deleting rows in `admin_sessions`, or rotating the Worker secret. Apply D1 migrations before deploying this version.

Login attempts are limited to five per minute per IP address in each Cloudflare location. This limit is approximate and does not replace broader abuse protection for the public endpoints. The rate-limit namespace ID in `wrangler.jsonc` must be unique within your Cloudflare account unless sharing a limit is intentional.
