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
    cargo build --locked --release -p query-api -p trade-log-query -p trade-collector -p trade-parser-publisher && \
    mkdir -p /out && cp target/release/query-api target/release/trade-log-query target/release/trade-collector target/release/trade-parser-publisher /out/

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

FROM runtime-base AS trade-parser-publisher
COPY --from=builder /out/trade-parser-publisher /usr/local/bin/trade-parser-publisher
EXPOSE 8083
ENTRYPOINT ["/usr/local/bin/trade-parser-publisher"]
CMD ["--config", "/etc/robotech/trade-parser-publisher.toml"]

# Optional isolated fixture image; ordinary application builds do not compile test binaries.
FROM builder AS fixture-builder
COPY config/ config/
RUN --mount=type=cache,id=robotech-cargo-registry,target=/usr/local/cargo/registry,sharing=locked \
    --mount=type=cache,id=robotech-cargo-git,target=/usr/local/cargo/git,sharing=locked \
    --mount=type=cache,id=robotech-cargo-target,target=/workspace/target,sharing=locked \
    cargo test --locked --release -p trade-log-adapters --test publishing_chain --no-run --message-format=json > /tmp/v05-fixture-build.json && \
    v05_test_executable=$(sed -n 's/.*"executable":"\([^" ]*publishing_chain-[^" ]*\)".*/\1/p' /tmp/v05-fixture-build.json | tail -n 1) && \
    test -n "$v05_test_executable" && cp "$v05_test_executable" /out/publishing-fixtures
FROM python:3.13-slim-bookworm AS verify-v05
COPY --from=fixture-builder /out/publishing-fixtures /app/publishing-fixtures
COPY --from=builder /out/trade-parser-publisher /app/trade-parser-publisher
COPY scripts/webhook-receiver.py scripts/run-v05-fixtures.sh /app/
ENV PYTHONDONTWRITEBYTECODE=1 ROBOTECH_RECEIVER_SCRIPT=/app/webhook-receiver.py ROBOTECH_TEST_DATABASE_EXACT=1 ROBOTECH_PUBLISHER_BINARY=/app/trade-parser-publisher
ENTRYPOINT ["/bin/sh", "/app/run-v05-fixtures.sh"]
