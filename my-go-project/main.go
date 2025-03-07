package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

var tasks []Task
var filename = "tasks.json"

func loadTasks() {
	file, err := os.ReadFile(filename)
	if err == nil {
		json.Unmarshal(file, &tasks)
	}
}

func saveTasks() {
	data, _ := json.MarshalIndent(tasks, "", "  ")
	os.WriteFile(filename, data, 0644)
}

func addTask(name string) {
	id := len(tasks) + 1
	task := Task{ID: id, Name: name, Complete: false}
	tasks = append(tasks, task)
	saveTasks()
	fmt.Println("✅ Tugas berhasil ditambahkan!")
}

func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("📌 Tidak ada tugas saat ini.")
		return
	}
	fmt.Println("\n📋 Daftar Tugas:")
	for _, task := range tasks {
		status := "❌"
		if task.Complete {
			status = "✅"
		}
		fmt.Printf("[%d] %s %s\n", task.ID, task.Name, status)
	}
}

func completeTask(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Complete = true
			saveTasks()
			fmt.Println("🎉 Tugas selesai!")
			return
		}
	}
	fmt.Println("⚠️ Tugas tidak ditemukan.")
}

func deleteTask(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTasks()
			fmt.Println("🗑️ Tugas dihapus!")
			return
		}
	}
	fmt.Println("⚠️ Tugas tidak ditemukan.")
}

func main() {
	loadTasks()
	for {
		fmt.Println("\n📌 To-Do List CLI")
		fmt.Println("1. Tambah Tugas")
		fmt.Println("2. Lihat Tugas")
		fmt.Println("3. Tandai Selesai")
		fmt.Println("4. Hapus Tugas")
		fmt.Println("5. Keluar")
		fmt.Print("Pilih menu: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Print("Masukkan nama tugas: ")
			var name string
			fmt.Scanln(&name)
			addTask(name)
		case 2:
			listTasks()
		case 3:
			fmt.Print("Masukkan ID tugas yang selesai: ")
			var id int
			fmt.Scan(&id)
			completeTask(id)
		case 4:
			fmt.Print("Masukkan ID tugas yang akan dihapus: ")
			var id int
			fmt.Scan(&id)
			deleteTask(id)
		case 5:
			fmt.Println("👋 Keluar dari aplikasi. Sampai jumpa!")
			return
		default:
			fmt.Println("⚠️ Pilihan tidak valid.")
		}
	}
}
