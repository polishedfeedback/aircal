# Air Freight Calculator

A command-line interface (CLI) tool built in Go for calculating air freight charges. The calculator supports both UK (120x100cm) and EU (120x80cm) pallet types and provides an interactive experience for entering shipment details.

## Installation

```bash
go install github.com/polishedfeedback/aircalc@latest
```

This command will download and install the binary to your `$GOPATH/bin` directory. Make sure your `$GOPATH/bin` is added to your system's `PATH`.

### Verify Installation

```bash
aircalc --version
```

If you see a version number, the installation was successful.

## Quick Start

1. Run the calculator:
```bash
aircalc
```

2. Follow the interactive prompts:
   - Choose pallet type (UK/EU)
   - Enter pallet details as `height,weight`
   Example: `150,350 160,400` or `150,350;160,400`

