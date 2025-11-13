# Okteto Pods to Node report
This is an experiment and Okteto does not officially support it.

This is an example on how you can use Okteto's public API to generate a CSV report of the Nodes where Pods that Okteto manages are scheduled.

* Create an Okteto Admin Token

* Export the token to a local variable:
```
export OKTETO_TOKEN=<<your-token>>
```
* Export the URL of your Okteto instance to a local variable:
```
export OKTETO_URL=<<your-okteto-url>>
```
* Clone the repository:
```
git clone https://github.com/okteto-community/list-pod-nodes
```
* Build the binary:
```
make build
```
* Make sure you have the correct okteto context and kubectl context in your terminal
```
okteto context use <your context>
okteto kubeconfig
```
* Run the command:
```
./list_pod_nodes
```
Once finished, the program will generate a CSV with all your Okteto-managed Pods and their corresponding Nodes.