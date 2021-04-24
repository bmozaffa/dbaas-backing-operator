/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	databasev1 "github.com/bmozaffa/dbaas-backing-operator/api/v1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConnectionReconciler reconciles a Connection object
type ConnectionReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=database.openshift.io,resources=connections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=database.openshift.io,resources=connections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=database.openshift.io,resources=connections/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Connection object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *ConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("connection", req.NamespacedName)

	connection := &databasev1.Connection{}
	if err := r.Get(ctx, req.NamespacedName, connection); err != nil {
		return ctrl.Result{}, err
	}

	//r.Log.Info("Processing database connection")
	//r.Log.Info(connection.Spec.Database)
	r.Log.Info("Processing database connection", "provider", connection.Spec.Provider, "database", connection.Spec.Database)
	if connection.Status.DBConfigMap != "" {
		r.Log.Info("Database connection previously reconciled", "DBConfigMap", connection.Status.DBConfigMap, "DBCredentials", connection.Status.DBCredentials)
		return ctrl.Result{}, nil
	}

	connection.Status.DBConfigMap = fmt.Sprintf("%s-%s-%s", connection.Spec.Provider, connection.Spec.Database, "cm")
	configMap := &corev1.ConfigMap{
		ObjectMeta: ctrl.ObjectMeta{
			Namespace: connection.Namespace,
			Name:      connection.Status.DBConfigMap,
		},
		Data: map[string]string{
			"db.host": "example.mongodb.net",
			"db.port": "27017",
			"db.name": connection.Spec.Database,
		},
	}
	if err := r.Create(context.TODO(), configMap); err != nil {
		return ctrl.Result{}, err
	}

	connection.Status.DBCredentials = fmt.Sprintf("%s-%s-%s", connection.Spec.Provider, connection.Spec.Database, "creds")
	secret := &corev1.Secret{
		ObjectMeta: ctrl.ObjectMeta{
			Namespace: connection.Namespace,
			Name:      connection.Status.DBCredentials,
		},
		Data: map[string][]byte{
			"db.user":     []byte("username1"),
			"db.password": []byte("password1"),
		},
	}
	if err := r.Create(context.TODO(), secret); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Status().Update(context.TODO(), connection); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1.Connection{}).
		Complete(r)
}
