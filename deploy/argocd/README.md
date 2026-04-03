# Argo CD GitOps

This directory contains GitOps bootstrap and applications.

## Layout

- `bootstrap/root-app.yaml`: app-of-apps root Application
- `apps/project.yaml`: Argo CD AppProject
- `apps/app-dev.yaml`: dev environment Application
- `apps/app-staging.yaml`: staging environment Application

## 1) Install Argo CD

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

## 2) Access Argo CD

```bash
kubectl -n argocd port-forward svc/argocd-server 8080:443
```

Get admin password:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 --decode && echo
```

## 3) Bootstrap app-of-apps

```bash
kubectl apply -f deploy/argocd/bootstrap/root-app.yaml
```

This root app syncs `deploy/argocd/apps`, which then creates:

- `AppProject` (`uptime-platform`)
- `uptime-dev` Application
- `uptime-staging` Application

## 4) Auto-sync behavior

Both `uptime-dev` and `uptime-staging` are configured with:

- `automated.prune: true`
- `automated.selfHeal: true`

So Git is the source of truth and drift is corrected automatically.
