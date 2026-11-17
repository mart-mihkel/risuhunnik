import { html } from "hono/html";
import type { Comment, Conundrum } from "./types";

const indicator = (name: string) =>
  html`<span class="typewriter-wrapper"
    ><span id="${name}" class="htmx-indicator typewriter">....</span></span
  >`;

export const list = (items: Conundrum[]) =>
  html`<ol>
    ${items.map((item) => html`<li class="${item.verified ? "" : "not-verified"}">${item.text} <button hx-get="/conundrum/${item.id}" hx-target="main">-&gt;</button></li>`)}
  </ol>`;

export const uploadResult = (refresh: boolean) =>
  html`<span id="upload-form-result"
    >${refresh ? html`<span hx-get="/conundrums" hx-trigger="load" hx-target="main"></span>` : ""}</span
  >`;

export const commentResult = (id: number | null) =>
  html`<span id="comment-form-result"
    >${id !== null ? html`<span hx-get="/conundrum/${id}" hx-trigger="load" hx-target="main"></span>` : ""}</span
  >`;

export const conundrums = (items: Conundrum[]) =>
  html`${list(items)}<br />
    <form hx-post="/upload" hx-target="#upload-form-result">
      <span>post conundrum:</span>
      <input
        required
        type="text"
        name="conundrum"
        placeholder="conundrum..."
      />${indicator("upload-form-loading")}<button type="submit">post+</button>
    </form>
    ${uploadResult(false)}`;

const authorLink = (author: string) =>
  html`<p>
    author......:
    <button
      hx-get="/author/${encodeURIComponent(author)}"
      hx-target="main"
      hx-indicator="#conundrum-nav-loading"
    >
      ${author}
    </button>
  </p>`;

export const stars = (item: Conundrum, starred: boolean) =>
  html`<p id="conundrum-stars">
    stars.......: ${item.stars}
    <button
      hx-post="/star/${item.id}"
      hx-swap="outerHTML"
      hx-target="#conundrum-stars"
      hx-indicator="#conundrum-star-loading"
    >
      star${starred ? "-" : "+"}</button
    >${indicator("conundrum-star-loading")}
  </p>`;

export const verified = (item: Conundrum, allowed: boolean) =>
  html`<p id="conundrum-verified">
    goated......:
    <button
      ${allowed ? "" : html`disabled`}
      hx-post="/verify/${item.id}"
      hx-swap="outerHTML"
      hx-target="#conundrum-verified"
      hx-indicator="#conundrum-verified-loading"
    >
      ${Boolean(item.verified)}</button
    >${indicator("conundrum-verified-loading")}
  </p>`;

export const detail = (
  item: Conundrum,
  comments: Comment[],
  allowed: boolean,
  starred: boolean,
  prev: number,
  next: number,
) =>
  html` <div class="conundrum-nav">
      <button
        ${prev ? "" : html`disabled`}
        hx-get="/conundrum/${prev}"
        hx-target="main"
        hx-indicator="#conundrum-nav-loading"
      >
        &lt;-prev</button
      >${indicator("conundrum-nav-loading")}<button
        ${next ? "" : html`disabled`}
        hx-get="/conundrum/${next}"
        hx-target="main"
        hx-indicator="#conundrum-nav-loading"
      >
        next-&gt;
      </button>
    </div>
    <p>conundrum...: ${item.id}</p>
    <p>date........: <span class="date">${item.date}</span></p>
    ${authorLink(item.author)}${stars(item, starred)}${verified(item, allowed)}
    <p class="conundrum-text">
      conundrum...:
      <span class="${item.verified ? "" : "not-verified"}">${item.text}</span>
    </p>
    ${comments.map(
      (comment) =>
        html`<br />
          <p>comment.....: ${comment.id}</p>
          <p>date........: <span class="date">${comment.date}</span></p>
          ${authorLink(comment.author)}
          <p class="comment-text">comment.....: ${comment.comment}</p>`,
    )}
    <br />
    <form hx-post="/comment" hx-target="#comment-form-result">
      <span>post comment:</span>
      <input
        required
        type="text"
        name="comment"
        placeholder="comment..."
      /><input
        type="hidden"
        name="conundrum-id"
        value="${item.id}"
      />${indicator("comment-form-loading")}<button type="submit">post+</button>
    </form>
    ${commentResult(null)}`;

export const page = (content: ReturnType<typeof html>) =>
  html`<!doctype html>
    <html lang="en">
      <head>
        <title>risuhunnik.xyz</title>
        <meta charset="UTF-8" />
        <meta name="description" content="A collection of wicked conundrums." />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link href="/static/style.css" rel="stylesheet" />
        <link href="/static/favicon.ico" rel="icon" type="image/x-icon" />
      </head>
      <body>
        <pre>
         ,        <a href="/">risuhunnik.xyz-&gt;</a>
      . ~ ~ .
    .'@ )) .. '..~  <a href="https://github.com/mart-mihkel/risuhunnik" target="_blank" rel="noopener noreferrer">github/risuhunnik-&gt;</a>
    &gt; - ~\`\`  '~{
                     <button hx-get="/conundrums" hx-target="main">all-conundrums-&gt;</button>
</pre>
        <main>${content}</main>
        <script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.6/dist/htmx.min.js"></script>
      </body>
    </html>`;
