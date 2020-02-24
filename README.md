# ercole-agent-virtualization
[![Build Status](https://travis-ci.org/ercole-io/ercole-agent-virtualization.png)](https://travis-ci.org/ercole-io/ercole-agent-virtualization) [![Join the chat at https://gitter.im/ercole-server/community](https://badges.gitter.im/ercole-server/community.svg)](https://gitter.im/ercole-server/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
 
Agent for the Ercole project. Documentation available [here](https://ercole.io).

## Requirements

- [go](https://golang.org/)

## How to build

    go build ./main.go -o ercole-agent-virtualization

## How to run

Just run the binary: `./ercole-agent-virtualization`

You can customize the agent parameters by copying the `config.json` file in the same directory as your ercole-agent-virtualization or in `/opt/ercole-agent-virtualization/config.json`.
