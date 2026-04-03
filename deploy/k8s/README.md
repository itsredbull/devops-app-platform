# Kubernetes Deploy (Kustomize)

This project uses Kustomize (not Helm) for Phase 5.

## Layout

- `deploy/k8s/base`: shared app + postgres resources
- `deploy/k8s/overlays/dev`: dev environment customization
- `deploy/k8s/overlays/staging`: staging environment customization

## Secrets and Config Management

Config is managed with `configMapGenerator` in overlays.

Secrets are managed with `secretGenerator` in overlays and loaded from `secret.env`.

Files committed:

- `secret.env.example` (template)

Files ignored:

- `secret.env` (real values)

## Deploy Dev

```bash
cp deploy/k8s/overlays/dev/secret.env.example deploy/k8s/overlays/dev/secret.env
# edit secret values
kubectl apply -k deploy/k8s/overlays/dev
```

## Deploy Staging

```bash
cp deploy/k8s/overlays/staging/secret.env.example deploy/k8s/overlays/staging/secret.env
# edit secret values
kubectl apply -k deploy/k8s/overlays/staging
```

## Validate Rendered Manifests

```bash
kubectl kustomize deploy/k8s/overlays/dev
kubectl kustomize deploy/k8s/overlays/staging
```

## Update App Image Tag (staging)

Edit `deploy/k8s/overlays/staging/kustomization.yaml` under `images.newTag`.
