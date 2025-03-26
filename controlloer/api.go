package controlloer

import (
	"go-final/dto"
	"go-final/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var database *gorm.DB

func DemoController(router *gin.Engine, gormdb *gorm.DB) {
	database = gormdb
	routers := router.Group("/api")
	{
		routers.GET("/ping", ping)
		routers.POST("/login", login)
		routers.PUT("/users/:id", updateUserInfo)
		routers.PUT("/users/:id/change-password", changePassword)
		routers.GET("/products/search", searchProducts)
		routers.POST("/cart/add/:id", addToCart)
	}
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Pong!!",
	})
}

func login(c *gin.Context) {
	// login
	loginform := model.Login{}
	if err := c.ShouldBindJSON(&loginform); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := model.Customer{}
	result := database.Where("email = ?", loginform.Email).First(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "ไม่พบผู้ใช้งาน"})
		return
	}
	if !checkPassword(user.Password, loginform.Password) {
		c.JSON(400, gin.H{"error": "รหัสผ่านไม่ถูกต้อง"})
		return
	}
	response := dto.ConvertUserResponse(user)
	c.JSON(200, response)
}

func updateUserInfo(c *gin.Context) {
	id := c.Param("id")
	var updateRequest dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := model.Customer{}
	result := database.First(&user, id)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "ไม่พบผู้ใช้งาน"})
		return
	}

	// Update user fields
	user.FirstName = updateRequest.FirstName
	user.LastName = updateRequest.LastName
	user.PhoneNumber = updateRequest.PhoneNumber
	user.Address = updateRequest.Address

	if err := database.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "อัปเดตข้อมูลไม่สำเร็จ"})
		return
	}

	response := dto.ConvertUserResponse(user)
	c.JSON(200, response)
}

func changePassword(c *gin.Context) {
	id := c.Param("id")
	changePassRequest := dto.ChangePasswordRequest{}
	if err := c.ShouldBindJSON(&changePassRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := model.Customer{}
	result := database.First(&user, id)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "ไม่พบผู้ใช้งาน"})
		return
	}

	if !checkPassword(user.Password, changePassRequest.OldPassword) {
		c.JSON(400, gin.H{"error": "รหัสผ่านเก่าไม่ถูกต้อง"})
		return
	}

	hashedPassword, err := hashPassword(changePassRequest.NewPassword)
	if err != nil {
		c.JSON(500, gin.H{"error": "เข้ารหัสผ่านไม่สำเร็จ"})
		return
	}

	user.Password = hashedPassword
	if err := database.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "อัปเดตรหัสผ่านไม่สำเร็จ"})
		return
	}

	c.JSON(200, gin.H{"message": "เปลี่ยนรหัสผ่านสำเร็จ"})
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func searchProducts(c *gin.Context) {
	searchRequest := dto.ProductSearchRequest{}
	if err := c.ShouldBindJSON(&searchRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	products := []model.Product{}
	query := database.Where("description LIKE ?", "%"+searchRequest.Keyword+"%")

	if searchRequest.MinPrice > 0 {
		query = query.Where("CAST(price AS DECIMAL) >= ?", searchRequest.MinPrice)
	}
	if searchRequest.MaxPrice > 0 {
		query = query.Where("CAST(price AS DECIMAL) <= ?", searchRequest.MaxPrice)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "ค้นหาสินค้าไม่สำเร็จ"})
		return
	}

	c.JSON(200, products)
}

func addToCart(c *gin.Context) {
	userID := c.Param("id")
	addRequest := dto.AddToCartRequest{}
	if err := c.ShouldBindJSON(&addRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	cart := model.Cart{}
	result := database.Where("customer_id = ? AND cart_name = ?", userID, addRequest.CartName).First(&cart)
	if result.Error != nil {
		// Create new cart
		cart = model.Cart{
			CustomerID: parseInt(userID),
			CartName:   addRequest.CartName,
		}
		if err := database.Create(&cart).Error; err != nil {
			c.JSON(500, gin.H{"error": "สร้างรถเข็นไม่สำเร็จ"})
			return
		}
	}

	existingItem := model.CartItem{}
	result = database.Where("cart_id = ? AND product_id = ?", cart.CartID, addRequest.ProductID).First(&existingItem)
	if result.Error == nil {
		existingItem.Quantity += addRequest.Quantity
		if err := database.Save(&existingItem).Error; err != nil {
			c.JSON(500, gin.H{"error": "อัพเดทจำนวนสินค้าไม่สำเร็จ"})
			return
		}
	} else {
		cartItem := model.CartItem{
			CartID:    cart.CartID,
			ProductID: addRequest.ProductID,
			Quantity:  addRequest.Quantity,
		}
		if err := database.Create(&cartItem).Error; err != nil {
			c.JSON(500, gin.H{"error": "เพิ่มสินค้าลงรถเข็นไม่สำเร็จ"})
			return
		}
	}

	// Fetch cart items with product details
	cartItems := []dto.CartItemResponse{}
	err := database.Table("cart_item").
		Select("cart_item.cart_item_id, product.product_name, product.description, product.price, cart_item.quantity").
		Joins("JOIN product ON product.product_id = cart_item.product_id").
		Where("cart_item.cart_id = ?", cart.CartID).
		Scan(&cartItems).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "ไม่สามารถดึงข้อมูลสินค้าในรถเข็นได้"})
		return
	}

	for i := range cartItems {
		price, _ := strconv.ParseFloat(cartItems[i].Price, 64)
		cartItems[i].TotalAmount = price * float64(cartItems[i].Quantity)
	}

	response := dto.CartResponse{
		CartID:   cart.CartID,
		CartName: cart.CartName,
		Items:    cartItems,
	}

	c.JSON(200, gin.H{
		"message": "เพิ่มสินค้าลงรถเข็นสำเร็จ",
		"cart":    response,
	})
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
