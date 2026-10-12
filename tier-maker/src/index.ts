import { DurableObject } from "cloudflare:workers";
import { Hono } from "hono";

const PORT = 8080;
const INACTIVITY_TIMEOUT_MS = 2 * 60 * 1000;
const POOL_SIZE = 3;

export class MyContainer extends DurableObject<Env> {
  private starting: Promise<void> | undefined;
  private monitoring: Promise<void> | undefined;

  constructor(ctx: DurableObjectState, env: Env) {
    super(ctx, env);

    const container = ctx.container;

    if (container?.running) {
      void ctx.blockConcurrencyWhile(() =>
        container.setInactivityTimeout(INACTIVITY_TIMEOUT_MS),
      );

      this.observeExit();
    }
  }

  async fetch(request: Request): Promise<Response> {
    this.starting ??= this.startAndWaitForPort().finally(() => {
      this.starting = undefined;
    });

    await this.starting;

    const url = new URL(request.url);
    url.protocol = "http:";
    url.host = "container";

    const forwarded = new Request(url, request);
    forwarded.headers.delete("host");

    return this.ctx.container!.getTcpPort(PORT).fetch(forwarded);
  }

  private async startAndWaitForPort(): Promise<void> {
    const container = this.ctx.container!;

    if (!container.running) {
      container.start({
        image: container.images.base,
        instance: "lite",
        enableInternet: false,
        env: {
          MESSAGE: "I was passed in when the Durable Object started me!",
          INSTANCE_ID: this.ctx.id.toString(),
        },
      });
    }

    this.observeExit();

    await container.setInactivityTimeout(INACTIVITY_TIMEOUT_MS);

    const port = container.getTcpPort(PORT);

    // Wait until the Go Fiber server is ready.
    for (let attempt = 0; attempt < 100; attempt++) {
      try {
        const response = await port.fetch("http://container/health", {
          signal: AbortSignal.timeout(1000),
        });

        await response.body?.cancel();

        if (response.ok) {
          return;
        }
      } catch {
        // The server may still be starting.
      }

      await scheduler.wait(200);
    }

    throw new Error("Go backend did not become ready on port 8080");
  }

  private observeExit(): void {
    if (this.monitoring) {
      return;
    }

    this.monitoring = this.ctx
      .container!.monitor()
      .then(() => {
        console.log("Container exited successfully");
      })
      .catch((error: unknown) => {
        console.error("Container failed:", error);
      })
      .finally(() => {
        this.monitoring = undefined;
      });

    this.ctx.waitUntil(this.monitoring);
  }
}

const app = new Hono<{ Bindings: Env }>();

// Home page.
app.get("/", (c) =>
  c.text(
    [
      "Available endpoints:",
      "GET  /health",
      "GET  /tier-list",
      "POST /tier-maker",
    ].join("\n"),
  ),
);

// Worker health check.
app.get("/health", (c) => c.json({ status: "ok" }));

// Forward backend API routes, including /tier-list and /tier-maker,
// to the same container.
app.all("*", (c) =>
  c.env.MY_CONTAINER.getByName("tier-maker").fetch(c.req.raw),
);

// Error handler.
app.onError((error, c) => {
  console.error("Container request failed:", error);

  return c.json(
    {
      error: "Container request failed",
      message: "Check the Cloudflare Worker logs and Go backend.",
    },
    502,
  );
});

export default app;