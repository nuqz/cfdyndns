pipeline {
    agent {
        docker { image 'golang:alpine' }
    }
    stages {
        stage('Test') {
            steps {
                sh 'go version'
            }
        }
    }
}
