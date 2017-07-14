
# Fixity (unstable)

_fixity: the quality of being incapable of mutation; immutability_

Fixity is an experimental immutable personal data store.

This is a highly experimental pet project. API and Json formats are not
yet final.

## Project Goals

With much inspiration from [Camlistore](https://camlistore.org),
Fixity aims to be:

- Low maintenance, small in scope.
- Easily sync files between computers.
- Versioned and content addressable.
- Per-id customizable deduplication for deduping of large and small data.
- A simple schemaless database.
- [Distributed](#distributed), and [eventually consistent](#eventually-consistent).
- Fully usable if any or all nodes are disconnected.


## Motivation

I've long been in need of a system to provide access to my files/data
in a cloud-like, distributed manner. Eg, to access my files from
any computer i own regardless of where i am. While open source
solutions exist, many of them store files in such a way that
they're centralized, requiring access through the firewall if you're
outside of the network; I did not like this.
Furthermore, they often store the data in a way that i felt obscured it.
Meaning it was difficult for me to manage the data myself.. the underlying
formats did not feel "open" to me.

Along with my needs for file storage, i wanted to store emails, chatlogs,
wiki pages, home inventory, etc. I wanted a database for my life, schemaless
and easy to manage.

Fixity is my attempt to implement my desired features in one package.
Fixity should store data in a format that you can read with any
text editor and write a program easily to extract, should you so feel the
need to. It can be used as a simple schemaless database. See also:
[Project Goals](#project-goals).


## Distributed

Fixity aims to connect all of your computers together, to provide
seamless storage between computers similar to a NAS
*(Network Attacked Storage)*. Unlike a NAS however, when a computer
leaves the private network *(such as laptops leaving wifi)*,
that disconnected computer should still be able to [write and read
data](#eventually-consistent). There is no hot new tech here,
Fixity is simply aimed so that laptops and phones can access all
data of whatever other Fixity nodes are available at any given time.
The distributed features here are small in scope, boring.

Fixity's features focus on accessibility of data, and allowing data to be
eventually consistent between nodes.
For data safety and health, Fixity will focus on backing up data onto
external hard drives, cloud storage, etc. It will not replicate data
automatically to ensure data redundancy. Again, small scope.


## Eventually Consistent

Because Fixity expects nodes to go offline or go outside of firewalls,
writing to offline nodes must eventually sync up to the network. This
is achieved by writing to what Fixity calls a blockchain *ledger*. This
has nothing to do with cryptocoin, but is a term used to convey how
Fixity achieves consensus between out of sync nodes. If, for example,
the same file is written on two disconnected nodes then when they eventually
do connect they will compare blockchains and one will write it's changes to
the other.

Conflicts of data like this are automatically resolved. However, since no
data is lost on this automatic resolution, the user is free to compare
versions at any time and choose one over the other, or write new data with
merged values, etc. No conflicts will ever prevent merging previously
offline nodes.


## Important Todos

This project is in very early development, so to summarize obvious things it
may be lacking the following items are notable unfinished TODOs:

- [ ] Distributed node implementation. *(ie, Fixity is not distributed yet)*
- [ ] Build system with multi-os released binaries.
- [ ] Continuous integration.
- [ ] Finish [Snail indexer](https://github.com/leeola/fixity/tree/master/indexes/snail),
  or preferably find a schemaless indexer to replace Snail.
- [ ] Node Auth implementation/decisions. To ensure writes only from desired users.
- [ ] Write signing. To identify the users who wrote the block/content.


## License

MIT
