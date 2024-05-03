#!/bin/bash 
set -e

# Builds configure-alertmanager-operator && runs it in a cluster. 
# This should be run in the 'configure-alertmanager-operator' repo while logged into a cluster. 

# Remember to disable olm management first:
# oc scale deploy/cluster-version-operator --replicas=0 -n openshift-cluster-version
# oc scale deploy/olm-operator --replicas=0 -n openshift-operator-lifecycle-manager



export ALLOW_DIRTY_CHECKOUT=true
export REGISTRY_USER="quay.io/hbhushan"
# export IMAGE_REPOSITORY="${REGISTRY_USER}"

## Disallow deployment to non-stage cluster
# ocmlogin
# ocm backplane status | grep "Backplane Server" | grep -q "https://api.integration.backplane.openshift.com"
# rc=$?
# if [ ${rc} -ne 0 ];
# then
#     printf "Refusing to deploy to non-stage cluster. Please login to a stage cluster before continuing."
#     exit 1
# fi

# if [ -z "${REGISTRY_TOKEN}" ];
# then
# 	printf "\${REGISTRY_TOKEN} unset.\n"
# 	exit 1
# fi

printBold() {
	printf "$(tput bold)${@}$(tput sgr0)"
}

printBold "Building operator\n"
#ALLOW_DIRTY_CHECKOUT=true REGISTRY_USER="${QUAY_USER}" REGISTRY_TOKEN="${REGISTRY_TOKEN}" make IMAGE_REPOSITORY="${QUAY_USER}" docker-build docker-push
make docker-build docker-push
sleep 1
printf "\n"

printBold "Disabling default RMO deploy\n"
#oc scale deploy/cluster-version-operator --replicas=0 -n openshift-cluster-version
#oc scale deploy/olm-operator --replicas=0 -n openshift-operator-lifecycle-manager
oc scale deploy/route-monitor-operator-controller-manager --replicas=0 -n openshift-route-monitor-operator
printf "\n"

printBold "Deploying new RMO\n"
oc delete -f deploy.yaml --ignore-not-found=true --as backplane-cluster-admin
oc create -f rmo-deploy.yaml --as backplane-cluster-admin
oc wait pods --for=condition=Ready -l app=route-monitor-operator -n openshift-route-monitor-operator
sleep 1
printf "\n"

#printBold "Updating operator image\n"
#oc patch deployment configure-alertmanager-operator -n openshift-monitoring --patch "{\"spec\": {\"template\": {\"spec\": {\"containers\": [{\"name\": \"configure-alertmanager-operator\",\"image\": \"quay.io/${IMAGE_REPOSITORY}/configure-alertmanager-operator:latest\"}]}}}}"
#patch='{"spec": {"template": {"spec": {"containers": [{"name": "configure-alertmanager-operator","image": "quay.io/${IMAGE_REPOSITORY}/configure-alertmanager-operator:latest"}]}}}}'
#oc patch deployment configure-alertmanager-operator -n openshift-monitoring --patch '' --as backplane-cluster-admin
#printf "\n"

#printBold "Refreshing operator pod\n"
#ocmlogin
#oc delete pod -l name=configure-alertmanager-operator -n openshift-monitoring
#oc wait pods --for=condition=Ready -l name=configure-alertmanager-operator -n openshift-monitoring
#sleep 1
#printf "\n"

printBold "Retrieving operator logs\n"
cmd="oc logs deploy/route-monitor-operator-controller-manager -n openshift-route-monitor-operator -f | less +F"
echo "${cmd}"
eval "${cmd}"
printf "\n"
