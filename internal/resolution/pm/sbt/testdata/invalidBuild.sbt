name := "invalid-project"
version := "1.0.0"

libraryDependencies ++= Seq(
    "org.scala-lang" % "scala-library" % "2.13.8",
    // Missing closing parenthesis
    "com.typesafe.akka" %% "akka-http" % "10.2.9"
    "com.typesafe.akka" %% "akka-stream" % "10.2.9"
)

// Invalid syntax
scalaVersion = "2.13.8"