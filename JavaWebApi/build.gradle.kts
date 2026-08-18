plugins {
    java
    id("org.springframework.boot") version "4.1.0"
    id("org.graalvm.buildtools.native") version "1.1.9"
}

group = "com.laixhe"
version = "0.0.1"

java {
    toolchain {
        // 使用本机 GraalVM 25 (JDK 25) 工具链
        languageVersion = JavaLanguageVersion.of(25)
    }
}

repositories {
    mavenCentral()
}

dependencies {
    // Spring Boot 4.1 BOM 统一管理版本
    implementation(platform("org.springframework.boot:spring-boot-dependencies:4.1.0"))

    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("org.springframework.boot:spring-boot-starter-security")
    implementation("org.springframework.boot:spring-boot-starter-validation")
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")

    // OpenAPI 文档: 由代码/注解动态生成, 不再维护静态 swagger.yaml (3.x 支持 Spring Boot 4)
    implementation("org.springdoc:springdoc-openapi-starter-webmvc-ui:3.1.0")

    // JWT (jjwt 是目前最主流的 Java JWT 库)
    implementation("io.jsonwebtoken:jjwt-api:0.13.0")
    runtimeOnly("io.jsonwebtoken:jjwt-impl:0.13.0")
    runtimeOnly("io.jsonwebtoken:jjwt-jackson:0.13.0")

    // 接口限流 (Bucket4j 是目前最主流的 Java 令牌桶限流库)
    implementation("com.bucket4j:bucket4j_jdk17-core:8.19.0")

    // Lombok: 减少样板代码
    compileOnly("org.projectlombok:lombok:1.18.46")
    annotationProcessor("org.projectlombok:lombok:1.18.46")

    // 数据库驱动: MySQL (生产) + H2 (本地/测试免安装)
    runtimeOnly("com.mysql:mysql-connector-j")
    runtimeOnly("com.h2database:h2")

    testCompileOnly("org.projectlombok:lombok:1.18.46")
    testAnnotationProcessor("org.projectlombok:lombok:1.18.46")
    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("org.springframework.boot:spring-boot-starter-webmvc-test")
}

tasks.withType<Test> {
    useJUnitPlatform()
}

// GraalVM Native Image 支持: ./gradlew nativeCompile 可产出原生可执行文件
graalvmNative {
    binaries {
        named("main") {
            imageName.set("webapi")
        }
    }
}
