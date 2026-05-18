export class NotFoundError extends Error {
  readonly statusCode = 404;

  constructor(message = "Courier not found") {
    super(message);
  }
}

export class ValidationError extends Error {
  constructor(
    public readonly errors: Record<string, string[]>,
    public readonly statusCode = 422
  ) {
    super("Validation failed");
  }
}
