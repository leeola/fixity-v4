
# Fixity (prototype)

_fixity: the quality of being incapable of mutation; immutability_

Fixity is an experimental immutable personal data store.

This is a highly experimental pet project. It will not be performant,
it will not be efficient. Usage of Fixity to store data is not recommended
until the API and json formats have been finalized. **Do not use currently.**


## Project Goals

With much inspiration from [Camlistore](https://camlistore.org),
Fixity aims to be:

- Low maintenance. Both for Fixity dev(s) and users.
- Easy to migrate away from. The data format isn't binary, it's just Json.
- Append only, versioned and easy to reason about.
- Developer friendly. Fixity has a friendly Go and HTTP API.
- A simple schemaless database. Good for wikis, bookmarks, notes, clipboards, etc.
- Distributed. Each node has access to the sum of all connected nodes data.
- Fully usable if nodes are offline. Nodes read whatever data they have access to.
- Able to deduplicate binary data.


## Motivation

I've long been in need of a system to provide access to my files/data
in a cloud-like, distributed manner. Eg, to access my files from
any computer i own regardless of where i am. While open source
solutions exist, many of them store files in such a way that
they're centralized, requiring access through the firewall if you're
outside of the network.
Furthermore, they often store the data in a way that i felt obscured it.
Meaning it was difficult for me to manage the data myself.. the underlying
formats did not feel "open" to me.

Along with my needs for file storage, i wanted to store emails, chatlogs,
wiki pages, home inventory, etc. I wanted a database for my life, schemaless
and easy to manage.

Fixity is my attempt to implement my desired features in one low maintenance
package. Fixity should store data in a format that you can read, and write
a program easily to extract, should you so feel the need to. It can be used
as a simple database, or store  See
[Project Goals](#project-goals) also.




## License

MIT
