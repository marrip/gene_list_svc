# gene_list_svc

Database cli to get bed files from gene/region lists for different diagnostics

## :speech_balloon: Introduction

Cli implemented in golang to add and update databases for AML and ALL with relevant
genes, transcripts, exons or genetic regions, extract all potential fusion partners
for certain genes (http://atlasgeneticsoncology.org/), get genetic coordinates from
the ENSEMBL api and write `.bed` files for different analyses and diagnostics.

## :heavy_exclamation_mark: Dependencies

To use the cli, the following dependencies must be met:

[![docker](https://img.shields.io/badge/docker-20.10.10-blue)](https://www.docker.com/)

## :wrench: Configuration

Env | Required | Default
--- | --- | ---
DB_HOST | - | localhost
DB_PORT | - | 5432
DB_USER | x | -
DB_PASSWORD | x | -
DB_NAME | - | gene_list
ATLAS_ROOT_URL | - | http://atlasgeneticsoncology.org
ENSEMBL_38_REST_URL | - | https://rest.ensembl.org
ENSEMBL_37_REST_URL | - | https://grch37.rest.ensembl.org

## :checkered_flag: Flags

```bash
  update
  --tsv string
      tsv containg list of genetic regions

  extract
  --analysis string
      choose analysis (cnv, pindel, snv, sv)
  --bed string
      set individual bed file name
  --build string
      choose genome build
  --chr bool
      use chr-prefix for chromosome ids
  --tables string
      comma-separated list of tables to be included
```

## :rocket: Usage

### Start docker containers

To use this cli, a postgres database server is needed with a docker network connection
and the gene list database:

```bash
$ docker network create -d bridge gene_list_svc
$ docker run -d --name gene_list_svc_db -e POSTGRES_PASSWORD=notsosecret --network=gene_list_svc postgres
$ docker exec -it gene_list_svc_db bash
$ psql -U postgres
> CREATE DATABASE gene_list;
> \q
$ exit
```

The second step is to create an instance of the gene_list_svc container which is on the
same network and has an input/output directory mounted:

```bash
$ docker run -it --rm -e DB_USER=postgres -e DB_PASSWORD=notsosecret -e DB_HOST=gene_list_svc_db --network=gene_list_svc -v /Path/on/host/to/input/output:/data marrip/gene_list_svc:latest bash
$ cd /data
$ gene_list_svc extract --analysis snv --tables my_table
```

### Input `.tsv` file

To load new data into the database, compile a `.tsv` file like so:

```bash
id	class	analyses	tables	include_partners	coordinates
ABL1	gene	snv,sv	aml	true
NM_000964.4 transcript  snv aml false
ENSE00001422265 exon  snv aml false
my_region region  snv aml false 5:55000-60000
...
```

The coordinates column can be skipped if no region of class `region` is supplied.

### Adding data to database

To add data from a `.tsv` file, simply run:

```bash
$ gene_list_svc update --tsv my_genes.tsv
```

### Retrieve data and write to bed file

Choose analysis, db tables, genome build and desired `.bed` file
name to retrieve data from the database:

```bash
$ gene_list_svc extract --analysis snv --tables my_table,your_table --build 38 --bed my_genes.bed
```

The flag `--chr` can be set to false if the `chr` prefix for chromosomes
is not desired.
