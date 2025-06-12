# Aquarium Controller Deployment to Kubernetes
---
Use the included helm chart to deploy the provided app to Kubernetes. The helm chart uses Smarter Device Manager to expose the necessary GPIO pins and devices (such as an attached camera) as K8s resources that can then be passed into the running K8s deployment.

## Installation

#### Smarter Device Management
- Ensure the relevant node(s) are specified in the `smarter-devices.yaml` file.
- Apply this `kubectl apply -f smater-devices.yaml`.
- Verify the exposed GPIO pins and devices by describing the node(s) `kubectl describe node <node>`.

#### Application
- Update the relevant variables in the `values.yaml` file within the `Chart` directory
- Apply the chart using `helm upgrade --install aquarium-controller --namespace default -f values.yaml .  --atomic`.


