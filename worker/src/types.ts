export type Bindings = {
  DB: D1Database;
  ADMIN_SECRET?: string;
  ADMIN_LOGIN_LIMITER: {
    limit(options: { key: string }): Promise<{ success: boolean }>;
  };
};
export type Conundrum = {
  id: number;
  text: string;
  author: string;
  verified: number;
  stars: number;
  date: string;
};
export type Comment = {
  id: number;
  cid: number;
  comment: string;
  author: string;
  date: string;
};
export type Identity = { author: string; starred: number[] };
