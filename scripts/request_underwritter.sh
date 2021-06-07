for i in {1..10}
do
    export NEW_PATH=$(cd .. && echo $PWD)
    export CORE_PEER_ADDRESS=localhost:11051
    export CORE_PEER_LOCALMSPID=underwritterMSP
    export CORE_PEER_TLS_ROOTCERT_FILE=$(cd .. && echo $PWD/fabric-ca/underwritter/peers/underwritter-peer/tls/ca.crt)
    export ORDERER_CA=$(cd .. && echo $PWD/fabric-ca/underwritter/orderers/underwritter-orderer/tls/ca.crt)
    export ORDERER_ADDRESS=localhost:11050
    t=$NEW_PATH/fabric-ca/underwritter/clients/underwritter-client-$i/msp
    export CORE_PEER_MSPCONFIGPATH=$t
    peer chaincode invoke -o $ORDERER_ADDRESS --tls true --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME -c '{ "Args": ["RequestBond","500"] }' --peerAddresses $CORE_PEER_ADDRESS --peerAddresses localhost:7051 --peerAddresses localhost:9051 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --tlsRootCertFiles  $(cd .. && echo $PWD/fabric-ca/originator/peers/originator-peer/tls/ca.crt) --tlsRootCertFiles  $(cd .. && echo $PWD/fabric-ca/aggregator/peers/aggregator-peer/tls/ca.crt)
done
