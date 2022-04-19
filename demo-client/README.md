# Demo Client

Demo client using Gradle/Java and supply chain gateway as the maven central source for public library access.

## Requirements

* Java / JDK 11

## Usage

```bash
./gradlew assemble --refresh-dependencies
```

## Observations

Failure to fetch vulnerable `log4j` dependency

```
> Could not resolve all files for configuration ':app:compileClasspath'.
   > Could not resolve org.apache.logging.log4j:log4j:2.16.0.
     Required by:
         project :app
      > Could not resolve org.apache.logging.log4j:log4j:2.16.0.
         > Could not get resource 'http://localhost:10000/maven2/org/apache/logging/log4j/log4j/2.16.0/log4j-2.16.0.pom'.
            > Could not GET 'http://localhost:10000/maven2/org/apache/logging/log4j/log4j/2.16.0/log4j-2.16.0.pom'. Received status code 403 from server: Forbidden
```
