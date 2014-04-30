<!--
Converts Wikipedia articles in wiki format into an XML format. It might
segfault or go into an "infinite" loop sometimes.

Evan Jones <evanj@mit.edu>
April, 2008
Released under a BSD licence.
http://evanjones.ca/software/wikipedia2text.html
-->

<?php
error_reporting(E_ALL);
require_once("mediawiki_converter.php");

$stdin = fopen('php://stdin', 'r');
while (1)
{
  $file = fgets($stdin);
  $file = chop($file);

  if (strcmp($file, "") == 0)
  {
     break;
  }

  $wikitext = file_get_contents($file);
  if (strlen($wikitext) > 0) {
     echo "$file\n";
     $filename_parts = explode("/", $file);
     $title = $filename_parts[count($filename_parts)-1];
     $title = str_replace(".txt", "", $title);
     $title = urldecode($title);

     // Configures options for converting to XML
     $xmlg = array();
     $xmlg["usetemplates"] = "none";
     $xmlg["resolvetemplates"] = "none";
     $xmlg["templates"] = array();
     $xmlg['add_gfdl'] = false;
     $xmlg['keep_interlanguage'] = false;
     $xmlg['keep_categories'] = false;
     $xmlg['text_hide_images'] = true;
     $xmlg['text_hide_tables'] = true;
     $xmlg["useapi"] = false;
     $xmlg["xml_articles_header"] = "<articles>";

     // No idea what it does, but it makes it work
     $content_provider = new ContentProviderHTTP;

     $converter = new MediaWikiConverter;
     $xml = $converter->article2xml($title, $wikitext , $xmlg);
     file_put_contents(str_replace('.txt', '.xml', $file), $xml);
     echo str_replace('.txt', '.xml', $file) . "\n";
  }
}

exit(1);
?>
