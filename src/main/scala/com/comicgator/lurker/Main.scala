package com.comicgator.lurker

import scala.io.{BufferedSource, Source}
import org.json4s._
import org.json4s.native.JsonMethods._

// Keep everything in one file until it bursts
// TODO: Segment into own files
object Main {
  def main(args: Array[String]): Unit = {
    val comics = loadComics()
    etl(comics)
  }

  def loadComics(): Array[Comic] = {
    implicit val formats = DefaultFormats // Brings in default date formats etc.
    val file: BufferedSource = Source.fromFile("comics.json")
    val json: JValue = parse(file.mkString)
    val comics: Array[Comic] = json.extract[Array[Comic]]
    comics
  }

  def etl(comics: Array[Comic]): Unit = {
    for (comic <- comics) {
      println(comic.toString())
    }
  }
}

class Comic(url: String, headline: String, banner: String, start: String,
            next: Parser, image: Parser, title: Parser, bonus: Parser, alt: Parser) {
  override def toString: String =
    s"(url=$url, headline=$headline, banner=$banner, start=$start, title=$title)"

}
case class Parser(method: String, patterns: Array[String])