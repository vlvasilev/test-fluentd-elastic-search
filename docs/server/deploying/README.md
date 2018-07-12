# DEPLOYING OF SERVER_ONE

## Simple Deploy
`export KUBECONFIG=<path to kubeconfig>`
`make deploy-server-one`

## Namespace
Kubernetes supports multiple virtual clusters backed by the same physical cluster. 
These virtual clusters are called namespaces.
In this namespace our server_one will leave.
`kubectl apply -f namespace.yaml`
## ServiceAccount
A service account provides an identity for processes that run in a Pod.
We need this service account to bound a role to it and thus the server_one
will have permission to crate and delete jobs and pods.
`kubectl apply -f serviceaccount.yaml`
## ClusterRole
In the RBAC API, a role contains rules that represent a set of permissions.
By this role we give permissions to server_one to manipulate jobs and pods.
## ClusterRoleBinding
A role binding grants the permissions defined in a role to a user or set of users.
`kubectl apply -f clusterrolebinding.yaml`
## Deployment
The deployment resource will take care of the Pod on which server_one leaves.
`kubectl apply -f deployment.yaml`
## Service
By the service of type LoadBalancer we are going to expose the server_one API.
`kubectl apply -f loadbalancer.yaml`
