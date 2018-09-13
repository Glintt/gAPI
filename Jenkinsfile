node{
	env.BUILD_VERSION = "${params.version}".trim()
    env.DOCKER_USER = "${env.DOCKER_USER}"

	stage('Clone repository') {
		checkout scm
	}

	stage('Build docker images') {
		dir('api') {
			sh "docker image build -t $DOCKER_USER/gapi-backend:$BUILD_VERSION -f Dockerfile ."
			sh "docker image build -t $DOCKER_USER/gapi-backend -f Dockerfile ."
		}

		dir('api') {
			sh "docker image build -t $DOCKER_USER/gapi-rabbitlistener:$BUILD_VERSION -f Dockerfile-rabbitlistener ."
			sh "docker image build -t $DOCKER_USER/gapi-rabbitlistener -f Dockerfile-rabbitlistener ."			
		}

		dir('dashboard') {
			sh "docker image build -t $DOCKER_USER/gapi-dashboard:$BUILD_VERSION ."
			sh "docker image build -t $DOCKER_USER/gapi-dashboard ."
		}
	}

	stage('Publish docker images') {
		sh "docker push $DOCKER_USER/gapi-backend:$BUILD_VERSION"
		sh "docker push $DOCKER_USER/gapi-backend"
		
		sh "docker push $DOCKER_USER/gapi-rabbitlistener:$BUILD_VERSION"
		sh "docker push $DOCKER_USER/gapi-rabbitlistener"
		
		sh "docker push $DOCKER_USER/gapi-dashboard:$BUILD_VERSION"
		sh "docker push $DOCKER_USER/gapi-dashboard"	
	}

	stage('Remove docker images from build machine') {		
		sh "docker image rm -f $DOCKER_USER/gapi-backend:$BUILD_VERSION"
		sh "docker image rm -f $DOCKER_USER/gapi-backend"
				
		sh "docker image rm -f $DOCKER_USER/gapi-rabbitlistener:$BUILD_VERSION"
		sh "docker image rm -f $DOCKER_USER/gapi-rabbitlistener"
				
		sh "docker image rm -f $DOCKER_USER/gapi-dashboard:$BUILD_VERSION"
		sh "docker image rm -f $DOCKER_USER/gapi-dashboard"
	}
}