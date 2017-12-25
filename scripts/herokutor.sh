#!/bin/bash
set -x
set -eo pipefail

# This job is from jenkins. so kill it if it is a pull request
exit_message() {
    echo ${1:-'exiting...'}
    code=${2:-1}
    exit ${code}
}

APP=${APP:-}
BRANCH=${BRANCH_NAME:-master}
BRANCH=${CHANGE_TARGET:-${BRANCH}}
CACHE_ROOT=${CACHE_ROOT:-/var/lib/jenkins/cache}

[ -z ${APP} ] && exit_message "The APP is not defined." # WTF, the APP NAME is important
[ -z ${CHANGE_AUTHOR} ] || exit_message "It's a PR, bail out" 0
if [[ ( "${BRANCH}" != "master" ) && ( "${BRANCH}" != "dev" ) ]]; then
    exit_message "Its not on correct branch, bail out" 0
fi

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

SOURCE_DIR=${1:-}
[ -z ${SOURCE_DIR} ] && exit_message "Must pass the source directory as the first parameter" 1
SOURCE_DIR=$(cd "${SOURCE_DIR}/" && pwd)

BUILD_DIR=${2:-$(mktemp -d)}
# Make sure the cache is writable on the worker server
CACHE_DIR=${CACHE_ROOT}/${APP}-${BRANCH}
ENV_DIR=$(mktemp -d)

mkdir -p "${BUILD_DIR}" "${CACHE_DIR}" "${ENV_DIR}"
BUILD=$(cd "${BUILD_DIR}/" && pwd)
CACHE=$(cd "${CACHE_DIR}/" && pwd)
VARS=$(cd "${ENV_DIR}/" && pwd)

#chown $(id -u):$(id -g) $ENV_DIR
#chown $(id -u):$(id -g) $CACHE_DIR

BUILD_PACKS_DIR=$(mktemp -d)

# Extract build data
pushd ${SOURCE_DIR}
GIT_WORK_TREE=${BUILD} git checkout -f HEAD

export LONGHASH=$(git log -n1 --pretty="format:%H" | cat)
export SHORTHASH=$(git log -n1 --pretty="format:%h"| cat)
export COMMITDATE=$(git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
export IMPDATE=$(date +%Y%m%d)
export COMMITCOUNT=$(git rev-list HEAD --count| cat)
export BUILDDATE=$(date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
popd

# Populate env for herokuish
env -0 | while IFS='=' read -r -d '' n v; do
    echo "${v}">"${VARS}/${n}";
done< <(env -0)

TEMPORARY=$(mktemp -d)

# Create Rockerfile to build with rocker (the Dockerfile enhancer tool)
cat > ${TEMPORARY}/Rockerfile <<EOF
FROM alpine:3.6

MOUNT {{ .Build }}:/crane

ENV TZ=Asia/Tehran

# I don't need to set GOPATH since the Makefile takes care of that
RUN apk add --no-cache --virtual .build-deps git go libc-dev make tzdata \
    && cp /usr/share/zoneinfo/\$TZ /etc/localtime && echo \$TZ > /etc/timezone \
    && apk add --no-cache ca-certificates bash wget && update-ca-certificates \
    && mkdir -p /gopath/src/clickyab.com/ && cp -r /crane /gopath/src/clickyab.com/ \
    && cd /gopath/src/clickyab.com/crane && make \
    && apk del .build-deps \
    && mkdir -p /app/bin \
    && mv /gopath/src/clickyab.com/crane/bin/* /app/bin/ \
    && rm -rf /gopath /go

TAG registry.clickyab.ae/clickyab/{{ .App }}:{{ .Version }}
PUSH registry.clickyab.ae/clickyab/{{ .App }}:{{ .Version }}
TAG registry.clickyab.ae/clickyab/{{ .App }}:latest
PUSH registry.clickyab.ae/clickyab/{{ .App }}:latest
EOF

TARGET=$(mktemp -d)
pushd ${TEMPORARY}
# Actual build
PUSH="--push"
if [[ ( "${BRANCH}" != "master" ) && ( "${BRANCH}" != "dev" ) ]]; then
    PUSH=""
fi
rocker build --no-cache ${PUSH} -var Build=${SOURCE_DIR} -var EnvDir=${VARS} -var Cache=${CACHE} -var Target=${TARGET} -var Version=${COMMITCOUNT} -var App=${APP}_${BRANCH}
popd

NAMESPACE="${APP}"
VERSION="${COMMITCOUNT}"
if [[ "${BRANCH}" == "dev" ]]; then
    NAMESPACE=${APP}-staging
    #VERSION="latest"
fi

echo "${VARS}" >> /tmp/kill-me
echo "${TARGET}" >> /tmp/kill-me
echo "${TEMPORARY}" >> /tmp/kill-me
echo "${BUILD_DIR}" >> /tmp/kill-me
echo "${BUILD_PACKS_DIR}" >> /tmp/kill-me

for WRK_TYP in demand-server supplier-server impression-worker click-worker
do
   kubectl -n ${NAMESPACE} set image deployment  ${APP}-${WRK_TYP} ${APP}-${BRANCH}=registry.clickyab.ae/clickyab/${APP}_${BRANCH}:${VERSION} --record
done
