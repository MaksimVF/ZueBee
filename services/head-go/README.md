




Head Service (Go) — gRPC gateway for LLM.

Setup:
1. Generate protobuf Go files:
   protoc -I proto --go_out=./gen --go-grpc_out=./gen proto/*.proto

2. Build:
   docker build -t head-go .

3. Run (env):
   CHAT_ADDR=":50055" METRICS_PORT=9001 PROVIDER_KEYS='{"openai":"sk-...","local":""}' FALLBACKS='{"gpt-4o":["openai","local"]}' CACHE_ENABLED=true CACHE_ADDR="cache:50053" docker run ...

Notes:
- Configure MTLS via env MTLS_ENABLED and file mounts for certs.
- Ensure billing & cache services are reachable.




