# titulizacion
Esta red está basada en un ejemplo de https://github.com/braduf/curso-hyperledger-fabric
Para desplegar la red se ejecuta el comando ./up.sh y para instalar el chaincode se ejecuta chaincode.sh,  cambiando el identificador del chaincode.

Una vez conseguido esto se pueden ejecutar diversas funciones: issue_bond emite bonos, sim_from_python realiza pagos a traves de un archivo .csv, poolpayments agrupa
los pagos validos cuando periodo==tiempo_cupon, etc

En este caso la organizacion aggregator seria el equivalente al mercado secundario. Habria que modificar los nombres de los archivos .yaml y de los enrollpeers y msp
