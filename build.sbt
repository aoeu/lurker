name := "lurker"

organization := "com.comicgator"

version := "0.0.1"

scalaVersion := "2.11.6"

libraryDependencies ++= Seq(
  "org.scalatest" % "scalatest_2.11" % "2.2.4" % "test" withSources() withJavadoc(),
  "org.scalacheck" %% "scalacheck" % "1.12.1" % "test" withSources() withJavadoc(),
  "org.json4s" %% "json4s-native" % "3.2.11",
  "ch.qos.logback" % "logback-classic" % "1.1.2",
  "org.scalaj" %% "scalaj-http" % "1.1.4"
)

initialCommands := "import com.comicgator.lurker._"

