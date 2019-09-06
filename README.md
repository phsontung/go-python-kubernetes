## Summary
- Frontend (Reacjts)
- Producer (Rabbitmq) use python
- Consumer (Rabbitmq) use golang and send to web client via websocket
- Traefik to routing traffic
- Kubernetes deployment

## Tips/Tricks
### Start local registry
```
docker run -d -p 5555:5000 --restart=always --name registry registry:2
```

### Debug kubernetes network
```
kubectl run --generator=run-pod/v1 tmp-shell --rm -i --tty --image nicolaka/netshoot -- /bin/bash
```