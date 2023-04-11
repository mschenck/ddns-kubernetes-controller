/*
Copyright 2023.

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

package controller

import (
	"context"
	"fmt"
	"net"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ddnsv1 "github.com/mschenck/ddns-kubernetes-controller/api/v1"
	dnslookup "github.com/mschenck/ddns-kubernetes-controller/internal/dnslookup"
	dnsprovider "github.com/mschenck/ddns-kubernetes-controller/internal/dnsprovider"
	iplookup "github.com/mschenck/ddns-kubernetes-controller/internal/iplookup"
)

// DdnsRecordReconciler reconciles a DdnsRecord object
type DdnsRecordReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ddns.mschenck.com,resources=ddnsrecords,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ddns.mschenck.com,resources=ddnsrecords/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ddns.mschenck.com,resources=ddnsrecords/finalizers,verbs=update

// DDNS Record Reconciler
//
// 1. Query for public IP (the IP we request from i.e. NAT gateway)
// 2. TODO(mschenck): Check what the record currently resolves to
// 3. Update DNS Zone Record to public IP
//
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DdnsRecord object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DdnsRecordReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	var err error

	var ddnsRecord ddnsv1.DdnsRecord
	if err = r.Get(ctx, req.NamespacedName, &ddnsRecord); err != nil {
		log.Log.Error(err, "unable to fetch DdnsRecord")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Duration for both record TTL as well as re-check interval
	recordDuration := ddnsRecord.Spec.Ttl.Duration
	recordSeconds := int64(ddnsRecord.Spec.Ttl.Seconds())
	ctrlResult := ctrl.Result{RequeueAfter: recordDuration} // Retry every TTL

	// Query Public IP
	var ip string
	ip, err = iplookup.Ipify()
	if err != nil {
		log.Log.Error(err, "Error fetching IP")
		return ctrl.Result{}, err
	}
	log.Log.Info(fmt.Sprintf("IP is: %q", ip))

	// Check what the zone record resolves to.
	var dnsIp string
	dnsIp, err = dnslookup.DnsLookup(ctx, ddnsRecord.Spec.Record, ddnsRecord.Spec.Zone)
	fqdn := fmt.Sprintf("%s.%s", ddnsRecord.Spec.Record, ddnsRecord.Spec.Zone)
	if err == nil {
		log.Log.Info(fmt.Sprintf("%q resolves to %q", fqdn, dnsIp))
		if dnsIp == ip {
			return ctrlResult, nil
		}
	} else if err.(*net.DNSError).IsNotFound {
		log.Log.Info("Record does not exist.")
	} else if err != nil {
		log.Log.Error(err, fmt.Sprintf("Failed DNS lookup for %s", fqdn))
		return ctrl.Result{}, err
	}

	// Update zone record
	a := dnsprovider.Aws{}
	err = a.UpdateRecord(ddnsRecord.Spec.Record, ddnsRecord.Spec.Zone, ip, recordSeconds)
	if err != nil {
		log.Log.Error(err, "Error Updating DNS Record")
		return ctrl.Result{}, err
	}
	log.Log.Info("Updated zone record.")

	return ctrlResult, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DdnsRecordReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ddnsv1.DdnsRecord{}).
		Complete(r)
}
