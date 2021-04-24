Can build and deploy with `make manifests generate docker-build docker-push deploy`. This will deploy the operator to the `dbaas-backing-operator-system` namespace.

The full build will currently also regenerate and overwrite the CRD. Revert the [CRD](config/crd/bases/database.openshift.io_connections.yaml) and make sure the original one is installed:
`git checkout config/crd/bases/database.openshift.io_connections.yaml`
`oc apply -f git checkout config/crd/bases/database.openshift.io_connections.yaml`

Create a mock [request](config/samples/database_v1_connection.yaml) to import a database: `oc create -f config/samples/database_v1_connection.yaml`. The operator will (mock) create both a configmap and a secret for the imported database and point to them from the status of the Connection object you created. To clean this up and run again later, make sure to remove all 3 objects: `oc delete Connection connection-sample; oc delete cm atlas-db1-cm; oc delete secret atlas-db1-creds`

Create a simple shell application [deployment](servicebinding/simple-busybox.yaml) based on `busybox` that will allow you to `rsh` into it and validate the binding has taken place: `oc create -f servicebinding/simple-busybox.yaml`. Wait for this to start, and verify that database connection info is not yet mounted by inspecting the directories off the root: `ls -l /`

Finally, create a [service binding](servicebinding/service-binding.yaml) that specifices `oc create -f servicebinding/service-binding.yaml`. This binding CR should get reconciled by the cluster-wide Service Binding Operator, and that status will reflect the progress. Once successful, a new deployment will roll out for busybox and it will have a directory with the values written: `oc get pods -o name | grep simple-busybox | xargs -I @ oc rsh @ ls -l /bindings/example-servicebinding`
