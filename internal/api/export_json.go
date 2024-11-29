package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type ExportData struct {
	Filter   FilterData               `json:"filter"`
	Products []map[string]interface{} `json:"products"`
	Events   []map[string]interface{} `json:"events"`
	Orders   []map[string]interface{} `json:"orders"`
	Payments []map[string]interface{} `json:"payments"`
}

type FilterData struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// ExportJSONHandler returns an Echo handler function that exports JSON based on start and end datetime
func ExportJSONHandler(app core.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse and validate query parameters
		startTime, endTime, err := parseQueryParams(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
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
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		exportData.Products = products

		// Fetch events within the specified timeframe
		events, err := fetchEvents(app, startTime, endTime)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		exportData.Events = events

		// Fetch orders within the specified timeframe
		orders, err := fetchOrders(app, startTime, endTime)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		exportData.Orders = orders

		// Fetch payments within the specified timeframe
		payments, err := fetchPayments(app, startTime, endTime)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		exportData.Payments = payments

		// Serialize the export data to JSON with indentation for readability
		return sendJSONResponse(c, exportData)
	}
}

// parseQueryParams parses and validates 'start' and 'end' query parameters
func parseQueryParams(c echo.Context) (time.Time, time.Time, error) {
	startStr := c.QueryParam("start")
	endStr := c.QueryParam("end")

	if startStr == "" || endStr == "" {
		return time.Time{}, time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "Missing 'start' or 'end' query parameters")
	}

	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid 'start' datetime format. Use RFC3339 format.")
	}

	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return time.Time{}, time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid 'end' datetime format. Use RFC3339 format.")
	}

	if startTime.After(endTime) {
		return time.Time{}, time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "'start' datetime must be before 'end' datetime")
	}

	return startTime, endTime, nil
}

// fetchAndEnrichProducts fetches all products and enriches them with related data
func fetchAndEnrichProducts(app core.App) ([]map[string]interface{}, error) {
	productRecords, err := app.Dao().FindRecordsByExpr("product", nil)
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

// fetchEvents fetches events within the specified timeframe
func fetchEvents(app core.App, startTime, endTime time.Time) ([]map[string]interface{}, error) {
	expr := dbx.NewExp("created BETWEEN {:start} AND {:end}", dbx.Params{
		"start": startTime,
		"end":   endTime,
	})
	eventRecords, err := app.Dao().FindRecordsByExpr("event", expr)
	if err != nil {
		return nil, err
	}

	events := make([]map[string]interface{}, 0, len(eventRecords))
	for _, record := range eventRecords {
		eventMap, err := getCleanRecordMap(record)
		if err != nil {
			return nil, err
		}
		events = append(events, eventMap)
	}
	return events, nil
}

// fetchOrders fetches orders within the specified timeframe
func fetchOrders(app core.App, startTime, endTime time.Time) ([]map[string]interface{}, error) {
	expr := dbx.NewExp("created BETWEEN {:start} AND {:end}", dbx.Params{
		"start": startTime,
		"end":   endTime,
	})
	orderRecords, err := app.Dao().FindRecordsByExpr("order", expr)
	if err != nil {
		return nil, err
	}

	orders := make([]map[string]interface{}, 0, len(orderRecords))
	for _, record := range orderRecords {
		orderMap, err := getCleanRecordMap(record)
		if err != nil {
			return nil, err
		}
		orders = append(orders, orderMap)
	}
	return orders, nil
}

// fetchOrders fetches orders within the specified timeframe
func fetchPayments(app core.App, startTime, endTime time.Time) ([]map[string]interface{}, error) {
	expr := dbx.NewExp("created BETWEEN {:start} AND {:end}", dbx.Params{
		"start": startTime,
		"end":   endTime,
	})
	paymentRecords, err := app.Dao().FindRecordsByExpr("payment", expr)
	if err != nil {
		return nil, err
	}

	payments := make([]map[string]interface{}, 0, len(paymentRecords))
	for _, record := range paymentRecords {
		paymentMap, err := getCleanRecordMap(record)
		if err != nil {
			return nil, err
		}
		payments = append(payments, paymentMap)
	}
	return payments, nil
}

// sendJSONResponse sends the JSON response as a downloadable file
func sendJSONResponse(c echo.Context, data ExportData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to serialize export data to JSON")
	}

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Content-Disposition", `attachment; filename="export.json"`)

	return c.Blob(http.StatusOK, "application/json", jsonData)
}

// getCleanRecordMap converts a Record to a clean map without collection metadata
func getCleanRecordMap(record *models.Record) (map[string]interface{}, error) {
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
	dao := app.Dao()

	// Enrich Attribute
	if attributes, ok := productMap["attribute"].([]interface{}); ok && len(attributes) > 0 {
		enrichedAttributes := make([]map[string]interface{}, 0, len(attributes))
		for _, attrID := range attributes {
			attrRecord, err := dao.FindRecordById("product_attribute", attrID.(string))
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
		categoryRecord, err := dao.FindRecordById("product_categ", categoryID)
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
		stationRecord, err := dao.FindRecordById("station", stationID)
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
