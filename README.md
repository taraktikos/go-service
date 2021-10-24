# Go service template

## Run docker compose on local env


```bash
docker compose -f zarf/docker-compose.dev.yml up --build
#or
make compose
```

## Run in local kubernetes cluster

```bash
make docker-build
make kind-up
make kind-load
make kind-apply
```

## Provision aws env

```bash
export TF_VAR_ghcr_token=$GHCR_TOKEN

terraform apply
```

```bash
chmod 400 .ssh/key.pem 
ssh -i .ssh/key.pem ec2-user@IP
```

