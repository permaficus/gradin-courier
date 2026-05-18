export class NotFoundError extends Error {
  constructor(message = "Courier not found") {
    super(message);
  }
}

export class ValidationError extends Error {
  constructor(public readonly errors: Record<string, string[]>) {
    super("Validation failed");
  }
}
