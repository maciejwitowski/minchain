package services

/*
Sync
* Protocol Overview*
Initial Sync Request
  - Node X broadcasts a sync request on the "sync" topic
  - The request includes X's current finalized block height

Peer Responses
  - Other nodes respond with:
  - Their finalized block height
  - Their latest block height
  - A hash of their finalized block

Consensus Determination
  - Node X collects responses for a set time period
  - It determines the network consensus on finalized height

Finalized Block Sync
  - X requests blocks from finalized height to its current height
  - It verifies the chain of blocks matches the consensus hash

Latest Block Sync
  - X requests blocks from finalized to latest heights
  - It applies standard consensus rules to accept these blocks

*Key Considerations*
Finalized vs Latest Blocks
  - Finalized blocks have reached consensus and won't be reverted
  - Latest blocks may still be subject to change

Handling Disagreements
  - If peers disagree on finalized height, X should:
  - Choose the majority view
  - If no clear majority, use the most recent common ancestor

Bridging the Gap
  - For blocks between finalized and latest:
  - Apply normal consensus rules (e.g., longest chain)
  - Be prepared to handle short-term forks

Optimizations
  - Use block headers for initial sync to reduce bandwidth
  - Implement parallel block downloads from multiple peers

Security Measures
  - Verify block signatures and difficulty adjustments
  - Implement safeguards against malicious peers providing false data
*/
type Sync struct {
}

func NewSync() {

}
