package com.comicgator.lurker

import org.json4s.{DefaultFormats, _}
import org.json4s.native.JsonMethods._
import org.json4s.native.Serialization.write
import org.slf4j.LoggerFactory
import scalaj.http._

import scala.io.{BufferedSource, Source}
// TODO: Split into separate files

/**
 * Main does all the things
 */
object Main {
  /**
   * Logging using logback settings in lurker/src/main/resources/logback.xml
   * If trouble try
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
    val argList = args.toList
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
    val (delta: Boolean, comicCodes: List[String]) = parseOptions(argList, delta = true, List[String]())
    // Empty comicCodes implies that all comics are targeted to be scraped.
    val comics: List[Comic] = if(comicCodes.isEmpty) comicMap.values.toList else comicCodes flatMap comicMap.get
    logger.info(s"Running${if(delta) " delta " else " "}ETL for comics: ${comics.mkString(", ")}")
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
    val comicMap: Map[String, Comic] = comics.map(c => (c.id, c)).toMap
    comicMap
  }

  /**
   * Extraction Transform and Load control flow for the scraping of the comic websites.
   * @param delta Boolean indicating the starting point of the scrape, true means start at the last strip scraped,
   *              false means start from the beginning
   * @param comics List of targeted Comic objects to be scraped.
   */
  def etl(delta: Boolean, comics: List[Comic]): Unit = {
    for(comic <- comics) {
      logger.info(s"Scraping comic: $comic")
      comic.save()
    }
  }
}


/**
 * Comic Class representing comic meta data and strategies for parsing strips from comic website.
 * @param id String predetermined in comics.json configuration file.
 * @param hostname String usually the domain and tld of the comic website, useful for matching or creating urls.
 * @param title String represents the general name of the comic.
 * @param creator String the name of the comic creator.
 * @param banner_image String location of the pre-generated image representing the comic in-app.
 * @param start_url String a url of the very first strip of the comic. ETL starts here if delta is false.
 * @param next_parser Parser strategy for getting the next endpoint from a strip page.
 * @param image_parser Parser strategy for getting the url of the image from a strip page.
 * @param title_parser Parser strategy for getting the title of the comic from the strip page.
 * @param bonus_image_parser Parser strategy for getting the bonus image from the strip page.
 * @param alt_text_parser Parser strategy for getting the alt_text from the strip page.
 */
class Comic(val id: String, hostname: String, title: String, creator: String, banner_image: String, start_url: String,
            next_parser: Parser, image_parser: Parser, title_parser: Parser, bonus_image_parser: Parser,
            alt_text_parser: Parser) {
  def logger = LoggerFactory.getLogger(this.getClass)
  /* underscores used for json import, recast them for consistency */
  private val bannerImage = banner_image
  private val startUrl = start_url
  val nextParser = next_parser
  val imageParser = image_parser
  val titleParser = title_parser
  val bonusImageParser = bonus_image_parser
  val altTextParser = alt_text_parser

  /**
   * Exports relevant variables to a json string suitable for http post/patch requests.
   * @return JSON formatted string
   */
  def export = {
    val jsonMap = Map("id" -> id, "hostname" -> hostname, "title" -> title, "creator" -> creator,
    "banner_image" -> bannerImage, "image_parser" -> imageParser.export)
    // Passing implicit DefaultFormats
    write(jsonMap)(DefaultFormats)
  }

  /**
   * Post/Patch Comic to Maestro and into the database.
   */
  def save() = {
    logger.info(s"Saving comic ${this.id}")
    logger.debug(this.export)
    val url = "http://192.168.34.10:3000/comic"
    val postResponse: HttpResponse[String] = Http(url).
      method("POST").
      auth("c11z", "c11z").
      header("User-Agent", "ComicGator (http://github.com/comicgator/lurker").
      postData(this.export).asString
    if(postResponse.is4xx) {
      logger.error(s"Unable to POST ${this.id}... attempting to PATCH ")
      val patchResponse: HttpResponse[String] = Http(url).
        param("id", s"eq.$this.id").
        method("PATCH").
        auth("c11z", "c11z").
        header("User-Agent", "ComicGator (http://github.com/comicgator/lurker").
        postData(this.export).asString
      if(patchResponse.is4xx) {
        logger.error(s"Unable to PATCH ${this.id}")
//        logger.error(s"${patchResponse.}")
//        throw new MaestroException("When in trouble or in doubt, run in circles, scream and shout.")
      }
      if(patchResponse.is5xx) logger.error("Someone has killed the server in mid path.")
    }
    if(postResponse.is5xx) logger.error("Someone forgot to turn on the server.")


  }

  /**
    * Override and simplify toString method
   * @return String of the Comic id
   */
  override def toString: String = s"$id"
}

/**
 * Parser strategy, right now only supports xpath.
 * @param method String a method for scrapping a piece of content out of the html.
 * @param patterns List of Strings representing an xpath or some unique pattern to get the content out of the html.
 */
case class Parser(method: String, patterns: List[String]) {
  def export = {
    val jsonMap = Map("method" -> method, "patterns" -> patterns)
    // Passing implicit DefaultFormats
    write(jsonMap)(DefaultFormats)
  }
}

/**
 * Strip Class represents a specific page of a web comic i.e. comic 'strip'.
 * @param comicId String representing the comic this strip is related to.
 * @param checksum String a md5 hash of the image html tag found on the strip page. Used to positively identify the
 *                 comic strip.
 * @param title String the title of the strip.
 * @param number Integer used to order the strips in time. Unique index on (comic_id, number) combination.
 * @param url String url of the comic strip.
 * @param image String image url of the comic strip.
 * @param thumbnailImage String a url of the thumbnail image for the comic strip generated by lurker.
 * @param bonusImage String a url of the bonus image of the comic strip.
 * @param altText String the title or alt text found with the main image.
 */
class Strip(comicId: String, checksum: String, title: String, number: Integer, url: String, image: String,
            thumbnailImage: String, bonusImage: String, altText: String) {
  def logger = LoggerFactory.getLogger(this.getClass)
  var id: Integer = 0
  def getId: Integer = this.id
  def setId(id: Integer) = this.id = id

  def export: String = {
    val jsonMap = Map("comic_id" -> comicId, "checksum" -> checksum, "title" -> title,
      "number" -> number, "url" -> url, "image" -> image, "thumbnail_image" -> thumbnailImage,
      "bonus_image" -> bonusImage, "alt_text" -> altText)
    // Passing implicit DefaultFormats
    write(jsonMap)(DefaultFormats)
  }

}
