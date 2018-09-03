# drunkard server

The idea behind that project, is simulate a way to put data into a database, a lot of data,
as a drunkard do with beers =)


## Run locally

```
make run
```

## Testing

```
make check-integration
```

## Releasing

```
make release version=<version>
```

It will create a git tag with the provided **<version>**
and build and publish a docker image.

## Vendoring

To update vendored dependencies run:

```
make vendor
```
