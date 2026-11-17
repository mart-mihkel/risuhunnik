import { Hono } from "hono";
import { deleteCookie, getCookie, setCookie } from "hono/cookie";
import { html } from "hono/html";
import { runDatabase } from "./database";
import type { Bindings } from "./types";
import { page } from "./views";

const sessionCookie = "__Host-risuhunnik-admin";
const sessionLifetime = 12 * 60 * 60;
const admin = new Hono<{ Bindings: Bindings }>();

function now(): number {
  return Math.floor(Date.now() / 1000);
}

async function hashToken(token: string): Promise<string> {
  const digest = await crypto.subtle.digest(
    "SHA-256",
    new TextEncoder().encode(token),
  );
  return Array.from(new Uint8Array(digest), (value) =>
    value.toString(16).padStart(2, "0"),
  ).join("");
}

export function isSameOrigin(request: Request): boolean {
  const origin = request.headers.get("Origin");
  return origin !== null && origin === new URL(request.url).origin;
}

export async function isAdmin(
  db: D1Database,
  token: string | undefined,
  secret: string | undefined,
): Promise<boolean> {
  if (
    !token ||
    !/^[0-9a-f]{64}$/.test(token) ||
    !secret ||
    !/^[0-9a-f]{64}$/.test(secret)
  )
    return false;

  const hash = await hashToken(token);
  const credentialHash = await hashToken(secret);
  const session = await runDatabase(() =>
    db
      .prepare(
        "SELECT 1 FROM admin_sessions WHERE token_hash = ? AND expires_at > ? AND credential_hash = ?",
      )
      .bind(hash, now(), credentialHash)
      .first(),
  );
  return session !== null;
}

export function adminCookie(header: string | undefined): string | undefined {
  return header
    ?.split(";")
    .map((item) => item.trim())
    .find((item) => item.startsWith(`${sessionCookie}=`))
    ?.slice(sessionCookie.length + 1);
}

admin.get("/admin", async (c) => {
  c.header("Cache-Control", "no-store");
  const loggedIn = await isAdmin(
    c.env.DB,
    getCookie(c, sessionCookie),
    c.env.ADMIN_SECRET,
  );
  return c.html(
    page(
      loggedIn
        ? html`<h1>Admin</h1><p>Signed in.</p><form method="post" action="/admin/logout"><button type="submit">Sign out</button></form>`
        : html`<h1>Admin sign in</h1><form method="post" action="/admin/login"><label>Admin secret <input type="password" name="token" autocomplete="current-password" required /></label><button type="submit">Sign in</button></form>`,
    ),
  );
});

admin.post("/admin/login", async (c) => {
  c.header("Cache-Control", "no-store");
  if (!isSameOrigin(c.req.raw)) return c.text("Forbidden", 403);

  const limit = await c.env.ADMIN_LOGIN_LIMITER.limit({
    key: c.req.header("CF-Connecting-IP") ?? "unknown",
  });
  if (!limit.success) return c.text("Too many attempts", 429);

  const form = await c.req.parseBody();
  if (typeof form.token !== "string" || !form.token || form.token.length > 512)
    return c.text("Invalid credentials", 401);

  const secret = c.env.ADMIN_SECRET;
  if (
    !secret ||
    !/^[0-9a-f]{64}$/.test(secret) ||
    !/^[0-9a-f]{64}$/.test(form.token)
  )
    return c.text("Invalid credentials", 401);

  const encoder = new TextEncoder();
  const algorithm = { name: "HMAC", hash: "SHA-256" };
  const expectedKey = await crypto.subtle.importKey(
    "raw",
    encoder.encode(secret),
    algorithm,
    false,
    ["verify"],
  );
  const candidateKey = await crypto.subtle.importKey(
    "raw",
    encoder.encode(form.token),
    algorithm,
    false,
    ["sign"],
  );
  const message = encoder.encode("risuhunnik-admin-login");
  const signature = await crypto.subtle.sign("HMAC", candidateKey, message);
  if (!(await crypto.subtle.verify("HMAC", expectedKey, signature, message)))
    return c.text("Invalid credentials", 401);

  const bytes = crypto.getRandomValues(new Uint8Array(32));
  const token = Array.from(bytes, (value) =>
    value.toString(16).padStart(2, "0"),
  ).join("");
  const hash = await hashToken(token);
  const credentialHash = await hashToken(secret);
  await runDatabase(() =>
    c.env.DB.prepare("DELETE FROM admin_sessions WHERE expires_at <= ?")
      .bind(now())
      .run(),
  );
  await runDatabase(() =>
    c.env.DB.prepare(
      "INSERT INTO admin_sessions (token_hash, expires_at, credential_hash) VALUES (?, ?, ?)",
    )
      .bind(hash, now() + sessionLifetime, credentialHash)
      .run(),
  );
  setCookie(c, sessionCookie, token, {
    path: "/",
    secure: true,
    httpOnly: true,
    sameSite: "Strict",
    maxAge: sessionLifetime,
  });
  return c.redirect("/admin", 303);
});

admin.post("/admin/logout", async (c) => {
  c.header("Cache-Control", "no-store");
  if (!isSameOrigin(c.req.raw)) return c.text("Forbidden", 403);

  const token = getCookie(c, sessionCookie);
  if (token && /^[0-9a-f]{64}$/.test(token)) {
    const hash = await hashToken(token);
    await runDatabase(() =>
      c.env.DB.prepare("DELETE FROM admin_sessions WHERE token_hash = ?")
        .bind(hash)
        .run(),
    );
  }
  deleteCookie(c, sessionCookie, {
    path: "/",
    secure: true,
    httpOnly: true,
    sameSite: "Strict",
  });
  return c.redirect("/admin", 303);
});

export default admin;
