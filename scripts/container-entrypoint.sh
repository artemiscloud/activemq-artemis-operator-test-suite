#!/usr/bin/env bash

set -e -o pipefail

if [[ "${DEBUG_ENTRYPOINT}" == "true" ]]; then
    set -x
    echo "User Id: $(id)"
    echo "Home dir contents (${HOME}):"
    ls -la "${HOME}"
    echo "Root dir contents (/):"
    ls -la /
fi

WORKDIR_VALUE="${WORKDIR:-/workdir}"
export KUBECONFIG="${WORKDIR_VALUE}/kubeconfig"
OUTPUT_DIR_PARAM_VALUE="${OUTPUT_DIR:-${WORKDIR_VALUE}/output}"
mkdir -p "${OUTPUT_DIR_PARAM_VALUE}"
REPORT_DIR_PARAM_VALUE="${REPORT_DIR:-${WORKDIR_VALUE}/reports}"
mkdir -p "${REPORT_DIR_PARAM_VALUE}"

# Get operator repo
REPOSITORY_ZIP_URL_VALUE="${REPOSITORY_ZIP_URL:-https://github.com/artemiscloud/activemq-artemis-operator/archive/refs/heads/main.zip}"
OPERATOR_REPO_ZIP_FILE="operator-repo.zip"
OPERATOR_REPO_DIR="operator-repo"
mkdir -p "${WORKDIR_VALUE}" && cd "${WORKDIR_VALUE}"
curl -s -L -o "${OPERATOR_REPO_ZIP_FILE}" "${REPOSITORY_ZIP_URL_VALUE}"
unzip -q -o -d "${OPERATOR_REPO_DIR}" "${OPERATOR_REPO_ZIP_FILE}"
rm -rf "${OPERATOR_REPO_ZIP_FILE}"
OPERATOR_REPO_DEPLOY_DIR="$(readlink -f "$(find "${OPERATOR_REPO_DIR}" -type d -name deploy)")"
cd - >/dev/null 2>&1

# Ginkgo params
if [[ -n ${GINKGO_PARAMS} ]]; then
    IFS=' ' read -ra GINKGO_DEFAULT_BASIC_PARAMS <<< "${GINKGO_PARAMS}"
else
    GINKGO_DEFAULT_BASIC_PARAMS=(-r -keepGoing -noColor -trace -v -randomizeSuites)
fi
ginkgo_params+=("${GINKGO_DEFAULT_BASIC_PARAMS[@]}")
ginkgo_params+=("-slowSpecThreshold" "${SLOW_SPEC_THRESHOLD:-300}")
ginkgo_params+=("-timeout" "${TIMEOUT:-1h}")
ginkgo_params+=("-flakeAttempts" "${FLAKE_ATTEMPTS:-0}")
ginkgo_params+=("-nodes" "${NODES:-1}")
ginkgo_params+=("-outputdir" "${OUTPUT_DIR:-${WORKDIR_VALUE}/output}")
[[ -n "${SKIP_TESTS}" ]] && ginkgo_params+=("-skip" "${SKIP_TESTS}")
[[ -n "${FOCUS_TESTS}" ]] && ginkgo_params+=("-focus" "${FOCUS_TESTS}")
ginkgo_params+=("test/..." "--")

# Test suite params
ginkgo_params+=("-operator-image" "${OPERATOR_IMAGE:-quay.io/artemiscloud/activemq-artemis-operator:latest}")
ginkgo_params+=("-broker-image" "${BROKER_IMAGE:-quay.io/artemiscloud/activemq-artemis-broker-kubernetes:1.0.7}")
ginkgo_params+=("-broker-image-second" "${BROKER_PREVIOUS_IMAGE:-quay.io/artemiscloud/activemq-artemis-broker-kubernetes:1.0.6}")
ginkgo_params+=("-repository" "${OPERATOR_REPO_DEPLOY_DIR}")
ginkgo_params+=("-report-dir" "${REPORT_DIR:-${WORKDIR_VALUE}/reports}")
ginkgo_params+=("-timeoutMult" "${TIMEOUT_MULT:-1}")
[[ "${GLOBAL_OPERATOR:-true}" == "true" ]] && ginkgo_params+=("-global")
[[ "${V2:-true}" == "true" ]] && ginkgo_params+=("-v2")
[[ "${V3:-true}" == "true" ]] && ginkgo_params+=("-v3")
[[ "${OPENSHIFT:-false}" == "true" ]] && ginkgo_params+=("-openshift")
[[ "${NO_ADMIN:-false}" == "true" ]] && ginkgo_params+=("-no-admin-available")
[[ "${IBMZ:-false}" == "true" ]] && ginkgo_params+=("-ibmz")
[[ "${PPC:-false}" == "true" ]] && ginkgo_params+=("-ppc")
ginkgo_params+=("-delete-namespace-on-failure" "${DELETE_NS_ON_FAILURE:-false}")
ginkgo_params+=("-logtostderr")
[[ "${DEBUG:-false}" == "true" ]] && ginkgo_params+=("-debug-run")

# log in into cluster
CLUSTER_NAME="$(echo "${CLUSTER_API_URL}" | sed -E -e 's/http.?:\/\/api\.//' -e 's/\:[0-9]+$//')"

if [[ "${OPENSHIFT}" == "true" ]] ; then
    if [[ "${CLUSTER_NAME}" =~ crc.testing ]]; then
        CLUSTER_APPS_PREFIX="apps-"
    else
        CLUSTER_APPS_PREFIX="apps."
    fi
    CLUSTER_APPS="${CLUSTER_APPS_PREFIX}${CLUSTER_NAME}"
    OAUTH_URL="https://oauth-openshift.${CLUSTER_APPS}/oauth/authorize?response_type=token&client_id=openshift-challenging-client"
    KUBEADMIN_TOKEN=$(curl -v --insecure --user "${CLUSTER_USERNAME}:${CLUSTER_PASSWORD}" --header "X-CSRF-Token: xxx" --url "${OAUTH_URL}" 2>&1 | grep -oP "access_token=\K[^&]*")
    kubectl --kubeconfig="${KUBECONFIG}" config set-credentials "${CLUSTER_USERNAME}/${CLUSTER_NAME}" --token="${KUBEADMIN_TOKEN}"
else
    kubectl --kubeconfig="${KUBECONFIG}" config set-credentials "${CLUSTER_USERNAME}/${CLUSTER_NAME}" --username="${CLUSTER_USERNAME}" --password="${CLUSTER_PASSWORD}"
fi

kubectl --kubeconfig="${KUBECONFIG}" config set-cluster "${CLUSTER_NAME}" --insecure-skip-tls-verify=true --server="${CLUSTER_API_URL}"
kubectl --kubeconfig="${KUBECONFIG}" config set-context "default/${CLUSTER_NAME}/${CLUSTER_USERNAME}" --user="${CLUSTER_USERNAME}/${CLUSTER_NAME}" --namespace=default --cluster="${CLUSTER_NAME}"
kubectl --kubeconfig="${KUBECONFIG}" config use-context "default/${CLUSTER_NAME}/${CLUSTER_USERNAME}"

# Display the execution command line and execute the test suite
echo "ginkgo params: ${ginkgo_params[*]}"

cd "${TEST_SUITE_DIR}"
ginkgo "${ginkgo_params[@]}"

