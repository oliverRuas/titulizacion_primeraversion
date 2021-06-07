export FABRIC_CFG_PATH=$(cd ../ && pwd)
configtxgen -profile FiveOrgsOrdererGenesis -channelID system-channel -outputBlock ../network/channel-artifacts/genesis.block
configtxgen -profile FiveOrgsChannel -outputCreateChannelTx ../network/channel-artifacts/channel.tx -channelID securitization

configtxgen -outputAnchorPeersUpdate ../network/channel-artifacts/originatorMSPanchors.tx -profile FiveOrgsChannel -asOrg originatorMSP -channelID securitization

configtxgen -outputAnchorPeersUpdate ../network/channel-artifacts/farmerMSPanchors.tx -profile FiveOrgsChannel -asOrg farmerMSP -channelID securitization

configtxgen -outputAnchorPeersUpdate ../network/channel-artifacts/aggregatorMSPanchors.tx -profile FiveOrgsChannel -asOrg aggregatorMSP -channelID securitization

configtxgen -outputAnchorPeersUpdate ../network/channel-artifacts/spvMSPanchors.tx -profile FiveOrgsChannel -asOrg spvMSP -channelID securitization

configtxgen -outputAnchorPeersUpdate ../network/channel-artifacts/underwritterMSPanchors.tx -profile FiveOrgsChannel -asOrg underwritterMSP -channelID securitization