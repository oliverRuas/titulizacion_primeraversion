#igual lo que hay que cambiar es el OU para la chaincode

export NEW_PATH=$(cd ../ && pwd)

# Originator
export FABRIC_CA_CLIENT_HOME=$NEW_PATH/network/fabric-ca/originator/ca/clients/admin
export CSR_NAMES_ORIGINATOR="C=SP,ST=OriginatorState,L=OriginatorLocation,O=Originator,OU=Hyperledger Fabric"
export FABRIC_CA_CLIENT_TLS_CERTFILES=$NEW_PATH/network/fabric-ca/originator/ca/ca-cert.pem
# Inscribimos la identidad de arranque de la CA
fabric-ca-client enroll -u https://rootca-originator-admin:rootca-originator-adminpw@localhost:7055 --id.type admin --tls.certfiles $FABRIC_CA_CLIENT_TLS_CERTFILES --csr.names "$CSR_NAMES_ORIGINATOR"
fabric-ca-client affiliation add originator
fabric-ca-client affiliation add originator.client1
fabric-ca-client affiliation add originator.peer1
fabric-ca-client affiliation add originator.orderer1
# Inscribimos la identidad de arranque de la TLS CA
export FABRIC_CA_CLIENT_TLS_CERTFILES=$NEW_PATH/network/fabric-ca/originator/tls-ca/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$NEW_PATH/network/fabric-ca/originator/tls-ca/clients/admin
fabric-ca-client enroll -u https://tls-ca-originator-admin:tls-ca-originator-adminpw@localhost:7054 --id.type admin --tls.certfiles $FABRIC_CA_CLIENT_TLS_CERTFILES --csr.names "$CSR_NAMES_ORIGINATOR"


# Farmer
export FABRIC_CA_CLIENT_HOME=$NEW_PATH/network/fabric-ca/farmer/ca/clients/admin
export FABRIC_CA_CLIENT_TLS_CERTFILES=$NEW_PATH/network/fabric-ca/farmer/ca/ca-cert.pem
export CSR_NAMES_FARMER="C=SP,ST=FarmerState,L=FarmerLocation,O=Farmer,OU=Hyperledger Fabric"
# Inscribimos la identidad de arranque de la CA
fabric-ca-client enroll -u https://rootca-farmer-admin:rootca-farmer-adminpw@localhost:8055 --id.type admin --tls.certfiles $FABRIC_CA_CLIENT_TLS_CERTFILES --csr.names "$CSR_NAMES_FARMER"
fabric-ca-client affiliation add farmer
fabric-ca-client affiliation add farmer.client1
fabric-ca-client affiliation add farmer.peer1
fabric-ca-client affiliation add farmer.orderer1
# Inscribimos la identidad de arranque de la TLS CA
export FABRIC_CA_CLIENT_TLS_CERTFILES=$NEW_PATH/network/fabric-ca/farmer/tls-ca/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$NEW_PATH/network/fabric-ca/farmer/tls-ca/clients/admin
fabric-ca-client enroll -u https://tls-ca-farmer-admin:tls-ca-farmer-adminpw@localhost:8054 --id.type admin --tls.certfiles $FABRIC_CA_CLIENT_TLS_CERTFILES --csr.names "$CSR_NAMES_FARMER"


# Aggregator
export FABRIC_CA_CLIENT_HOME=$NEW_PATH/network/fabric-ca/aggregator/ca/clients/admin
export FABRIC_CA_CLIENT_TLS_CERTFILES=$NEW_PATH/network/fabric-ca/aggregator/ca/ca-cert.pem
export CSR_NAMES_AGGREGATOR="C=SP,ST=AggregatorState,L=AggregatorLocation,O=Aggregator,OU=Hyperledger Fabric"
# Inscribimos la identidad de arranque de la CA
fabric-ca-client enroll -u https://rootca-aggregator-admin:rootca-aggregator-adminpw@localhost:9055 --id.type admin --tls.certfiles $FABRIC_CA_CLIENT_TLS_CERTFILES --csr.names "$CSR_NAMES_AGGREGATOR"
fabric-ca-client affiliation add aggregator
fabric-ca-client affiliation add aggregator.client1
fabric-ca-client affiliation add aggregator.peer1
fabric-ca-client affiliation add aggregator.orderer1
# Inscribimos la identidad de arranque de la TLS CA
export FABRIC_CA_CLIENT_TLS_CERTFILES=$NEW_PATH/network/fabric-ca/aggregator/tls-ca/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$NEW_PATH/network/fabric-ca/aggregator/tls-ca/clients/admin
fabric-ca-client enroll -u https://tls-ca-aggregator-admin:tls-ca-aggregator-adminpw@localhost:9054 --id.type admin --tls.certfiles $FABRIC_CA_CLIENT_TLS_CERTFILES --csr.names "$CSR_NAMES_AGGREGATOR"

#SPV
export FABRIC_CA_CLIENT_HOME=$NEW_PATH/network/fabric-ca/spv/ca/clients/admin
export CSR_NAMES_SPV="C=SP,ST=SPVState,L=SPVLocation,O=SPV,OU=Hyperledger Fabric"
export FABRIC_CA_CLIENT_TLS_CERTFILES=$NEW_PATH/network/fabric-ca/spv/ca/ca-cert.pem
# Inscribimos la identidad de arranque de la CA
fabric-ca-client enroll -u https://rootca-spv-admin:rootca-spv-adminpw@localhost:10055 --id.type admin --tls.certfiles $FABRIC_CA_CLIENT_TLS_CERTFILES --csr.names "$CSR_NAMES_SPV"
fabric-ca-client affiliation add spv
fabric-ca-client affiliation add spv.client1
fabric-ca-client affiliation add spv.peer1
fabric-ca-client affiliation add spv.orderer1
# Inscribimos la identidad de arranque de la TLS CA
export FABRIC_CA_CLIENT_TLS_CERTFILES=$NEW_PATH/network/fabric-ca/spv/tls-ca/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$NEW_PATH/network/fabric-ca/spv/tls-ca/clients/admin
fabric-ca-client enroll -u https://tls-ca-spv-admin:tls-ca-spv-adminpw@localhost:10054 --id.type admin --tls.certfiles $FABRIC_CA_CLIENT_TLS_CERTFILES --csr.names "$CSR_NAMES_SPV"

#Underwritter
export FABRIC_CA_CLIENT_HOME=$NEW_PATH/network/fabric-ca/underwritter/ca/clients/admin
export CSR_NAMES_UNDERWRITTER="C=SP,ST=UnderwritterState,L=UnderwritterLocation,O=Underwritter,OU=Hyperledger Fabric"
export FABRIC_CA_CLIENT_TLS_CERTFILES=$NEW_PATH/network/fabric-ca/underwritter/ca/ca-cert.pem
# Inscribimos la identidad de arranque de la CA
fabric-ca-client enroll -u https://rootca-underwritter-admin:rootca-underwritter-adminpw@localhost:11055 --id.type admin --tls.certfiles $FABRIC_CA_CLIENT_TLS_CERTFILES --csr.names "$CSR_NAMES_UNDERWRITTER"
fabric-ca-client affiliation add underwritter
fabric-ca-client affiliation add underwritter.client1
fabric-ca-client affiliation add underwritter.peer1
fabric-ca-client affiliation add underwritter.orderer1
# Inscribimos la identidad de arranque de la TLS CA
export FABRIC_CA_CLIENT_TLS_CERTFILES=$NEW_PATH/network/fabric-ca/underwritter/tls-ca/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$NEW_PATH/network/fabric-ca/underwritter/tls-ca/clients/admin
fabric-ca-client enroll -u https://tls-ca-underwritter-admin:tls-ca-underwritter-adminpw@localhost:11054 --id.type admin --tls.certfiles $FABRIC_CA_CLIENT_TLS_CERTFILES --csr.names "$CSR_NAMES_UNDERWRITTER"



# Register intermediate CA in the root CA
#fabric-ca-client register --id.name int.ca.org1.acme.com --id.secret password --id.attrs 'hf.IntermediateCA=true' -u http://admin:adminpw@localhost:7054
# Register intermediate CA in the tls root CA
#fabric-ca-client register --id.name tls.int.ca.org1.acme.com --id.secret password --id.attrs 'hf.IntermediateCA=true' -u http://admin:adminpw@localhost:7055


