#Con esto consultamos el numero de pagos realizados en el periodo
periodo=$1
export NEW_PATH=$(cd .. && echo $PWD)
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_LOCALMSPID=originatorMSP
export CORE_PEER_TLS_ROOTCERT_FILE=$(cd .. && echo $PWD/fabric-ca/originator/peers/originator-peer/tls/ca.crt)
export ORDERER_CA=$(cd .. && echo $PWD/fabric-ca/originator/orderers/originator-orderer/tls/ca.crt)
export ORDERER_ADDRESS=localhost:7050
t=$NEW_PATH/fabric-ca/originator/clients/originator-client/msp
export CORE_PEER_MSPCONFIGPATH=$t
export NUMERO=\'\{\"Args\"\:\[\"QueryAssetNumberByPeriod\",\"$periodo\"\]\}\'
# x=echo $NUMERO
echo peer chaincode invoke -o $ORDERER_ADDRESS --tls true --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME -c $NUMERO --peerAddresses $CORE_PEER_ADDRESS --peerAddresses localhost:8051 --peerAddresses localhost:9051 --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE --tlsRootCertFiles  $(cd .. && echo $PWD/fabric-ca/farmer/peers/farmer-peer/tls/ca.crt) --tlsRootCertFiles  $(cd .. && echo $PWD/fabric-ca/aggregator/peers/aggregator-peer/tls/ca.crt) > transfer_tokens.sh
chmod +x ./transfer_tokens.sh
./transfer_tokens.sh

