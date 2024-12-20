package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ExportData struct {
	Filter   FilterData               `json:"filter"`
	Products []map[string]interface{} `json:"products"`
	Orders   []map[string]interface{} `json:"orders"`
	Payments []map[string]interface{} `json:"payments"`
}

type FilterData struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// ExportJSONHandler returns an Echo handler function that exports JSON based on start and end datetime
func ExportJSONHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		// Parse and validate query parameters
		startTime, endTime, err := parseQueryParams(e)
		if err != nil {
			return e.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		// Initialize the export data
		exportData := ExportData{
			Filter: FilterData{
				Start: startTime,
				End:   endTime,
			},
		}

		// Fetch and enrich products
		products, err := fetchAndEnrichProducts(app)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		exportData.Products = products

		// Fetch orders with order_items and build maps
		orders, ordersMap, orderItemsMap, err := fetchOrdersWithItems(app, startTime, endTime)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		exportData.Orders = orders

		// Fetch payments and build paymentsMap
		payments, paymentsMap, err := fetchAndEnrichPayments(app, startTime, endTime)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		exportData.Payments = payments

		// Fetch events and assign them to the appropriate objects
		err = processEventsAndAssign(app, startTime, endTime, ordersMap, orderItemsMap, paymentsMap)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		// Serialize the export data to JSON with indentation for readability
		return sendJSONResponse(e, exportData)
	}
}

// parseQueryParams parses and validates 'start' and 'end' query parameters
func parseQueryParams(e *core.RequestEvent) (time.Time, time.Time, error) {
	startStr := e.Request.URL.Query().Get("start")
	endStr := e.Request.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("Missing 'start' or 'end' query parameters")
	}

	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("Invalid 'start' datetime format. Use RFC3339 format.")
	}

	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("Invalid 'end' datetime format. Use RFC3339 format.")
	}

	if startTime.After(endTime) {
		return time.Time{}, time.Time{}, fmt.Errorf("'start' datetime must be before 'end' datetime")
	}

	return startTime, endTime, nil
}

// fetchAndEnrichProducts fetches all products and enriches them with related data
func fetchAndEnrichProducts(app core.App) ([]map[string]interface{}, error) {
	productRecords, err := app.FindAllRecords("product")
	if err != nil {
		return nil, err
	}

	products := make([]map[string]interface{}, 0, len(productRecords))
	for _, record := range productRecords {
		productMap, err := getCleanRecordMap(record)
		if err != nil {
			return nil, err
		}

		enrichedProductMap, err := enrichProductData(app, productMap)
		if err != nil {
			return nil, err
		}

		products = append(products, enrichedProductMap)
	}
	return products, nil
}

func stringSliceToInterfaceSlice(strings []string) []interface{} {
	interfaces := make([]interface{}, len(strings))
	for i, v := range strings {
		interfaces[i] = v
	}
	return interfaces
}

// fetchOrdersWithItems fetches orders and their associated order_items
func fetchOrdersWithItems(app core.App, startTime, endTime time.Time) ([]map[string]interface{}, map[string]map[string]interface{}, map[string]map[string]interface{}, error) {
	filter := "created >= {:start} && created <= {:end}"
	params := dbx.Params{
		"start": startTime,
		"end":   endTime,
	}
	orderRecords, err := app.FindRecordsByFilter("order", filter, "", 0, 0, params)
	if err != nil {
		return nil, nil, nil, err
	}

	orders := make([]map[string]interface{}, 0, len(orderRecords))
	ordersMap := make(map[string]map[string]interface{})

	// Collect order IDs
	orderIDs := make([]string, len(orderRecords))
	for i, record := range orderRecords {
		orderIDs[i] = record.Id
	}

	// Convert orderIDs to []interface{}
	orderIDsInterface := stringSliceToInterfaceSlice(orderIDs)

	// Fetch order_items associated with the orders
	orderItemRecords, err := app.FindAllRecords("order_item", dbx.In("order", orderIDsInterface...))
	if err != nil {
		return nil, nil, nil, err
	}

	// Create maps for easy lookup
	orderItemsByOrderID := make(map[string][]map[string]interface{})
	orderItemsMap := make(map[string]map[string]interface{})

	for _, record := range orderItemRecords {
		orderItemMap, err := getCleanRecordMap(record)
		if err != nil {
			return nil, nil, nil, err
		}

		// Enrich "menu_item"
		if menuItemID, ok := orderItemMap["menu_item"].(string); ok && menuItemID != "" {
			menuItemRecord, err := app.FindRecordById("menu_item", menuItemID)
			if err != nil {
				return nil, nil, nil, err
			}
			menuItemMap, err := getCleanRecordMap(menuItemRecord)
			if err != nil {
				return nil, nil, nil, err
			}
			orderItemMap["menu_item"] = menuItemMap
		}

		orderItemID := record.Id
		orderItemsMap[orderItemID] = orderItemMap

		orderID := record.GetString("order")
		orderItemsByOrderID[orderID] = append(orderItemsByOrderID[orderID], orderItemMap)
	}

	// Attach order_items to orders
	for _, record := range orderRecords {
		orderMap, err := getCleanRecordMap(record)
		if err != nil {
			return nil, nil, nil, err
		}

		orderID := record.Id
		if items, ok := orderItemsByOrderID[orderID]; ok {
			orderMap["order_items"] = items
		} else {
			orderMap["order_items"] = []map[string]interface{}{}
		}

		orders = append(orders, orderMap)
		ordersMap[orderID] = orderMap
	}

	return orders, ordersMap, orderItemsMap, nil
}

// fetchAndEnrichPayments fetches payments and enriches them with related data
func fetchAndEnrichPayments(app core.App, startTime, endTime time.Time) ([]map[string]interface{}, map[string]map[string]interface{}, error) {
	filter := "created >= {:start} && created <= {:end}"
	params := dbx.Params{
		"start": startTime,
		"end":   endTime,
	}
	paymentRecords, err := app.FindRecordsByFilter("payment", filter, "", 0, 0, params)
	if err != nil {
		return nil, nil, err
	}

	payments := make([]map[string]interface{}, 0, len(paymentRecords))
	paymentsMap := make(map[string]map[string]interface{})

	for _, record := range paymentRecords {
		paymentMap, err := getCleanRecordMap(record)
		if err != nil {
			return nil, nil, err
		}

		// Enrich 'payment_option' in payment
		enrichedPaymentMap, err := enrichPaymentData(app, paymentMap)
		if err != nil {
			return nil, nil, err
		}

		payments = append(payments, enrichedPaymentMap)
		paymentsMap[record.Id] = enrichedPaymentMap
	}
	return payments, paymentsMap, nil
}

// processEventsAndAssign processes events and assigns them to orders, order_items, or payments
func processEventsAndAssign(app core.App, startTime, endTime time.Time, ordersMap, orderItemsMap, paymentsMap map[string]map[string]interface{}) error {
	filter := "created >= {:start} && created <= {:end}"
	params := dbx.Params{
		"start": startTime,
		"end":   endTime,
	}
	eventRecords, err := app.FindRecordsByFilter("event", filter, "", 0, 0, params)
	if err != nil {
		return err
	}

	for _, record := range eventRecords {
		eventMap, err := getCleanRecordMap(record)
		if err != nil {
			return err
		}

		// Get 'type' and 'content' from event
		eventType, _ := eventMap["type"].(string)
		contentRaw, _ := eventMap["content"]
		content, _ := contentRaw.(map[string]interface{})

		if content == nil {
			continue // Skip events without content
		}

		switch eventType {
		case "order_item":
			orderItemID, _ := content["order_item_id"].(string)
			if orderItem, ok := orderItemsMap[orderItemID]; ok {
				// Append event to orderItem["events"]
				appendEvent(orderItem, eventMap)
			}
		case "order":
			orderID, _ := content["order_id"].(string)
			if order, ok := ordersMap[orderID]; ok {
				// Append event to order["events"]
				appendEvent(order, eventMap)
			}
		case "payment":
			paymentID, _ := content["payment_id"].(string)
			if payment, ok := paymentsMap[paymentID]; ok {
				// Append event to payment["events"]
				appendEvent(payment, eventMap)
			}
		default:
			// Unknown type, ignore or handle as needed
		}
	}

	return nil
}

// appendEvent appends an event to the 'events' slice of the object
func appendEvent(obj map[string]interface{}, eventMap map[string]interface{}) {
	events, ok := obj["events"].([]interface{})
	if !ok {
		events = []interface{}{}
	}
	events = append(events, eventMap)
	obj["events"] = events
}

// sendJSONResponse sends the JSON response as a downloadable file
func sendJSONResponse(e *core.RequestEvent, data ExportData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return e.InternalServerError("Failed to serialize export data to JSON", err)
	}

	e.Response.Header().Set("Content-Disposition", `attachment; filename="export.json"`)

	return e.Blob(http.StatusOK, "application/json", jsonData)
}

// getCleanRecordMap converts a Record to a clean map without collection metadata
func getCleanRecordMap(record *core.Record) (map[string]interface{}, error) {
	var recordMap map[string]interface{}

	// Marshal the record to JSON
	jsonBytes, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON back into a map
	err = json.Unmarshal(jsonBytes, &recordMap)
	if err != nil {
		return nil, err
	}

	return cleanRecordMap(recordMap), nil
}

// cleanRecordMap removes collection metadata from the record map
func cleanRecordMap(recordMap map[string]interface{}) map[string]interface{} {
	delete(recordMap, "collectionId")
	delete(recordMap, "collectionName")
	return recordMap
}

// enrichProductData adds detailed information for attributes, station, and category
func enrichProductData(app core.App, productMap map[string]interface{}) (map[string]interface{}, error) {
	// Enrich Attribute
	if attributes, ok := productMap["attribute"].([]interface{}); ok && len(attributes) > 0 {
		enrichedAttributes := make([]map[string]interface{}, 0, len(attributes))
		for _, attrID := range attributes {
			attrRecord, err := app.FindRecordById("product_attribute", attrID.(string))
			if err == nil {
				attrMap, err := getCleanRecordMap(attrRecord)
				if err != nil {
					return nil, err
				}
				enrichedAttributes = append(enrichedAttributes, attrMap)
			}
		}
		productMap["attribute"] = enrichedAttributes
	}

	// Enrich Category
	if categoryID, ok := productMap["category"].(string); ok {
		categoryRecord, err := app.FindRecordById("product_categ", categoryID)
		if err == nil {
			categoryMap, err := getCleanRecordMap(categoryRecord)
			if err != nil {
				return nil, err
			}
			productMap["category"] = categoryMap
		}
	}

	// Enrich Station
	if stationID, ok := productMap["station"].(string); ok {
		stationRecord, err := app.FindRecordById("station", stationID)
		if err == nil {
			stationMap, err := getCleanRecordMap(stationRecord)
			if err != nil {
				return nil, err
			}
			productMap["station"] = stationMap
		}
	}

	return productMap, nil
}

// enrichPaymentData enriches 'payment_option' in payment with full details
func enrichPaymentData(app core.App, paymentMap map[string]interface{}) (map[string]interface{}, error) {
	// Enrich 'payment_option'
	if paymentOptionID, ok := paymentMap["payment_option"].(string); ok {
		paymentOptionRecord, err := app.FindRecordById("payment_option", paymentOptionID)
		if err == nil {
			paymentOptionMap, err := getCleanRecordMap(paymentOptionRecord)
			if err != nil {
				return nil, err
			}
			paymentMap["payment_option"] = paymentOptionMap
		}
	}

	return paymentMap, nil
}
