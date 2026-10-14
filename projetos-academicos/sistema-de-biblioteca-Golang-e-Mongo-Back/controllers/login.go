package controller // Especifica o nome do pacote

//// Função para autenticar usuário
//func LoginUsuario(c *gin.Context, login models.LoginRequest) {
//
//	usersCollection := dataBase.DB.Collection("Users")
//
//	filter := bson.D{{"username", login.User}}
//	var result models.User
//
//	// Verificar se o usuário existe pelo nome de usuário
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	err := usersCollection.FindOne(ctx, filter).Decode(&result)
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			log.Printf("Usuário não encontrado: %v", err)
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
//			return
//		}
//		log.Printf("Erro ao encontrar usuário: %v", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while finding user"})
//		return
//	}
//
//	// Compare the password with the hashed password stored in the database
//	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(login.Password))
//	if err != nil {
//		log.Printf("Erro ao comparar senha: %v", err)
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
//		return
//	}
//
//	// Definir o cookie do usuário
//	c.SetCookie("username", login.User, 3600, "/", "localhost", false, true)
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "Login successful",
//		"user":    result,
//	})
//}
//
//// Função para autenticar administrador
//func LoginAdmin(c *gin.Context, login models.LoginRequest) {
//
//	usersCollection := dataBase.DB.Collection("Admins")
//
//	filter := bson.D{{"username", login.User}}
//	var result models.Admin
//
//	// Verificar se o usuário existe pelo nome de usuário
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	err := usersCollection.FindOne(ctx, filter).Decode(&result)
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			log.Printf("Usuário não encontrado: %v", err)
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
//			return
//		}
//		log.Printf("Erro ao encontrar usuário: %v", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while finding user"})
//		return
//	}
//
//	// Compare the password with the hashed password stored in the database
//	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(login.Password))
//	if err != nil {
//		log.Printf("Erro ao comparar senha: %v", err)
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
//		return
//	}
//
//	// Definir o cookie do usuário
//	c.SetCookie("username", login.User, 3600, "/", "localhost", false, true)
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "Login successful",
//		"user":    result,
//	})
//}
