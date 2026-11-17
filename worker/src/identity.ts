import { runDatabase } from "./database";
import type { Identity } from "./types";

export const cookieName = "risuhunnik-cookie";

export function readCookie(header: string | undefined): string | undefined {
  const entry = header
    ?.split(";")
    .map((part) => part.trim())
    .find((part) => part.startsWith(`${cookieName}=`));
  if (!entry) return undefined;

  try {
    return decodeURIComponent(
      entry.slice(cookieName.length + 1).replace(/\+/g, " "),
    );
  } catch {
    return undefined;
  }
}

export async function identity(
  db: D1Database,
  encoded: string | undefined,
): Promise<{ value: Identity; fresh: boolean }> {
  if (encoded) {
    try {
      const parsed: unknown = JSON.parse(encoded);
      if (
        typeof parsed === "object" &&
        parsed !== null &&
        "author" in parsed &&
        typeof parsed.author === "string" &&
        "starred" in parsed &&
        (parsed.starred === null || Array.isArray(parsed.starred)) &&
        (!("token" in parsed) || typeof parsed.token === "string")
      ) {
        const starred: unknown = parsed.starred;

        if (
          starred === null ||
          (Array.isArray(starred) &&
            starred.every(
              (id: unknown) =>
                typeof id === "number" && Number.isSafeInteger(id),
            ))
        ) {
          return {
            value: {
              author: parsed.author,
              starred: starred ?? [],
              token:
                "token" in parsed && typeof parsed.token === "string"
                  ? parsed.token
                  : "",
            },
            fresh: false,
          };
        }
      }
    } catch {
      /* Invalid cookies are replaced with a new identity. */
    }
  }

  const row = await runDatabase(() =>
    db
      .prepare("SELECT author FROM authors ORDER BY RANDOM() LIMIT 1")
      .first<{ author: string }>(),
  );
  if (!row) throw new Error("No authors available");

  return { value: { author: row.author, starred: [], token: "" }, fresh: true };
}

export async function validToken(
  db: D1Database,
  token: string,
): Promise<boolean> {
  if (!token) return false;

  return (
    (await runDatabase(() =>
      db
        .prepare("SELECT 1 AS valid FROM tokens WHERE token = ?")
        .bind(token)
        .first(),
    )) !== null
  );
}
