pipeline {
    agent {
        node {
            label 'master'
            customWorkspace "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/grugrut/beet-backend"
        }
    }
    tools {
        go 'Go1.8'
    }
    environment {
        GOPATH="${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/"
    }
        
    stages {
        stage ('Checkout') {
            steps {
                git url: 'https://github.com/grugrut/beet-backend.git'
            }
        }
        stage ('preTest') {
            steps {
                sh 'go version'
                sh 'go get -u github.com/golang/dep/...'
                sh '${GOPATH}/bin/dep init'
                sh 'go get github.com/golang/lint/golint'

            }
        }
        stage ('Test') {
            steps {
                sh 'go vet ./...'
                sh '${GOPATH}/bin/golint ./...'
                stepcounter settings: [[encoding: 'UTF-8', filePattern: '**/*.go', filePatternExclude: 'vendor/**/*.go', key: 'Go']]
            }
        }
        stage ('Build') {
            steps {
                sh 'go build -o beet .'
            }
        }
        stage ('Deploy') {
            when {
                branch 'master'
            }
            steps {
                withCredentials([string(credentialsId: 'DEPLOY_PATH', variable: 'DEPLOY_PATH')]) {
                    sh 'cp -fp ${WORKSPACE}/beet ${DEPLOY_PATH}/bin/'
                }
                sh 'sudo /sbin/service beet restart'
            }
        }
    }
    post {
        success {
            notifySlack("good")
        }
        failure {
            notifySlack("danger")
        }
    }
}

def notifySlack(color) {
    slackSend color: color, message: "job ${env.JOB_NAME}[No.${env.BUILD_NUMBER}] was builded ${currentBuild.result}. (<${env.BUILD_URL}|Open>)"
}
