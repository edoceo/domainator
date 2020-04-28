# domainator

This tool searches all the domains for your name

## Dependencies

* https://github.com/miekg/dns


## Operation

Just run `./domainator --name=example` and it will check for all the domains.

You may have to set `ulimit -n 2048` to make this work, since it opens a lot of ports all at once.

## DNS and TLD files

This tool uses two files `./dns-list.txt` and `tld-list.txt` which have been manually constructed.
The dns-list is the list of servers we'll query, round-robin.
The tld-list is the set of TLDs that will be queried.

You can get some details on domains [https://www.iana.org/domains/root/files](from IANA). 

## Domain Stats

* https://ntldstats.com/tld
* https://hostingtribunal.com/blog/tld-statistics/#gref
* https://www.statista.com/statistics/265677/number-of-internet-top-level-domains-worldwide/


https://ops.tips/blog/udp-client-and-server-in-go/



https://en.wikipedia.org/wiki/Country_code_second-level_domain

https://en.wikipedia.org/wiki/Second-level_domain

https://en.wikipedia.org/wiki/Country_code_top-level_domain

https://en.wikipedia.org/wiki/List_of_Internet_top-level_domains#Internationalized_geographic_top-level_domains

https://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml

DNS Lists:
  https://www.lifewire.com/free-and-public-dns-servers-2626062

