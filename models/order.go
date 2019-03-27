package models

import (
	"database/sql"
	"log"
	"math/big"

	"sto_order_manager/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB_mysql *sql.DB

type Order struct {
	ExchangeAddress        string `json:"exchangeAddress"`
	MakerAddress           string `json:"makerAddress"`
	TakerAddress           string `json:"takerAddress"`
	FeeRecipientAddress    string `json:"feeRecipientAddress"`
	SenderAddress          string `json:"senderAddress"`
	MakerAssetAmount       string `json:"makerAssetAmount"`
	TakerAssetAmount       string `json:"takerAssetAmount"`
	MakerAssetFilledAmount string `json:"makerAssetFilledAmount"`
	TakerAssetFilledAmount string `json:"takerAssetFilledAmount"`
	MakerFee               string `json:"makerFee"`
	TakerFee               string `json:"takerFee"`
	ExpirationTimeSeconds  string `json:"expirationTimeSeconds"`
	Salt                   string `json:"salt"`
	MakerAssetData         string `json:"makerAssetData"`
	TakerAssetData         string `json:"takerAssetData"`
	Signature              string `json:"signature"`
}

type OrderSlice []Order

func (o OrderSlice) Len() int      { return len(o) }
func (o OrderSlice) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (o OrderSlice) Less(i, j int) bool {
	x := new(big.Float)
	y := new(big.Float)
	z := new(big.Float)
	x.SetString(o[i].MakerAssetAmount)
	y.SetString(o[i].TakerAssetAmount)
	z.Quo(y, x)

	u := new(big.Float)
	v := new(big.Float)
	w := new(big.Float)
	u.SetString(o[j].MakerAssetAmount)
	v.SetString(o[j].TakerAssetAmount)
	w.Quo(v, u)
	return z.Cmp(w) <= 0
}

func CreateOrder(order *Order) (int, error) {
	result, err := DB_mysql.Exec(`insert into orders(exchangeAddress, 
	makerAddress, takerAddress, feeRecipientAddress, senderAddress, 
	makerAssetAmount, takerAssetAmount, makerAssetFilledAmount, 
	takerAssetFilledAmount, makerFee, takerFee, expirationTimeSeconds, 
	salt, makerAssetData, takerAssetData, signature) values (?, ?, ?, ?, ?, ?, 
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		order.ExchangeAddress,
		order.MakerAddress,
		order.TakerAddress,
		order.FeeRecipientAddress,
		order.SenderAddress,
		order.MakerAssetAmount,
		order.TakerAssetAmount,
		order.MakerAssetFilledAmount,
		order.TakerAssetFilledAmount,
		order.MakerFee,
		order.TakerFee,
		order.ExpirationTimeSeconds,
		order.Salt,
		order.MakerAssetData,
		order.TakerAssetData,
		order.Signature)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func UpdateOrder(makerAssetFilledAmount, takerAssetFilledAmount, signature string) error {
	var MakerAssetAmount string
	row := DB_mysql.QueryRow(`select MakerAssetAmount from orders where signature = ?`,
		signature)
	err := row.Scan(&MakerAssetAmount)
	if err != nil {
		return err
	}
	if MakerAssetAmount == makerAssetFilledAmount {
		result, err := DB_mysql.Exec("delete from orders where signature = ?",
			signature)
		if err != nil {
			return err
		}
		_, err = result.LastInsertId()
		if err != nil {
			return err
		}
		return nil
	}

	result, err := DB_mysql.Exec(`update orders set makerAssetFilledAmount = ?, 
	takerAssetFilledAmount = ? where signature = ?`,
		makerAssetFilledAmount,
		takerAssetFilledAmount,
		signature)

	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func GetOrders(makerAssetData, takerAssetData string) (orders OrderSlice, err error) {
	rows, err := DB_mysql.Query(`select exchangeAddress, 
	makerAddress, takerAddress, feeRecipientAddress, senderAddress, 
	makerAssetAmount, takerAssetAmount, makerAssetFilledAmount, 
	takerAssetFilledAmount, makerFee, takerFee, expirationTimeSeconds, 
	salt, makerAssetData, takerAssetData, signature from orders 
	where makerAssetData=? and takerAssetData=?`,
		makerAssetData, takerAssetData)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tmp Order
	for rows.Next() {
		err = rows.Scan(&tmp.ExchangeAddress, &tmp.MakerAddress, &tmp.TakerAddress,
			&tmp.FeeRecipientAddress, &tmp.SenderAddress, &tmp.MakerAssetAmount,
			&tmp.TakerAssetAmount, &tmp.MakerAssetFilledAmount, &tmp.TakerAssetFilledAmount,
			&tmp.MakerFee, &tmp.TakerFee, &tmp.ExpirationTimeSeconds, &tmp.Salt,
			&tmp.MakerAssetData, &tmp.TakerAssetData, &tmp.Signature)
		order := tmp
		orders = append(orders, order)
		if err != nil {
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return
}

func init() {
	var err error
	username := config.V.Mysql.Username
	password := config.V.Mysql.Password
	host := config.V.Mysql.Host
	port := config.V.Mysql.Port
	dbname := config.V.Mysql.Dbname
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname
	DB_mysql, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
}
