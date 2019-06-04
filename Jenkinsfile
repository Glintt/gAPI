node{
	env.DB = "mongo".trim()
	env.LOGS_TYPE = "Elastic".trim()
	env.QUEUE_TYPE = "Rabbit".trim()
    env.DOCKER_USER = "${env.DOCKER_USER}"
	env.IMAGE_NAME = "gapi"

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

	if (params.version != "") {
		env.BUILD_VERSION_NAME = params.version
	}

	if (params.image_name != "") {
		env.IMAGE_NAME = params.image_name
	}
	
	stage('Build docker images') {
		dir('api') {
			sh "docker image build --build-arg db=$DB --build-arg logs_type=$LOGS_TYPE --build-arg queue_type=$QUEUE_TYPE -t $DOCKER_USER/$IMAGE_NAME-backend:$BUILD_VERSION_NAME -f Dockerfile ."
			sh "docker image build --build-arg db=$DB --build-arg logs_type=$LOGS_TYPE --build-arg queue_type=$QUEUE_TYPE -t $DOCKER_USER/$IMAGE_NAME-backend:latest -f Dockerfile ."
		}

		dir('api') {
			sh "docker image build --build-arg db=$DB --build-arg logs_type=$LOGS_TYPE --build-arg queue_type=Rabbit -t $DOCKER_USER/$IMAGE_NAME-rabbitlistener:$BUILD_VERSION_NAME -f Dockerfile-rabbitlistener ."
			sh "docker image build --build-arg db=$DB --build-arg logs_type=$LOGS_TYPE --build-arg queue_type=Rabbit -t $DOCKER_USER/$IMAGE_NAME-rabbitlistener:latest -f Dockerfile-rabbitlistener ."
		}

		dir('dashboard') {
			sh "docker image build -t $DOCKER_USER/$IMAGE_NAME-dashboard:$BUILD_VERSION_NAME ."
			sh "docker image build -t $DOCKER_USER/$IMAGE_NAME-dashboard:latest ."
		}
	}

	stage('Publish docker images') {
		sh "docker push $DOCKER_USER/$IMAGE_NAME-backend:$BUILD_VERSION_NAME"
		sh "docker push $DOCKER_USER/$IMAGE_NAME-backend:latest"
		
		sh "docker push $DOCKER_USER/$IMAGE_NAME-rabbitlistener:$BUILD_VERSION_NAME"
		sh "docker push $DOCKER_USER/$IMAGE_NAME-rabbitlistener:latest"
		
		sh "docker push $DOCKER_USER/$IMAGE_NAME-dashboard:$BUILD_VERSION_NAME"
		sh "docker push $DOCKER_USER/$IMAGE_NAME-dashboard:latest"
	}

	stage('Remove docker images from build machine') {		
		sh "docker image rm -f $DOCKER_USER/$IMAGE_NAME-backend:$BUILD_VERSION_NAME"
		sh "docker image rm -f $DOCKER_USER/$IMAGE_NAME-backend:latest"
				
		sh "docker image rm -f $DOCKER_USER/$IMAGE_NAME-rabbitlistener:$BUILD_VERSION_NAME"
		sh "docker image rm -f $DOCKER_USER/$IMAGE_NAME-rabbitlistener:latest"
				
		sh "docker image rm -f $DOCKER_USER/$IMAGE_NAME-dashboard:$BUILD_VERSION_NAME"
		sh "docker image rm -f $DOCKER_USER/$IMAGE_NAME-dashboard:latest"
	}
}