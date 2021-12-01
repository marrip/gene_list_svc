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
  -analysis string
    	Select analysis (snv, cnv, sv, pindel).
  -bed string
    	Output bed file name (default: [list]_[analysis]_[build].bed).
  -build string
    	Select genome build (37, 38; default: 38).
  -list string
    	Select gene list (aml, aml_ext, all).
  -tsv string
    	Path to tsv file containing gene list.
```

## :rocket: Usage

### Start docker containers - *Still under development - no docker available yet!*

To use this cli, a postgres database server is needed and start the cli docker
with an interactive session:

```bash
docker run -d -e POSTGRES_PASSWORD=notsosecret postgres
docker run -it --rm marrip/gene_list_svc:latest bash
```

### Input `.tsv` file

To load new data into the database, compile a `.tsv` file like so:

```bash
BCR	Gene	snv	aml	Specific
KMT2A	Gene	sv	all	All
RUNX1	Gene	snv,sv,cnv	aml,all	Specific
...
```

### Adding data to database

To add data from a `.tsv` file, simply run:

```bash
gene_list_svc -tsv /path/to/data.tsv
```

### Retrieve data and write to bed file

Choose analysis, diagnostic route, genome build and desired `.bed` file
name to retrieve data from the database:

```bash
gene_list_svc -list aml -analysis snv -build 38 -bed /path/to/aml_snv_38.bed
```
