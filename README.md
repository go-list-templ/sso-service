# SSO service

Template sso service on go

---

## Create secret and login in docker registry

Create namespace

```bash
kubectl create namespace sso-service
```

Create secret

```bash
kubectl create secret docker-registry registrysecret \
  --docker-server=ghcr.io \
  --docker-username=GH_USERNAME \
  --docker-password=GH_TOKEN \
  -n sso-service
```

Login

```bash
werf cr login ghcr.io -u GH_USERNAME -p GH_TOKEN
```

---

## Install dependency helm

```bash
werf helm dependency update .helm
```

---

## Generate secret

Generate private key

```bash
make generate-private-key
```

Create secret

```bash
kubectl create secret generic token-secret \
  --from-file=jwt_private_key=./secrets/jwt_private_key \
  -n sso-service
```

---

## Run and build App

Run and deploy to from Helm to Kuber

```bash
werf converge --repo=ghcr.io/go-list-templ/sso-service --platform=linux/amd64
```

Stop and remove release in kuber

```bash
werf dismiss
```

Forward port on localhost from app

```bash
werf kubectl port-forward svc/sso-service 8080:8080 -n sso-service
werf kubectl port-forward svc/sso-service 8081:8081 -n sso-service
```

Get events in kuber

```bash
kubectl get events -n sso-service --sort-by='.lastTimestamp'
```

Delete all images from container registry (token with rules on write+delete packages)

```bash
werf purge --repo ghcr.io/go-list-templ/sso-service --repo-github-token GH_TOKEN
```