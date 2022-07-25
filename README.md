## etty

Etymology dataset scraper and visualizer, written in Go. This is a WIP.

![Screenshot of Web Interface](https://raw.githubusercontent.com/jamespwilliams/etty/master/_assets/web.png)

### Table of Contents

- [Command-Line](#command-line)
- [Web Interface](#web-interface)
- [Wordnet Generator](#wordnet-generator)

### Command-Line

Inspired by https://github.com/jmsv/ety-python. You'll need a wordnet in TSV
format, for example Gerard de Melo's [Etymological Wordnet](http://etym.org/):

```console
$ wget https://archive.org/download/etymwn-20130208/etymwn-20130208.zip
$ unzip etymwn-20130208.zip
```

Then, for example:

```console
$ go install github.com/jamespwilliams/etty/cmd/etty@latest
$ etty etymwn.tsv psychoneuroendocrinological eng

    psychoneuroendocrinological (eng)
    ├── psychoneuroendocrinology (eng)
    │   ├── neuro- (eng)
    │   │   └── νευρο- (grc)
    │   ├── psycho- (eng)
    │   │   └── ψυχή (grc)
    │   │       └── ψύχω (grc)
    │   └── endocrinology (eng)
    │       ├── endocrine (eng)
    │       ├── -logy (eng)
    │       └── -ology (eng)
    └── -ical (eng)

$ etty etymwn.tsv microgyfrifiadur cym

    microgyfrifiadur (cym)
    ├── micro- (cym)
    └── cyfrifiadur (cym)
        ├── cyfrif (cym)
            └── -adur (cym)

$ etty etymwn.tsv 'ალერსიანი' kat

    ალერსიანი (kat)
    ├── ალერსი (kat)
    │   └── աղերս (xcl)
    └── -იანი (kat)
```

The output is coloured nicely, too:

![Screenshot of CLI](https://raw.githubusercontent.com/jamespwilliams/etty/master/_assets/cli.png)

### Web Interface

You'll need a wordnet in TSV format, as mentioned in the Command-Line section.

First, start the API:

```console
$ go run ./cmd/etty-api/ path/to/wordnet.txt tcp ':3000'
```

The web interface itself is just a static HTML file, can be served (for example)
with Python:

```console
$ python3 -m http.server
```

### Wordnet Generator

Gerard de Melo's [Etymological Wordnet](http://etym.org/)
is quite out of date (last updated in 2013). The wordnet is mostly sourced from
Wiktionary, which has much more detailed etymology information these days.

The binary in `cmd/wiktionary-parse/` accepts an XML dump of Wiktionary, and
will output a wordnet similar to de Melo's, which the other binaries in this
project can then accept.

Something along the lines of the following should work:

```console
$ wget https://dumps.wikimedia.org/enwiktionary/latest/enwiktionary-latest-pages-articles-multistream.xml.bz2
$ pv enwiktionary-latest-pages-articles-multistream.xml.bz2 | bzcat | go run ./cmd/wiktionary-parse > wordnet.txt
```

`wiktionary-parse` attempts to extract etymological information using the
templates in the Etymology sections of words.
