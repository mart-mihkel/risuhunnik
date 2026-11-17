export type Bindings = { DB: D1Database };
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
export type Identity = { author: string; starred: number[]; token: string };
