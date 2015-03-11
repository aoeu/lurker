package com.comicgator.lurker

import org.json4s.{DefaultFormats, _}
import org.json4s.native.JsonMethods._
import org.json4s.native.Serialization.write
import org.slf4j.LoggerFactory
import dispatch._, Defaults._

import scala.io.{BufferedSource, Source}
import scala.util.{Failure, Success, Try}

// TODO: Split into separate files

/**
 * Main does all the things
 */
object Main {
  /**
   * Logging using logback settings in lurker/src/main/resources/logback.xml
   * If trouble try:
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
      Maestro.saveComic(comic)
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
class Strip(val id: Integer, comicId: String, val checksum: String, title: String, number: Integer, val url: String,
            image: String, thumbnailImage: String, bonusImage: String, altText: String, updated_at: String, created_at: String) {
  def logger = LoggerFactory.getLogger(this.getClass)

  def export(comicId: String, checksum: String, title: String, number: Integer, url: String, image: String,
  thumbnailImage: String, bonusImage: String, altText: String): String = {
    val jsonMap = Map("comic_id" -> comicId, "checksum" -> checksum, "title" -> title,
      "number" -> number, "url" -> url, "image" -> image, "thumbnail_image" -> thumbnailImage,
      "bonus_image" -> bonusImage, "alt_text" -> altText)
    // Passing implicit DefaultFormats
    write(jsonMap)(DefaultFormats)
  }

  def export: String = {
    val jsonMap = Map("comic_id" -> this.comicId, "checksum" -> this.checksum, "title" -> this.title,
      "number" -> this.number, "url" -> this.url, "image" -> this.image, "thumbnail_image" -> this.thumbnailImage,
      "bonus_image" -> this.bonusImage, "alt_text" -> this.altText)
    // Passing implicit DefaultFormats
    write(jsonMap)(DefaultFormats)
  }

}

object Maestro {
  def logger = LoggerFactory.getLogger(this.getClass)
  val mHost = host("192.168.34.10", 3000)
    .addHeader("User-Agent", "Comic Gator (http://github.com/comicgator/lurker)")
    .setContentType("application/json","utf-8")
    .as_!("c11z", "c11z")

  val asStatus: (Res) => Int = as.Response {res =>
    res.getStatusCode
  }

  val asStrip: (Res) => Strip = as.Response {res =>
    val json: JValue = parse(res.getResponseBody)
    implicit val formats = DefaultFormats
    json.extract[Strip]
  }

  val asRes: (Res) => Res = as.Response {res => res}

//  def get() = Http(this.mHost.GET)
  def get(endpoint: String): Future[Res] = Http((this.mHost / endpoint).GET)
  def get(endpoint: String, params: Map[String, String]): Future[Res] = {
    Http((this.mHost / endpoint).GET <<? params)
  }
  def post(endpoint: String, body: String): Future[Res] = {
    val req = (this.mHost / endpoint).POST.setBody(body)
    logger.debug(s"POST body=$body}")
    Http(req OK asRes)
  }
  def patch(endpoint: String, params: Map[String, String], body: String): Future[Res] = {
    val req = (this.mHost / endpoint).PATCH.setBody(body) <<? params
    logger.debug(s"PATCH body= $body")
    Http(req OK asRes)
  }

  def saveComic(c: Comic) = {
    logger.info(s"Saving comic ${c.id}")
    val res = this.post("comic", c.export) recover {
      case e => this.patch("comic", Map("id" -> s"eq.${c.id}"), c.export)
    }
    res onComplete {
      case Success(r) => logger.info(s"Successfully updated comic ${c.id}")
      case Failure(ex) => logger.error(s"Failed to update comic ${c.id}")
    }
  }

//  def saveStrip(s: Strip): Integer = {
//    logger.info(s"Saving Strip ${s.url}")
//    val res = this.post("strip", s.export) recover {
//      case e => this.patch("strip", Map("checksum" -> s"eq.${s.checksum}"), s.export)
//    }
//  }

//  def postNow(endpoint: String, body: String): Int = {
//    val res = this.post(endpoint, body).result(Duration(5, "seconds"))
//    res.getStatusCode
//  }
//  def patchNow(endpoint: String, params: Map[String, Seq[String]], body: String): Int = {
//    val res = this.patch(endpoint, params, body).result(Duration(5,"seconds"))
//    res.getStatusCode
//  }
}
