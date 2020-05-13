node {
	stage('代码检出'){
	    git credentialsId: '51048f57-c4b0-40c0-a83b-61586069665c', url: 'http://192.168.1.173/hqj/hqj.git'
	}

	stage('mvn 打包'){
		def POM_PATH = "project/euerka/pom.xml"
		sh "mvn clean package -Dautoconfig.skip=true -Dmaven.test.skip=true -f ${POM_PATH}"
	}

	stage('docker build'){
		def EUERKA_PORT = 8761
		def IMAGE = "192.168.1.172/hqj-dev/euerka:v1.${BUILD_NUMBER}"
		def BUILD_ARG = "--build-arg JAR_FILE=project/${JOB_NAME}/target/euerka-server.jar --build-arg SERVER_PORT=${EUERKA_PORT} ."
		docker.withRegistry("http://192.168.1.172" , 'ded20f3c-0eed-481e-97cb-f043cc22affc') {
		    def CUSTOMIMAGE = docker.build(
		        "${IMAGE}", "${BUILD_ARG}  -f ./Dockerfile"
		        )
			CUSTOMIMAGE.push()
		}
	}
	stage('发布应用'){
		sshPublisher(
			publishers: [
				sshPublisherDesc(configName: 'fz_server_175', 
					transfers: [
						sshTransfer(cleanRemote: false,
							execCommand: '''
                                echo "正在发布..."
                                IMAGE_ID=`docker ps -a |grep ${JOB_NAME} | awk 'NR==1''{print $1}'` ; \
                                DOCKER_CMD="docker run -d  --name=${JOB_NAME} 192.168.1.172/hqj-dev/${JOB_NAME}:v1.${BUILD_NUMBER}" ; \
                                if [[ $IMAGE_ID ]];then docker rm -f $IMAGE_ID &&  $DOCKER_CMD ; \
                                else $DOCKER_CMD; \
                                fi
                                sleep 5
                                echo "正在检测容器状态: ${JOB_NAME}" ; \
                                IMAGE_ID=`docker ps |grep ${JOB_NAME} |awk '{print $1}'` ;\
                                if [[ $IMAGE_ID ]];then echo "${JOB_NAME} 已发布" ;\
                                else ERROR_MESSAGE=`docker logs ${JOB_NAME}` ;\
                                echo ${JOB_NAME} no running! \
                                echo "容器启动失败，回滚中..." ; \
                                ROLLBACK_IMAGE_ID=`expr ${BUILD_NUMBER} - 1` ; \
                                ROLLBACK_DOCKER="docker run -d  --name=${JOB_NAME} 192.168.1.172/hqj-dev/${JOB_NAME}:v1.$ROLLBACK_IMAGE_ID" ; \
                                docker rm -f ${JOB_NAME} && $ROLLBACK_DOCKER ;\
                                echo "回退后版本号： 192.168.1.172/hqj-dev/${JOB_NAME}:v1.$ROLLBACK_IMAGE_ID" ;\
                                fi
							''', 
							execTimeout: 600000
						)
					]
				)
		    ]
		)
	}
}