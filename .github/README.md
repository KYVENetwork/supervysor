<div align="center">
  <h1>@supervysor</h1>
</div>

<p align="center">
<strong></strong>
</p>

## Content

- [What is the supervysor?](#what-is-the-supervysor)
- [How does it work?](#how-does-it-work)

## What is the supervysor?

Currently a protocol validator needs to sync
a whole chain as a full node to participate successfully in a data pool.
In parallel, Cosmos full nodes require a lot of storage (~ 10TB for Osmosis),
which leads to high operation costs and less efficient funding usage.
The supervysor manages the sync process of the node. Therefore, it must sync the blocks from currentKey to currentKey + n blocks and stop running when it reaches the upper limit. In addition it needs to make the synced data accessible for the protocol node even the node stopped the syncing process.

## How does it work?

<p align="center">
  <img width="70%" src="assets/supverysor.png" />
</p>