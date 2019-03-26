package controllers

import (
	"log"
	"net/http"
	"sort"
	"sto_order_manager/config"
	"sto_order_manager/models"

	"github.com/gin-gonic/gin"
)

type OrderBook struct {
	Bids       models.OrderSlice `json:"bids"`
	Asks       models.OrderSlice `json:"asks"`
	BidsAmount int               `json:"bidsAmount"`
	AsksAmount int               `json:"asksAmount"`
}

func AddOrder(c *gin.Context) {
	order := new(models.Order)
	if err := c.BindJSON(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order.MakerAssetFilledAmount = "0"
	order.TakerAssetFilledAmount = "0"
	_, err := models.CreateOrder(order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": "Add Order Success!"})
}

func FillOrder(c *gin.Context) {
	/*
		_signature, ok := c.Get("signature")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request non signature"})
			return
		}
		signature := _signature.(string)
		_makerAssetFilledAmount, ok := c.Get("makerAssetFilledAmount")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Requesr non makerAssetFilledAmount"})
			return
		}
		makerAssetFilledAmount := _makerAssetFilledAmount.(string)
		_takerAssetFilledAmount, ok := c.Get("takerAssetFilledAmount")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Requesr non takerAssetFilledAmount"})
			return
		}
		takerAssetFilledAmount := _takerAssetFilledAmount.(string)
	*/
	order := new(models.Order)
	if err := c.BindJSON(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(order.MakerAssetFilledAmount)
	log.Println(order.TakerAssetFilledAmount)
	err := models.UpdateOrder(order.MakerAssetFilledAmount, order.TakerAssetFilledAmount,
		order.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": "Fill Order Success!"})
}

func GetOrderBook(c *gin.Context) {
	tokenSymbol := c.Query("tokenSymbol")
	if tokenSymbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request non tokenSymbol"})
		return
	}
	// tokenSymbol := _tokenSymbol.(string)
	_tokenAssetData := config.V.Token[tokenSymbol]
	tokenAssetData := _tokenAssetData.(string)
	_wethAssetData := config.V.Token["WETH"]
	wethAssetData := _wethAssetData.(string)
	asks, err := models.GetOrders(tokenAssetData, wethAssetData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bids, err := models.GetOrders(wethAssetData, tokenAssetData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sort.Sort(bids)
	sort.Sort(asks)
	response := OrderBook{
		bids,
		asks,
		bids.Len(),
		asks.Len(),
	}
	c.JSON(http.StatusOK, response)
}
