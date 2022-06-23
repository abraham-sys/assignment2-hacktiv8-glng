package controllers

import (
	"assignment2/db/connection"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ItemID      uint   `json:"lineItemId"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderID     uint   `json:"order_id"`
}

type Order struct {
	OrderID      uint   `json:"order_id"`
	OrderedAt    string `json:"orderedAt"`
	CustomerName string `json:"customerName"`
	Quantity     uint   `json:"quantity"`
	Items        []Item `json:"items"`
}

var orderData = []Order{}

func CreateOrder(ctx *gin.Context) {
	var newOrder Order
	var (
		sqlInsertOrderStatement string = `
		INSERT INTO orders (customer_name, ordered_at, quantity)
		VALUES ($1, $2, $3) 
		RETURNING order_id;`

		sqlInsertItemStatement string = `
		INSERT INTO items (item_code, description, quantity, order_id)
		VALUES ($1, $2, $3, $4);`
	)

	connection := connection.ConnectDB()

	err := ctx.ShouldBindJSON(&newOrder)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, eachItem := range newOrder.Items {
		newOrder.Quantity += uint(eachItem.Quantity)
	}

	var dest uint

	resultQueryOrder := connection.QueryRow(
		sqlInsertOrderStatement,
		newOrder.CustomerName,
		newOrder.OrderedAt,
		newOrder.Quantity,
	)

	err = resultQueryOrder.Scan(&dest)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	for _, eachItem := range newOrder.Items {

		_, err := connection.Exec(
			sqlInsertItemStatement,
			eachItem.Description,
			eachItem.ItemCode,
			eachItem.Quantity,
			dest,
		)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "ok",
		"message": "Order successfully created",
	})
}

func GetOrders(ctx *gin.Context) {
	var queryGetOrders string = `
	SELECT * FROM orders o 	
	LEFT JOIN items i 
	ON o.order_id = i.order_id;`
	var orders []Order

	connection := connection.ConnectDB()

	getOrdersResult, err := connection.Query(queryGetOrders)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	var currentOrder Order

	for getOrdersResult.Next() {
		var order Order
		var item Item

		if err := getOrdersResult.Scan(
			&order.OrderID,
			&order.CustomerName,
			&order.OrderedAt,
			&order.Quantity,
			&item.ItemID,
			&item.ItemCode,
			&item.Description,
			&item.Quantity,
			&item.OrderID,
		); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		if currentOrder.OrderID != order.OrderID && currentOrder.OrderID != 0 {
			orders = append(orders, currentOrder)
			currentOrder = order
		} else if currentOrder.OrderID == 0 {
			currentOrder = order
		}

		if currentOrder.OrderID == item.OrderID {
			currentOrder.Items = append(currentOrder.Items, item)
		}
	}
	orders = append(orders, currentOrder)

	ctx.JSON(http.StatusAccepted, gin.H{
		"status": "ok",
		"data":   orders,
	})

}

func UpdateOrder(ctx *gin.Context) {
	orderId := ctx.Param("orderId")
	var orderToBeUpdated Order

	var (
		queryUpdateOrder string = `
		UPDATE orders 
		
		SET 
			customer_name=($1),
			ordered_at=($2),
			quantity=($3)

		WHERE
			order_id=($4)

		RETURNING
			order_id;
		`
		queryUpdateItem string = `
		UPDATE items 
		
		SET 
			item_code=($1),
			description=($2),
			quantity=($3),
			order_id=($4)

		WHERE
			item_id=($5) AND order_id=($4)
		`
	)

	connection := connection.ConnectDB()

	err := ctx.ShouldBindJSON(&orderToBeUpdated)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, eachItem := range orderToBeUpdated.Items {
		orderToBeUpdated.Quantity += uint(eachItem.Quantity)
	}

	var dest uint

	resultQueryOrder := connection.QueryRow(
		queryUpdateOrder,
		orderToBeUpdated.CustomerName,
		orderToBeUpdated.OrderedAt,
		orderToBeUpdated.Quantity,
		orderId,
	)

	err = resultQueryOrder.Scan(&dest)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	for _, eachItem := range orderToBeUpdated.Items {

		_, err := connection.Exec(
			queryUpdateItem,
			eachItem.ItemCode,
			eachItem.Description,
			eachItem.Quantity,
			dest,
			eachItem.ItemID,
		)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "ok",
		"message": "Order successfully updated",
	})
}

func DeleteById(ctx *gin.Context) {
	orderId := ctx.Param("orderId")

	connection := connection.ConnectDB()

	for _, table := range []string{"items", "orders"} {
		_, err := connection.Exec(fmt.Sprintf("DELETE FROM %s WHERE order_id = %s;", table, orderId))
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"status":  "ok",
		"message": "Order successfully deleted",
	})

}
