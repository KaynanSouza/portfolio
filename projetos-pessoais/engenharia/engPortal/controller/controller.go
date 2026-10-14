package controller

import (
	"context"
	"encoding/json"
	"engPortal/models"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetProjects(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading configuration"})
		return
	}

	token := os.Getenv("GITLAB_PRIVATE_TOKEN")
	groupID := os.Getenv("GITLAB_GROUP_ID")

	if token == "" || groupID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing GITLAB_PRIVATE_TOKEN or GITLAB_GROUP_ID in .env file"})
		return
	}

	url := fmt.Sprintf("https://gitlab.com/api/v4/groups/%s/projects", groupID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Add("PRIVATE-TOKEN", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects from GitLab"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from GitLab"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, "application/json; charset=utf-8", body)
		log.Printf("GitLab API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
		return
	}

	// ... dentro da função GetProjects ...

	var gitlabProjects []map[string]interface{}
	if err := json.Unmarshal(body, &gitlabProjects); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse GitLab response"})
		return
	}

	var formattedProjects []models.ProjectResponse

	for _, p := range gitlabProjects {
		// Extraímos os dados básicos
		project := models.ProjectResponse{
			Name:        p["name"].(string),
			Description: "",
			Branch:      p["default_branch"].(string),
			WebURL:      p["web_url"].(string),
		}

		if p["description"] != nil {
			project.Description = p["description"].(string)
		}

		// Tratamento de Data (simplificado para YYYY-MM-DD como na imagem)
		if activity, ok := p["last_activity_at"].(string); ok && len(activity) >= 10 {
			project.LastCommitDate = activity[:10]
		}

		// Mock do status de CI (O GitLab só envia isso se você pedir o campo 'last_pipeline')
		// Para fins de visualização igual à imagem:
		project.StatusCI = "Success"

		formattedProjects = append(formattedProjects, project)
	}

	c.JSON(http.StatusOK, formattedProjects)

	/*c.Data(http.StatusOK, "application/json; charset=utf-8", body)*/
}

func GetKanbanData(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading configuration"})
		return
	}

	tenantID := os.Getenv("TENANT_ID")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	planID := os.Getenv("PLAN_IDS")

	if tenantID == "" || clientID == "" || clientSecret == "" || planID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing Microsoft Graph credentials in .env"})
		return
	}

	accessToken, err := getMicrosoftGraphToken(tenantID, clientID, clientSecret)
	if err != nil {
		log.Printf("Failed to get MS Graph token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Fetch Plan Title
	planTitle, err := fetchPlanTitle(planID, accessToken)
	if err != nil {
		log.Printf("Failed to fetch plan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plan details"})
		return
	}

	// Fetch Buckets
	buckets, err := fetchBuckets(planID, accessToken)
	if err != nil {
		log.Printf("Failed to fetch buckets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch buckets"})
		return
	}

	// Fetch Tasks
	tasks, err := fetchTasks(planID, accessToken)
	if err != nil {
		log.Printf("Failed to fetch tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	// Organize tasks into buckets
	userCache := make(map[string]string)
	bucketMap := make(map[string]*models.KanbanBucket)
	for i := range buckets {
		buckets[i].Tasks = []models.KanbanTask{}
		bucketMap[buckets[i].ID] = &buckets[i]
	}

	for _, task := range tasks {
		// Resolve User ID to Name
		if name, ok := userCache[task.CreatedBy]; ok {
			task.CreatedBy = name
		} else {
			displayName := fetchUserName(task.CreatedBy, accessToken)
			userCache[task.CreatedBy] = displayName
			task.CreatedBy = displayName
		}

		if bucket, ok := bucketMap[task.BucketID]; ok {
			bucket.Tasks = append(bucket.Tasks, task)
		}
	}

	kanbanBoard := models.KanbanBoard{
		PlanID:  planID,
		Title:   planTitle,
		Buckets: buckets,
	}

	c.JSON(http.StatusOK, kanbanBoard)
}

func fetchUserName(userID, token string) string {
	if userID == "" {
		return "Unknown"
	}

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", userID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return userID
	}
	defer resp.Body.Close()

	var result struct {
		DisplayName string `json:"displayName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result.DisplayName == "" {
		return userID
	}

	return result.DisplayName
}

func getMicrosoftGraphToken(tenantID, clientID, clientSecret string) (string, error) {
	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	data := fmt.Sprintf("client_id=%s&scope=https://graph.microsoft.com/.default&client_secret=%s&grant_type=client_credentials", clientID, clientSecret)

	req, _ := http.NewRequest("POST", url, strings.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}

	return "", fmt.Errorf("could not find access_token in response")
}

func fetchPlanTitle(planID, token string) (string, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/planner/plans/%s", planID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if title, ok := result["title"].(string); ok {
		return title, nil
	}
	return "Untitled Plan", nil
}

func fetchBuckets(planID, token string) ([]models.KanbanBucket, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/planner/plans/%s/buckets", planID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Value []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"value"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	var buckets []models.KanbanBucket
	for _, b := range result.Value {
		buckets = append(buckets, models.KanbanBucket{
			ID:   b.ID,
			Name: b.Name,
		})
	}
	return buckets, nil
}

func fetchTasks(planID, token string) ([]models.KanbanTask, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/planner/plans/%s/tasks", planID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Value []struct {
			Title           string `json:"title"`
			BucketID        string `json:"bucketId"`
			PercentComplete int    `json:"percentComplete"`
			CreatedBy       struct {
				User struct {
					ID string `json:"id"`
				} `json:"user"`
			} `json:"createdBy"`
			DueDateTime string `json:"dueDateTime"`
		} `json:"value"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	var tasks []models.KanbanTask
	for _, t := range result.Value {
		status := "Em progresso"
		if t.PercentComplete == 0 {
			status = "A fazer"
		} else if t.PercentComplete == 100 {
			status = "Concluído"
		}

		tasks = append(tasks, models.KanbanTask{
			Title:       t.Title,
			CreatedBy:   t.CreatedBy.User.ID,
			DueDateTime: t.DueDateTime,
			Status:      status,
			BucketID:    t.BucketID,
		})
	}
	return tasks, nil
}

func GetDocs(c *gin.Context) {
	dir := "./documentos"

	// Cria a pasta automaticamente se ela não existir para evitar erros
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível ler a pasta de documentos"})
		return
	}

	var docs []map[string]interface{}

	for i, f := range files {
		if f.IsDir() {
			continue // Ignora subpastas
		}

		info, err := f.Info()
		if err != nil {
			continue
		}

		// Pega a extensão sem o ponto (ex: "pdf")
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(f.Name()), "."))

		// Formata o tamanho do arquivo
		size := info.Size()
		var sizeStr string
		if size > 1024*1024 {
			sizeStr = fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
		} else {
			sizeStr = fmt.Sprintf("%.2f KB", float64(size)/1024)
		}

		docs = append(docs, map[string]interface{}{
			"id":       fmt.Sprintf("doc-%d", i),
			"title":    f.Name(),
			"filename": f.Name(),
			"type":     ext,
			"size":     sizeStr,
		})
	}

	// Se a pasta estiver vazia, retorna array vazio para o React não quebrar
	if docs == nil {
		docs = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, docs)
}

func DownloadFile(c *gin.Context) {
	// MUDE AQUI: Pega o nome do arquivo da query string (ex: /download?file=nome.pdf)
	filename := c.Query("file")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome do arquivo não fornecido"})
		return
	}

	cleanName := filepath.Base(filename)
	filePath := filepath.Join("./documentos", cleanName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado no disco"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+cleanName+"\"")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")

	c.File(filePath)
}

type Training struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title    string             `bson:"title" json:"title"`
	Category string             `bson:"category" json:"category"`
	Status   string             `bson:"status" json:"status"`
	Duration string             `bson:"duration" json:"duration"`
	Enrolled int                `bson:"enrolled" json:"enrolled"`
}

func GetTrainings(c *gin.Context) {
	// Carrega o .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Aviso: .env não encontrado ou erro ao ler.")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MONGO_URI não configurada"})
		return
	}

	// Conecta ao MongoDB
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar no MongoDB"})
		return
	}
	defer client.Disconnect(context.TODO())

	// Seleciona o banco de dados ("portal") e a coleção ("trainings")
	collection := client.Database("portal").Collection("trainings")

	// Busca todos os treinamentos do banco
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar treinamentos"})
		return
	}
	defer cursor.Close(context.TODO())

	var trainings []Training
	if err = cursor.All(context.TODO(), &trainings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados do banco"})
		return
	}

	// Se o banco estiver vazio, retorna um array vazio em vez de null
	if trainings == nil {
		trainings = []Training{}
	}

	c.JSON(http.StatusOK, trainings)
}

func CreateTraining(c *gin.Context) {
	// Pega os dados enviados pelo React
	var newTraining Training
	if err := c.ShouldBindJSON(&newTraining); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("Aviso: .env não encontrado ou erro ao ler.")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		fmt.Println("Erro: MONGO_URI está vazia")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MONGO_URI não configurada no servidor"})
		return
	}

	// Conecta no banco
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Erro ao conectar no Mongo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar no banco"})
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("portal").Collection("trainings")

	// Insere no banco
	result, err := collection.InsertOne(context.TODO(), newTraining)
	if err != nil {
		fmt.Println("Erro ao inserir:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar treinamento no banco"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Treinamento salvo com sucesso",
		"id":      result.InsertedID,
	})
}

func UpdateTraining(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updateData Training
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	godotenv.Load()
	uri := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar no banco"})
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("portal").Collection("trainings")

	// Prepara os dados para atualizar apenas o que foi enviado
	update := bson.M{
		"$set": bson.M{
			"title":    updateData.Title,
			"category": updateData.Category,
			"duration": updateData.Duration,
			"status":   updateData.Status,
			// Não atualizamos o 'enrolled' aqui para não zerar os inscritos sem querer ao editar
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar treinamento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Treinamento atualizado com sucesso"})
}

func DeleteTraining(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	godotenv.Load()
	uri := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar no banco"})
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("portal").Collection("trainings")

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar treinamento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Treinamento deletado com sucesso"})
}
