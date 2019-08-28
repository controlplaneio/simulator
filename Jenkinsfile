#!/usr/bin/env groovy

def getDockerImageTag() {
  if (env.GIT_COMMIT == "") {
    error "GIT_COMMIT value was empty at usage. "
  }
  return "${env.BUILD_ID}-${env.GIT_COMMIT}"
}

pipeline {
  agent none

    post {
      failure {
        emailext (
            subject: "Simulator build on '${env.GIT_BRANCH} branch failed",
            body: "<p>Simulator build failed.</p> <ul><li>Branch: '${env.GIT_BRANCH}'</li><li>Commit: '${env.GIT_COMMIT}'</li><li>Build No: '${env.BUILD_NUMBER}'</li><ul><p>${env.CHANGES}</p>",
            to: "kubernetes-simulator-build-notifications@googlegroups.com",
            from: "jenkins@control-plane.io"
            )
      }
    }

  environment {
    ENVIRONMENT = 'ops'
    DOCKER_IMAGE_TAG = "${getDockerImageTag()}"
    AWS_REGION = "eu-west-1"
    AWS_DEFAULT_REGION = "eu-west-1"
  }

  stages {

    stage('Test') {
      agent {
        docker {
          image 'docker.io/controlplane/gcloud-sdk:latest'
            args '-v /var/run/docker.sock:/var/run/docker.sock ' +
            '--user=root ' +
            '--cap-drop=ALL ' +
            '--cap-add=DAC_OVERRIDE'
        }
      }

      options {
        timeout(time: 15, unit: 'MINUTES')
          retry(1)
          timestamps()
      }

      steps {
        ansiColor('xterm') {
          sh 'env && make docker-test'
        }
      }
    }
  }
}
