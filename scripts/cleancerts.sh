function cleanCA(){
    org=$1
    ca=$2

    CA_PATH=../network/fabric-ca/$org/$ca
    sudo rm -r $CA_PATH/clients
    sudo rm -r $CA_PATH/msp
    sudo rm $CA_PATH/ca-cert.pem
    sudo rm $CA_PATH/tls-cert.pem
    sudo rm $CA_PATH/fabric-ca-server.db
    sudo rm $CA_PATH/IssuerPublicKey
    sudo rm $CA_PATH/IssuerRevocationPublicKey
    CA_CHAIN_FILE=$CA_PATH/ca-chain.pem
    if test -f "$CA_CHAIN_FILE"; then
        sudo rm $CA_CHAIN_FILE
    fi
}

function cleanOrgMSP() {
    org=$1

    MSP_PATH=../network/fabric-ca/$org/msp
    sudo rm -r $MSP_PATH/cacerts
    # sudo rm -r $MSP_PATH/intermediatecerts
    sudo rm -r $MSP_PATH/tlscacerts
    # sudo rm -r $MSP_PATH/tlsintermediatecerts
}

function cleanLocalMSP() {
    org=$1
    name=$2
    type=$3

    LOCAL_MSP_PATH=../network/fabric-ca/$org/${type}s/$name/msp
    TLS_FOLDER_PATH=../network/fabric-ca/$org/${type}s/$name/tls

    sudo rm -r $LOCAL_MSP_PATH
    sudo rm -r $TLS_FOLDER_PATH
}

cleanCA originator ca
cleanCA originator tls-ca
cleanCA farmer ca
cleanCA farmer tls-ca
cleanCA aggregator ca
cleanCA aggregator tls-ca
cleanCA spv ca
cleanCA spv tls-ca
cleanCA underwritter ca
cleanCA underwritter tls-ca


# cleanOrgMSP org1.acme.com
# cleanOrgMSP org2.acme.com
# cleanOrgMSP org3.acme.com
# #cleanOrgMSP acme.com

# cleanLocalMSP org1.acme.com orderer.org1.acme.com orderer
# cleanLocalMSP org2.acme.com orderer.org2.acme.com orderer
# cleanLocalMSP org3.acme.com orderer.org3.acme.com orderer

# cleanLocalMSP org1.acme.com peer0.org1.acme.com peer
# cleanLocalMSP org2.acme.com peer0.org2.acme.com peer
# cleanLocalMSP org3.acme.com peer0.org3.acme.com peer
# #cleanLocalMSP acme.com orderer.acme.com orderer

# cleanLocalMSP org1.acme.com admin@org1.acme.com user
# cleanLocalMSP org2.acme.com admin@org2.acme.com user
# cleanLocalMSP org3.acme.com admin@org3.acme.com user
# #cleanLocalMSP acme.com admin@acme.com user
