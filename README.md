# Bardock
Monorepo Dockerfile management tool

```
Usage:
  bardock [flags]

Flags:
  -f, --file string   path to the bardock YAML file (default "bardock.yaml")
  -h, --help          help for bardock
  -p, --push          override push image to registry
  -t, --tag string    override image tag
  -v, --verbose       verbose output
      --version       print the version number of bardock
```

### Install

```
bash -c "$(curl -fsSL https://raw.githubusercontent.com/alexxi19/bardock/main/install.sh)"
```

### Creating your Bardock spec file   

In your repo root, create a `bardock.yaml` file.


```yaml
services:
  # Required: Service Name
  service1:
    # Required: Dockerfile path relative to bardock.yaml
    dockerfile: ./service1/Dockerfile
    # Required: Image name for the service
    image_name: service1
    # Optional(default: "."): Docker build context relative to bardock.yaml, the docker build command will be run in this directory 
    build_context: ./service1
    # Optional(default: true): Boolean for pushing the image 
    push: true

config:
  # Required: Registry URL (e.g., Docker Hub username)
  registry_url: alexxi19
  # Optional(default: "latest"): Docker image tag
  # Options: latest, git_sha (uses the prefix of the current git commit sha)
  image_tag: latest
```

You can find more examples in the `examples/` directory. 

### Using Bardock 

To build and push an image for a service using bardock, run

```
bardock [service name]
```

Based on the previous example bardock file, the corresponding command to build and push the image for service1 is:

```
bardock service1
```

### Github actions

You can find an example for using bardock in github actions in `.github/workflows/install.yaml`.

### Override options

You can also use bardock's cli flags to override the tag or push options.

For example: 

```
bardock service1 --tag dev --push false
```
