# Archived

This (Go) Version of Fixity is being archived and the format being altered.
Failings from this version of the fixity "spec" were centered around the
`embeddable database` goal of fixity.

This implementation does not map cleanly to SQL-oriented databases. The
schemaless design has proven more of a burdon than a benefit. Schemaless
ends up translating to requring custom indexers to understand the database
oriented data. This requirement ends up meaning that for an embedded database
Fixity would have to use its own implementation of indexing and the resulting
output would never compare to more SQL-centric real databases.

This implementation has proven good in many respects, but falls short of that
schema problem. The new format will have the goal of being easily translatable
to SQL oriented databases.

As an aside, the primary author (Lee Olayvar) has switched primarily from Go to
Rust, and the new version is being written in Rust - hence the repo change.

For future developments, see [https://github.com/leeola/fixity](https://github.com/leeola/fixity).

## License

MIT
