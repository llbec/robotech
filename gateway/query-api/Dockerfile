FROM rust:1.96.1-bookworm@sha256:a339861ae23e9abb272cea45dfafde21760d2ce6577a70f8a926153677902663 AS builder
WORKDIR /workspace
COPY Cargo.toml Cargo.lock rust-toolchain.toml ./
COPY gateway/query-api/ gateway/query-api/
COPY tests/ tests/
RUN cargo build --locked --release -p query-api

FROM debian:bookworm-slim@sha256:3783cc01769c7b2b1b83a5c5ad96c815348e28ed7da68e2e3687004faa906251 AS runtime
RUN groupadd --gid 10001 robotech && useradd --uid 10001 --gid 10001 --no-create-home robotech
COPY --from=builder /workspace/target/release/query-api /usr/local/bin/query-api
USER 10001:10001
EXPOSE 8080
STOPSIGNAL SIGTERM
ENTRYPOINT ["/usr/local/bin/query-api"]
CMD ["--config", "/etc/robotech/query-api.toml"]
