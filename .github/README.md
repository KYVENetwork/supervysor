![banner](../assets/banner.png)

<p align="center">
<strong>Run your KYVE Protocol node as efficient as possible</strong>
</p>

<div align="center">
  <img alt="License: Apache-2.0" src="https://badgen.net/github/license/KYVENetwork/supervysor?color=green" />

  <img alt="License: Apache-2.0" src="https://badgen.net/github/stars/KYVENetwork/supervysor?color=green" />

  <img alt="License: Apache-2.0" src="https://badgen.net/github/contributors/KYVENetwork/supervysor?color=green" />

  <img alt="License: Apache-2.0" src="https://badgen.net/github/releases/KYVENetwork/supervysor?color=green" />
</div>

<div align="center">
  <a href="https://twitter.com/KYVENetwork" target="_blank">
    <img alt="Twitter" src="https://badgen.net/badge/icon/twitter?icon=twitter&label" />
  </a>
  <a href="https://discord.com/invite/kyve" target="_blank">
    <img alt="Discord" src="https://badgen.net/badge/icon/discord?icon=discord&label" />
  </a>
  <a href="https://t.me/kyvenet" target="_blank">
    <img alt="Telegram" src="https://badgen.net/badge/icon/telegram?icon=telegram&label" />
  </a>
</div>

<br>

> [!IMPORTANT]
> In this README you will find information on contribution guidelines and
> detailed documentation about the low level implementation of the supervysor.
>
> You can find the complete documentation on installation and usage
> here: **[https://docs.kyve.network/tools/supervysor](https://docs.kyve.network/tools/supervysor)**
> 
# Build from Source

You can install the supervysor from source by pulling the supervysor repository and switching to the correct version and building
as follows:

```bash
git clone git@github.com:KYVENetwork/supervysor.git
cd supervysor
git checkout tags/vx.x.x -b vx.x.x
make
```

This will build supervysor in `/build` directory. Afterwards you may want to put it into your machine's PATH like
as follows:

```bash
cp build/supervysor ~/go/bin/supervysor
```

# How to contribute

Generally, you can contribute to the supervysor via Pull Requests. The following branch conventions are required:

- **feat/\***: for a new feature
- **fix/\***: for fixing a bug
- **refactor/\***: for improving code maintainability without changing any logic
- **docs/\***: for upgrading documentation around the supervysor
- **test/\***: for addition or changes in tests

For committing new changes [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) have
to be used.

Once the Pull Request is ready, it can be opened against the `main` branch. Once the owners have approved
your Pull Request, your changes get merged.

# Releases

For creating release versions [Semantic Versioning](https://semver.org/) is used. A release is created
by manually creating the release over GitHub. For that the latest commit gets tagged with the new version,
which will also be the title of the release. Once the release is created it gets automatically published to
https://pkg.go.dev/github.com/KYVENetwork/supervysor.

# What is the supervysor?

Participating in a KYVE data pool such as CosmosHub or Osmosis requires running two nodes: the KYVE protocol node and the data source node (e.g., full node of CosmosHub, Osmosis, etc.). However, running these full nodes in parallel can result in high storage requirements (approximately >10TB for Osmosis), leading to increased operational costs and inefficient resource utilization. This inefficiency arises because the node begins synchronizing from the start, even though it only requires storage for a certain range of blocks. Additionally, the node lacks information about the progress of the KYVE pool and the already validated data, making pruning impractical when running a node as a KYVE data source.

However, if the synchronization process is halted, the node cannot fulfill its responsibilities as data source effectively. To overcome this challenge, the supervysor is introduced as a solution. The supervysor manages the data source node process based on the requirements of a KYVE data pool. It ensures that the node synchronizes only up to the necessary extent and continues to provide data even when the synchronization process is paused and prunes data that already has been validated.

By implementing the supervysor, the synchronization process and the disk storage requirements are optimized, reducing unnecessary operational costs. The node can focus on synchronizing up to the required point, thus efficiently utilizing resources while fulfilling its role as a data source for the KYVE pool.

# Structure

The supervysor is a process manager that is wrapped around a Tendermint node or the Cosmovisor. After the initial start, the node-height and the pool-height of the KYVE data pool are queried at a specified interval, after which the difference between the two values is calculated. If the difference is higher than ```height_difference_max``` , the node is set to the `Ghost Mode`. In this mode, the synchronization process is stopped by making the address book inaccessible and by starting the node without seeds and with a modified laddr. This ensures that the node cannot reach other peers and thus cannot synchronize new blocks. If the difference is smaller than ```height_difference_min```, the address book is made accessible again and the node is started with specified seeds so that peers can be found and the synchronization process can continue. If the difference is smaller than ```height_difference_max``` and larger than ```height_difference_min``` the current mode is kept. In both modes, the endpoints are accessible to the protocol node, so the required data remains accessible even if the node does not synchronize.

<p align="center">
  <img width="70%" src="../assets/supervysor.png" />
</p>

To keep memory requirements as low as possible, we need to specify a maximum value for how far the data source node can synchronize beyond the current pool height, that is calculated as followed:

* `height_difference_max = max_bundle_size / upload_interval * 60 * 60 * 24 * 2` (maximum bundles for 2 days)
* `height_difference_min = height_difference_max / 2` (maximum bundles for 1 day)

These values ensure that
* the data source node will always be 1 day ahead to the latest pool-height,
* the data source node will not sync to the latest height, because it will stop syncing when the required blocks for the next 2 days are stored locally,
* the data source node has a time window of 1 day to connect to peers to continue syncing before the pool catches up.

## Pruning

Aside from the optimized syncing process, pruning already validated data is the second role of the supervysor to fulfill its goal of reducing disk storage requirements. Therefore, a custom pruning method is used, which relies on the provided Tendermint functionality of pruning all blocks and the state until a specified height. In the context of the supervysor, this until-height should always be lower than the latest validated height of the KYVE data pool to ensure no data is pruned that needs validation. Unfortunately, the node has to be stopped to execute the pruning process, while a pruning-interval needs specification in hours. During this interval, the supervysor halts the current node process, prunes blocks and state until the already validated height, and restarts the node. Due to the required time to connect with peers and to prevent the pool from catching up with the node, the pruning process is only initiated if the node is in GhostMode. If the node is in NormalMode, even if the interval reaches the pruning threshold, pruning will be enabled immediately after the node enters GhostMode. Additionally, it is recommended to set the pruning-interval to a value of at least six hours to ensure there is enough time to find peers before the pool catches up.