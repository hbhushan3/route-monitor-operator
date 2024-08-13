/*


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

package hostedcontrolplane

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"

	routev1 "github.com/openshift/api/route/v1"
	hypershiftv1beta1 "github.com/openshift/hypershift/api/hypershift/v1beta1"
	"github.com/openshift/route-monitor-operator/api/v1alpha1"
	"github.com/openshift/route-monitor-operator/pkg/util/finalizer"
	utilreconcile "github.com/openshift/route-monitor-operator/pkg/util/reconcile"

	v1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	dynatrace "github.com/openshift/route-monitor-operator/pkg/dynatrace"
)

const (
	// hostedcontrolplaneFinalizer defines the finalizer used by this controller's objects
	hostedcontrolplaneFinalizer = "hostedcontrolplane.routemonitoroperator.monitoring.openshift.io/finalizer"

	// watchResourceLabel is a label key indicating which objects this controller should reconcile against
	watchResourceLabel = "hostedcontrolplane.routemonitoroperator.monitoring.openshift.io/managed"

	//httpMonitorLabel is added to hcp object to keep track of when to create and delete of dynatrace http monitor
	httpMonitorLabel = "dynatrace.http.monitor/id"

	//fetch dynatrace secret to get dynatrace api token and tennant url
	secretNamespace    = "openshift-route-monitor-operator"
	secretName         = "dynatrace-token"
	dynatraceApiKey    = "apiToken"
	dynatraceTenantKey = "apiUrl"
)

var logger logr.Logger = ctrl.Log.WithName("controllers").WithName("HostedControlPlane")

// HostedControlPlaneReconciler reconciles a HostedControlPlane object
type HostedControlPlaneReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// NewHostedControlPlaneReconciler creates a HostedControlPlaneReconciler
func NewHostedControlPlaneReconciler(mgr manager.Manager) *HostedControlPlaneReconciler {
	return &HostedControlPlaneReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
}

//+kubebuilder:rbac:groups=openshift.io,resources=hostedcontrolplanes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=openshift.io,resources=hostedcontrolplanes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=openshift.io,resources=hostedcontrolplanes/finalizers,verbs=update

// Reconcile responds to events against watched objects
func (r *HostedControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.WithName("Reconcile").WithValues("name", req.Name, "namespace", req.Namespace)
	log.Info("Reconciling HostedControlPlanes")
	defer log.Info("Finished reconciling HostedControlPlane")

	// Fetch the HostedControlPlane instance
	hostedcontrolplane := &hypershiftv1beta1.HostedControlPlane{}
	err := r.Client.Get(ctx, req.NamespacedName, hostedcontrolplane)
	if err != nil {
		if kerr.IsNotFound(err) {
			log.Info("HostedControlPlane not found, assumed deleted")
			return utilreconcile.Stop()
		}
		log.Error(err, "unable to fetch HostedControlPlane")
		return utilreconcile.RequeueWith(err)
	}

	//Create Dynatrace API client
	valueDynatraceApiToken, valueDynatraceTenant, err := r.getDynatraceSecrets(ctx)
	if err != nil {
		log.Error(err, "failed to get secret for Dynatrace API client")
		return utilreconcile.RequeueWith(err)
	}
	baseURL := fmt.Sprintf("%s/v1", valueDynatraceTenant)
	dynatraceApiClient := dynatrace.NewDynatraceApiClient(baseURL, valueDynatraceApiToken)

	// If the HostedControlPlane is marked for deletion, clean up
	shouldDelete := finalizer.WasDeleteRequested(hostedcontrolplane)
	if shouldDelete {
		err = deleteDynatraceHttpMonitorResources(dynatraceApiClient, log, hostedcontrolplane)
		if err != nil {
			log.Error(err, "failed to delete Dynatrace HTTP Monitor Resources")
			return utilreconcile.RequeueWith(err)
		}

		err := r.finalizeHostedControlPlane(ctx, log, hostedcontrolplane)
		if err != nil {
			log.Error(err, "failed to finalize HostedControlPlane")
			return utilreconcile.RequeueWith(err)
		}
		finalizer.Remove(hostedcontrolplane, hostedcontrolplaneFinalizer)
		err = r.Client.Update(ctx, hostedcontrolplane)
		if err != nil {
			return utilreconcile.RequeueWith(err)
		}
		return utilreconcile.Stop()
	}

	if !finalizer.Contains(hostedcontrolplane.Finalizers, hostedcontrolplaneFinalizer) {
		finalizer.Add(hostedcontrolplane, hostedcontrolplaneFinalizer)
		err := r.Client.Update(ctx, hostedcontrolplane)
		if err != nil {
			return utilreconcile.RequeueWith(err)
		}
	}

	// Check if the HostedControlPlane is ready
	if !hostedcontrolplane.Status.Ready {
		log.Info("skipped deploying monitoring objects: HostedControlPlane not ready")
		return utilreconcile.Stop()
	}

	log.Info("Deploying internal monitoring objects")
	err = r.deployInternalMonitoringObjects(ctx, log, hostedcontrolplane)
	if err != nil {
		log.Error(err, "failed to deploy internal monitoring components")
		return utilreconcile.RequeueWith(err)
	}

	log.Info("Deploying HTTP Monitor Resources")
	err = r.deployDynatraceHttpMonitorResources(ctx, dynatraceApiClient, log, hostedcontrolplane)
	if err != nil {
		log.Error(err, "failed to deploy Dynatrace HTTP Monitor Resources")
		return utilreconcile.RequeueWith(err)
	}

	return ctrl.Result{}, err
}

// deployInternalMonitoringObjects creates or updates the objects needed to monitor the kube-apiserver using cluster-internal routes
func (r *HostedControlPlaneReconciler) deployInternalMonitoringObjects(ctx context.Context, log logr.Logger, hostedcontrolplane *hypershiftv1beta1.HostedControlPlane) error {
	// Create or update route object
	expectedRoute := r.buildInternalMonitoringRoute(hostedcontrolplane)
	err := r.Client.Create(ctx, &expectedRoute)
	if err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
		// Object already exists: update it
		actualRoute := routev1.Route{}
		err := r.Client.Get(ctx, types.NamespacedName{Name: expectedRoute.Name, Namespace: expectedRoute.Namespace}, &actualRoute)
		if err != nil {
			return err
		}
		expectedRoute.ObjectMeta = buildMetadataForUpdate(expectedRoute.ObjectMeta, actualRoute.ObjectMeta)
		err = r.Client.Update(ctx, &expectedRoute)
		if err != nil {
			return err
		}
	}

	// Quick fix to discover the API server port from the service resource
	apiServerService := v1.Service{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: "kube-apiserver", Namespace: hostedcontrolplane.Namespace}, &apiServerService)
	if err != nil {
		return fmt.Errorf("couldn't query API server service resource: %w", err)
	}
	apiServerPort := int64(6443)
	if len(apiServerService.Spec.Ports) > 0 {
		apiServerPort = int64(apiServerService.Spec.Ports[0].Port)
	}

	// Create or update RouteMonitor object
	expectedRouteMonitor := r.buildInternalMonitoringRouteMonitor(expectedRoute, hostedcontrolplane, apiServerPort)
	err = r.Client.Create(ctx, &expectedRouteMonitor)
	if err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
		// Object already exists: update it
		actualRouteMonitor := v1alpha1.RouteMonitor{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: expectedRouteMonitor.Name, Namespace: expectedRouteMonitor.Namespace}, &actualRouteMonitor)
		if err != nil {
			return err
		}
		expectedRouteMonitor.ObjectMeta = buildMetadataForUpdate(expectedRouteMonitor.ObjectMeta, actualRouteMonitor.ObjectMeta)
		err = r.Client.Update(ctx, &expectedRouteMonitor)
		if err != nil {
			return err
		}
	}

	return nil
}

// buildMetadataForUpdate is a helper function to generate valid metadata for an Update request by combining the expected object's Metadata and the actual (on-cluster) object's Metadata
func buildMetadataForUpdate(expected, actual metav1.ObjectMeta) metav1.ObjectMeta {
	actual.Labels = expected.Labels
	actual.OwnerReferences = expected.OwnerReferences
	return actual
}

// buildInternalMonitoringRoute constructs the Route needed to monitor a HostedControlPlane's kube-apiserver via cluster-internal routes
func (r *HostedControlPlaneReconciler) buildInternalMonitoringRoute(hostedcontrolplane *hypershiftv1beta1.HostedControlPlane) routev1.Route {
	weight := int32(100)
	route := routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf("%s-kube-apiserver-monitoring", hostedcontrolplane.Name),
			Namespace:       hostedcontrolplane.Namespace,
			OwnerReferences: buildOwnerReferences(hostedcontrolplane),
			Labels: map[string]string{
				watchResourceLabel: "true",
			},
		},
		Spec: routev1.RouteSpec{
			Host: fmt.Sprintf("kube-apiserver.%s.svc.cluster.local", hostedcontrolplane.Namespace),
			TLS: &routev1.TLSConfig{
				InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyNone,
				Termination:                   routev1.TLSTerminationPassthrough,
			},
			To: routev1.RouteTargetReference{
				Kind:   "Service",
				Name:   "kube-apiserver",
				Weight: &weight,
			},
			WildcardPolicy: routev1.WildcardPolicyNone,
		},
	}

	return route
}

// buildInternalMonitoringRouteMonitor constructs the expected RouteMonitor needed to probe a HostedControlPlane's kube-apiserver using cluster-internal routes
func (r *HostedControlPlaneReconciler) buildInternalMonitoringRouteMonitor(route routev1.Route, hostedcontrolplane *hypershiftv1beta1.HostedControlPlane, apiServerPort int64) v1alpha1.RouteMonitor {
	routemonitor := v1alpha1.RouteMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:            route.Name,
			Namespace:       route.Namespace,
			OwnerReferences: buildOwnerReferences(hostedcontrolplane),
			Labels: map[string]string{
				watchResourceLabel: "true",
			},
		},
		Spec: v1alpha1.RouteMonitorSpec{
			Route: v1alpha1.RouteMonitorRouteSpec{
				Name:      route.Name,
				Namespace: route.Namespace,
				Port:      apiServerPort,
				Suffix:    "/livez",
			},
			SkipPrometheusRule: false,
			Slo: v1alpha1.SloSpec{
				TargetAvailabilityPercent: "99.5",
			},
			InsecureSkipTLSVerify: true,
			ServiceMonitorType:    v1alpha1.ServiceMonitorTypeRHOBS,
		},
	}
	return routemonitor
}

// buildOwnerReferences generates a set OwnerReferences indicating the HostedControlPlane is the owner+controller of the object. This is used
// to trigger reconciles against non-HCP objects (ie - the route & routemonitor generated by this controller)
func buildOwnerReferences(hostedcontrolplane *hypershiftv1beta1.HostedControlPlane) []metav1.OwnerReference {
	return []metav1.OwnerReference{*metav1.NewControllerRef(&hostedcontrolplane.ObjectMeta, hostedcontrolplane.GroupVersionKind())}
}

// finalizeHostedControlPlane cleans up HostedControlPlane-related objects managed by the HostedControlPlaneReconciler
func (r *HostedControlPlaneReconciler) finalizeHostedControlPlane(ctx context.Context, log logr.Logger, hostedcontrolplane *hypershiftv1beta1.HostedControlPlane) error {
	err := r.deleteInternalMonitoringObjects(ctx, log, hostedcontrolplane)
	if err != nil {
		return fmt.Errorf("failed to cleanup internal monitoring resources: %w", err)
	}
	return nil
}

// deleteInternalMonitoringObjects removes the internal monitoring objects for the provided HostedControlPlane
func (r *HostedControlPlaneReconciler) deleteInternalMonitoringObjects(ctx context.Context, log logr.Logger, hostedcontrolplane *hypershiftv1beta1.HostedControlPlane) error {
	// Delete Route
	expectedRoute := r.buildInternalMonitoringRoute(hostedcontrolplane)
	err := r.Client.Delete(ctx, &expectedRoute)
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		log.Info(fmt.Sprintf("Skipped deleting Route %s/%s: already deleted", expectedRoute.Namespace, expectedRoute.Name))
	}

	// Delete routemonitor, port is not relevant for deletion
	expectedRouteMonitor := r.buildInternalMonitoringRouteMonitor(expectedRoute, hostedcontrolplane, 6443)
	err = r.Client.Delete(ctx, &expectedRouteMonitor)
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		log.Info(fmt.Sprintf("Skipped deleting RouteMonitor %s/%s: already deleted", expectedRoute.Namespace, expectedRoute.Name))
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HostedControlPlaneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	selector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      watchResourceLabel,
				Operator: metav1.LabelSelectorOpExists,
			},
		},
	}
	selectorPredicate, err := predicate.LabelSelectorPredicate(selector)
	if err != nil {
		return fmt.Errorf("failed to build label selector predicate for routes: %w", err)
	}

	// The following:
	// - Reconciles against all HostedControlPlane objects
	// - Additionally watches against route & routemonitor objects with the 'watchResourceLabel' present.
	//   When these objects are modified, the HCP specified in the objects' .metadata.OwnerReferences is
	//   reconciled
	return ctrl.NewControllerManagedBy(mgr).
		For(&hypershiftv1beta1.HostedControlPlane{}).
		Watches(
			&routev1.Route{},
			handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &hypershiftv1beta1.HostedControlPlane{}, handler.OnlyControllerOwner()),
			builder.WithPredicates(selectorPredicate),
		).
		Watches(
			&v1alpha1.RouteMonitor{},
			handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &hypershiftv1beta1.HostedControlPlane{}, handler.OnlyControllerOwner()),
			builder.WithPredicates(selectorPredicate),
		).
		Complete(r)
}

// ------------------------------synthetic-monitoring--------------------------

func (r *HostedControlPlaneReconciler) getDynatraceSecrets(ctx context.Context) (string, string, error) {

	secret := &v1.Secret{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: secretName, Namespace: secretNamespace}, secret)
	if err != nil {
		return "", "", fmt.Errorf("error getting Kubernetes secret: %v", err)
	}

	valueBytesDynatraceApiToken, ok := secret.Data[dynatraceApiKey]
	if !ok {
		return "", "", fmt.Errorf("secret did not contain key %s", dynatraceApiKey)
	}
	if len(valueBytesDynatraceApiToken) == 0 {
		return "", "", fmt.Errorf("%s is empty", dynatraceApiKey)
	}
	valueDynatraceApiToken := string(valueBytesDynatraceApiToken)

	valueBytesDynatraceTenant, ok := secret.Data[dynatraceTenantKey]
	if !ok {
		return "", "", fmt.Errorf("secret did not contain key %s", dynatraceTenantKey)
	}
	if len(valueBytesDynatraceTenant) == 0 {
		return "", "", fmt.Errorf("%s is empty", dynatraceTenantKey)
	}
	valueDynatraceTenant := string(valueBytesDynatraceTenant)

	return valueDynatraceApiToken, valueDynatraceTenant, nil
}

func getDynatraceHttpMonitorId(hostedcontrolplane *hypershiftv1beta1.HostedControlPlane) (string, bool) {
	labels := hostedcontrolplane.GetLabels()
	dynatraceHttpMonitorId, ok := labels[httpMonitorLabel]
	return dynatraceHttpMonitorId, ok
}

func (r *HostedControlPlaneReconciler) UpdateHostedControlPlaneLabels(ctx context.Context, hostedcontrolplane *hypershiftv1beta1.HostedControlPlane, key, value string) error {
	labels := hostedcontrolplane.GetLabels()
	labels[key] = value
	hostedcontrolplane.SetLabels(labels)

	err := r.Client.Update(ctx, hostedcontrolplane)
	if err != nil {
		return fmt.Errorf("error updating hostedcontrolplane monitor: %v", err)
	}
	return nil
}

func GetAPIServerHostname(hostedcontrolplane *hypershiftv1beta1.HostedControlPlane) (string, error) {
	for _, service := range hostedcontrolplane.Spec.Services {
		if service.Service == "APIServer" {
			return service.ServicePublishingStrategy.Route.Hostname, nil
		}
	}
	return "", fmt.Errorf("APIServer service not found in the hostedcontrolplane")
}

func (r *HostedControlPlaneReconciler) deployDynatraceHttpMonitorResources(ctx context.Context, dynatraceApiClient *dynatrace.DynatraceApiClient, log logr.Logger, hostedcontrolplane *hypershiftv1beta1.HostedControlPlane) error {
	//if http monitor does not exist, and hcp is not marked for deletion, and hcp is ready, then create http monitor
	//get apiserver
	apiServerHostname, err := GetAPIServerHostname(hostedcontrolplane)
	if err != nil {
		return fmt.Errorf("failed to get APIServer hostname %v", err)
	}
	monitorName := strings.Replace(apiServerHostname, "api.", "", 1)
	// apiServerHostname := hostedcontrolplane.Spec.Services[1].ServicePublishingStrategy.Route.Hostname
	monitorLocation := hostedcontrolplane.Spec.Platform.AWS.EndpointAccess

	//in hcp, spec.services.service["APIServer"].servicePublishingStrategy.route.hostname is api.test-rs1.dgcj.i3.devshift.org
	// apiUrl := "https://api.hb-testing.j1b6.i3.devshift.org/livez"

	apiUrl := fmt.Sprintf("https://%s/livez", apiServerHostname)

	dynatraceHttpMonitorId, ok := getDynatraceHttpMonitorId(hostedcontrolplane)
	if ok {
		log.Info(fmt.Sprintf("HTTP monitor label found. Skipping creating a monitor. Monitor id is %s", dynatraceHttpMonitorId))
		return nil
	}

	// determine location and create monitor
	//public
	if monitorLocation == "PublicAndPrivate" {

		clusterId := hostedcontrolplane.Spec.ClusterID
		clusterRegion := hostedcontrolplane.Spec.Platform.AWS.Region

		dynatraceEquivalentClusterRegionId, err := dynatraceApiClient.GetDynatraceEquivalentClusterRegionId(clusterRegion)
		if err != nil {
			return fmt.Errorf("error getting DynatraceEquivalentClusterRegionId: %v", err)
		}

		monitorId, err := dynatraceApiClient.CreateDynatraceHttpMonitor(monitorName, apiUrl, clusterId, dynatraceEquivalentClusterRegionId)
		if err != nil {
			return fmt.Errorf("error creating HTTP monitor: %v", err)
		}

		err = r.UpdateHostedControlPlaneLabels(ctx, hostedcontrolplane, httpMonitorLabel, monitorId)
		if err != nil {
			return fmt.Errorf("failed to update hostedcontrolplane monitor labels %v", err)
		}

		log.Info("Created HTTP monitor ", monitorId, clusterId)
	}

	//private not supported yet
	if monitorLocation == "Private" {
		log.Info("Private API - Not supported yet.")
		return nil
	}

	return nil
}

func deleteDynatraceHttpMonitorResources(dynatraceApiClient *dynatrace.DynatraceApiClient, log logr.Logger, hostedcontrolplane *hypershiftv1beta1.HostedControlPlane) error {
	//check if monitor exists - has label/monitor on hcp, then delete it
	// key := "dynatrace.http.monitor/id"
	dynatraceHttpMonitorId, ok := getDynatraceHttpMonitorId(hostedcontrolplane)
	if !ok {
		log.Info("HTTP monitor label not found. Skipping deleting HTTP monitor")
		return nil
	}

	err := dynatraceApiClient.DeleteDynatraceHttpMonitor(dynatraceHttpMonitorId)
	if err != nil {
		return fmt.Errorf("error deleting HTTP monitor. Status Code: %v", err)
	}
	log.Info("Successfully deleted HTTP monitor")
	return nil
}
