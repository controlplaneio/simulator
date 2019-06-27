#!/usr/bin/env groovy

def getDockerImageTag() {
  if (env.GIT_COMMIT == "") {
    error "GIT_COMMIT value was empty at usage. "
  }
  return "${env.BUILD_ID}-${env.GIT_COMMIT}"
}

pipeline {
  agent none

  environment {
    ENVIRONMENT = 'ops'
    DOCKER_IMAGE_TAG = "${getDockerImageTag()}"
  }

  stages {

    stage('Test') {
      agent {
        docker { image 'docker.io/controlplane/gcloud-sdk:latest' }
      }

      options {
        timeout(time: 15, unit: 'MINUTES')
        retry(1)
        timestamps()
      }

      steps {
        ansiColor('xterm') {
          sh 'make test'
        }
      }
    }
  }
}
