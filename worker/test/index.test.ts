import { env } from "cloudflare:workers";
import { beforeEach, describe, expect, it } from "vitest";
import schema from "../migrations/001-schema.sql?raw";
import authors from "../migrations/002-authors.sql?raw";
import app from "../src/index";

type TestEnv = { DB: D1Database };
const db = (env as unknown as TestEnv).DB;
const secret = "0123456789abcdef".repeat(4);
let configuredSecret: string | undefined = secret;
let allowLogin = true;

async function request(path: string, init?: RequestInit): Promise<Response> {
  return app.request(`https://example.com${path}`, init, {
    DB: db,
    ADMIN_SECRET: configuredSecret,
    ADMIN_LOGIN_LIMITER: { limit: async () => ({ success: allowLogin }) },
  });
}

beforeEach(async () => {
  allowLogin = true;
  configuredSecret = secret;
  for (const table of ["comments", "conundrums", "authors", "admin_sessions"]) {
    await db.prepare(`DROP TABLE IF EXISTS ${table}`).run();
  }
  for (const sql of [schema, authors]) {
    for (const statement of sql
      .split(";")
      .map((part) => part.trim())
      .filter(Boolean)) {
      await db.prepare(statement).run();
    }
  }
});

describe("worker", () => {
  it("creates only the application tables and a complete session schema", async () => {
    const tables = await db
      .prepare(
        "SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name",
      )
      .all<{ name: string }>();
    expect(tables.results.map((table) => table.name)).toEqual([
      "admin_sessions",
      "authors",
      "comments",
      "conundrums",
    ]);
    const columns = await db
      .prepare("PRAGMA table_info(admin_sessions)")
      .all<{ name: string; notnull: number }>();
    expect(
      columns.results.find((column) => column.name === "credential_hash")
        ?.notnull,
    ).toBe(1);
  });

  it("serves verified conundrums and escapes user content", async () => {
    await db
      .prepare(
        "INSERT INTO conundrums (text, author, verified) VALUES (?, ?, ?)",
      )
      .bind("<script>alert(1)</script>", "test", 1)
      .run();
    await db
      .prepare("INSERT INTO conundrums (text, author) VALUES (?, ?)")
      .bind("hidden", "test")
      .run();
    const response = await request("/");
    const body = await response.text();
    expect(body).toContain("&lt;script&gt;");
    expect(body).not.toContain("hidden");
    expect(await (await request("/conundrums")).text()).toContain("hidden");
  });

  it("posts, stars, comments and requires a valid token to verify", async () => {
    const upload = await request("/upload", {
      method: "POST",
      body: new URLSearchParams({ conundrum: "A puzzle" }),
    });
    expect(upload.status).toBe(200);
    const cookie = upload.headers.get("set-cookie")?.split(";")[0];
    if (!cookie) throw new Error("Missing cookie");
    const star = await request("/star/1", {
      method: "POST",
      headers: { cookie },
    });
    expect(await star.text()).toContain("star-");
    const starredCookie = star.headers.get("set-cookie")?.split(";")[0];
    if (!starredCookie) throw new Error("Missing star cookie");
    const unstar = await request("/star/1", {
      method: "POST",
      headers: { cookie: starredCookie },
    });
    expect(await unstar.text()).toContain("star+");
    expect(
      (
        await request("/comment", {
          method: "POST",
          headers: { cookie },
          body: new URLSearchParams({ comment: "hello", "conundrum-id": "1" }),
        })
      ).status,
    ).toBe(200);
    expect(
      await (await request("/conundrum/1", { headers: { cookie } })).text(),
    ).toContain("hello");
    expect(
      (
        await request("/verify/1", {
          method: "POST",
          headers: { cookie, origin: "https://example.com" },
        })
      ).status,
    ).toBe(401);
    const login = await request("/admin/login", {
      method: "POST",
      headers: { origin: "https://example.com" },
      body: new URLSearchParams({ token: secret }),
    });
    expect(login.status).toBe(303);
    const setCookie = login.headers.get("set-cookie");
    expect(setCookie).toContain("HttpOnly");
    expect(setCookie).toContain("Secure");
    const admin = setCookie?.split(";")[0];
    if (!admin) throw new Error("Missing admin session");
    const stored = await db
      .prepare("SELECT token_hash FROM admin_sessions")
      .first<{ token_hash: string }>();
    expect(stored?.token_hash).not.toBe(admin.split("=")[1]);
    expect(
      await (await request("/admin", { headers: { cookie: admin } })).text(),
    ).toContain("Signed in.");
    expect(
      (
        await request("/verify/1", {
          method: "POST",
          headers: { cookie: admin, origin: "https://evil.example" },
        })
      ).status,
    ).toBe(403);
    expect(
      (
        await request("/verify/1", {
          method: "POST",
          headers: { cookie: admin, origin: "https://example.com" },
        })
      ).status,
    ).toBe(200);
    expect(
      (
        await db
          .prepare("SELECT verified FROM conundrums WHERE id = 1")
          .first<{ verified: number }>()
      )?.verified,
    ).toBe(1);
    expect(
      (
        await request("/admin/logout", {
          method: "POST",
          headers: { cookie: admin, origin: "https://example.com" },
        })
      ).status,
    ).toBe(303);
    expect(
      (
        await request("/verify/1", {
          method: "POST",
          headers: { cookie: admin, origin: "https://example.com" },
        })
      ).status,
    ).toBe(401);
  });

  it("rejects cross-origin requests, incorrect credentials, legacy cookies and throttled logins", async () => {
    const legacy = `risuhunnik-cookie=${encodeURIComponent(JSON.stringify({ author: "Admin", starred: [], token: "secret" }))}`;
    expect(
      (
        await request("/verify/1", {
          method: "POST",
          headers: { cookie: legacy, origin: "https://example.com" },
        })
      ).status,
    ).toBe(401);
    expect(
      (
        await request("/admin/login", {
          method: "POST",
          headers: { origin: "https://evil.example" },
          body: new URLSearchParams({ token: secret }),
        })
      ).status,
    ).toBe(403);
    expect(
      (
        await request("/admin/login", {
          method: "POST",
          headers: { origin: "https://example.com" },
          body: new URLSearchParams({ token: "f".repeat(64) }),
        })
      ).status,
    ).toBe(401);
    allowLogin = false;
    expect(
      (
        await request("/admin/login", {
          method: "POST",
          headers: { origin: "https://example.com" },
          body: new URLSearchParams({ token: secret }),
        })
      ).status,
    ).toBe(429);
  });

  it("fails closed without a configured secret and invalidates sessions after rotation", async () => {
    configuredSecret = undefined;
    const loginRequest = (): RequestInit => ({
      method: "POST",
      headers: { origin: "https://example.com" },
      body: new URLSearchParams({ token: secret }),
    });
    expect((await request("/admin/login", loginRequest())).status).toBe(401);
    configuredSecret = secret;
    const login = await request("/admin/login", loginRequest());
    expect(login.status).toBe(303);
    const cookie = login.headers.get("set-cookie")?.split(";")[0];
    if (!cookie) throw new Error("Missing admin session");
    configuredSecret = "a".repeat(64);
    expect(
      (
        await request("/verify/1", {
          method: "POST",
          headers: { cookie, origin: "https://example.com" },
        })
      ).status,
    ).toBe(401);
  });

  it("rejects expired sessions", async () => {
    const login = await request("/admin/login", {
      method: "POST",
      headers: { origin: "https://example.com" },
      body: new URLSearchParams({ token: secret }),
    });
    const cookie = login.headers.get("set-cookie")?.split(";")[0];
    if (!cookie) throw new Error("Missing admin session");
    await db.prepare("UPDATE admin_sessions SET expires_at = 0").run();
    expect(
      (
        await request("/verify/1", {
          method: "POST",
          headers: { cookie, origin: "https://example.com" },
        })
      ).status,
    ).toBe(401);
  });
});
