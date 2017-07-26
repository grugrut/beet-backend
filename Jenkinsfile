pipeline {
    agent {
        node {
            customWorkspace "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/grugrut/beet-backend"
        }
    }
    tools {
        Go 'Go1.8'
    }
    environment {
        GOROOT="${root}"
        GOPATH="${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/"
        PATH="${GOROOT}/bin:${GOPATH}/bin:$PATH"
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
                sh 'dep init'
            }
        }
        stage ('Test') {
            steps {
                sh 'go vet .'
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
            detail_link = "(<${env.BUILD_URL}|Open>)"
            slack_color = "good"
            slack_msg = "job ${env.JOB_NAME}[No.${env.BUILD_NUMBER}] was builded ${currentBuild.result}. ${detail_link}"
            slackSend color: "${slack_color}", message: "${slack_msg}"
        }
        failure {
            detail_link = "(<${env.BUILD_URL}|Open>)"
            slack_color = "danger"
            slack_msg = "job ${env.JOB_NAME}[No.${env.BUILD_NUMBER}] was builded ${currentBuild.result}. ${detail_link}"
            slackSend color: "${slack_color}", message: "${slack_msg}"
        }
    }
}
