// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package controllers is a generated GoMock package.
package controllers

import (
	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/openshift/route-monitor-operator/api/v1alpha1"
	blackboxexporter "github.com/openshift/route-monitor-operator/pkg/consts/blackboxexporter"
	reconcile "github.com/openshift/route-monitor-operator/pkg/util/reconcile"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	v10 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	reflect "reflect"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// MockMonitorResourceHandler is a mock of MonitorResourceHandler interface
type MockMonitorResourceHandler struct {
	ctrl     *gomock.Controller
	recorder *MockMonitorResourceHandlerMockRecorder
}

// MockMonitorResourceHandlerMockRecorder is the mock recorder for MockMonitorResourceHandler
type MockMonitorResourceHandlerMockRecorder struct {
	mock *MockMonitorResourceHandler
}

// NewMockMonitorResourceHandler creates a new mock instance
func NewMockMonitorResourceHandler(ctrl *gomock.Controller) *MockMonitorResourceHandler {
	mock := &MockMonitorResourceHandler{ctrl: ctrl}
	mock.recorder = &MockMonitorResourceHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMonitorResourceHandler) EXPECT() *MockMonitorResourceHandlerMockRecorder {
	return m.recorder
}

// SetErrorStatus mocks base method
func (m *MockMonitorResourceHandler) SetErrorStatus(errorStatus *string, err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetErrorStatus", errorStatus, err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// SetErrorStatus indicates an expected call of SetErrorStatus
func (mr *MockMonitorResourceHandlerMockRecorder) SetErrorStatus(errorStatus, err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetErrorStatus", reflect.TypeOf((*MockMonitorResourceHandler)(nil).SetErrorStatus), errorStatus, err)
}

// ParseMonitorSLOSpecs mocks base method
func (m *MockMonitorResourceHandler) ParseMonitorSLOSpecs(routeURL string, sloSpec v1alpha1.SloSpec) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseMonitorSLOSpecs", routeURL, sloSpec)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseMonitorSLOSpecs indicates an expected call of ParseMonitorSLOSpecs
func (mr *MockMonitorResourceHandlerMockRecorder) ParseMonitorSLOSpecs(routeURL, sloSpec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseMonitorSLOSpecs", reflect.TypeOf((*MockMonitorResourceHandler)(nil).ParseMonitorSLOSpecs), routeURL, sloSpec)
}

// SetResourceReference mocks base method
func (m *MockMonitorResourceHandler) SetResourceReference(reference *v1alpha1.NamespacedName, target types.NamespacedName) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetResourceReference", reference, target)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetResourceReference indicates an expected call of SetResourceReference
func (mr *MockMonitorResourceHandlerMockRecorder) SetResourceReference(reference, target interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetResourceReference", reflect.TypeOf((*MockMonitorResourceHandler)(nil).SetResourceReference), reference, target)
}

// UpdateMonitorResource mocks base method
func (m *MockMonitorResourceHandler) UpdateMonitorResource(cr client.Object) (reconcile.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMonitorResource", cr)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMonitorResource indicates an expected call of UpdateMonitorResource
func (mr *MockMonitorResourceHandlerMockRecorder) UpdateMonitorResource(cr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMonitorResource", reflect.TypeOf((*MockMonitorResourceHandler)(nil).UpdateMonitorResource), cr)
}

// UpdateMonitorResourceStatus mocks base method
func (m *MockMonitorResourceHandler) UpdateMonitorResourceStatus(cr client.Object) (reconcile.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMonitorResourceStatus", cr)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMonitorResourceStatus indicates an expected call of UpdateMonitorResourceStatus
func (mr *MockMonitorResourceHandlerMockRecorder) UpdateMonitorResourceStatus(cr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMonitorResourceStatus", reflect.TypeOf((*MockMonitorResourceHandler)(nil).UpdateMonitorResourceStatus), cr)
}

// SetFinalizer mocks base method
func (m *MockMonitorResourceHandler) SetFinalizer(o v10.Object, finalizerKey string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetFinalizer", o, finalizerKey)
	ret0, _ := ret[0].(bool)
	return ret0
}

// SetFinalizer indicates an expected call of SetFinalizer
func (mr *MockMonitorResourceHandlerMockRecorder) SetFinalizer(o, finalizerKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFinalizer", reflect.TypeOf((*MockMonitorResourceHandler)(nil).SetFinalizer), o, finalizerKey)
}

// DeleteFinalizer mocks base method
func (m *MockMonitorResourceHandler) DeleteFinalizer(o v10.Object, finalizerKey string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFinalizer", o, finalizerKey)
	ret0, _ := ret[0].(bool)
	return ret0
}

// DeleteFinalizer indicates an expected call of DeleteFinalizer
func (mr *MockMonitorResourceHandlerMockRecorder) DeleteFinalizer(o, finalizerKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFinalizer", reflect.TypeOf((*MockMonitorResourceHandler)(nil).DeleteFinalizer), o, finalizerKey)
}

// GetClusterID mocks base method
func (m *MockMonitorResourceHandler) GetClusterID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusterID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetClusterID indicates an expected call of GetClusterID
func (mr *MockMonitorResourceHandlerMockRecorder) GetClusterID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterID", reflect.TypeOf((*MockMonitorResourceHandler)(nil).GetClusterID))
}

// MockServiceMonitorHandler is a mock of ServiceMonitorHandler interface
type MockServiceMonitorHandler struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMonitorHandlerMockRecorder
}

// MockServiceMonitorHandlerMockRecorder is the mock recorder for MockServiceMonitorHandler
type MockServiceMonitorHandlerMockRecorder struct {
	mock *MockServiceMonitorHandler
}

// NewMockServiceMonitorHandler creates a new mock instance
func NewMockServiceMonitorHandler(ctrl *gomock.Controller) *MockServiceMonitorHandler {
	mock := &MockServiceMonitorHandler{ctrl: ctrl}
	mock.recorder = &MockServiceMonitorHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockServiceMonitorHandler) EXPECT() *MockServiceMonitorHandlerMockRecorder {
	return m.recorder
}

// GetServiceMonitor mocks base method
func (m *MockServiceMonitorHandler) GetServiceMonitor(namespacedName types.NamespacedName) (v1.ServiceMonitor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceMonitor", namespacedName)
	ret0, _ := ret[0].(v1.ServiceMonitor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceMonitor indicates an expected call of GetServiceMonitor
func (mr *MockServiceMonitorHandlerMockRecorder) GetServiceMonitor(namespacedName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceMonitor", reflect.TypeOf((*MockServiceMonitorHandler)(nil).GetServiceMonitor), namespacedName)
}

// UpdateServiceMonitorDeployment mocks base method
func (m *MockServiceMonitorHandler) UpdateServiceMonitorDeployment(template v1.ServiceMonitor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateServiceMonitorDeployment", template)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateServiceMonitorDeployment indicates an expected call of UpdateServiceMonitorDeployment
func (mr *MockServiceMonitorHandlerMockRecorder) UpdateServiceMonitorDeployment(template interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateServiceMonitorDeployment", reflect.TypeOf((*MockServiceMonitorHandler)(nil).UpdateServiceMonitorDeployment), template)
}

// DeleteServiceMonitorDeployment mocks base method
func (m *MockServiceMonitorHandler) DeleteServiceMonitorDeployment(serviceMonitorRef v1alpha1.NamespacedName) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteServiceMonitorDeployment", serviceMonitorRef)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteServiceMonitorDeployment indicates an expected call of DeleteServiceMonitorDeployment
func (mr *MockServiceMonitorHandlerMockRecorder) DeleteServiceMonitorDeployment(serviceMonitorRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteServiceMonitorDeployment", reflect.TypeOf((*MockServiceMonitorHandler)(nil).DeleteServiceMonitorDeployment), serviceMonitorRef)
}

// MockPrometheusRuleHandler is a mock of PrometheusRuleHandler interface
type MockPrometheusRuleHandler struct {
	ctrl     *gomock.Controller
	recorder *MockPrometheusRuleHandlerMockRecorder
}

// MockPrometheusRuleHandlerMockRecorder is the mock recorder for MockPrometheusRuleHandler
type MockPrometheusRuleHandlerMockRecorder struct {
	mock *MockPrometheusRuleHandler
}

// NewMockPrometheusRuleHandler creates a new mock instance
func NewMockPrometheusRuleHandler(ctrl *gomock.Controller) *MockPrometheusRuleHandler {
	mock := &MockPrometheusRuleHandler{ctrl: ctrl}
	mock.recorder = &MockPrometheusRuleHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPrometheusRuleHandler) EXPECT() *MockPrometheusRuleHandlerMockRecorder {
	return m.recorder
}

// UpdatePrometheusRuleDeployment mocks base method
func (m *MockPrometheusRuleHandler) UpdatePrometheusRuleDeployment(template v1.PrometheusRule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePrometheusRuleDeployment", template)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePrometheusRuleDeployment indicates an expected call of UpdatePrometheusRuleDeployment
func (mr *MockPrometheusRuleHandlerMockRecorder) UpdatePrometheusRuleDeployment(template interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePrometheusRuleDeployment", reflect.TypeOf((*MockPrometheusRuleHandler)(nil).UpdatePrometheusRuleDeployment), template)
}

// DeletePrometheusRuleDeployment mocks base method
func (m *MockPrometheusRuleHandler) DeletePrometheusRuleDeployment(prometheusRuleRef v1alpha1.NamespacedName) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePrometheusRuleDeployment", prometheusRuleRef)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePrometheusRuleDeployment indicates an expected call of DeletePrometheusRuleDeployment
func (mr *MockPrometheusRuleHandlerMockRecorder) DeletePrometheusRuleDeployment(prometheusRuleRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePrometheusRuleDeployment", reflect.TypeOf((*MockPrometheusRuleHandler)(nil).DeletePrometheusRuleDeployment), prometheusRuleRef)
}

// MockBlackBoxExporterHandler is a mock of BlackBoxExporterHandler interface
type MockBlackBoxExporterHandler struct {
	ctrl     *gomock.Controller
	recorder *MockBlackBoxExporterHandlerMockRecorder
}

// MockBlackBoxExporterHandlerMockRecorder is the mock recorder for MockBlackBoxExporterHandler
type MockBlackBoxExporterHandlerMockRecorder struct {
	mock *MockBlackBoxExporterHandler
}

// NewMockBlackBoxExporterHandler creates a new mock instance
func NewMockBlackBoxExporterHandler(ctrl *gomock.Controller) *MockBlackBoxExporterHandler {
	mock := &MockBlackBoxExporterHandler{ctrl: ctrl}
	mock.recorder = &MockBlackBoxExporterHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBlackBoxExporterHandler) EXPECT() *MockBlackBoxExporterHandlerMockRecorder {
	return m.recorder
}

// EnsureBlackBoxExporterResourcesExist mocks base method
func (m *MockBlackBoxExporterHandler) EnsureBlackBoxExporterResourcesExist() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnsureBlackBoxExporterResourcesExist")
	ret0, _ := ret[0].(error)
	return ret0
}

// EnsureBlackBoxExporterResourcesExist indicates an expected call of EnsureBlackBoxExporterResourcesExist
func (mr *MockBlackBoxExporterHandlerMockRecorder) EnsureBlackBoxExporterResourcesExist() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsureBlackBoxExporterResourcesExist", reflect.TypeOf((*MockBlackBoxExporterHandler)(nil).EnsureBlackBoxExporterResourcesExist))
}

// EnsureBlackBoxExporterResourcesAbsent mocks base method
func (m *MockBlackBoxExporterHandler) EnsureBlackBoxExporterResourcesAbsent() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnsureBlackBoxExporterResourcesAbsent")
	ret0, _ := ret[0].(error)
	return ret0
}

// EnsureBlackBoxExporterResourcesAbsent indicates an expected call of EnsureBlackBoxExporterResourcesAbsent
func (mr *MockBlackBoxExporterHandlerMockRecorder) EnsureBlackBoxExporterResourcesAbsent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsureBlackBoxExporterResourcesAbsent", reflect.TypeOf((*MockBlackBoxExporterHandler)(nil).EnsureBlackBoxExporterResourcesAbsent))
}

// ShouldDeleteBlackBoxExporterResources mocks base method
func (m *MockBlackBoxExporterHandler) ShouldDeleteBlackBoxExporterResources() (blackboxexporter.ShouldDeleteBlackBoxExporter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldDeleteBlackBoxExporterResources")
	ret0, _ := ret[0].(blackboxexporter.ShouldDeleteBlackBoxExporter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShouldDeleteBlackBoxExporterResources indicates an expected call of ShouldDeleteBlackBoxExporterResources
func (mr *MockBlackBoxExporterHandlerMockRecorder) ShouldDeleteBlackBoxExporterResources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldDeleteBlackBoxExporterResources", reflect.TypeOf((*MockBlackBoxExporterHandler)(nil).ShouldDeleteBlackBoxExporterResources))
}

// GetBlackBoxExporterNamespace mocks base method
func (m *MockBlackBoxExporterHandler) GetBlackBoxExporterNamespace() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlackBoxExporterNamespace")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetBlackBoxExporterNamespace indicates an expected call of GetBlackBoxExporterNamespace
func (mr *MockBlackBoxExporterHandlerMockRecorder) GetBlackBoxExporterNamespace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlackBoxExporterNamespace", reflect.TypeOf((*MockBlackBoxExporterHandler)(nil).GetBlackBoxExporterNamespace))
}
