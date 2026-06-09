package api

import (
    "net/http"
    "strconv"

    "pra-exchange/database"

    "github.com/gin-gonic/gin"
)

type Server struct {
    router *gin.Engine
    port   int
}

func NewServer(port int) *Server {
    router := gin.Default()
    
    s := &Server{
        router: router,
        port:   port,
    }
    
    s.setupRoutes()
    return s
}

func (s *Server) setupRoutes() {
    // Health check
    s.router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
    
    // دریافت لیست معاملات
    s.router.GET("/api/trades", s.getTrades)
    
    // دریافت یک معامله
    s.router.GET("/api/trades/:id", s.getTradeByID)
    
    // دریافت تراکنش‌های یک کاربر
    s.router.GET("/api/transfers/:address", s.getTransfersByAddress)
    
    // آمار کلی
    s.router.GET("/api/stats", s.getStats)
}

func (s *Server) getTrades(c *gin.Context) {
    var trades []database.Trade
    database.DB.Order("created_at desc").Find(&trades)
    c.JSON(http.StatusOK, trades)
}

func (s *Server) getTradeByID(c *gin.Context) {
    id := c.Param("id")
    var trade database.Trade
    if err := database.DB.Where("trade_id = ?", id).First(&trade).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Trade not found"})
        return
    }
    c.JSON(http.StatusOK, trade)
}

func (s *Server) getTransfersByAddress(c *gin.Context) {
    address := c.Param("address")
    var transfers []database.TransferEvent
    database.DB.Where("from = ? OR to = ?", address, address).
        Order("created_at desc").
        Find(&transfers)
    c.JSON(http.StatusOK, transfers)
}

func (s *Server) getStats(c *gin.Context) {
    var totalTrades int64
    var completedTrades int64
    var totalVolume float64
    
    database.DB.Model(&database.Trade{}).Count(&totalTrades)
    database.DB.Model(&database.Trade{}).Where("status = ?", "completed").Count(&completedTrades)
    database.DB.Model(&database.Trade{}).Select("COALESCE(SUM(pra_amount), 0)").Scan(&totalVolume)
    
    c.JSON(http.StatusOK, gin.H{
        "total_trades":      totalTrades,
        "completed_trades":  completedTrades,
        "total_volume_pra":  totalVolume,
    })
}

func (s *Server) Start() {
    s.router.Run(":" + strconv.Itoa(s.port))
}