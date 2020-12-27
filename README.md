## etymology

Etymology library for Go, heavily inspired by
https://github.com/jmsv/ety-python.

This is a WIP.

Uses Gerard de Melo's [Etymological Wordnet](http://etym.org/).
This is quite out of date, though (2013). The wordnet is mostly sourced from
Wiktionary, which has much more detailed etymology information these days.

### Usage

You'll need a copy of de Melo's Wordnet:

```console
[jpw@xyz:~] $ wget https://cs.rutgers.edu/~gd343/downloads/etymwn-20130208.zip
[jpw@xyz:~] $ unzip etymwn-20130208.zip
```

Then, for example:

```console
[jpw@xyz:~] $ go run github.com/jamespwilliams/etymology/cmd/ety etymwn.tsv psychoneuroendocrinological eng
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

### Future

I have the vague idea of making a web interface for this, which accepts a word
and displays an etymological tree for it.

Also, it'd be nice to use more up-to-date etymology information, as mentioned
earlier.
