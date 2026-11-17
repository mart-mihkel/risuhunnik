import { Hono } from "hono";
import { setCookie } from "hono/cookie";
import { html } from "hono/html";
import admin, { adminCookie, isAdmin, isSameOrigin } from "./admin";
import { runDatabase } from "./database";
import { cookieName, identity, readCookie } from "./identity";
import type { Bindings, Comment, Conundrum } from "./types";
import {
  commentResult,
  conundrums,
  detail,
  list,
  page,
  stars,
  uploadResult,
  verified,
} from "./views";

const app = new Hono<{ Bindings: Bindings }>();

app.route("/", admin);

function parseId(value: string): number | null {
  if (!/^[1-9]\d*$/.test(value)) return null;

  const id = Number(value);

  return Number.isSafeInteger(id) ? id : null;
}

app.get("/", async (c) => {
  const rows = await runDatabase(() =>
    c.env.DB.prepare(
      "SELECT * FROM conundrums WHERE verified = 1",
    ).all<Conundrum>(),
  );

  return c.html(page(conundrums(rows.results)));
});

app.get("/conundrums", async (c) => {
  const rows = await runDatabase(() =>
    c.env.DB.prepare("SELECT * FROM conundrums").all<Conundrum>(),
  );

  return c.html(conundrums(rows.results));
});

app.get("/author/:author", async (c) => {
  const author = c.req.param("author");
  const items = await runDatabase(() =>
    c.env.DB.prepare("SELECT * FROM conundrums WHERE author = ?")
      .bind(author)
      .all<Conundrum>(),
  );
  const comments = await runDatabase(() =>
    c.env.DB.prepare("SELECT COUNT(*) AS count FROM comments WHERE author = ?")
      .bind(author)
      .first<{ count: number }>(),
  );
  const total = items.results.reduce((sum, item) => sum + item.stars, 0);
  return c.html(
    html`<p>author......: ${author}</p>
      <p>stars.......: ${total}</p>
      <p>comments....: ${comments?.count ?? 0}</p>
      <p>conundrums..:</p>
      <br />${list(items.results)}`,
  );
});

app.get("/conundrum/:id", async (c) => {
  const id = parseId(c.req.param("id"));
  if (id === null) return c.text("Invalid id", 400);

  const item = await runDatabase(() =>
    c.env.DB.prepare("SELECT * FROM conundrums WHERE id = ?")
      .bind(id)
      .first<Conundrum>(),
  );
  if (!item) return c.notFound();

  const comments = await runDatabase(() =>
    c.env.DB.prepare("SELECT * FROM comments WHERE cid = ?")
      .bind(id)
      .all<Comment>(),
  );
  const user = await identity(c.env.DB, readCookie(c.req.header("cookie")));
  if (user.fresh)
    setCookie(c, cookieName, JSON.stringify(user.value), {
      path: "/",
      httpOnly: true,
      sameSite: "Lax",
    });

  c.header("Cache-Control", "no-store");
  const allowed = await isAdmin(
    c.env.DB,
    adminCookie(c.req.header("cookie")),
    c.env.ADMIN_SECRET,
  );
  const prev = await runDatabase(() =>
    c.env.DB.prepare("SELECT MAX(id) AS id FROM conundrums WHERE id < ?")
      .bind(id)
      .first<{ id: number | null }>(),
  );
  const next = await runDatabase(() =>
    c.env.DB.prepare("SELECT MIN(id) AS id FROM conundrums WHERE id > ?")
      .bind(id)
      .first<{ id: number | null }>(),
  );

  return c.html(
    detail(
      item,
      comments.results,
      allowed,
      user.value.starred.includes(id),
      prev?.id ?? 0,
      next?.id ?? 0,
    ),
  );
});

app.post("/upload", async (c) => {
  const form = await c.req.parseBody();
  if (typeof form.conundrum !== "string" || !form.conundrum.trim())
    return c.text("Invalid conundrum", 400);

  const user = await identity(c.env.DB, readCookie(c.req.header("cookie")));
  await runDatabase(() =>
    c.env.DB.prepare("INSERT INTO conundrums (text, author) VALUES (?, ?)")
      .bind(form.conundrum, user.value.author)
      .run(),
  );
  if (user.fresh)
    setCookie(c, cookieName, JSON.stringify(user.value), {
      path: "/",
      httpOnly: true,
      sameSite: "Lax",
    });

  return c.html(uploadResult(true));
});

app.post("/comment", async (c) => {
  const form = await c.req.parseBody();
  const id =
    typeof form["conundrum-id"] === "string"
      ? parseId(form["conundrum-id"])
      : null;
  if (id === null || typeof form.comment !== "string" || !form.comment.trim())
    return c.text("Invalid comment", 400);

  const user = await identity(c.env.DB, readCookie(c.req.header("cookie")));
  await runDatabase(() =>
    c.env.DB.prepare(
      "INSERT INTO comments (cid, comment, author) VALUES (?, ?, ?)",
    )
      .bind(id, form.comment, user.value.author)
      .run(),
  );
  if (user.fresh)
    setCookie(c, cookieName, JSON.stringify(user.value), {
      path: "/",
      httpOnly: true,
      sameSite: "Lax",
    });

  return c.html(commentResult(id));
});

app.post("/star/:id", async (c) => {
  const id = parseId(c.req.param("id"));
  if (id === null) return c.text("Invalid id", 400);

  const user = await identity(c.env.DB, readCookie(c.req.header("cookie")));
  const starred = user.value.starred.includes(id);
  const item = await runDatabase(() =>
    c.env.DB.prepare(
      "UPDATE conundrums SET stars = stars + ? WHERE id = ? RETURNING *",
    )
      .bind(starred ? -1 : 1, id)
      .first<Conundrum>(),
  );
  if (!item) return c.notFound();

  user.value.starred = starred
    ? user.value.starred.filter((value) => value !== id)
    : [...user.value.starred, id];
  setCookie(c, cookieName, JSON.stringify(user.value), {
    path: "/",
    httpOnly: true,
    sameSite: "Lax",
  });

  return c.html(stars(item, !starred));
});

app.post("/verify/:id", async (c) => {
  const id = parseId(c.req.param("id"));
  if (id === null) return c.text("Invalid id", 400);

  if (!isSameOrigin(c.req.raw)) return c.text("Forbidden", 403);

  if (
    !(await isAdmin(
      c.env.DB,
      adminCookie(c.req.header("cookie")),
      c.env.ADMIN_SECRET,
    ))
  )
    return c.text("Unauthorized", 401);

  const item = await runDatabase(() =>
    c.env.DB.prepare(
      "UPDATE conundrums SET verified = 1 - verified WHERE id = ? RETURNING *",
    )
      .bind(id)
      .first<Conundrum>(),
  );
  if (!item) return c.notFound();

  return c.html(verified(item, true));
});

export default app;
