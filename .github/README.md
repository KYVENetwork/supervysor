<div align="center">
  <h1>@supervysor</h1>
</div>

![banner](../assets/banner.png)

## Content

- [What is the supervysor?](#what-is-the-supervysor)
- [How does it work?](#how-does-it-work)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)

## What is the supervysor?

To successfully participate in a KYVE data pool such as Cosmoshub or Osmosis, you need to run two nodes: the KYVE protocol node and the data source node (full node of Cosmoshub, Osmosis, etc.). In parallel, these full nodes require a lot of storage (~ 10TB for Osmosis), which leads to high operation costs and therefore less efficient funding usage although pruning can be actived. This is due to the fact that after the start, the node immediately begins to synchronize up to the last block, although only the storage of a certain range of blocks is necessary. However, if this synchronization process is stopped, the node cannot fulfill its tasks as a data source. The supervysor solves this problem by managing the state of the node depending on the requirements of the data pool. This ensures that the node synchronizes only as far as needed, while still providing data when the synchronization process stops.

## How does it work?

The supervysor is a process manager that is wrapped around a node or the cosmovisor. After the initial start, the node-height and the pool-height of the KYVE data pool are queried at a specified interval, after which the difference between the two values is calculated. If the difference is higher than ```height_difference_max``` , the node is set to the `Ghost Mode`. In this mode, the synchronization process is stopped by making the address book inaccessible and starting the node without seeds and with a modified laddr. This ensures that the node cannot reach other peers and thus cannot synchronize new blocks. If the difference is smaller than ```height_difference_min```, the address book is made accessible again and the node is started normally so that peers can be found and the synchronization process can continue. If the difference is smaller than ```height_difference_max``` and larger than ```height_difference_min``` the current mode is kept. In both modes, the endpoints are accessible to the protocol node, so the required data remains accessible even if the node does not synchronize.

<p align="center">
  <img width="70%" src="../assets/supervysor.png" />
</p>

To keep memory requirements as low as possible, we need to specify both a maximum value for how far the data source node can synchronize beyond the current pool height and the matching pruning settings to make sure that not validated data can be pruned. Derived from this, these values were calculated as followed:

* `min_retain_blocks = max_bundle_size / upload_interval * 60 * 60 * 24 * 7` (maximum bundles for 7 days)
* `height_difference_max = max_bundle_size / upload_interval * 60 * 60 * 24 * 7` (maximum bundles for 5 days)
* `height_difference_min = height_difference_max` (maximum bundles for 2.5 days)

These values ensure that
* the data source node will always be 2.5 days ahead to the latest pool-height
* the data source node will not sync to the latest height, because it will stop syncing when the required blocks for the next 5 days are stored locally
* only the required blocks for the next 7 days are kept locally, everything else will be pruned
* because `min_retain_blocks > height_difference_max`, nothing will be pruned before it was validated in the data pool 

_Note: Currently these settings are static and not usable for integrations with state requests. This will be part of one of the next versions and will be published in a few weeks._
## Requirements

The supervysor manages the process of the data source node. First of all, it should be ensured that this node can run successfully. In addition, to successfully participate in a KYVE data pool, it is necessary to create a protocol validator and join a data pool. Further information can be found here: https://docs.kyve.network/validators/protocol_nodes/overview

Make sure your Go version is at least ```1.20```.
## Installation

To install the latest version of `supervysor`, run the following command:

```bash
go install github.com/KYVENetwork/supervysor/cmd/supervysor@latest
```

To install a previous version, you can specify the version:

```bash
go install github.com/KYVENetwork/supervysor/cmd/supervysor@v0.1.0
```

_Optional:_ If you have issues to successfully run the `go install` command, make sure to export the following to your environment:

```bash
env GIT_TERMINAL_PROMPT=1
```

Run `supervysor version` to check the installed version.

You can also install from source by pulling the supervysor repository and switching to the correct version and building
as follows:

```bash
git clone git@github.com:KYVENetwork/supervysor.git
cd supervysor
git checkout tags/vx.x.x -b vx.x.x
make supervysor
```

This will build supervysor in `/build` directory. Afterwards you may want to put it into your machine's PATH like
as follows:

```bash
cp build/supervysor ~/go/bin/supervysor
```

## Usage

To use the supervysor, you first need to initialize it:

```bash
supervysor init
--address-book-path string   'path to address book (e.g. /root/.osmosisd/config/addrbook.json)'
--binary-path       string   'path to chain binaries (e.g. /root/go/bin/osmosisd)'
--chain-id          string   'KYVE chain-id'
--pool-id           int      'KYVE pool-id'
--seeds             string   'seeds for the node to connect'
```

This command creates a config file at ```~/.supervysor/config.toml``` which is editable and required to start the supervysor.

To start the supervysor after the successful initialisation, run the following command:

```bash
supervysor start
```

Then the supervisor starts the chain binaries or cosmovisor to manage the syncing process depending on the required data of the KYVE pool.

## Examples

### 1. Run a Cosmovisor Osmosis node with the supervysor

To run an Osmosis node with the Cosmovisor you have to download and set up the correct binaries. You can see a more detailed
introduction [here](https://docs.osmosis.zone/networks/join-mainnet/).

Verify the correct installation and setup with the successful start of the node:

```bash
cosmovisor run start [flags]
```

With your node being able to run using Cosmovisor, you can install the supervysor to start optimize this process for KYVE purposes. After the [installation](#installation), you can initialize the supervysor with the following command:

```bash
supervysor init
--address-book-path '/root/.osmosisd/config/addrbook.json'
--binary-path       '/root/go/bin/cosmovisor'
--chain-id          'kaon-1'
--pool-id           1
--seeds             '21d7539792ee2e0d650b199bf742c56ae0cf499e@162.55.132.230:2000,44ff091135ef2c69421eacfa136860472ac26e60@65.21.141.212:2000,ec4d3571bf709ab78df61716e47b5ac03d077a1a@65.108.43.26:2000'
```

After the successful initialisation you can start your node with:

```bash
supervysor start
```

The supervysor then will start an Osmosis node as data source for the pool with the ID 1 of the KYVE Kaon Network.

_Note: The shown pool for the example haven't been created yet._