# SysKit

> A modern, Linux-first command-line toolkit built with Go for system inspection, resource monitoring, and diagnostics.

## Overview

SysKit is an open-source CLI application designed for developers, backend engineers, DevOps professionals, and Linux enthusiasts who need fast, reliable, and consistent access to system information.

Rather than simply wrapping existing Linux utilities, SysKit aims to collect information directly from native Linux interfaces such as `/proc` and `/sys` whenever possible, providing both a practical daily-use tool and a hands-on exploration of Linux internals.

The project is built with a strong focus on performance, clean architecture, maintainability, and an excellent terminal user experience.

## Goals

* Build a unified interface for common Linux inspection tasks.
* Learn advanced Go through a real-world project.
* Explore Linux internals by interacting directly with kernel interfaces.
* Create a fast, modular, and extensible CLI application.
* Serve as a long-term open-source learning project.

## Planned Features

* System information
* CPU monitoring
* Memory monitoring
* Disk usage analysis
* Process inspection and management
* Network monitoring
* Port inspection
* Filesystem information
* System health diagnostics
* Interactive terminal dashboard
* JSON and YAML output
* Plugin system
* Docker and Kubernetes integration
* Remote monitoring over SSH

## Design Principles

* Linux-first
* Performance-oriented
* Native APIs whenever possible
* Modular architecture
* Consistent user experience
* Minimal dependencies
* Comprehensive testing
* Extensible by design

## Technology Stack

* Go
* Linux Kernel Interfaces (`/proc`, `/sys`)
* Cobra (CLI)
* Bubble Tea (Terminal UI)
* Lip Gloss (Terminal Styling)

## Project Status

🚧 **Early development**

SysKit is currently in the design and planning phase. The project is being developed using a specification-driven approach, with a strong emphasis on software architecture, Linux internals, and modern Go engineering practices before implementation begins.

## Vision

SysKit aims to become a reliable daily companion for backend engineers and Linux users by providing a modern command-line experience for system inspection, monitoring, and diagnostics.

## License

This project is licensed under the MIT License.
