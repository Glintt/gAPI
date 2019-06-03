node{
	env.DB = "mongo".trim()
	env.LOGS_TYPE = "Elastic".trim()
	env.QUEUE_TYPE = "Rabbit".trim()
    env.DOCKER_USER = "${env.DOCKER_USER}"

	stage('Clone repository') {
		checkout scm
	}

	last_version = sh (
      returnStdout: true,
      script: 'git describe --abbrev=0 --tag'
    ).trim()

	env.BUILD_VERSION_NAME = "$last_version"
	if (params.close_version == false) {
		env.BUILD_VERSION_NAME = env.BUILD_VERSION_NAME + ".$BUILD_NUMBER"
	}
	
	stage('Build docker images') {
		dir('api') {
			sh "docker image build --build-arg db=$DB --build-arg logs_type=$LOGS_TYPE --build-arg queue_type=$QUEUE_TYPE -t $DOCKER_USER/gapi-backend:$BUILD_VERSION_NAME -f Dockerfile ."
			sh "docker image build --build-arg db=$DB --build-arg logs_type=$LOGS_TYPE --build-arg queue_type=$QUEUE_TYPE -t $DOCKER_USER/gapi-backend:latest -f Dockerfile ."
		}

		dir('api') {
			sh "docker image build --build-arg db=$DB --build-arg logs_type=$LOGS_TYPE --build-arg queue_type=Rabbit -t $DOCKER_USER/gapi-rabbitlistener:$BUILD_VERSION_NAME -f Dockerfile-rabbitlistener ."
			sh "docker image build --build-arg db=$DB --build-arg logs_type=$LOGS_TYPE --build-arg queue_type=Rabbit -t $DOCKER_USER/gapi-rabbitlistener:latest -f Dockerfile-rabbitlistener ."
		}

		dir('dashboard') {
			sh "docker image build -t $DOCKER_USER/gapi-dashboard:$BUILD_VERSION_NAME ."
			sh "docker image build -t $DOCKER_USER/gapi-dashboard:latest ."
		}
	}

	stage('Publish docker images') {
		sh "docker push $DOCKER_USER/gapi-backend:$BUILD_VERSION_NAME"
		sh "docker push $DOCKER_USER/gapi-backend:latest"
		
		sh "docker push $DOCKER_USER/gapi-rabbitlistener:$BUILD_VERSION_NAME"
		sh "docker push $DOCKER_USER/gapi-rabbitlistener:latest"
		
		sh "docker push $DOCKER_USER/gapi-dashboard:$BUILD_VERSION_NAME"
		sh "docker push $DOCKER_USER/gapi-dashboard:latest"
	}

	stage('Remove docker images from build machine') {		
		sh "docker image rm -f $DOCKER_USER/gapi-backend:$BUILD_VERSION_NAME"
		sh "docker image rm -f $DOCKER_USER/gapi-backend:latest"
				
		sh "docker image rm -f $DOCKER_USER/gapi-rabbitlistener:$BUILD_VERSION_NAME"
		sh "docker image rm -f $DOCKER_USER/gapi-rabbitlistener:latest"
				
		sh "docker image rm -f $DOCKER_USER/gapi-dashboard:$BUILD_VERSION_NAME"
		sh "docker image rm -f $DOCKER_USER/gapi-dashboard:latest"
	}
}