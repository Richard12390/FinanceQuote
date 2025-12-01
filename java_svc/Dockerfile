# Build
FROM maven:3.9-eclipse-temurin-21 AS builder
WORKDIR /app
COPY pom.xml mvnw ./
COPY .mvn ./.mvn
RUN ./mvnw -q dependency:go-offline
COPY . .
RUN ./mvnw -DskipTests package

# Runtime
FROM eclipse-temurin:21-jre
WORKDIR /app
COPY --from=builder /app/target/*.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java","-jar","/app/app.jar"]