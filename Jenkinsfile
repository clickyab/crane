node {
    stage('Prepare') {
        checkout scm
        sh "make clean && make prepare"
    }
    stage('Build') {
        checkout scm
        sh "./scripts/ci-test.sh all"
    }
    stage('Test') {
        checkout scm
        def REDIS_NAME = sh(script: 'cat /dev/urandom | tr -dc "a-zA-Z0-9" | fold -w 32 | head -n 1', returnStdout: true).trim()
        sh "docker rm -f $REDIS_NAME || true"
        sh "docker run -d --rm --name $REDIS_NAME redis:alpine"
        def TEST_SERVICES_REDIS_ADDRESS=sh( script: "docker inspect -f '{{.NetworkSettings.IPAddress}}' $REDIS_NAME", returnStdout: true).trim()
        sh "./scripts/ci-test.sh test"
        sh "docker rm -f $REDIS_NAME"
    }
    stage('Lint') {
        checkout scm
        sh "./scripts/ci-test.sh lint"
    }
    stage('Deploy into k8s') {
        checkout scm
        def OUT_LOG = sh(script: 'mktemp', returnStdout: true).trim()
        def OUT_LOG_COLOR = sh(script: 'mktemp', returnStdout: true).trim()
        sh "APP=crane OUT_LOG=$OUT_LOG OUT_LOG_COLOR=$OUT_LOG_COLOR bash -x ./scripts/herokutor.sh `pwd`"
        def color = readFile OUT_LOG_COLOR
        def message = readFile OUT_LOG
        slackSend color: color, message: message
    }
}
