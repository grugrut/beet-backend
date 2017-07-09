def err_msg = ""

node {
    try {
        def root = tool name: 'Go1.8', type: 'go'
        ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/grugrut/beet-backend") {
            withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
                env.PATH="${GOPATH}/bin:$PATH"
                
                stage 'Checkout'
            
                git url: 'https://github.com/grugrut/beet-backend.git'
            
                stage 'preTest'
                sh 'go version'
                sh 'go get -u github.com/golang/dep/...'
                sh 'dep init'
                
                stage 'Test'
                sh 'go vet .'
                
                stage 'Build'
                sh 'go build -o beet .'

                stage 'Deploy'
                withCredentials([string(credentialsId: 'DEPLOY_PATH', variable: 'DEPLOY_PATH')]) {
                    sh 'cp -fp ${WORKSPACE}/src/beet ${DEPLOY_PATH}/bin/'
                }
                sh '/sbin/service beet restart'
            }
        }
    } catch (e) {
        err_msg = "${e}"
        currentBuild.result = "FAILURE"
    } finally {
        if (currentBuild.result != "FAILURE") {
            currentBuild.result = "SUCCESS"
        }
        notify(err_msg)
    }
}

def notify(msg) {
    def detail_link = "(<${env.BUILD_URL}|Open>)"
    def slack_color = "good"
    if(currentBuild.result == "FAILURE") {
        slack_color = "danger"
    }
    def slack_msg = "job ${env.JOB_NAME}[No.${env.BUILD_NUMBER}] was builded ${currentBuild.result}. ${detail_link}\n\n${msg}"
    slackSend color: "${slack_color}", message: "${slack_msg}"
}
