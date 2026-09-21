FROM rust:1.88-bookworm AS builder
WORKDIR /workspace
COPY Cargo.toml Cargo.lock ./
COPY services/trade-log/Cargo.toml services/trade-log/Cargo.toml
COPY services/trade-log/src services/trade-log/src
COPY services/trade-log/migrations services/trade-log/migrations
RUN cargo build --locked --release --bin trade-log-server --bin trade-log-migrate

FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /workspace/target/release/trade-log-server /usr/local/bin/trade-log-server
COPY --from=builder /workspace/target/release/trade-log-migrate /usr/local/bin/trade-log-migrate
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["trade-log-server"]
