node {
    stage('Prepare') {
        checkout scm
        sh "make clean && make prepare"
    }
    stage('Build') {
        checkout scm
        sh "./scripts/ci-test.sh all"
    }
    stage('Lint') {
        checkout scm
        sh "./scripts/ci-test.sh lint"
    }
    stage('Deploy into k8s') {
        checkout scm
        def OUT_LOG = sh(script: 'mktemp', returnStdout: true).trim()
        def OUT_LOG_COLOR = sh(script: 'mktemp', returnStdout: true).trim()
        sh "APP=crane OUT_LOG=$OUT_LOG OUT_LOG_COLOR=$OUT_LOG_COLOR bash ./scripts/herokutor.sh `pwd`"
        def color = readFile OUT_LOG_COLOR
        def message = readFile OUT_LOG
        slackSend color: color, message: message
    }
}
