node {
    def EUERKA_PORT = 8061
    def IMAGE = "192.168.1.172/hqj-dev/config-server:v1.k8s.${BUILD_NUMBER}"
	stage('checkout code'){
	    git credentialsId: '51048f57-c4b0-40c0-a83b-61586069665c', url: 'http://192.168.1.173/hqj/hqj.git'
	}
    
    stage('package and scan'){
        withSonarQubeEnv(credentialsId: 'sonarqube') {
            def POM_PATH = "project/config/pom.xml"
            sh "mvn clean package -Dautoconfig.skip=true -Dmaven.test.skip=true -f ${POM_PATH} sonar:sonar  "
        }
	}
	
    stage("sonarqube quality testing"){
        timeout(time: 1, unit: 'HOURS') {
            def qg = waitForQualityGate()
            if (qg.status != 'OK') {
                error "未通过sonarqube质量检测，流程停止: ${qg.status}"
            }
        }
    }
    
	stage('docker building'){
		def BUILD_ARG = "--build-arg JAR_FILE=project/config/target/config-server.jar --build-arg SERVER_PORT=${EUERKA_PORT} ."
		docker.withRegistry("http://192.168.1.172" , 'ded20f3c-0eed-481e-97cb-f043cc22affc') {
            def CUSTOMIMAGE = docker.build("${IMAGE}", "${BUILD_ARG}  -f ./Dockerfile")
			CUSTOMIMAGE.push()
		}
	}
    stage('deployment app to k8s'){
		sshPublisher(
			publishers: [
				sshPublisherDesc(configName: 'fz_server_175', 
					transfers: [
						sshTransfer(cleanRemote: false,
							execCommand: '''
                                IMAGE=192.168.1.172/hqj-dev/config-server:v1.k8s.${BUILD_NUMBER} ;\
                                source .bashrc ;\
                                kubectl set image deployment/hqj-config-server hqj-config-server="$IMAGE" --record -n wwdt ;\
                                kubectl rollout status deployment/hqj-config-server -n wwdt
							''', 
							execTimeout: 1200000
						)
					]
				)
		    ]
		)
	}
	stage('automated testing'){
		sh ('''curl -o index.html http://192.168.1.171:8080/api/open/run_auto_test?id=23\\&token=bff99cca744a0955e0f3c8b60e131c04de6db0abfaef3d88d0f6134bebda26a5\\&mode=html\\&email=false\\&download=false''')
		publishHTML([
		    allowMissing: false, 
		    alwaysLinkToLastBuild: false, includes: 'index.html', 
		    keepAll: true, reportDir: '', reportFiles: 'index.html', 
		    reportName: 'HTML Report', reportTitles: '接口测试报告'
		    ]
        )
	}
}