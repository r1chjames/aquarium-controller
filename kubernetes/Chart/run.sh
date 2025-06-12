# helm del ddns-client --namespace infra
helm upgrade --install aquarium-controller --namespace default -f values.yaml .  --atomic
