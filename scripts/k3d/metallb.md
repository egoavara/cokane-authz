# 구성과정

```bash
k3d cluster create -c ./k3d.yaml
```

```bash
docker network inspect k3d-egoavara-net | jq '.[0].IPAM.Config[0].Subnet'
```
https://metallb.universe.tf/installation/