# Cloud Slicer
HTTP API to slice models for 3D printing.

## About
Around two months ago I wondered if I could control my non-networked 3D printer from a Discord bot using some 
programming and an old Chromebox turned server. CloudSlicer is the slicing portion of that project, it took me 
embarrassingly long to even consider using the CLI and even longer to figure out how to get it working. CloudSlicer is
currently in active development and not recommended for anything but light hobby use.

## Usage
If you want to use CloudSlicer then you are going to need to self-host it, it is pretty barebones, so you will need to
add features that you need yourself. Authentication, rate-limiting, or just security in general are not present, 
CloudSlicer is not meant for production
### Prerequisites
- [Docker](https://docs.docker.com/install/)
### Generate Config
You will need a config for your printer and filament, this is also where you control the settings for the print. The 
easiest way to do this is to download [Prusa Slicer](https://www.prusa3d.com/page/prusaslicer_424/) to your machine, 
set your printer, filament, settings, etc. Then click file > export to generate the .ini. The pla.ini is for a Prusa
MK4 with PLA filament.
### Run
```bash
docker build -t <image_name> .
docker run -p 8080:8080 <image_name>
```


