node{
	env.BUILD_VERSION = "${params.version}".trim()
    env.DOCKER_USER = "${env.DOCKER_USER}"

	stage('Clone repository') {
		checkout scm
	}

	stage('Build docker images') {
		dir('api') {
			sh "docker image build -t $DOCKER_USER/gAPI-backend:$BUILD_VERSION -f Dockerfile ."
			sh "docker image build -t $DOCKER_USER/gAPI-backend -f Dockerfile ."
		}

		dir('api') {
			sh "docker image build -t $DOCKER_USER/gAPI-rabbitlistener:$BUILD_VERSION -f Dockerfile-rabbitlistener ."
			sh "docker image build -t $DOCKER_USER/gAPI-rabbitlistener -f Dockerfile-rabbitlistener ."			
		}

		dir('dashboard') {
			sh "docker image build -t $DOCKER_USER/gAPI-dashboard:$BUILD_VERSION ."
			sh "docker image build -t $DOCKER_USER/gAPI-dashboard ."
		}
	}

	stage('Publish docker images') {
		sh "docker push $DOCKER_USER/gAPI-backend:$BUILD_VERSION"
		sh "docker push $DOCKER_USER/gAPI-backend"
		
		sh "docker push $DOCKER_USER/gAPI-rabbitlistener:$BUILD_VERSION"
		sh "docker push $DOCKER_USER/gAPI-rabbitlistener"
		
		sh "docker push $DOCKER_USER/gAPI-dashboard:$BUILD_VERSION"
		sh "docker push $DOCKER_USER/gAPI-dashboard"	
	}

	stage('Remove docker images from build machine') {		
		sh "docker image rm -f $DOCKER_USER/gAPI-backend:$BUILD_VERSION"
		sh "docker image rm -f $DOCKER_USER/gAPI-backend"
				
		sh "docker image rm -f $DOCKER_USER/gAPI-rabbitlistener:$BUILD_VERSION"
		sh "docker image rm -f $DOCKER_USER/gAPI-rabbitlistener"
				
		sh "docker image rm -f $DOCKER_USER/gAPI-dashboard:$BUILD_VERSION"
		sh "docker image rm -f $DOCKER_USER/gAPI-dashboard"
	}
}