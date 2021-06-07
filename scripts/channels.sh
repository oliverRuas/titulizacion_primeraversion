export FABRIC_CFG_PATH=$(cd ../ && pwd)
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/originator/admins/originator-admin/msp)
export CLIENTAUTH_CERTFILE=$(cd ../ && echo $PWD/fabric-ca/originator/admins/originator-admin/tls/server.crt)
export CLIENTAUTH_KEYFILE=$(cd ../ && echo $PWD/fabric-ca/originator/admins/originator-admin/tls/server.key)
export CORE_PEER_LOCALMSPID=originatorMSP
export ORDERER_CA=$(cd ../ && echo $PWD/fabric-ca/originator/orderers/originator-orderer/tls/ca.crt)
#INVENTADO
# export ORDERER_TLS=$(cd ../ && echo $PWD/fabric-ca/originator/orderers/originator-orderer/tls/ca.crt)

# Create the application channel
peer channel create -o localhost:7050 -c securitization -f ../network/channel-artifacts/channel.tx --outputBlock ../network/channel-artifacts/securitization.genesis.block --tls --cafile $ORDERER_CA --clientauth --certfile $CLIENTAUTH_CERTFILE --keyfile $CLIENTAUTH_KEYFILE
sleep 2

# Let the peers join the channel
#originator-peer
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_LOCALMSPID=originatorMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=$(cd ../ && echo $PWD/fabric-ca/originator/peers/originator-peer/tls/ca.crt)
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/originator/admins/originator-admin/msp)
peer channel join -b ../network/channel-artifacts/securitization.genesis.block

sleep 2

#farmer-peer
export CORE_PEER_ADDRESS=localhost:8051
export CORE_PEER_LOCALMSPID=farmerMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=$(cd ../ && echo $PWD/fabric-ca/farmer/peers/farmer-peer/tls/ca.crt)
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/farmer/admins/farmer-admin/msp)
peer channel join -b ../network/channel-artifacts/securitization.genesis.block

sleep 2

#aggregator-peer
export CORE_PEER_ADDRESS=localhost:9051
export CORE_PEER_LOCALMSPID=aggregatorMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=$(cd ../ && echo $PWD/fabric-ca/aggregator/peers/aggregator-peer/tls/ca.crt)
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/aggregator/admins/aggregator-admin/msp)
peer channel join -b ../network/channel-artifacts/securitization.genesis.block

sleep 2

#spv-peer
export CORE_PEER_ADDRESS=localhost:10051
export CORE_PEER_LOCALMSPID=spvMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=$(cd ../ && echo $PWD/fabric-ca/spv/peers/spv-peer/tls/ca.crt)
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/spv/admins/spv-admin/msp)
peer channel join -b ../network/channel-artifacts/securitization.genesis.block

sleep 2

#underwritter-peer
export CORE_PEER_ADDRESS=localhost:11051
export CORE_PEER_LOCALMSPID=underwritterMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=$(cd ../ && echo $PWD/fabric-ca/underwritter/peers/underwritter-peer/tls/ca.crt)
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/underwritter/admins/underwritter-admin/msp)
peer channel join -b ../network/channel-artifacts/securitization.genesis.block

sleep 2

# Set the anchor peers in the network. This let the peers communicate each other

#Originator


export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_LOCALMSPID=originatorMSP
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/originator/admins/originator-admin/msp)
export ORDERER_CA=$(cd ../ && echo $PWD/fabric-ca/originator/orderers/originator-orderer/tls/ca.crt)

peer channel update -c securitization -f ../network/channel-artifacts/originatorMSPanchors.tx -o localhost:7050 --tls --cafile $ORDERER_CA

sleep 2

#farmer


export CORE_PEER_ADDRESS=localhost:8051
export CORE_PEER_LOCALMSPID=farmerMSP
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/farmer/admins/farmer-admin/msp)
export ORDERER_CA=$(cd ../ && echo $PWD/fabric-ca/farmer/orderers/farmer-orderer/tls/ca.crt)

peer channel update -c securitization -f ../network/channel-artifacts/farmerMSPanchors.tx -o localhost:8050 --tls --cafile $ORDERER_CA


sleep 2

#aggregator


export CORE_PEER_ADDRESS=localhost:9051
export CORE_PEER_LOCALMSPID=aggregatorMSP
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/aggregator/admins/aggregator-admin/msp)
export ORDERER_CA=$(cd ../ && echo $PWD/fabric-ca/aggregator/orderers/aggregator-orderer/tls/ca.crt)

peer channel update -c securitization -f ../network/channel-artifacts/aggregatorMSPanchors.tx -o localhost:9050 --tls --cafile $ORDERER_CA


sleep 2

#spv


export CORE_PEER_ADDRESS=localhost:10051
export CORE_PEER_LOCALMSPID=spvMSP
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/spv/admins/spv-admin/msp)
export ORDERER_CA=$(cd ../ && echo $PWD/fabric-ca/spv/orderers/spv-orderer/tls/ca.crt)

peer channel update -c securitization -f ../network/channel-artifacts/spvMSPanchors.tx -o localhost:10050 --tls --cafile $ORDERER_CA

sleep 2

#underwritter


export CORE_PEER_ADDRESS=localhost:11051
export CORE_PEER_LOCALMSPID=underwritterMSP
export CORE_PEER_MSPCONFIGPATH=$(cd ../ && echo $PWD/fabric-ca/underwritter/admins/underwritter-admin/msp)
export ORDERER_CA=$(cd ../ && echo $PWD/fabric-ca/underwritter/orderers/underwritter-orderer/tls/ca.crt)

peer channel update -c securitization -f ../network/channel-artifacts/underwritterMSPanchors.tx -o localhost:11050 --tls --cafile $ORDERER_CA




