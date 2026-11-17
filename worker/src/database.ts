import { Data, Effect } from "effect";

export class DatabaseError extends Data.TaggedError("DatabaseError")<{
  readonly cause: unknown;
}> {}

export function databaseOperation<A>(
  operation: () => Promise<A>,
): Effect.Effect<A, DatabaseError> {
  return Effect.tryPromise({
    try: operation,
    catch: (cause) => new DatabaseError({ cause }),
  });
}

export function runDatabase<A>(operation: () => Promise<A>): Promise<A> {
  return Effect.runPromise(databaseOperation(operation));
}
