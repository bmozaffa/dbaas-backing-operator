Can build and deploy with `make manifests generate docker-build docker-push deploy`. This will deploy the operator to the `dbaas-backing-operator-system` namespace.

The above command will also fully generate and patch the CRD, before installing it. The CRD left on the file system locally is ignored by git, as it is incomplete and should be ignored. 

Create a mock [request](config/samples/database_v1_connection.yaml) to import a database: `oc create -f config/samples/database_v1_connection.yaml`. The operator will (mock) create both a configmap and a secret for the imported database and point to them from the status of the Connection object you created. To clean this up and run again later, make sure to remove all 3 objects: `oc delete Connection connection-sample; oc delete cm atlas-db1-cm; oc delete secret atlas-db1-creds`

Create a simple shell application [deployment](servicebinding/simple-busybox.yaml) based on `busybox` that will allow you to `rsh` into it and validate the binding has taken place. Wait for the pod to start, and verify that database connection info is not available as a set of environment variables yet, by using `oc get pods -o name | grep simple-busybox | xargs -I @ oc rsh @ printenv`.

Finally, create a [service binding](servicebinding/service-binding.yaml) that specifices `oc create -f servicebinding/service-binding.yaml`. This binding CR should get reconciled by the cluster-wide Service Binding Operator, and that status will reflect the progress. Once successful, a new deployment will roll out for busybox and it will inject environment variables with the derived values: `oc get pods -o name | grep simple-busybox | xargs -I @ oc rsh @ printenv`
