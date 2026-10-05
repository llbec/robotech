# syntax=docker/dockerfile:1
FROM rust:1.96.1-bookworm@sha256:a339861ae23e9abb272cea45dfafde21760d2ce6577a70f8a926153677902663 AS builder
WORKDIR /workspace
ENV RUSTUP_TOOLCHAIN=1.96.1
COPY Cargo.toml Cargo.lock rust-toolchain.toml ./
COPY crates/ crates/
COPY gateway/ gateway/
COPY services/ services/
COPY protocols/ protocols/
COPY tests/ tests/
RUN --mount=type=cache,id=robotech-cargo-registry,target=/usr/local/cargo/registry,sharing=locked \
    --mount=type=cache,id=robotech-cargo-git,target=/usr/local/cargo/git,sharing=locked \
    --mount=type=cache,id=robotech-cargo-target,target=/workspace/target,sharing=locked \
    cargo build --locked --release -p query-api -p trade-log-query -p trade-collector && \
    mkdir -p /out && cp target/release/query-api target/release/trade-log-query target/release/trade-collector /out/

FROM debian:bookworm-slim@sha256:3783cc01769c7b2b1b83a5c5ad96c815348e28ed7da68e2e3687004faa906251 AS runtime-base
RUN groupadd --gid 10001 robotech && useradd --uid 10001 --gid 10001 --no-create-home robotech
USER 10001:10001
STOPSIGNAL SIGTERM

FROM runtime-base AS trade-log-query
COPY --from=builder /out/trade-log-query /usr/local/bin/trade-log-query
EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/trade-log-query"]
CMD ["--config", "/etc/robotech/trade-log.toml"]

FROM runtime-base AS query-api
COPY --from=builder /out/query-api /usr/local/bin/query-api
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/query-api"]
CMD ["--config", "/etc/robotech/query-api.toml"]

FROM runtime-base AS trade-collector
COPY --from=builder /out/trade-collector /usr/local/bin/trade-collector
EXPOSE 8082
ENTRYPOINT ["/usr/local/bin/trade-collector"]
CMD ["--config", "/etc/robotech/trade-collector.toml"]
