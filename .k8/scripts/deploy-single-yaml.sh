#!/bin/bash

set -x

if [[ -z $1 ]]; then
  echo "Please provide DOCKERREPO"
  exit 1
fi

if [[ -z $2 ]]; then
  echo "Please provide IMAGETAG"
  exit 1
fi

if [[ -z $3 ]]; then
  echo "Please provide NAMESPACE"
  exit 1
fi

if [[ -z $4 ]]; then
  echo "Please provide number of replicas"
  exit 1
fi

if [[ -z $5 ]]; then
  echo "Please provide deployment type: prod|stage"
  exit 1
fi

if [[ -z $6 ]]; then
  echo "Please provide k8s context"
  exit 1
fi

if [[ -z $7 ]]; then
  echo "Pleae provide domain suffix"
  exit 1
fi

if [[ -z $8 ]]; then
  echo "Please provide k8s app version"
  exit 1
fi


DOCKERREPO=$1
IMAGETAG=$2
NAMESPACE=$3
YAML_FILE=$4
DEPL_TYPE=$5
K8S_CONTEXT=$6
DOMAIN_SUFFIX=$7
K8S_API_VERSION=$8

echo $DOCKERREPO
echo $IMAGETAG
echo $NAMESPACE
echo $YAML_FILE
echo $DEPL_TYPE
echo $K8S_CONTEXT
echo $DOMAIN_SUFFIX
echo $K8S_API_VERSION

export POD_POSTFIX=${RANDOM}
echo $POD_POSTFIX

# create namespace if not exists
if ! kubectl --context=$K8S_CONTEXT get ns | grep -q $NAMESPACE; then
  echo "$NAMESPACE created"
  kubectl --context=$K8S_CONTEXT create namespace $NAMESPACE
fi

if [[ -f .k8/${YAML_FILE} ]]; then
   sed -i -e "s|REPLACE_NAMESPACE|${NAMESPACE}|g" \
          -e "s|REPLACE_JOB_POSTFIX|${POD_POSTFIX}|g" \
          -e "s|REPLACE_DOCKER_REPO|${DOCKERREPO}|g" \
          -e "s|REPLACE_DOMAIN_SUFFIX|${DOMAIN_SUFFIX}|g" \
          -e "s|REPLACE_API_VERSION|${K8S_API_VERSION}|g" \
          -e "s|REPLACE_IMAGETAG|${IMAGETAG}|g" ".k8/${YAML_FILE}" || exit 1
fi

#deploy
kubectl --context=$K8S_CONTEXT apply -f .k8/${YAML_FILE} --wait=true || exit 1
kubectl wait --for=condition=terminated -n validators pod vault-plugin-secrets-test-$(kubectl get pods --selector=job-name=vault-plugin-secrets-test-$IMAGETAG  -o=jsonpath='{.items[0].metadata.name}' -n validators) || exit 1