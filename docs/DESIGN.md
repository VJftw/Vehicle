# Vehicle Design

Vehicle is a tool which enables you to run arbitrary commands on a **temporary** instance on supported cloud providers. This tool aims to fill the following gaps:
 - Programmatic use of a temporary bastion instance to provision infrastructure in air-gapped cloud environments.
 - Running tests on on/against a temporary instance.

## Goals

- Portable: The single binary 

## Background

## Design

### Configuration

```yaml
---
vehicles:
    my_vehicle:
        provider: docker

tasks:
- when: 
  command: 
  files: 
```

## Future Considerations
