node {
    stage('Prepare') {
        checkout scm
        sh "make clean && make prepare"
    }
    stage('Build') {
        checkout scm
        sh "./bin/ci-test.sh all"
    }
    stage('Test') {
        checkout scm
        def REDIS_NAME = sh(script: 'cat /dev/urandom | tr -dc "a-zA-Z0-9" | fold -w 32 | head -n 1', returnStdout: true).trim()
        sh "docker rm -f $REDIS_NAME || true"
        sh "docker run -d --rm --name $REDIS_NAME redis:alpine"
        def TEST_SERVICES_REDIS_ADDRESS=sh( script: "docker inspect -f '{{.NetworkSettings.IPAddress}}' $REDIS_NAME", returnStdout: true).trim()
        sh "./bin/ci-test.sh test"
        sh "docker rm -f $REDIS_NAME"
    }
    stage('Lint') {
        checkout scm
        sh "./bin/ci-test.sh lint"
    }
    stage('Deploy into k8s') {
        checkout scm
        sh "APP=crane bash ./bin/herokutor.sh `pwd`"
    }
}
