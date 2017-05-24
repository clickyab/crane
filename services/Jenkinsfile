node { 
    stage('Test') {
        checkout scm
        sh "docker rm -f redis-jenkins-service-test || true"
        sh "docker run -d --name redis-jenkins-service-test -p 6379:6379 redis:alpine"
        sh "GOPATH=`mktemp -d` make -f ./Makefile.mk services_test"
        sh "docker rm -f redis-jenkins-service-test"
    }
}
