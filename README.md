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

### Run
```bash
docker build -t <image_name> .
docker run -p 8080:8080 <image_name>
```

### Generate Config
You will need a config for your printer and filament, this is also where you control the settings for the print. The 
easiest way to do this is to download [Prusa Slicer](https://www.prusa3d.com/page/prusaslicer_424/) to your machine, 
set your printer, filament, settings, etc. Then click file > export to generate the .ini. You will include this
file in your request along with your model.

### Use the tool
The easiest way to use CloudSlicer is with a GUI API client like [Postman](https://www.postman.com/),
[HTTPie Desktop](https://httpie.io/desktop), or [Insomnia](https://insomnia.rest/). They are pretty
straightforward to use, just include a Multipart Form and model and config fields with their respective
files.


