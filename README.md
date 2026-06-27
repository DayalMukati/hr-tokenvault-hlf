# hr-tokenvault-hlf

Solution repository for the **Tokenized Card Vault** Hyperledger Fabric chaincode challenge (NPCI / HackerRank, Intermediate).

This repo is a standard Hyperledger Fabric **test-network** plus a chaincode skeleton at
[`chaincode/tokenvault.go`](chaincode/tokenvault.go). It is cloned into the candidate's
environment by the HackerRank **Setup Script** (via [`setup.sh`](setup.sh)).

## What ships here

| Path | Purpose |
|------|---------|
| `test-network/` | Fabric test network (`network.sh up createChannel`, two orgs, `mychannel`). |
| `chaincode/tokenvault.go` | Chaincode **skeleton** with empty function bodies — the candidate implements these. |
| `chaincode/go.mod`, `go.sum`, `vendor/` | Pinned `fabric-contract-api-go v1.2.0` dependencies (vendored). |
| `setup.sh` | `git clone` this repo into the challenge dir. |

## Candidate task

1. Implement the functions in `chaincode/tokenvault.go`.
2. Deploy: `cd test-network && ./network.sh deployCC -ccn tokenvaultcc -ccp ../chaincode -ccl go`
3. `tok1`: `IssueToken` → `SuspendToken` → `ResumeToken` (ends `ACTIVE`).
4. `tok2`: `IssueToken` → `DeleteToken` (ends absent).

See the full problem statement and scoring breakdown in the question folder.

---

Authored by **Dayal Mukati** — [dayalmukati.com](https://dayalmukati.com)
