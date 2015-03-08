name := "lurker"

organization := "com.comicgator"

version := "0.0.1"

scalaVersion := "2.11.6"

libraryDependencies ++= Seq(
  "org.scalatest" % "scalatest_2.11" % "2.2.4" % "test" withSources() withJavadoc(),
  "org.scalacheck" %% "scalacheck" % "1.12.1" % "test" withSources() withJavadoc(),
  "org.json4s" %% "json4s-native" % "3.2.11"
)

initialCommands := "import com.comicgator.lurker._"

