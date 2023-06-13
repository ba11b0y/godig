# godig

godig is an implementation of [DNS in a weekend](https://implement-dns.wizardzines.com/index.html) in Go!

## Usage

```shell
go run cmd/main.go
 
Querying 198.41.0.4 for twitter.com
Querying 192.5.6.30 for twitter.com
Querying 198.41.0.4 for a.r06.twtrdns.net
Querying 192.5.6.30 for a.r06.twtrdns.net
Querying 205.251.195.207 for a.r06.twtrdns.net
Querying 205.251.192.179 for twitter.com
Resolved IP for twitter.com is 104.244.42.1

```


#### TODO

 - Update README on differences between the blog's implementation and this one.
 - Add tests.
 - Support CLI usage.
 - Add support for more records, refer https://implement-dns.wizardzines.com/book/exercises.html
