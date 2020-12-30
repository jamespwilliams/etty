## etymology

Etymology dataset scraper and visualizer, written in Go. This is WIP.

### Table of Contents

- [Command-Line](#command-line)
- [Web Interface](#web-interface)
- [Wordnet Generator](#wordnet-generator)

### Command-Line

Inspired by https://github.com/jmsv/ety-python. You'll need a wordnet, for
example Gerard de Melo's [Etymological Wordnet](http://etym.org/):

```console
[jpw@xyz:~] $ wget https://cs.rutgers.edu/~gd343/downloads/etymwn-20130208.zip
[jpw@xyz:~] $ unzip etymwn-20130208.zip
```

Then, for example:

```console
[jpw@xyz:~] $ go get github.com/jamespwilliams/etymology/cmd/ety
[jpw@xyz:~] $ ety etymwn.tsv psychoneuroendocrinological eng
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

[jpw@xyz:~] $ go run github.com/jamespwilliams/etymology/cmd/ety etymwn.tsv microgyfrifiadur cym
microgyfrifiadur (cym)
├── micro- (cym)
└── cyfrifiadur (cym)
    ├── cyfrif (cym)
        └── -adur (cym)

[jpw@xyz:~] $ go run github.com/jamespwilliams/etymology/cmd/ety etymwn.tsv 'ალერსიანი' kat
ალერსიანი (kat)
├── ალერსი (kat)
│   └── աղերս (xcl)
└── -იანი (kat)
```

The output is coloured nicely, too:

![Screenshot of CLI](https://raw.githubusercontent.com/jamespwilliams/etymology/master/_assets/cli.png)

### Web Interface

You'll need a wordnet, as in the Command-Line section.

First, start the API:

```console
[jpw@xyz:ety] $ go run ./cmd/ety-api/ path/to/wordnet.txt
```

The web interface itself is just a static HTML file, can be served (for example)
with Python:

```console
[jpw@xyz:ety/_www] $ python3 -m http.server
```

![Screenshot of Web Interface](https://raw.githubusercontent.com/jamespwilliams/etymology/master/_assets/web.png)

### Wordnet Generator

Gerard de Melo's [Etymological Wordnet](http://etym.org/)
is quite out of date (last updated in 2013). The wordnet is mostly sourced from
Wiktionary, which has much more detailed etymology information these days.

The binary in `cmd/wiktionary-parse/` accepts an XML dump of Wiktionary, and
will output a wordnet similar to de Melo's, which the other binaries in this
project can then accept.

Something along the lines of the following should work:

```console
[jpw@xyz:ety] $ wget https://dumps.wikimedia.org/enwiktionary/latest/enwiktionary-latest-pages-articles-multistream.xml.bz2
[jpw@xyz:ety] $ pv enwiktionary-latest-pages-articles-multistream.xml.bz2 | bzcat | go run ./cmd/wiktionary-parse > wordnet.txt
```

`wiktionary-parse` attempts to extract etymological information using the
templates in the Etymology sections of words.
