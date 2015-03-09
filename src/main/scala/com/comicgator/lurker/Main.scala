package com.comicgator.lurker

import scala.io.{BufferedSource, Source}
import org.slf4j.LoggerFactory
import org.json4s._
import org.json4s.native.JsonMethods._

// TODO: Split into separate files

/**
 * Main object reads the initialization file comics.json, extracts them to comic objects. It takes a
 */

object Main {
  /**
   * Logging using logback settings in lurker/src/main/resources/logback.xml
   * If err try
   * StatusPrinter.print((LoggerFactory.getILoggerFactory).asInstanceOf[LoggerContext])
   * @return
   */
  def logger = LoggerFactory.getLogger(this.getClass)
  val usage = """Usage: [--no-delta] [<code>...]
                 |--no-delta
                 |    Runs ETL from first page of comic strip.
                 |    Default is to run from last strip saved to the database.
                 |<code>...
                 |    List of target comics to run etl against. Use "name" from comics.json.
                 |    Default is to run all comics.""".stripMargin
  def main(args: Array[String]): Unit = {
    val arglist = args.toList
    logger.info("Spinning up Lurker")
    /**
     * parseOptions: a closure that recursively examines the args list. It accepts defaults for the delta boolean and
     * comics list which are then updated and returned.
     * @param args List of command line arguments
     * @param delta Boolean default true starts etl from last strip scraped otherwise starts at the first strip.
     * @param comicCodes List of comics to be scraped described by their code name. If returned empty indicates that
     *                   all comics are to be scraped.
     * @return Tuple of augmented delta boolean and comic list.
     */
    def parseOptions(args: List[String], delta: Boolean, comicCodes: List[String]): (Boolean, List[String]) = {
      args match {
        case Nil => (delta, comicCodes)
        case ("-h" | "--help") :: tail =>
          println(usage)
          sys.exit(0)
        case "--no-delta" :: tail =>
          parseOptions(tail, delta = false, comicCodes)
        case value :: tail =>
          parseOptions(tail, delta, value :: comicCodes)
      }
    }
    val comicMap = loadComics()
    val (delta: Boolean, comicCodes: List[String]) = parseOptions(arglist, delta = true, List[String]())
    // Empty comicCodes implies that all comics are targeted to be scraped.
    val comics: List[Comic] = if(comicCodes.isEmpty) comicMap.values.toList else comicCodes flatMap comicMap.get
    logger.info(s"Running ${if(delta) "delta etl" else "full etl"} for comics ${comics.mkString(", ")}")
    etl(delta, comics)
  }

  /**
   * Load comics.json file into List of Comic objects.
   * @return Map keying on the code name of the comic to the Comic object.
   */
  def loadComics(): Map[String, Comic] = {
    implicit val formats = DefaultFormats // Brings in default date formats etc.
    val file: BufferedSource = Source.fromFile("comics.json")
    val json: JValue = parse(file.mkString)
    val comics: List[Comic] = json.extract[List[Comic]]
    val comicMap: Map[String, Comic] = comics.map(c => (c.code, c)).toMap
    comicMap
  }

  /**
   * Extraction Transform and Load control flow for the scraping of the comic websites.
   * @param delta Boolean indicating the starting point of the scrape, true means start at the last strip scraped,
   *              false means start from the beginning
   * @param comics List of targeted Comic objects to be scraped.
   */
  def etl(delta: Boolean, comics: List[Comic]): Unit = {
    logger.debug("this is where I etl")
  }
}

class Comic(val code: String, url: String, title: String, banner: String, start: String, next: Parser, image: Parser, titleParser: Parser,
            bonus: Parser, alt: Parser) {
  override def toString: String = s"$code"
}
case class Parser(method: String, patterns: Array[String])