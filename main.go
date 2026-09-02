package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

// Zamanı 2 saat önce, 3 dakika önce falan yapmak için yardımcı fonksiyon
func formatRelativeTime(t time.Time) string {
	duration := time.Since(t)
	if duration.Minutes() < 1 {
		return "az önce"
	}
	if duration.Hours() < 1 {
		return fmt.Sprintf("%d dakika önce", int(duration.Minutes()))
	}
	if duration.Hours() < 24 {
		return fmt.Sprintf("%d saat önce", int(duration.Hours()))
	}
	return fmt.Sprintf("%d gün önce", int(duration.Hours()/24))
}

func main() {
	// Argüman girilmediyse (sadece "todo" yazıldıysa) yardım bas (E1 kuralı)
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	command := os.Args[1]

	// Program çalışırken önce veriyi disktken yükle
	// Load fonksiyonu storage.go içinde, ama aynı package oldukları için import etmeden görebiliyor!
	store, err := Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
		os.Exit(1)
	}

	// Yönlendirme (Switch)
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Hata: Görev metni gerekli. Örnek: todo add \"Görev adı\"")
			os.Exit(1)
		}

		// TrimSpace: Sadece boşluk girilmesini engelliyor (E5)
		title := strings.TrimSpace(os.Args[2])
		if title == "" {
			fmt.Fprintln(os.Stderr, "Hata: Görev metni boş olamaz.")
			os.Exit(1)
		}

		// task.go'daki fonksiyona adres (&store) gönderiyoruz ki asıl liste güncellensin
		newTask := AddTask(&store, title)

		// try-except yok, Go'da hataları her seferinde böyle yakalamak gerekiyor
		if err := Save(store); err != nil {
			fmt.Fprintf(os.Stderr, "Hata: Görev kaydedilemedi: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Görev eklendi (ID: %d)\n", newTask.ID)

	case "list":
		if len(store.Tasks) == 0 {
			fmt.Println("Henüz görev yok. Eklemek için:\ntodo add \"görev metni\"")
			os.Exit(0)
		}

		// --all flag kontrolü
		showAll := false
		if len(os.Args) >= 3 && os.Args[2] == "--all" {
			showAll = true
		}

		// Tablo düzgün dursun diye tabwriter
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tDURUM\tGÖREV\tOLUŞTURULMA\t")

		pendingCount, doneCount := 0, 0
		for _, t := range store.Tasks {
			if t.Done {
				doneCount++
			} else {
				pendingCount++
			}

			// Eğer all denilmediyse ve görev bittiyse atla
			if !showAll && t.Done {
				continue
			}

			status := "[ ]"
			if t.Done {
				status = "[x]"
			}

			// formatRelativeTime ile süreyi Türkçeleştirip ekrana bas
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n", t.ID, status, t.Title, formatRelativeTime(t.CreatedAt))
		}
		w.Flush() // Tabloyu döktük
		fmt.Printf("\nToplam: %d görev  Tamamlanan: %d  Bekleyen: %d\n", len(store.Tasks), doneCount, pendingCount)

	case "done":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Hata: Görev ID'si gerekli. Örnek: todo done 1")
			os.Exit(1)
		}

		// Terminalden gelen string'i Integer'a çevir (harf girilirse diye koruma var)
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hata: Geçersiz ID: %q bir sayı değil\n", os.Args[2])
			os.Exit(1)
		}

		task, err := CompleteTask(&store, id)
		if err != nil {
			if err.Error() == "zaten tamamlanmış" { // task.go'dan gelen özel hata
				fmt.Printf("Bilgi: Görev zaten tamamlanmış: %q\n", task.Title)
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "Hata: %d numaralı görev bulunamadı\n", id)
			os.Exit(1)
		}

		if err := Save(store); err != nil {
			fmt.Fprintf(os.Stderr, "Hata: Kaydedilemedi: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Görev tamamlandı: %q\n", task.Title)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Hata: Görev ID'si gerekli. Örnek: todo delete 1")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hata: Geçersiz ID: %q bir sayı değil\n", os.Args[2])
			os.Exit(1)
		}

		task, err := DeleteTask(&store, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hata: %d numaralı görev bulunamadı\n", id)
			os.Exit(1)
		}

		if err := Save(store); err != nil {
			fmt.Fprintf(os.Stderr, "Hata: Silinme kaydedilemedi: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Görev silindi: %q\n", task.Title)

	case "help":
		printHelp()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Hata: Bilinmeyen komut: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

// Kullanım menüsü
func printHelp() {
	fmt.Println(`todo - basit bir görev yöneticisi

Kullanım:
  todo add <metin>     Yeni görev ekler
  todo list [--all]    Görevleri listeler
  todo done <id>       Görevi tamamlandı işaretler
  todo delete <id>     Görevi siler
  todo help            Bu mesajı gösterir`)
}
