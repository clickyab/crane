#!/bin/bash
set -x
set -eo pipefail

# This job is from jenkins. so kill it if it is a pull request
exit_message() {
    echo ${1:-'exiting...'}
    code=${2:-1}
    if [[ "${code}" == "0" ]]; then
        echo "${APP}:${BRANCH}.${COMMIT_COUNT}" >> ${OUT_LOG}
        echo "Build was OK, but its not the correct branch(${APP}:${BRANCH}.${COMMIT_COUNT} By ${CHANGE_AUTHOR}). ignore this" >> ${OUT_LOG}
        echo "green" > ${OUT_LOG_COLOR}
    else
        echo "${APP}:${BRANCH}.${COMMIT_COUNT}" >> ${OUT_LOG}
        echo "Build was NOT OK (${APP}:${BRANCH}.${COMMIT_COUNT} By ${CHANGE_AUTHOR}). Verify with dev team." >> ${OUT_LOG}
        echo "red" > ${OUT_LOG_COLOR}
    fi;
    exit ${code}
}

OUT_LOG=${OUT_LOG:-/dev/null}
OUT_LOG_COLOR=${OUT_LOG_COLOR:-/dev/null}
echo "" > ${OUT_LOG}
echo "red" > ${OUT_LOG_COLOR}
APP=${APP:-}
BRANCH=${BRANCH_NAME:-master}
BRANCH=${CHANGE_TARGET:-${BRANCH}}
CACHE_ROOT=${CACHE_ROOT:-/var/lib/jenkins/cache}

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

export LONG_HASH=$(git log -n1 --pretty="format:%H" | cat)
export SHORT_HASH=$(git log -n1 --pretty="format:%h"| cat)
export COMMIT_DATE=$(git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
export IMP_DATE=$(date +%Y%m%d)
export COMMIT_COUNT=$(git rev-list HEAD --count| cat)
export BUILD_DATE=$(date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
popd

[ -z ${APP} ] && exit_message "The APP is not defined." # WTF, the APP NAME is important
[ -z ${CHANGE_AUTHOR} ] || exit_message "It's a PR, bail out" 0
if [[ ( "${BRANCH}" != "master" ) && ( "${BRANCH}" != "dev" ) && ( "${BRANCH}" != "revert" ) ]]; then
    exit_message "Its not on correct branch, bail out" 0
fi

# Populate env for herokuish
env -0 | while IFS='=' read -r -d '' n v; do
    echo "${v}">"${VARS}/${n}";
done< <(env -0)

TEMPORARY=$(mktemp -d)

# Create Rockerfile to build with rocker (the Dockerfile enhancer tool)
cat > ${TEMPORARY}/Rockerfile <<EOF
FROM alpine:3.6

ENV LONG_HASH ${LONG_HASH}
ENV SHORT_HASH ${SHORT_HASH}
ENV COMMIT_DATE ${COMMIT_DATE}
ENV IMP_DATE ${IMP_DATE}
ENV COMMIT_COUNT ${COMMIT_COUNT}
ENV BUILD_DATE ${BUILD_DATE}

MOUNT {{ .Build }}:/crane

ENV TZ=Asia/Tehran

# I don't need to set GOPATH since the Makefile takes care of that
RUN apk add --no-cache --virtual .build-deps git go libc-dev make \
    && apk add --no-cache ca-certificates bash wget tzdata && update-ca-certificates \
    && cp /usr/share/zoneinfo/\$TZ /etc/localtime && echo \$TZ > /etc/timezone \
    && mkdir -p /gopath/src/clickyab.com/ && cp -r /crane /gopath/src/clickyab.com/ \
    && cd /gopath/src/clickyab.com/crane && make \
    && apk del .build-deps \
    && mkdir -p /app/bin /app/statics \
    && mv /gopath/src/clickyab.com/crane/bin/* /app/bin/ \
    && mv /gopath/src/clickyab.com/crane/statics/* /app/statics/ \
    && mkdir -p /app/statics \
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
if [[ ( "${BRANCH}" != "master" ) && ( "${BRANCH}" != "dev" ) && ( "${BRANCH}" != "revert" ) ]]; then
    PUSH=""
fi
rocker build --no-cache ${PUSH} -var Build=${BUILD} -var EnvDir=${VARS} -var Cache=${CACHE} -var Target=${TARGET} -var Version=${COMMIT_COUNT} -var App=${APP}_${BRANCH}
popd

NAMESPACE="${APP}"
VERSION="${COMMIT_COUNT}"
if [[ "${BRANCH}" == "dev" ]]; then
    NAMESPACE=${APP}-staging
    #VERSION="latest"
fi

if [[ "${BRANCH}" == "revert" ]]; then
    NAMESPACE=${APP}
    BRANCH="master"
    #VERSION="latest"
fi

echo "${VARS}" >> /tmp/kill-me
echo "${TARGET}" >> /tmp/kill-me
echo "${TEMPORARY}" >> /tmp/kill-me
echo "${BUILD_DIR}" >> /tmp/kill-me
echo "${BUILD_PACKS_DIR}" >> /tmp/kill-me

echo "${APP}:${BRANCH}.${COMMIT_COUNT}" >> ${OUT_LOG}
echo "The branch ${BRANCH} build finished, try to deploy it" >> ${OUT_LOG}
echo "If there is no report after this for successful deploy, it means the deply failed. report it please." >> ${OUT_LOG}
for WRK_TYP in demand-server supplier-server impression-worker click-worker
do
   kubectl -n ${NAMESPACE} set image deployment  ${APP}-${WRK_TYP} ${APP}-${BRANCH}=registry.clickyab.ae/clickyab/${APP}_${BRANCH}:${VERSION} --record
done
echo "..." >> ${OUT_LOG}
echo "Deploy done successfully to image registry.clickyab.ae/clickyab/${APP}:${BRANCH}.${COMMIT_COUNT}" >> ${OUT_LOG}
echo "green" >> ${OUT_LOG_COLOR}
