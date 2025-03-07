package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Struktur data untuk tugas
type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

// Slice untuk menyimpan daftar tugas
var tasks = []Task{
	{ID: 1, Name: "Belajar Golang", Complete: false},
	{ID: 2, Name: "Membuat API", Complete: false},
}

// 🔹 Ambil semua tugas
func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

// 🔹 Tambah tugas baru
func addTask(c *gin.Context) {
	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	c.JSON(http.StatusCreated, newTask)
}

// 🔹 Tandai tugas sebagai selesai
func completeTask(c *gin.Context) {
	id := c.Param("id")
	for i, task := range tasks {
		if id == strconv.Itoa(task.ID) {
			tasks[i].Complete = true
			c.JSON(http.StatusOK, tasks[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Tugas tidak ditemukan"})
}

// 🔹 Hapus tugas
func deleteTask(c *gin.Context) {
	id := c.Param("id")
	for i, task := range tasks {
		if id == strconv.Itoa(task.ID) {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Tugas dihapus"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Tugas tidak ditemukan"})
}

func main() {
	r := gin.Default()

	// Routing API
	r.GET("/tasks", getTasks)          // 🔹 Ambil semua tugas
	r.POST("/tasks", addTask)          // 🔹 Tambah tugas baru
	r.PUT("/tasks/:id", completeTask)  // 🔹 Tandai selesai
	r.DELETE("/tasks/:id", deleteTask) // 🔹 Hapus tugas

	// Menjalankan server di port 8080
	r.Run(":8080")
}
