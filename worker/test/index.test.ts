import { env } from "cloudflare:workers";
import { beforeEach, describe, expect, it } from "vitest";
import schema from "../migrations/20250803223004.sql?raw";
import authors from "../migrations/20250803230701.sql?raw";
import tokens from "../migrations/20250821225330.sql?raw";
import app from "../src/index";

type TestEnv = { DB: D1Database };
const db = (env as unknown as TestEnv).DB;

async function request(path: string, init?: RequestInit): Promise<Response> {
  return app.request(`https://example.com${path}`, init, { DB: db });
}

beforeEach(async () => {
  for (const table of ["comments", "conundrums", "authors", "tokens"]) {
    await db.prepare(`DROP TABLE IF EXISTS ${table}`).run();
  }
  for (const sql of [schema, authors, tokens]) {
    for (const statement of sql
      .split(";")
      .map((part) => part.trim())
      .filter(Boolean)) {
      await db.prepare(statement).run();
    }
  }
});

describe("worker", () => {
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
          headers: { cookie },
        })
      ).status,
    ).toBe(401);
    await db
      .prepare("INSERT INTO tokens (token) VALUES (?)")
      .bind("secret")
      .run();
    const admin = `risuhunnik-cookie=${encodeURIComponent(JSON.stringify({ author: "Admin", starred: [], token: "secret" }))}`;
    expect(
      (
        await request("/verify/1", {
          method: "POST",
          headers: { cookie: admin },
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
  });
});
