FROM woshitiancai/openjdk8:v1.0.0

MAINTAINER JENKINS

ARG JAR_FILE

ARG SERVER_PORT

WORKDIR /home/java

COPY  ${JAR_FILE} app.jar

ENTRYPOINT ["java", "-XX:+UnlockExperimentalVMOptions", "-XX:+UseCGroupMemoryLimitForHeap", "-Djava.security.egd=file:/dev/./urandom",  "-jar", "app.jar", "--eureka.client.serviceUrl.defaultZone=http://hqj:hqj_pass@hqj-euerka:8761/eureka/", "--eureka.instance.prefer-ip-address=false","--eureka.instance.instance-id=${spring.cloud.client.ipAddress}:${server.port}"]

EXPOSE ${SERVER_PORT}

