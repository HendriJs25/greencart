FROM oven/bun:1 AS builder

WORKDIR /app

COPY package.json bun.lockb ./

RUN bun install --frozen-lockfile --ignore-scripts

COPY . .

RUN bun --bun run build


FROM oven/bun:1 AS runtime

WORKDIR /app

ENV NODE_ENV=production
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000

COPY --from=builder /app/.output ./.output

EXPOSE 3000

CMD ["bun", "--bun", "run", ".output/server/index.mjs"]