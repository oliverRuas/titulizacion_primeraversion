/*
oliver
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"math/rand"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const media float64 = 1000.0
const desv float64 = 10.0

//Defino e inicializo variables que usaré más adelante, que son fijas y que se aceptaran con la aprobación del chaincode
//Años que va a durar el contrato
const years int = 1

//periodos por año
const months int = 6

//número de contratos, que se usará para introducir un default del 1% por año en el modelo del SPV
const numero_contratos int = 100

const rate int = 1

const tiempo_cupon int = 6

const cantidad_bonos int = 5000

const valor_emision float64 = 88.14

const cupon float64 = 2

var numeros_defaulters []int = numeros_aleatorios(rate, years)

//se puede quitar
var periodo int

//se puede quitar
//var rand.Seed(time.Now().UTC().UnixNano())

//defino esta estructura que implementará la lógica del modulo Contract del paquete contractapi
type SmartContract struct {
	contractapi.Contract
}

//Defino una estructura de un activo llamado PPA que tiene estas
//estas propiedades (ojo, atributos de una struct en go con primera
//letra en Mayusc) y definimos la representacion json de estos atributos
//que es como se va a guardar en el ledger
type PPA struct {
	DocType  string  `json:"docType"`
	Client   string  `json:"client"`
	Energy   float64 `json:"energy"`
	Default  bool    `json:"default"`
	Payments float64 `json:"payments"`
	Fecha    Datos
	Period   int `json:"periodo"`
}

type Datos struct {
	Day   int
	Month time.Month
	Year  int
}

type Pagos struct {
	Doctype string  `json:"docType"`
	Total   float64 `json:"total"`
	Owner   string  `json:"propietario"`
}

type Pool struct {
	Doctype string  `json:"docType"`
	Total   float64 `json:"total"`
	Balance float64 `json:"balance"`
	Client  string  `json:"cliente"`
	State   string  `json:"estado"`
}

//Estructura que se usará para obtener las ID de los clients de la organizacion farmer
type FarmerID struct {
	Doctype   string `json:"docType"`
	Identidad string `json:"identidad"`
}

//estructura que se usará para calcular el total de los payments del modelo del SPV
type ValorTotal struct {
	Doctype string  `json:"docType"`
	Total   float64 `json:"total"`
}

type PagosImpagos struct {
	Payments float64 `json:"pagos"`
	Default  bool    `json:"impago"`
}

type BondUTXO struct {
	Key               string  `json:"utxo_key"`
	Issuer            string  `json:"issuer"`
	Investor          string  `json:"owner"`
	Amount            int     `json:"amount"`
	CouponRate        float64 `json:"couponrate"`
	Value             string  `json:"value"` // Expected payout
	MaturityDate      string  `json:"maturitydate"`
	AmountPaid        string  `json:"amountpaid"`
	OriginalValue     string  `json:"originalvalue"`
	MonthlyPayout     string  `json:"monthlypayout"`
	RemainingPayments string  `json:"remainingpayments"`
}

type Request struct {
	Doctype    string `json:"docType"`
	Requester  string `json:"bonista"`
	Amount     int    `json:"cantidad"`
	KeyRequest string `json:"clave"`
	//	RequestedAt	*time.Location	`json:"peticion"`
}

type State struct {
	DocType string `json:"docType"`
	Value   interface{}
}

type Bond struct {
	Doctype      string  `json:"docType"`
	ValorEmision float64 `json:"valorEmision"`
	Cupon        float64 `json:"cupon"`
	Vencimiento  int     `json:"vencimiento"`
	TiempoPagos  int     `json:"tiempoPagos"`
}

// CurrencyEventNames specifies the names of the events that should be fired after the txs
var CurrencyEventNames = map[string]string{
	"Mint":     "Minted",
	"Transfer": "Transfered",
	"Redeem":   "Redeemed",
}

// MintedPayload is the payload of the Minted Events
type MintedPayload struct {
	Minter       string `json:"minter"`
	UTXOID       string `json:"UtxoId"`
	Receiver     string `json:"receiver"`
	CurrencyCode string `json:"currencyCode"`
}

// TransferedPayload is the payload of the Transfered Events
type TransferedPayload struct {
	TransferedBy string `json:"transferedBy"`
	//SpentUTXOIDSet   []string `json:"spentUtxoIdSet"`
	ChangeUTXOID     string `json:"changeUtxoId"`
	TransferedUTXOID string `json:"transferedUtxoId"`
	Receiver         string `json:"receiver"`
	CurrencyCode     string `json:"currencyCode"`
}

// RedeemPayload is the payload of the RedemptionRequested Events
type RedeemPayload struct {
	Requestor    string `json:"requestor"`
	Redeemer     string `json:"redeemer"`
	UTXOID       string `json:"utxoID"`
	CurrencyCode string `json:"currencyCode"`
}

// RedeemPrivateData is the private data that should be shared in a redeem transaction
type RedeemPrivateData struct {
	UtxoID        string `json:"utxoID"`
	AccountNumber string `json:"accountNumber"`
	Bank          string `json:"bank"`
}

func (s *SmartContract) BondInfo(ctx contractapi.TransactionContextInterface) *Bond {
	x := &Bond{
		Doctype:      "Bono",
		ValorEmision: valor_emision,
		Cupon:        cupon,
		Vencimiento:  years * months,
		TiempoPagos:  tiempo_cupon}
	return x
}

func fecha_inicio() Datos {
	return Datos{
		Day:   12,
		Month: time.Month(5),
		Year:  2021,
	}
}

func queryperiod() Datos {
	fecha := Datos{
		Day:   time.Now().Day(),
		Month: time.Now().Month(),
		Year:  time.Now().Year(),
	}
	return fecha
}

func numeros_aleatorios(rate int, anhos int) []int {
	//rand.Seed(time.Now().UTC().UnixNano())
	x := rate * numero_contratos / 100
	xx := x * anhos
	var m []int
	var r int
	//var size int
	for i := 0; len(m) < xx; i++ {
		r = rand.Intn(numero_contratos + 1)
		for _, l := range m {
			if r == l || r == 0 {
				r = rand.Intn(numero_contratos + 1)
			}
		}
		m = append(m, r)
	}
	return m
}

//buscar que esto solo se pueda ejecutar una sola vez
func issue(ctx contractapi.TransactionContextInterface) error {
	hasOU, err := cid.HasOUValue(ctx.GetStub(), "cliente")
	if err != nil {
		return err
	}
	if !hasOU {
		return ErrNoSPV
	}
	identity := ctx.GetClientIdentity()
	spv, err := identity.GetID()
	if err != nil {
		return err
	}
	org, err := identity.GetMSPID()
	if err != nil {
		return err
	}
	if org != "spvMSP" {
		return ErrNoSPV
	}
	spvID := &FarmerID{
		Doctype:   "SPVidentidad",
		Identidad: spv,
	}
	spvIDasBytes, err := json.Marshal(spvID)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}
	spvKey := ctx.GetStub().GetTxID()
	return ctx.GetStub().PutState(spvKey, spvIDasBytes)
}

//faltan comprobaciones de identidad
// func (s *SmartContract) IssueBond(ctx contractapi.TransactionContextInterface) error {
// 	hasOU, err := cid.HasOUValue(ctx.GetStub(), "cliente")
// 	if err != nil {
// 		return err
// 	}
// 	if !hasOU {
// 		return ErrNoSPV
// 	}
// 	identity := ctx.GetClientIdentity()
// 	// spv, err := identity.GetID()
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	org, err := identity.GetMSPID()
// 	if err != nil {
// 		return err
// 	}
// 	if org != "spvMSP" {
// 		return ErrNoSPV
// 	}
// 	fechainicio := fecha_inicio()
// 	fecha := &Datos{
// 		Day:   time.Now().Day(),
// 		Month: time.Now().Month(),
// 		Year:  time.Now().Year(),
// 	}
// 	log.Printf("fecha actual: %v", fecha)
// 	log.Printf("fecha inicio: %v", fechainicio)

// 	log.Printf("anho: %v", fecha.Year)
// 	log.Printf("anho inicio: %v", fechainicio.Year)
// 	log.Printf("mes: %v", fecha.Month)
// 	log.Printf("mes inicio: %v", fechainicio.Month)
// 	if fecha.Month != fechainicio.Month || fecha.Year != fechainicio.Year {
// 		return fmt.Errorf("Not possible to issue bond")
// 	}
// 	identidades, err := s.QueryIdentitiesSPV(ctx)
// 	log.Printf("identidades: %v", identidades)
// 	if identidades == nil {
// 		mint(ctx, cantidad_bonos)
// 		err = issue(ctx)
// 		if err != nil {
// 			fmt.Printf("error: %s", err.Error())
// 			return err
// 		}
// 	}else{
// 		return ErrSPVBond
// 	}
// 	return nil
// }

func (s *SmartContract) PoolPaymentsOriginator(ctx contractapi.TransactionContextInterface) error {
	//Comprobamos que no se hayan agrupado los pagos
	//añadido 27 abril
	fecha1 := fecha_inicio()
	fecha2 := queryperiod()
	log.Printf("fecha actual: %v", fecha2)
	var y1 int = int(fecha2.Month)
	var y0 int = int(fecha1.Month)
	periodo := (12*fecha2.Year + y1) - (12*fecha1.Year + y0) + 1

	log.Printf("numero de periodo: %v", periodo)
	///////////////////////

	periodoAsString := strconv.Itoa(periodo)
	oldpoolKey, err := ctx.GetStub().CreateCompositeKey("pool", []string{periodoAsString})
	if err != nil {
		return fmt.Errorf("failed to read")
	}
	valueBytes, err := ctx.GetStub().GetState(oldpoolKey)
	if err != nil {
		return fmt.Errorf("failed to read")
	}
	if valueBytes != nil {
		return fmt.Errorf("These payments have already been pooled")
	}
	hasOU, err := cid.HasOUValue(ctx.GetStub(), "cliente")
	if err != nil {
		return err
	}
	if !hasOU {
		return ErrNoOriginator
	}
	identity := ctx.GetClientIdentity()
	spv, err := identity.GetID()
	if err != nil {
		return err
	}
	org, err := identity.GetMSPID()
	if err != nil {
		return err
	}
	if org != "originatorMSP" {
		return ErrNoOriginator
	}
	if periodo%tiempo_cupon != 0 {
		return ErrNoPeriod
	}
	//consulto que para ese periodo se hayan emitido todos los pagos
	num, _ := s.QueryAssetNumberByPeriod(ctx, periodo)
	if num < numero_contratos {
		return ErrNoPeriod
	}
	cociente := periodo / tiempo_cupon
	periodo_anterior := (tiempo_cupon * (cociente - 1)) + 1
	var nuevoPPA []*PagosImpagos
	var valor float64
	var nuevoValor float64 = 0
	var total float64
	for j := periodo_anterior; j < periodo+1; j++ {
		total = 0
		nuevoPPA, err = s.QueryPaymentsAndDefaultByPeriod(ctx, j)
		for _, k := range nuevoPPA {
			if !k.Default {
				valor = k.Payments
				total = total + valor
			}
		}
		nuevoValor = nuevoValor + total
	}
	log.Printf("Valor total: %v", nuevoValor)

	pool := &Pool{
		Doctype: "pool",
		Total:   nuevoValor,
		Balance: nuevoValor,
		Client:  spv,
		State:   "ISSUED",
	}
	//periodoAsString,_:=strconv.Atoi(periodo)

	poolAsBytes, err := json.Marshal(&pool)
	//si no quiero validar el err, defino como elemento, _ :=json.Marshal()
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		// return err
	}
	poolKey, err := ctx.GetStub().CreateCompositeKey("pool", []string{periodoAsString})
	ctx.GetStub().PutState(poolKey, poolAsBytes)
	return nil
}

func (s *SmartContract) QueryClientUTXOs(ctx contractapi.TransactionContextInterface, client string) ([]string, error) {
	// since utxos have a composite key of owner:utxoKey, we can query for all utxos matching owner:*
	utxoResultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("utxo", []string{client})
	if err != nil {
		return nil, err
	}
	defer utxoResultsIterator.Close()

	var newUTXO UTXO
	var utxos []*UTXO
	for utxoResultsIterator.HasNext() {
		utxoRecord, err := utxoResultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// composite key is expected to be owner:utxoKey
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(utxoRecord.Key)
		if err != nil {
			return nil, err
		}

		if len(compositeKeyParts) != 2 {
			return nil, fmt.Errorf("expected composite key with two parts (owner:utxoKey)")
		}

		utxoKey := compositeKeyParts[1] // owner is at [0], utxoKey is at[1]

		if utxoRecord.Value == nil {
			return nil, fmt.Errorf("utxo %s has no value", utxoKey)
		}

		amount, _ := strconv.Atoi(string(utxoRecord.Value)) // Error handling not needed since Itoa() was used when setting the utxo amount, guaranteeing it was an integer.

		utxo := &UTXO{
			Key:    utxoKey,
			Owner:  client,
			Amount: amount,
		}
		newUTXO.Key = utxoKey
		newUTXO.Amount = amount
		utxos = append(utxos, utxo)
	}
	//log.Printf("Este cliente tiene un UTXO")
	return []string{newUTXO.Key, strconv.Itoa(newUTXO.Amount)}, nil
}

func (s *SmartContract) DistributePayments(ctx contractapi.TransactionContextInterface, client string) error {
	///////////////////////////
	//añadido 27 de abril
	fecha1 := fecha_inicio()
	fecha2 := queryperiod()
	// log.Printf("fecha actual: %v", fecha2)
	var y1 int = int(fecha2.Month)
	var y0 int = int(fecha1.Month)
	periodo := (12*fecha2.Year + y1) - (12*fecha1.Year + y0) + 1

	log.Printf("numero de periodo: %v", periodo)
	////////////////////////////////////////
	periodoAsString := strconv.Itoa(periodo)
	oldpoolKey, _ := ctx.GetStub().CreateCompositeKey("pool", []string{periodoAsString})
	// if err != nil {
	// 	return fmt.Errorf("failed to read")
	// }
	valueBytes, _ := ctx.GetStub().GetState(oldpoolKey)
	// if err != nil {
	// 	return fmt.Errorf("failed to read")
	// }
	if valueBytes == nil {
		return fmt.Errorf("These payments havent been pooled")
	}

	oldpagosKey, _ := ctx.GetStub().CreateCompositeKey("pagos", []string{periodoAsString, client})
	// if err!=nil{
	// 	return fmt.Errorf("failed to read")
	// }
	pagosBytes, _ := ctx.GetStub().GetState(oldpagosKey)
	// if err != nil {
	// 	return fmt.Errorf("failed to read")
	// }
	if pagosBytes != nil {
		return fmt.Errorf("These payments have already been distributed")
	}
	hasOU, err := cid.HasOUValue(ctx.GetStub(), "cliente")
	if err != nil {
		return err
	}
	if !hasOU {
		return ErrNoSPV
	}
	identity := ctx.GetClientIdentity()
	//spv, err := identity.GetID()
	if err != nil {
		return err
	}
	org, _ := identity.GetMSPID()
	// if err != nil {
	// 	return err
	// }
	if org != "spvMSP" {
		return ErrNoSpv
	}
	var pool Pool
	json.Unmarshal(valueBytes, &pool)
	// if err != nil {
	// 	return err
	// }
	//Compruebo cuantos bonos tiene el cliente en el instante en el que se produce la transaccion
	ss, _ := s.QueryClientUTXOs(ctx, client)
	// if err!=nil{
	// 	return fmt.Errorf("Error: %v",err)
	// }
	//comprobar si los bonos estan en distintos documentos (el vector de utxos)
	cant, _ := strconv.Atoi(ss[1])
	cant1 := float64(cant)
	log.Printf("Bonos en posesion: %v", cant)

	// valor_cupon:=cupon*valor_nominal/100
	// valor_final:=((cupon/100)+1)*valor_nominal
	pagos := &Pagos{}
	if periodo != years*months {
		pagos = &Pagos{
			Doctype: "pagos",
			Owner:   client,
			Total:   cant1 * (pool.Total / float64(cantidad_bonos)),
			//Total:		cant1*valor_cupon,
		}
	} else {
		pagos = &Pagos{
			Doctype: "pagos",
			Owner:   client,
			Total:   cant1 * (pool.Total / float64(cantidad_bonos)),
			//Total:		cant1*valor_final,
		}
	}

	balance := pool.Balance - pagos.Total
	if balance < 0 {
		return fmt.Errorf("Imposible de realizar la tx")
	} else if balance == 0 {
		pool.State = "REDEEMED"
		pool.Balance = balance
	} else {
		pool.State = "TRANSFERRED"
		pool.Balance = balance
	}
	//cant=float64(cant)
	pagosAsBytes, _ := json.Marshal(&pagos)
	pagosKey, _ := ctx.GetStub().CreateCompositeKey("pagos", []string{periodoAsString, client})
	ctx.GetStub().PutState(pagosKey, pagosAsBytes)
	// if err!=nil{
	// 	return err
	// }
	poolAsBytes, _ := json.Marshal(&pool)
	//si no quiero validar el err, defino como elemento, _ :=json.Marshal()
	// if err != nil {
	// 	fmt.Printf("Marshal error: %s", err.Error())
	// 	return err
	// }
	poolKey, _ := ctx.GetStub().CreateCompositeKey("pool", []string{periodoAsString})
	return ctx.GetStub().PutState(poolKey, poolAsBytes)
	// jsonAsBytes, _ := json.Marshal(`"{sending payments to receipient: ` + client + ` }"`)
	// return ctx.GetStub().SetEvent("Payments distributed: ",jsonAsBytes)
}

func (s *SmartContract) RequestBond(ctx contractapi.TransactionContextInterface, cantidad int) error {
	if cantidad > cantidad_bonos {
		return ErrWrongNumber
	}
	hasOU, err := cid.HasOUValue(ctx.GetStub(), "cliente")
	if err != nil {
		return err
	}
	if !hasOU {
		return ErrNoUnderwritter
	}
	identity := ctx.GetClientIdentity()
	underwritter, err := identity.GetID()
	if err != nil {
		return err
	}
	org, err := identity.GetMSPID()
	if err != nil {
		return err
	}
	x := 0
	orgs := []string{"underwritterMSP", "aggregatorMSP"}
	for _, b := range orgs {
		if org != b {
			x = x + 1
		}
	}
	if x == 2 {
		return ErrNoUnderwritter
	}
	requestKey := ctx.GetStub().GetTxID()
	finalKey, _ := ctx.GetStub().CreateCompositeKey("request", []string{underwritter, requestKey})
	request := &Request{
		Doctype:    "request",
		Requester:  underwritter,
		Amount:     cantidad,
		KeyRequest: finalKey,
		//RequestedAt: 	time.UTC,
	}
	requestAsBytes, err := json.Marshal(&request)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}

	//request.KeyRequest=requestKey
	return ctx.GetStub().PutState(finalKey, requestAsBytes)
}

func (s *SmartContract) TransferTokens(ctx contractapi.TransactionContextInterface, receiver string) error {
	//chequeamos que el que recibe se ha inscrito (falta)
	idUTXO, err := s.ClientUTXOs(ctx)
	//log.Printf("obtiene algo el clientUTXO? %v",idUTXO)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	var total int = 0
	var keys []string
	requests, _ := s.QueryRequests(ctx, receiver)
	for _, k := range requests {
		total = total + k.Amount
		keys = append(keys, k.KeyRequest)
	}
	if total == 0 {
		return ErrWrongNumber
	}
	//log.Printf("Total de bonos que voy a asignar: %v", total)
	value := idUTXO[0]
	cant, _ := strconv.Atoi(idUTXO[1])

	if cant-total < 0 || cant == 0 {
		return ErrWrongNumber
	}
	_, err = transfer(ctx, []string{value}, cant-total, receiver, total)
	if err != nil {
		return err
	}
	// jsonAsBytes, _ := json.Marshal(`"{sending: ` + strconv.Itoa(total) + ` tokens to receipient: ` + receiver + ` }"`)
	// er:=ctx.GetStub().SetEvent("token transferred: ",jsonAsBytes)
	// if er!=nil{
	// 	return er
	// }
	for _, r := range keys {
		ctx.GetStub().DelState(r)
		log.Printf("estado elimininado: %v", r)
	}
	return nil
}

func transfer(ctx contractapi.TransactionContextInterface, utxoInputKeys []string, otro_amount int, receiver string, amount int) (*TokenUTXO, error) {
	//var issuer string
	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client id: %v", err)
	}

	// Validate and summarize utxo inputs
	utxoInputs := make(map[string]*TokenUTXO)
	var totalInputAmount int
	for _, utxoInputKey := range utxoInputKeys {
		if utxoInputs[utxoInputKey] != nil {
			return nil, fmt.Errorf("the same utxo input can not be spend twice")
		}

		utxoInputCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{clientID, utxoInputKey})
		if err != nil {
			return nil, fmt.Errorf("failed to create composite key: %v", err)
		}
		// validate that client has a utxo matching the input key
		valueBytes, err := ctx.GetStub().GetState(utxoInputCompositeKey)
		if err != nil {
			return nil, fmt.Errorf("failed to read utxoInputCompositeKey %s from world state: %v", utxoInputCompositeKey, err)
		}

		if valueBytes == nil {
			return nil, fmt.Errorf("utxoInput %s not found for client %s", utxoInputKey, clientID)
		}
		//	log.Printf("valor en bytes: %v",valueBytes)
		//amount, _ := strconv.Atoi(string(valueBytes)) // Error handling not needed since Itoa() was used when setting the utxo amount, guaranteeing it was an integer.

		utxoInput := &TokenUTXO{
			Doctype: "token",
			Key:     utxoInputKey,
			Owner:   clientID,
			Amount:  amount,
		}

		totalInputAmount += amount
		utxoInputs[utxoInputKey] = utxoInput
	}
	// Since the transaction is valid, now delete utxo inputs from owner's state
	//issuer=utxoInput.Issuer
	new_UTXO := new(TokenUTXO)
	new_UTXO.Key = ctx.GetStub().GetTxID() + ".0"
	//utxoOutput.Issuer:=issuer
	new_UTXO.Doctype = "token"
	new_UTXO.Owner = clientID
	new_UTXO.Amount = otro_amount
	new_UTXOCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{new_UTXO.Owner, new_UTXO.Key})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %v", err)
	}
	err = ctx.GetStub().PutState(new_UTXOCompositeKey, []byte(strconv.Itoa(new_UTXO.Amount)))
	if err != nil {
		return nil, err
	}

	for _, utxoInput := range utxoInputs {
		utxoInputCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{utxoInput.Owner, utxoInput.Key})
		if err != nil {
			return nil, fmt.Errorf("failed to create composite key: %v", err)
		}
		err = ctx.GetStub().DelState(utxoInputCompositeKey)
		if err != nil {
			return nil, err
		}
		//	log.Printf("utxoInput deleted: %+v", utxoInput)
	}
	// Create utxo outputs using a composite key based on the owner and utxo key
	//for _, utxoOutput := range utxoOutputs {
	utxoOutput := new(TokenUTXO)
	utxoOutput.Key = ctx.GetStub().GetTxID() + ".0"
	//utxoOutput.Issuer:=issuer
	utxoOutput.Doctype = "token"
	utxoOutput.Owner = receiver
	utxoOutput.Amount = totalInputAmount
	utxoOutputCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{utxoOutput.Owner, utxoOutput.Key})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %v", err)
	}
	err = ctx.GetStub().PutState(utxoOutputCompositeKey, []byte(strconv.Itoa(utxoOutput.Amount)))
	if err != nil {
		return nil, err
	}
	//log.Printf("utxoOutput created: %+v", utxoOutput)
	//}
	return utxoOutput, nil
}

//ejemplo fabric-sample
//esta función crea "amount" UTXOs que tiene que ser entero (para hacer una equivalencia con los pagos: pagos*100)
//la creación se hace en el momento en que se registra el pago por parte del client del farmer si el default es false. El client no controla que
//cantidad emite, es directamente el pago*100 y no va a poder eliminar el UTXO (falta por implementar DELETEUTXO)
//Se asegura que la ID del utxo es unica, pues se crea a partir del usuario y de la ID de la transaccion
//seria interesante cambiar a que la funcion no se pudiese ejecutar, que fuese una interfaz



func IssueBond(ctx contractapi.TransactionContextInterface, receiver string) (err error) {
// func mint(ctx contractapi.TransactionContextInterface, amount int, receiver string) (payload MintedPayload, err error) {
	// Check minter authorization - this sample assumes Org1 is the central banker with privilege to mint new tokens
	hasOU, err := cid.HasOUValue(ctx.GetStub(), "cliente")
	if err != nil {
		return err
	}
	if !hasOU {
		return ErrNoSPV
	}
	identity := ctx.GetClientIdentity()
	org, err := identity.GetMSPID()
	if err != nil {
		return err
	}
	if org != "spvMSP" {
		return ErrNoSPV
	}
	// Get ID of submitting client identity
	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return
		// return nil, fmt.Errorf("failed to get client id: %v", err)
	}
	fechainicio := fecha_inicio()
	fecha := &Datos{
		Day:   time.Now().Day(),
		Month: time.Now().Month(),
		Year:  time.Now().Year(),
	}
	log.Printf("fecha actual: %v", fecha)
	log.Printf("fecha inicio: %v", fechainicio)

	log.Printf("anho: %v", fecha.Year)
	log.Printf("anho inicio: %v", fechainicio.Year)
	log.Printf("mes: %v", fecha.Month)
	log.Printf("mes inicio: %v", fechainicio.Month)
	if fecha.Month != fechainicio.Month || fecha.Year != fechainicio.Year {
		return fmt.Errorf("Not possible to issue bond")
	}
	if amount <= 0 {
		return
		// return nil, fmt.Errorf("mint amount must be a positive integer")
	}
	identidades, err := s.QueryIdentitiesSPV(ctx)
	log.Printf("identidades: %v", identidades)
	amount:=numero_contratos
	if identidades == nil {
		utxo := UTXOBond{
			utxo.Key : ctx.GetStub().GetTxID() + ".0",
			utxo.Issuer : ctx.GetClientIdentity().GetMSPID(),
			utxo.Investor : receiver,
			utxo.Amount : amount,
			utxo.CouponRate : couponrate,
			utxo.Value      : ,
			utxo.MaturityDate : ,
			utxo.AmountPaid    : ,
			utxo.OriginalValue : ,
			utxo.MonthlyPayout : ,
			utxo.RemainingPayments : ,
			}
		// the utxo has a composite key of owner:utxoKey, this enables ClientUTXOs() function to query for an owner's utxos.
		utxoCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{utxo.Key})
	
		doctype:= "utxoBond"
		state:=State{
			doctype,
			utxo,
		}
	
		value, err := json.Marshal(&state)
		if err != nil {
			return
		}
	
		err = ctx.GetStub().PutState(utxoCompositeKey, value)
		if err != nil {
			return nil, err
		}
		err = issue(ctx)
		if err != nil {
			fmt.Printf("error: %s", err.Error())
			return err
		}
	}else{
		return ErrSPVBond
	}
	// ctx.GetStub().SetEvent()	
	return
}

//ejemplo fabric-samples
//esta función de usa para obtener la ID del usuario (no es necesaria ya que son 3 lineas de codigo)
func (s *SmartContract) ClientID(ctx contractapi.TransactionContextInterface) (string, error) {

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client id: %v", err)
	}

	return clientID, nil
}

//Ejemplo fabric-sample
// Transfer transfers UTXOs containing tokens from client to recipient(s)
//Cambiar a que la funcion pase a ser una interfaz
//Esta función está cambiada respecto a la del ejemplo, no se transfiere una cantidad, sino que se transfiere el total de cada client del
//farmer, por lo que si el mint es correcto no necesita comprobar los fondos de los UTXO del client

//esta funcion va sumando todos los UTXO de todos los clientes y los va eliminando. Cuando ha terminado emite un nuevo token
//a nombre del originador con la cantidad=suma(utxos)
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, utxoInputKeys []string, amount int) (*UTXO, error) {
	hasOU, err := cid.HasOUValue(ctx.GetStub(), "cliente")
	if err != nil {
		return nil, err
	}
	if !hasOU {
		return nil, ErrNoFarmer
	}
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != "farmerMSP" {
		return nil, fmt.Errorf("client is not authorized to receive new tokens")
	}

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client id: %v", err)
	}

	// Validate and summarize utxo inputs
	utxoInputs := make(map[string]*UTXO)
	var totalInputAmount int
	for _, utxoInputKey := range utxoInputKeys {
		if utxoInputs[utxoInputKey] != nil {
			return nil, fmt.Errorf("the same utxo input can not be spend twice")
		}

		utxoInputCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{clientID, utxoInputKey})
		if err != nil {
			return nil, fmt.Errorf("failed to create composite key: %v", err)
		}

		// validate that client has a utxo matching the input key
		valueBytes, err := ctx.GetStub().GetState(utxoInputCompositeKey)
		if err != nil {
			return nil, fmt.Errorf("failed to read utxoInputCompositeKey %s from world state: %v", utxoInputCompositeKey, err)
		}

		if valueBytes == nil {
			return nil, fmt.Errorf("utxoInput %s not found for client %s", utxoInputKey, clientID)
		}

		//amount, _ := strconv.Atoi(string(valueBytes)) // Error handling not needed since Itoa() was used when setting the utxo amount, guaranteeing it was an integer.

		utxoInput := &UTXO{
			Key:    utxoInputKey,
			Owner:  clientID,
			Amount: amount,
		}

		totalInputAmount += amount
		utxoInputs[utxoInputKey] = utxoInput
	}
	// Validate and summarize utxo outputs
	//var totalOutputAmount int
	//txID := ctx.GetStub().GetTxID()
	// for i, utxoOutput := range utxoOutputs {

	//  	if utxoOutput.Amount <= 0 {
	//  		return nil, fmt.Errorf("utxo output amount must be a positive integer")
	//  	}

	//  	utxoOutputs[i].Key = fmt.Sprintf("%s.%d", txID, i)

	//  	totalOutputAmount += utxoOutput.Amount
	//  }

	//  // Validate total inputs equals total outputs
	//  if totalInputAmount != totalOutputAmount {
	//  	return nil, fmt.Errorf("total utxoInput amount %d does not equal total utxoOutput amount %d", totalInputAmount, totalOutputAmount)
	// }

	// Since the transaction is valid, now delete utxo inputs from owner's state
	for _, utxoInput := range utxoInputs {

		utxoInputCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{utxoInput.Owner, utxoInput.Key})
		if err != nil {
			return nil, fmt.Errorf("failed to create composite key: %v", err)
		}

		err = ctx.GetStub().DelState(utxoInputCompositeKey)
		if err != nil {
			return nil, err
		}
		log.Printf("utxoInput deleted: %+v", utxoInput)
	}

	// Create utxo outputs using a composite key based on the owner and utxo key
	//for _, utxoOutput := range utxoOutputs {
	utxoOutput := new(UTXO)
	utxoOutput.Key = ctx.GetStub().GetTxID() + ".0"
	mspid := "originatorMSP"
	utxoOutput.Owner = mspid
	utxoOutput.Amount = totalInputAmount
	utxoOutputCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{utxoOutput.Owner, utxoOutput.Key})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %v", err)
	}

	err = ctx.GetStub().PutState(utxoOutputCompositeKey, []byte(strconv.Itoa(utxoOutput.Amount)))
	if err != nil {
		return nil, err
	}
	log.Printf("utxoOutput created: %+v", utxoOutput)

	return utxoOutput, nil
}

//ejemplo fabric-samples
// ClientUTXOs returns all UTXOs owned by the calling client
//cambiar a interfaz
func (s *SmartContract) ClientUTXOs(ctx contractapi.TransactionContextInterface) ([]string, error) {

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client id: %v", err)
	}

	// since utxos have a composite key of owner:utxoKey, we can query for all utxos matching owner:*
	utxoResultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("utxo", []string{clientID})
	if err != nil {
		return nil, err
	}
	defer utxoResultsIterator.Close()

	var newUTXO UTXO
	var utxos []*UTXO
	for utxoResultsIterator.HasNext() {
		utxoRecord, err := utxoResultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// composite key is expected to be owner:utxoKey
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(utxoRecord.Key)
		if err != nil {
			return nil, err
		}

		if len(compositeKeyParts) != 2 {
			return nil, fmt.Errorf("expected composite key with two parts (owner:utxoKey)")
		}

		utxoKey := compositeKeyParts[1] // owner is at [0], utxoKey is at[1]

		if utxoRecord.Value == nil {
			return nil, fmt.Errorf("utxo %s has no value", utxoKey)
		}

		amount, _ := strconv.Atoi(string(utxoRecord.Value)) // Error handling not needed since Itoa() was used when setting the utxo amount, guaranteeing it was an integer.

		utxo := &UTXO{
			Key:    utxoKey,
			Owner:  clientID,
			Amount: amount,
		}
		newUTXO.Key = utxoKey
		newUTXO.Amount = amount
		utxos = append(utxos, utxo)
	}
	//log.Printf("Este cliente tiene un UTXO")
	return []string{newUTXO.Key, strconv.Itoa(newUTXO.Amount)}, nil
}

//adaptacion fabric-samples
//funcion que suma todos los UTXO emitidos por los clients del farmer y los pone a nombre del originador.
// func (s *SmartContract) ClientUTXOoriginator(ctx contractapi.TransactionContextInterface) (*UTXO, error) {

// 	// // Get ID of submitting client identity
// 	identity := ctx.GetClientIdentity()
// 	clientID, err := identity.GetID()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get client id: %v", err)
// 	}

// 	org, err := identity.GetMSPID()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if org != "originatorMSP" {
// 		return nil, ErrNoOriginator
// 	}
// 	// since utxos have a composite key of owner:utxoKey, we can query for all utxos matching owner:*
// 	utxoResultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("utxo", []string{"originatorMSP"})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer utxoResultsIterator.Close()

// 	var total int = 0
// 	var oldClientID string = "originatorMSP"
// 	//var newUTXO UTXO
// 	//var utxos []*UTXO
// 	for utxoResultsIterator.HasNext() {
// 		utxoRecord, err := utxoResultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		// composite key is expected to be owner:utxoKey
// 		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(utxoRecord.Key)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if len(compositeKeyParts) != 2 {
// 			return nil, fmt.Errorf("expected composite key with two parts (owner:utxoKey)")
// 		}

// 		utxoKey := compositeKeyParts[1] // owner is at [0], utxoKey is at[1]

// 		if utxoRecord.Value == nil {
// 			return nil, fmt.Errorf("utxo %s has no value", utxoKey)
// 		}

// 		amount, _ := strconv.Atoi(string(utxoRecord.Value)) // Error handling not needed since Itoa() was used when setting the utxo amount, guaranteeing it was an integer.

// 		total = total + amount
// 		utxo := &UTXO{
// 			Key:    utxoKey,
// 			Owner:  oldClientID,
// 			Amount: total,
// 		}

// 		utxoCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{utxo.Owner, utxo.Key})
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create composite key: %v", err)
// 		}

// 		err = ctx.GetStub().DelState(utxoCompositeKey)
// 		if err != nil {
// 			return nil, err
// 		}
// 		log.Printf("utxoInput deleted: %+v", utxoCompositeKey)
// 	}

// 	utxoKey := ctx.GetStub().GetTxID() + ".0"
// 	utxo := &UTXO{
// 		Key:    utxoKey,
// 		Owner:  clientID,
// 		Amount: total,
// 	}

// 	utxoCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{utxo.Owner, utxo.Key})

// 	log.Printf("total de los UTXOs: %v", total)

// 	err = ctx.GetStub().PutState(utxoCompositeKey, []byte(strconv.Itoa(utxo.Amount)))
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Printf("utxo created: %+v", utxo)

// 	return utxo, nil
// }

//función que permite a los clients del farmer emitir una transaccion con su identidad (publica) al ledger
//puede ser interesante que esto se transfiera como PrivateData entre spv, originator y farmer, o con un canal aparte
func (s *SmartContract) RegisteringFarmers(ctx contractapi.TransactionContextInterface) error {
	hasOU, err := cid.HasOUValue(ctx.GetStub(), "cliente")
	if err != nil {
		return err
	}
	if !hasOU {
		return ErrNoFarmer
	}

	identity := ctx.GetClientIdentity()
	farmer, err := identity.GetID()
	if err != nil {
		return err
	}
	org, err := identity.GetMSPID()
	if err != nil {
		return err
	}
	if org != "farmerMSP" {
		return ErrNoFarmer
	}
	log.Printf("Pasa todos los filtros de identidad")

	farmerID := &FarmerID{
		Doctype:   "identidad",
		Identidad: farmer,
	}
	farmerIDasBytes, err := json.Marshal(farmerID)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}
	farmerKey := ctx.GetStub().GetTxID()
	log.Printf("ppakey: %v", farmerKey)

	return ctx.GetStub().PutState(farmerKey, farmerIDasBytes)
}

//con esta función consultamos el ID de los clients del farmer que se han registrado. Esto solo se puede hacer si se pertenece a la org
//SPV, que usara esta lista para asignar valores de energia y de default para usar en su modelo
func (s *SmartContract) QueryFarmerID(ctx contractapi.TransactionContextInterface) ([]*FarmerID, error) {
	var identidades []*FarmerID
	identity := ctx.GetClientIdentity()

	org, err := identity.GetMSPID()
	if err != nil {
		return nil, err
	}
	if org != "spvMSP" {
		return nil, ErrNoFarmer
	}
	log.Printf("Pasa todos los filtros de identidad")
	identidades, err = s.QueryIdentities(ctx)
	//log.Printf("identidades: %v",identidades)
	return identidades, err
}

//ESta función sirve, una vez registrados el ID de los clients del farmer para asignar al periodo 1 pagos y defaulters a las ID registradas
//segun la N(100,10)
//En este caso, como seria el modelo del spv no se emiten UTXOs
// func (s *SmartContract) InitPaymentsForSPV(ctx contractapi.TransactionContextInterface) error {

// 	identity := ctx.GetClientIdentity()
// 	farmer, err := identity.GetID()
// 	if err != nil {
// 		return err
// 	}
// 	org, err := identity.GetMSPID()
// 	if err != nil {
// 		return err
// 	}
// 	if org != "spvMSP" {
// 		return ErrNoFarmer
// 	}
// 	log.Printf("Pasa todos los filtros de identidad")
// 	periodo := 1

// 	var valor_total float64
// 	numeroAleatorio := rand.Intn(numero_contratos)
// 	// var defaulters []string
// 	// var nuevoPPA []*PPA
// 	lista, err := s.QueryFarmerID(ctx)
// 	//log.Printf("lista: %v", lista)
// 	contador := 1
// 	for _, n := range lista {
// 		de_fault := false
// 		if contador == numeroAleatorio {
// 			log.Printf("defaulter: %v", contador)
// 			de_fault = true
// 		}
// 		farmer = n.Identidad
// 		//log.Printf("Identidad: %v", farmer)
// 		media := 100.0
// 		desv := 10.0
// 		Z := rand.NormFloat64()
// 		X := media + math.Sqrt(desv)*Z
// 		energy := math.Round(X*1000) / 1000
// 		payments := math.Round(0.10*energy*100) / 100
// 		ppa := &PPA{
// 			DocType:  "ppa",
// 			Client:   farmer,
// 			Energy:   energy,
// 			Payments: payments,
// 			Default:  de_fault,
// 			Period:   periodo,
// 		}
// 		//log.Printf("Contrato: %v",ppa)
// 		ppaKey, _ := ctx.GetStub().CreateCompositeKey("ppa", []string{farmer})
// 		//log.Printf("ppakey: %v",ppaKey)
// 		ppaAsBytes, _ := json.Marshal(ppa)
// 		ctx.GetStub().PutState(ppaKey, ppaAsBytes)
// 		contador = contador + 1
// 		if !de_fault {
// 			valor_total = valor_total + payments
// 		}
// 	}
// 	valorr := &ValorTotal{
// 		Doctype: "cantidad",
// 		Total:   valor_total,
// 	}
// 	valorrAsBytes, _ := json.Marshal(valorr)
// 	valorKey, _ := ctx.GetStub().CreateCompositeKey("cantidad", []string{strconv.Itoa(periodo)})
// 	ctx.GetStub().PutState(valorKey, valorrAsBytes)
// 	return nil
// }

//Esta función sirve para consultar todos los pagos que ha realizado una determinada ID. Se usan indexes, por lo que solo busca los que tengan esa ID
func (s *SmartContract) GetHistoryFarmer(ctx contractapi.TransactionContextInterface, farmer string) ([]*PPA, error) {
	identity := ctx.GetClientIdentity()
	farmer, err := identity.GetID()
	if err != nil {
		return nil, err
	}
	org, err := identity.GetMSPID()
	if err != nil {
		return nil, err
	}
	if org != "spvMSP" {
		return nil, ErrNoSpv
	}
	var nuevoPPA []*PPA
	nuevoPPA, err = s.QueryIdentityHistory(ctx, farmer)
	return nuevoPPA, err
}

//Esta función la ejecutará alguien de la org SPV para determinar los pagos durante los 119 periodos restantes. Cada periodo 1 del año se asigna
//un 1% de default, si hay 100 contratos pues 1 y nunca más va a salir del default, sus contratos se siguen registrando pero los pagos
//no se suman al total
// func (s *SmartContract) SimulatedPaymentsForSPV(ctx contractapi.TransactionContextInterface, periodo int) error {
// 	identity := ctx.GetClientIdentity()
// 	farmer, err := identity.GetID()
// 	if err != nil {
// 		return err
// 	}
// 	org, err := identity.GetMSPID()
// 	if err != nil {
// 		return err
// 	}
// 	if org != "spvMSP" {
// 		return ErrNoSpv
// 	}
// 	//log.Printf("Pasa todos los filtros de identidad")

// 	var defaulters []string
// 	var nuevoPPA []*PPA

// 	var resto int
// 	var valor_total float64
// 	var payments float64
// 	//	for nn:=2;nn<49;nn++{
// 	resto = periodo % 12
// 	nuevoPPA, err = s.QueryAssetByPeriod(ctx, periodo-1)
// 	for _, j := range nuevoPPA {
// 		defaulter := j.Client
// 		//log.Printf("defaulter: %v",defaulter)
// 		defaulters = append(defaulters, defaulter)
// 		// log.Printf("defaulters: %v", defaulters)
// 	}
// 	log.Printf("defaulters antes: %v", defaulters)
// 	lista, _ := s.QueryFarmerID(ctx)
// 	if resto == 1 {
// 		num_aleatorio := rand.Intn(numero_contratos)
// 		log.Printf("numero aleatorio: %v", num_aleatorio)
// 		nuevo_default := lista[num_aleatorio]
// 		for _, newrange := range defaulters {
// 			nuevo_default = lista[num_aleatorio]
// 			if nuevo_default.Identidad != newrange {
// 				log.Printf("sigue buscando")
// 			} else {
// 				log.Printf("se ha repetido el default")
// 				num_aleatorio = rand.Intn(numero_contratos)
// 			}
// 		}
// 		defaulters = append(defaulters, nuevo_default.Identidad)
// 		log.Printf("lista de defaulters: %v", defaulters)
// 		for _, n := range lista {
// 			de_fault := false
// 			for _, k := range defaulters {
// 				if n.Identidad == k {
// 					log.Printf("DEFAULTER HAS APPEARED")
// 					de_fault = true
// 				}
// 				farmer = n.Identidad
// 				media := 100.0
// 				desv := 10.0
// 				Z := rand.NormFloat64()
// 				X := media + math.Sqrt(desv)*Z
// 				energy := math.Round(X*1000) / 1000
// 				payments = math.Round(0.10*energy*100) / 100
// 				ppa := &PPA{
// 					DocType:  "ppa",
// 					Client:   farmer,
// 					Energy:   energy,
// 					Payments: payments,
// 					Default:  de_fault,
// 					Period:   periodo,
// 				}
// 				// log.Printf("Contrato: %v",ppa)
// 				ppaAsBytes, _ := json.Marshal(ppa)
// 				ppaKey, _ := ctx.GetStub().CreateCompositeKey("ppa", []string{farmer, strconv.Itoa(periodo)})
// 				//log.Printf("ppakey: %v",ppaKey)
// 				ctx.GetStub().PutState(ppaKey, ppaAsBytes)
// 			}
// 			if !de_fault {
// 				valor_total = valor_total + payments
// 			}
// 		}

// 	} else {
// 		for _, n := range lista {
// 			de_fault := false
// 			for _, k := range defaulters {
// 				if n.Identidad == k {
// 					log.Printf("Hay defaulter")
// 					de_fault = true
// 				}
// 				farmer = n.Identidad
// 				media := 100.0
// 				desv := 10.0
// 				Z := rand.NormFloat64()
// 				X := media + math.Sqrt(desv)*Z
// 				energy := math.Round(X*1000) / 1000
// 				payments = math.Round(0.10*energy*100) / 100
// 				ppa := &PPA{
// 					DocType:  "ppa",
// 					Client:   farmer,
// 					Energy:   energy,
// 					Payments: payments,
// 					Default:  de_fault,
// 					Period:   periodo,
// 				}
// 				// log.Printf("Contrato: %v",ppa)
// 				ppaAsBytes, _ := json.Marshal(ppa)
// 				ppaKey, _ := ctx.GetStub().CreateCompositeKey("ppa", []string{farmer, strconv.Itoa(periodo)})
// 				//log.Printf("ppakey: %v",ppaKey)
// 				ctx.GetStub().PutState(ppaKey, ppaAsBytes)
// 			}
// 			if !de_fault {
// 				valor_total = valor_total + payments
// 			}

// 		}
// 	}
// 	valorr := &ValorTotal{
// 		Doctype: "cantidad",
// 		Total:   valor_total,
// 	}
// 	valorrAsBytes, _ := json.Marshal(valorr)
// 	valorKey, _ := ctx.GetStub().CreateCompositeKey("cantidad", []string{strconv.Itoa(periodo)})
// 	ctx.GetStub().PutState(valorKey, valorrAsBytes)
// 	return nil
// }

//Esta función se usa para calcular el total de los pagos. Cambiar para calcular segun el rate del 3%
//CAMBIAR!!!!!!!!
// func (s *SmartContract) CalculateSPV(ctx contractapi.TransactionContextInterface) error {
// 	rate := 0.03
// 	identity := ctx.GetClientIdentity()
// 	org, err := identity.GetMSPID()
// 	if err != nil {
// 		return err
// 	}
// 	if org != "spvMSP" {
// 		return ErrNoSpv
// 	}
// 	//log.Printf("Pasa todos los filtros de identidad")

// 	var nuevoPPA []*ValorTotal
// 	var valor float64
// 	var total float64
// 	var nuevoValor float64
// 	nuevoPPA, err = s.QueryAssetByPeriodSPV(ctx)
// 	contador := 1
// 	for _, k := range nuevoPPA {
// 		valor = k.Total
// 		nuevoValor = valor / math.Pow(1+rate, float64(contador))
// 		total = total + nuevoValor
// 		contador = contador + 1
// 	}

// 	log.Printf("Valor total: %v", total)
// 	return nil
// }

//esta función sirve para que alguna ID que sea client de la org farmer pueda emitir sus pagos siguiendo una normal N(100,10). El mismo client
//decide default y periodo. Si no hay default emite el token con una cantidad 100*pagos, ya que no acepta decimales.
// func (s *SmartContract) WriteSimulatedPayments(ctx contractapi.TransactionContextInterface, periodo int) error {
// 	//en un entorno productivo habria que comprobar que al periodo le corresponden unas fechas determinadas
// 	hasOU, err := cid.HasOUValue(ctx.GetStub(), "cliente")
// 	if err != nil {
// 		return err
// 	}
// 	if !hasOU {
// 		return ErrNoFarmer
// 	}
// 	identity := ctx.GetClientIdentity()
// 	farmer, err := identity.GetID()
// 	if err != nil {
// 		return err
// 	}
// 	org, err := identity.GetMSPID()
// 	if err != nil {
// 		return err
// 	}
// 	if org != "farmerMSP" {
// 		return ErrNoFarmer
// 	}
// 	if periodo > years*months {
// 		return ErrNoPeriod
// 	}
// 	num, _ := s.QueryAssetNumberByPeriod(ctx, periodo)
// 	if num >= numero_contratos {
// 		return ErrNoPeriod
// 	}

// 	if num == 0 && periodo > 1 {
// 		num, _ := s.QueryAssetNumberByPeriod(ctx, periodo-1)
// 		if num != numero_contratos {
// 			return ErrNoPeriod
// 		}
// 	}

// 	ppa_client, err := s.QueryIdentityHistory(ctx, farmer)
// 	for _, rango_client := range ppa_client {
// 		if rango_client.Period == periodo && rango_client.Client == farmer {
// 			return ErrFarmerPeriod
// 		}
// 	}
// 	de_fault := false
// 	var resto int
// 	resto = periodo % 12
// 	lista, _ := s.QueryAssetByPeriod(ctx, periodo-1)
// 	anho := periodo / 12
// 	for _, ll := range lista {
// 		if ll.Client == farmer {
// 			de_fault = true
// 			break
// 		}
// 	}

// 	if resto == 1 {
// 		lista_identidades, _ := s.QueryIdentities(ctx)
// 		//	log.Printf("entra si es primer mes")
// 		nuevo_num := numeros_defaulters[anho]
// 		nuevo_defaulter := lista_identidades[nuevo_num]
// 		//	log.Printf("pasa al buscar")
// 		if nuevo_defaulter.Identidad == farmer {
// 			//		log.Printf("Hay default")
// 			de_fault = true
// 		}
// 	}
// 	Z := rand.NormFloat64()
// 	X := media + math.Sqrt(desv)*Z
// 	energy := math.Round(X*1000) / 1000
// 	payments := math.Round(0.10*energy*100) / 100
// 	log.Printf("Payments: %v", payments)
// 	ppa := &PPA{
// 		DocType:  "ppa",
// 		Client:   farmer,
// 		Energy:   energy,
// 		Payments: payments,
// 		Default:  de_fault,
// 		Fecha: Datos{
// 			Day:   time.Now().Day(),
// 			Month: time.Now().Month(),
// 			Year:  time.Now().Year(),
// 		},
// 		Period: periodo,
// 	}
// 	if !de_fault {
// 		// _,err=mint(ctx,int(100*payments))
// 		// if err != nil {
// 		// 	return fmt.Errorf("failed to mint utxo: %v", err)
// 		// }
// 		ppaAsBytes, err := json.Marshal(&ppa)
// 		// //si no quiero validar el err, defino como elemento, _ :=json.Marshal()
// 		if err != nil {
// 			fmt.Printf("Marshal error: %s", err.Error())
// 			return err
// 		}
// 		ppaKey := ctx.GetStub().GetTxID()
// 		ctx.GetStub().PutState(ppaKey, ppaAsBytes)
// 	} else {
// 		ppaAsBytes, err := json.Marshal(ppa)
// 		//si no quiero validar el err, defino como elemento, _ :=json.Marshal()
// 		if err != nil {
// 			fmt.Printf("Marshal error: %s", err.Error())
// 			return err
// 		}
// 		ppaKey := ctx.GetStub().GetTxID()

// 		ctx.GetStub().PutState(ppaKey, ppaAsBytes)
// 	}
// 	return nil
// }

//igual que la anterior pero deciciendo el propio client la cantidad
//sseria interesante en estas 2 pasar private data quefuncione como comprobante de que efectivamente se ha realizado el pago

func (s *SmartContract) WritePayments(ctx contractapi.TransactionContextInterface, payments float64, de_fault bool) error {
	fecha1 := fecha_inicio()
	fecha2 := queryperiod()
	// log.Printf("fecha actual: %v", fecha2)
	var y1 int = int(fecha2.Month)
	var y0 int = int(fecha1.Month)
	periodo := (12*fecha2.Year + y1) - (12*fecha1.Year + y0) + 1

	// log.Printf("numero de periodo: %v", periodo)

	hasOU, err := cid.HasOUValue(ctx.GetStub(), "cliente")
	if err != nil {
		return err
		// return shim.Error("Expecting integer value for asset holding")
	}
	if !hasOU {
		// return shim.Error("Expecting integer value for asset holding")
		return ErrNoFarmer
	}
	identity := ctx.GetClientIdentity()
	farmer, err := identity.GetID()
	if err != nil {
		// return shim.Error("Expecting integer value for asset holding")
		return err
	}
	org, err := identity.GetMSPID()
	if err != nil {
		// return shim.Error("Expecting integer value for asset holding")
		return err
	}
	if org != "farmerMSP" {
		// return shim.Error("Expecting integer value for asset holding")
		return ErrNoFarmer
	}
	num, _ := s.QueryAssetNumberByPeriod(ctx, periodo)
	if num >= numero_contratos {
		// return shim.Error("Expecting integer value for asset holding")
		return ErrNoPeriod
	}
	if num == 0 && periodo > 1 {
		num, _ := s.QueryAssetNumberByPeriod(ctx, periodo-1)
		if num != numero_contratos {
			return ErrNoPeriod
			// return shim.Error("Expecting integer value for asset holding")
		}
	}
	ppa_client, err := s.QueryIdentityHistory(ctx, farmer)
	// log.Printf("este cliente ya ha subido: %v", ppa_client)
	for _, rango_client := range ppa_client {
		if rango_client.Period == periodo && rango_client.Client == farmer {
			return ErrFarmerPeriod
			// return shim.Error("Expecting integer value for asset holding")
		}
	}
	energy := 10 * payments

	// de_fault:=false

	ppa := &PPA{
		DocType:  "ppa",
		Client:   farmer,
		Energy:   energy,
		Payments: payments,
		Default:  de_fault,
		Fecha: Datos{
			Day:   time.Now().Day(),
			Month: time.Now().Month(),
			Year:  time.Now().Year(),
		},
		Period: periodo,
	}
	ppaAsBytes, err := json.Marshal(&ppa)
	//si no quiero validar el err, defino como elemento, _ :=json.Marshal()
	if err != nil {
		// fmt.Printf("Marshal error: %s", err.Error())
		// return shim.Error("Expecting integer value for asset holding")
		return err
	}
	ppaKey := ctx.GetStub().GetTxID()
	err = ctx.GetStub().PutState(ppaKey, ppaAsBytes)
	if err != nil {
		// return shim.Error("Expecting integer value for asset holding")
		return fmt.Errorf("Error al subir: %v", err)
	}
	jsonAsBytes, _ := json.Marshal(`"{sending payments to receipient: ` + farmer + ` }"`)
	err = ctx.GetStub().SetEvent("Payments distributed: ", jsonAsBytes)
	return err
}

// func (s *SmartContract) WritePayments(ctx contractapi.TransactionContextInterface, payments float64, periodo int) ([]string,error) {
// 	hasOU, err := cid.HasOUValue(ctx.GetStub(), "client1")
// 	if err != nil {
// 		return nil,err
// 	}
// 	if !hasOU {
// 		return nil,ErrNoFarmer
// 	}
// 	identity := ctx.GetClientIdentity()
// 	farmer, err := identity.GetID()
// 	if err != nil {
// 		return nil, err
// 	}
// 	org, err := identity.GetMSPID()
// 	if err != nil {
// 		return nil,err
// 	}
// 	if org!="farmerMSP"{
// 		return nil,ErrNoFarmer
// 	}
// 	num,_:=s.QueryAssetNumberByPeriod(ctx, periodo)
// 	if num>=numero_contratos{
// 		return nil,ErrNoPeriod
// 	}
// 	if num==0 && periodo>1{
// 		num,_:=s.QueryAssetNumberByPeriod(ctx,periodo-1)
// 		if num!=numero_contratos{
// 			return nil,ErrNoPeriod
// 		}
// 	}
// 	ppa_client,err:=s.QueryIdentityHistory(ctx,farmer)
// 	for _,rango_client:=range ppa_client{
// 		if rango_client.Period==periodo && rango_client.Client==farmer{
// 			return nil,ErrFarmerPeriod
// 		}
// 	}
// 	energy:=10*payments

// 	de_fault:=false

// 	ppa := PPA{
// 		DocType:	"ppa",
// 		Client:     farmer,
// 		Energy:     energy,
// 		Payments:   payments,
// 		Default:    de_fault,
// 		Fecha: Datos{
// 			Day:   time.Now().Day(),
// 			Month: time.Now().Month(),
// 			Year:  time.Now().Year(),
// 		},
// 		Period: periodo,
// 	}
// 	newUTXO:=new(TokenUTXO)
// 	if !de_fault{
// 		newUTXO,err=mint(ctx,int(100*payments))
// 		if err != nil {
// 			return nil,fmt.Errorf("failed to mint utxo: %v", err)
// 		}
// 		clave:=newUTXO.Key

// 		log.Printf("la clave es: %v",clave)
// 	}
// 	//ppaAsBytes
// 	ppaAsBytes , err := json.Marshal(&ppa)
// 	//si no quiero validar el err, defino como elemento, _ :=json.Marshal()
// 	if err != nil {
// 		fmt.Printf("Marshal error: %s", err.Error())
// 		return nil,err
// 	}
// 	ppaKey:=ctx.GetStub().GetTxID()
// 	//  Save index entry to world state. Only the key name is needed, no need to store a duplicate copy of the asset.
// 	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value

// 	//Pero ademas, tambien voy a emitir una cantidad de UTXO equivalente a los payments (multiplicamos por 100 para que sea entero) y tambien se
// 	//enviarian datos privados como numero de cuenta etc para probar que efectivamente se ha realizado el pago, que iran al originador
// 	//(o al aggregator?)

// 	err= ctx.GetStub().PutState(ppaKey, ppaAsBytes)
// 	//clave:=newUTXO.Key
// 	if err!=nil{
// 	 	return nil,fmt.Errorf("Error al subir: %v",err)
// 	}
// 	return []string{newUTXO.Key}, nil
// }

//Esta función sirve para que los UTXOs creados se puedan transferir al originator
// func (s *SmartContract) AfterTransaction(ctx contractapi.TransactionContextInterface) error {
// 	idUTXO, err := s.ClientUTXOs(ctx)
// 	//log.Printf("obtiene algo el clientUTXO? %v",idUTXO)
// 	if err != nil {
// 		return fmt.Errorf("Error: %v", err)
// 	}
// 	value := idUTXO[0]
// 	//log.Printf("valor de la clave: %v",value)
// 	cant, _ := strconv.Atoi(idUTXO[1])
// 	_, err = s.Transfer(ctx, []string{value}, cant)
// 	return err
// }

//funcion para eliminar utxos emitidos por los clients del farmer
func (s *SmartContract) DeletePayment(ctx contractapi.TransactionContextInterface, utxoKey string) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}

	if clientMSPID != "originatorMSP" {
		return fmt.Errorf("client is not authorized to receive new tokens")
	}
	utxoOutputCompositeKey, err := ctx.GetStub().CreateCompositeKey("utxo", []string{clientMSPID, utxoKey})
	err = ctx.GetStub().DelState(utxoOutputCompositeKey)
	if err != nil {
		return err
	}
	log.Printf("utxoInput deleted: %+v", utxoOutputCompositeKey)

	return nil
}

//Funciones de consultas usando INDEXES

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*PPA, error) {
	var assets []*PPA
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset PPA
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func constructQueryRequestFromIteratorID(resultsIterator shim.StateQueryIteratorInterface) ([]*Request, error) {
	var assets []*Request
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Request
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func constructQueryResponseFromIteratorID(resultsIterator shim.StateQueryIteratorInterface) ([]*FarmerID, error) {
	var assets []*FarmerID
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset FarmerID
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func constructQueryResponseFromIteratorRequestID(resultsIterator shim.StateQueryIteratorInterface) ([]*FarmerID, error) {
	var assets []*FarmerID
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Request
		var asset2 FarmerID
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		asset2.Identidad = asset.Requester
		log.Printf("valor: %v", asset)
		assets = append(assets, &asset2)
	}
	return assets, nil
}

func constructQueryResponseFromIteratorSPV(resultsIterator shim.StateQueryIteratorInterface) ([]*ValorTotal, error) {
	var assets []*ValorTotal
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset ValorTotal
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func constructQueryResponseFromIteratorHistory(resultsIterator shim.StateQueryIteratorInterface) ([]*PPA, error) {
	var assets []*PPA
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset PPA
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func constructQueryResponseFromIteratorNumber(resultsIterator shim.StateQueryIteratorInterface) (int, error) {
	var contador int
	contador = 0
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return 0, err
		}
		var asset PPA
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return 0, err
		}
		contador = contador + 1
	}
	return contador, nil
}

func constructQueryResponseFromIteratorPayments(resultsIterator shim.StateQueryIteratorInterface) ([]*PagosImpagos, error) {
	var vector []*PagosImpagos
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset PPA
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		pagos := &PagosImpagos{
			Payments: asset.Payments,
			Default:  asset.Default,
		}
		vector = append(vector, pagos)
	}
	return vector, nil
}

func constructQueryUTXOFromIteratorRequest(resultsIterator shim.StateQueryIteratorInterface) ([]*TokenUTXO, error) {
	var vector []*TokenUTXO
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset TokenUTXO
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		token := &TokenUTXO{
			Owner:  asset.Owner,
			Amount: asset.Amount,
		}
		vector = append(vector, token)
	}
	return vector, nil
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*PPA, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return constructQueryResponseFromIterator(resultsIterator)
}

func getQueryResultForQueryStringID(ctx contractapi.TransactionContextInterface, queryString string) ([]*FarmerID, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return constructQueryResponseFromIteratorID(resultsIterator)
}

func getQueryRequestForQueryStringID(ctx contractapi.TransactionContextInterface, queryString string) ([]*Request, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return constructQueryRequestFromIteratorID(resultsIterator)
}

func getQueryResultForQueryStringHistory(ctx contractapi.TransactionContextInterface, queryString string) ([]*PPA, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return constructQueryResponseFromIteratorHistory(resultsIterator)
}

func getQueryResultForQueryStringSPV(ctx contractapi.TransactionContextInterface, queryString string) ([]*ValorTotal, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return constructQueryResponseFromIteratorSPV(resultsIterator)
}

func getQueryResultForQueryStringNumber(ctx contractapi.TransactionContextInterface, queryString string) (int, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return 0, err
	}
	defer resultsIterator.Close()
	return constructQueryResponseFromIteratorNumber(resultsIterator)
}

func getQueryResultForQueryStringPayments(ctx contractapi.TransactionContextInterface, queryString string) ([]*PagosImpagos, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return constructQueryResponseFromIteratorPayments(resultsIterator)
}

func getQueryResultForQueryRequestID(ctx contractapi.TransactionContextInterface, queryString string) ([]*FarmerID, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return constructQueryResponseFromIteratorRequestID(resultsIterator)
}

func getQueryUTXOForQueryRequest(ctx contractapi.TransactionContextInterface, queryString string) ([]*TokenUTXO, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return constructQueryUTXOFromIteratorRequest(resultsIterator)
}

func (s *SmartContract) QueryAssetNumberByPeriod(ctx contractapi.TransactionContextInterface, periodo int) (int, error) {
	// periodoAsString:=strconv.Itoa(periodo)
	// queryString := fmt.Sprintf(`{"selector":{"docType":"ppa","client":"%s"}}`,client)
	queryString := fmt.Sprintf(`{"selector":{"docType":"ppa","periodo":%d}}`, periodo)

	// queryString := fmt.Sprintf({\"selector\":{\"docType\":\"ppa\",\"periodo\":\"1\",\"default\":\"true\"},\"use_index\":[\"_design/indexPeriodDoc\", \"indexPeriod\"]}", periodo)
	//log.Printf("El string que le pasamos: %v",queryString)
	return getQueryResultForQueryStringNumber(ctx, queryString)
}

func (s *SmartContract) QueryAssetByPeriod(ctx contractapi.TransactionContextInterface, periodo int) ([]*PPA, error) {
	// periodoAsString:=strconv.Itoa(periodo)
	valor := true
	// queryString := fmt.Sprintf(`{"selector":{"docType":"ppa","client":"%s"}}`,client)
	queryString := fmt.Sprintf(`{"selector":{"docType":"ppa","default":%t,"periodo":%d}}`, valor, periodo)

	// queryString := fmt.Sprintf({\"selector\":{\"docType\":\"ppa\",\"periodo\":\"1\",\"default\":\"true\"},\"use_index\":[\"_design/indexPeriodDoc\", \"indexPeriod\"]}", periodo)
	//log.Printf("El string que le pasamos: %v",queryString)
	return getQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryAssetByPeriodSPV(ctx contractapi.TransactionContextInterface) ([]*ValorTotal, error) {
	// periodoAsString:=strconv.Itoa(periodo)
	// queryString := fmt.Sprintf(`{"selector":{"docType":"ppa","client":"%s"}}`,client)
	queryString := fmt.Sprintf(`{"selector":{"docType":"cantidad"}}`)

	return getQueryResultForQueryStringSPV(ctx, queryString)
}

func (s *SmartContract) QueryIdentities(ctx contractapi.TransactionContextInterface) ([]*FarmerID, error) {
	// periodoAsString:=strconv.Itoa(periodo)
	// queryString := fmt.Sprintf(`{"selector":{"docType":"ppa","client":"%s"}}`,client)
	queryString := fmt.Sprintf(`{"selector":{"docType":"identidad"}}`)

	//queryString := fmt.Sprintf({\"selector\":{\"docType\":\"ppa\",\"periodo\":\"1\",\"default\":\"true\"},\"use_index\":[\"_design/indexPeriodDoc\", \"indexPeriod\"]}", periodo)
	//log.Printf("El string que le pasamos: %v",queryString)
	return getQueryResultForQueryStringID(ctx, queryString)
}

func (s *SmartContract) QueryRequests(ctx contractapi.TransactionContextInterface, client string) ([]*Request, error) {
	// periodoAsString:=strconv.Itoa(periodo)
	// queryString := fmt.Sprintf(`{"selector":{"docType":"ppa","client":"%s"}}`,client)
	queryString := fmt.Sprintf(`{"selector":{"docType":"request","bonista":"%s"}}`, client)
	return getQueryRequestForQueryStringID(ctx, queryString)
}

func (s *SmartContract) QueryIdentitiesSPV(ctx contractapi.TransactionContextInterface) ([]*FarmerID, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"SPVidentidad"}}`)
	return getQueryResultForQueryStringID(ctx, queryString)
}

func (s *SmartContract) QueryIdentityHistory(ctx contractapi.TransactionContextInterface, farmer string) ([]*PPA, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"ppa","client":"%s"}}`, farmer)

	//log.Printf("El string que le pasamos: %v",queryString)
	return getQueryResultForQueryStringHistory(ctx, queryString)
}

func (s *SmartContract) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*PPA, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryAssetByID(ctx contractapi.TransactionContextInterface, ppaID string) (*PPA, error) {
	ppaAsBytes, err := ctx.GetStub().GetState(ppaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ppa %s: %v", ppaID, err)
	}
	if ppaAsBytes == nil {
		return nil, fmt.Errorf("ppa %s does not exist", ppaID)
	}
	var ppa PPA
	err = json.Unmarshal(ppaAsBytes, &ppa)
	if err != nil {
		return nil, err
	}
	return &ppa, nil
}

func (s *SmartContract) QueryPaymentsAndDefaultByPeriod(ctx contractapi.TransactionContextInterface, periodo int) ([]*PagosImpagos, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"ppa","periodo":%d}}`, periodo)
	return getQueryResultForQueryStringPayments(ctx, queryString)
}

func (s *SmartContract) QueryRequesterID(ctx contractapi.TransactionContextInterface) ([]*FarmerID, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"request"}}`)
	return getQueryResultForQueryRequestID(ctx, queryString)
}

func (s *SmartContract) QueryAllUTXOs(ctx contractapi.TransactionContextInterface) ([]*TokenUTXO, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"token"}}`)
	return getQueryUTXOForQueryRequest(ctx, queryString)
}

//y ya tenemos la funcion que nos permite guardar en la blockchain

//creo el metodo main
func main() {
	//levantamos un nuevo chaincode y le enviamos la estructura
	//SmartContract, que devuelve 2 valores
	SC := new(SmartContract)
	// S.AfterTransaction
	// NewSC:=new(Contract)
	// NewSC.AfterTransaction=SC.AfterTransaction
	chaincode, err := contractapi.NewChaincode(SC)
	//chaincode, err := contractapi.NewChaincode(new(SmartContract))
	//verificamos si hay algun error
	if err != nil {
		fmt.Printf("Error create ppa chaincode: %s", err.Error())
		//y terminaria la ejecucion del codigo
		return
	}

	//verificamos si hay algun error al ejecutar esta funcion
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting ppa chaincode: %s", err.Error())
	}
}

//Variables en las que se guardan todos los tipos de errores posibles.
var (
	ErrOldID                 = errors.New("This PPA's ID already exists")
	ErrAtraso                = errors.New("This PPA will be considered in default")
	ErrNumMax                = errors.New("Not on correct period or achieved max number of contracts")
	ErrWrongPeriod           = errors.New("You are searching in a wrong period")
	ErrNotAValidFormatClient = errors.New("Client name hasnt a valid format")
	ErrNoFarmer              = errors.New("The identity should be a farmer to execute the transaction")
	ErrNoOriginator          = errors.New("The identity should be an originator to execute the transaction")
	ErrNoSpv                 = errors.New("The identity should be a SPV to execute the transaction")
	ErrNoPeriod              = errors.New("You are not allowed to write in this period")
	ErrFarmerPeriod          = errors.New("This client has already submit a payment for this period")
	ErrNoUnderwritter        = errors.New("The identity should be an underwritter to execute the transaction")
	ErrWrongNumber           = errors.New("There are 10000 bonds")
	ErrNoSPV                 = errors.New("The identity should belong to the SPV to execute the transaction")
	ErrSPVBond               = errors.New("Bonds have already been issued")
)

//con Stub es como accedo al world state y al ledger
