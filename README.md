<div align="center">
  <h1>
    URL Info
  </h1>
  <p>
    Simple URL information server
  </p>
  <p>
    <a href="https://github.com/aasaam/url-info/actions/workflows/build.yml" target="_blank"><img src="https://github.com/aasaam/url-info/actions/workflows/build.yml/badge.svg" alt="build" /></a>
    <a href="https://goreportcard.com/report/github.com/aasaam/url-info"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/aasaam/url-info"></a>
    <a href="https://hub.docker.com/r/aasaam/url-info" target="_blank"><img src="https://img.shields.io/docker/image-size/aasaam/url-info?label=docker%20image" alt="docker" /></a>
    <a href="https://github.com/aasaam/url-info/blob/master/LICENSE"><img alt="License" src="https://img.shields.io/github/license/aasaam/url-info"></a>
  </p>
</div>

## Guide

For see available options

```bash
docker run --rm ghcr.io/aasaam/url-info:latest -h
```

Default cache dir is `/tmp` for change set env variable `ASM_URL_INFO_TEMPORARY_PATH=/path/you/want`

```bash
$ curl -s -X POST -H 'Content-Type: application/json' -d '{"url":"https://news.yahoo.com//","image_resize":"720x240", "image_quality":85}' http://127.0.0.1:4000/info | jq
{
  "title": "Yahoo News - Latest News &amp; Headlines",
  "lang": "en",
  "dir": "ltr",
  "description": "The latest news and headlines from Yahoo! News. Get breaking news stories and in-depth coverage with videos and photos.",
  "url": "https://news.yahoo.com/",
  "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAAA7VBMVEVgAdJfANFnDNRuF9ZqEtVgAtJmCtO8lOvm2PeyhOljBtJeANFdANGFPdzy6/v69/2WWOF+Mdqvf+idYuNkCNNjB9KZXeKuf+iANdu6kuv////YwfRsFNV9MNru4/rx6ft/NNp6K9nt4/nt4/qaX+Lq3vn9/P6ia+RiBdK9leyse+e8leyZXuLi0vbRt/JyH9eGP9zz7Pvs4vnq3fj07vyVV+Gpd+atfOdyHtfKqu/MrvChauT6+P2GPtyUVOD59v1/M9rYwfNzINfcyPVuGNZqEtRoDtRjB9O5kOvv5fqkbuVhA9JtFtVxHNZlCtO5+6fEAAAAm0lEQVQY02NgwA0YQQCFz8TMwookwsjGzsHJBRLg5uFhZOTh4eHl4xfgBgpwCwoJi4iKiUtISknLgBXIyskrKCopq6iqqYMUMDBqaEpxaklp6+jq6XODTeM2MDQyNjE1M7dgYIQabyklZWVtY2jLDbWPx86eX8DB0coJJsBtx+/swujqBnMno7uHpxeyuxm9fXz9GFF9xY3CRwUAfpMOI9WCN3MAAAAASUVORK5CYII=",
  "image_url": "https://s.yimg.com/cv/apiv2/social/images/yahoo_default_logo.png",
  "image_optimized": "/tmp/uuu/5a47e476c530617c22e84e9d2f518ea2.jpg",
  "canonical": "https://news.yahoo.com/"
}
```

<div>
  <p align="center">
    <a href="https://aasaam.com" title="aasaam software development group">
      <img alt="aasaam software development group" width="64" src="https://raw.githubusercontent.com/aasaam/information/master/logo/aasaam.svg">
    </a>
    <br />
    aasaam software development group
  </p>
</div>
