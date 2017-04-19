node { 
    stage('Dependency') {
        checkout scm
        sh "git submodule update --init --recursive"
        sh "rm -rf ./src/services/tmp/src"
        sh "make clean"
        sh "make restore"
    }
    stage('Codegen') {
        checkout scm
        sh "make codegen"
        sh "make mockentity"
    }
    stage('Build') {
        checkout scm
        sh "make all"
    }
    stage('Lint') {
        checkout scm
        sh "make lint"
    }
    stage('Test') {
        checkout scm
        sh "docker rm -f redis-jenkins-service-test || true"
        sh "docker run -d --name redis-jenkins-test -p 6379:6379 redis:alpine"
        sh "make test"
        sh "docker rm -f redis-jenkins-test"
    }
}
