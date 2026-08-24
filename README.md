# pg5-rumbles-app

A tiny Go HTTP service used to demonstrate a full GitOps promotion + progressive
delivery pipeline on Amazon EKS (hub-and-spoke Argo CD, ECR, CodePipeline,
Argo Rollouts).

This is the **application repo**. Its only job is to produce an immutable
container image. Deployment configuration lives in the separate
[`pg5-rumbles-gitops`](https://github.com/devanshpoplii/pg5-rumbles-gitops) repo.

## Endpoints

| Path       | Purpose                                                        |
|------------|----------------------------------------------------------------|
| `/`        | Simple greeting with the running version                       |
| `/health`  | Liveness/readiness probe; signal source for Rollouts analysis  |
| `/version` | Reports the running version + hostname (useful during canary)  |

## Run locally

```bash
go run .
# in another terminal:
curl localhost:8080/health
curl localhost:8080/version
```

## Test

```bash
go test ./...
```

## Build the image

The version is injected at build time and baked into the binary:

```bash
docker build --build-arg VERSION=0.1.0 -t rumbles:0.1.0 .
```

In the pipeline, the image is pushed to ECR and promoted **by digest**
(`@sha256:...`) for true immutability — never by a mutable tag.
