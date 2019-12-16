FROM openjdk:8
COPY . /usr/src
WORKDIR /usr/src
RUN javac -cp lib/gson-2.8.6.jar *.java */**.java
CMD ["java", "-cp", ".:lib/gson-2.8.6.jar", "Entry_LALBFunction"]
