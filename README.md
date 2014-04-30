
    wget http://www.polishmywriting.com/download/wikipedia2text_rsm_mods.tgz
    tar zxvf wikipedia2text_rsm_mods.tgz
    cd wikipedia2text

Get the dump ending in `-pages-articles.xml.bz2` from http://dumps.wikimedia.org/:

    cd data
    wget http://dumps.wikimedia.org/afwiki/20140429/afwiki-20140429-pages-articles.xml.bz2
    bunzip2 afwiki-20140429-pages-articles.xml.bz2

Parse the XML into folders containing text files:

    ../tools/wikipedia2text/xmldump2files.py afwiki-20140429-pages-articles.xml afrikaans

Actually, just read this tutorial:

    http://www.evanjones.ca/software/wikipedia2text.html