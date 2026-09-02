package main

import (
	"fmt"
	"os"
	"strconv" // String'i Integer'a (metni sayıya) çevirmek için gerekli paket
	"strings"
	"text/tabwriter"
	"time"
)

// Görevlerin özelliklerini tutan yapı
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Done        bool       `json:"done"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

func main() {
	// Program sadece "todo" diye çalıştırılırsa (yanında komut yoksa)
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	command := os.Args[1] // İlk kelimemiz komut (add, list, done, delete)

	// Hangi komut olursa olsun önce JSON dosyasını okuyup verileri almalıyız
	store, err := Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "add":
		// add için görev metni lazım, yani en az 3 kelime olmalı ("todo" "add" "metin")
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Hata: Görev metni gerekli. Örnek: todo add \"Görev adı\"")
			os.Exit(1)
		}

		title := strings.TrimSpace(os.Args[2])
		if title == "" { // Kullanıcı sadece boşluk girdiyse iptal et (E4, E5 kuralları)
			fmt.Fprintln(os.Stderr, "Hata: Görev metni boş olamaz.")
			os.Exit(1)
		}

		newTask := Task{
			ID:          store.NextID,
			Title:       title,
			Done:        false,
			CreatedAt:   time.Now(),
			CompletedAt: nil,
		}

		// Yeni görevi listeye ekle ve bir sonraki ID'yi hazırla
		store.Tasks = append(store.Tasks, newTask)
		store.NextID++

		// JSON dosyasına kaydet, hata verirse programdan çık
		err = Save(store)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hata: Görev kaydedilemedi: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Görev eklendi (ID: %d)\n", newTask.ID)

	case "list":
		// Eğer listede hiç görev yoksa tabloyu hiç çizme
		if len(store.Tasks) == 0 {
			fmt.Println("Henüz görev yok. Eklemek için:\ntodo add \"görev metni\"")
			os.Exit(0)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tDURUM\tGÖREV\tOLUŞTURULMA\t")

		for _, t := range store.Tasks {
			status := "[ ]"
			if t.Done {
				status = "[x]"
			}

			timeStr := t.CreatedAt.Format("2006-01-02 15:04") // Go'nun tuhaf tarih formatı
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n", t.ID, status, t.Title, timeStr)
		}
		w.Flush()

	case "done":
		// done komutunun yanına ID numarası girilmiş mi kontrol et
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Hata: Görev ID'si gerekli. Örnek: todo done 1")
			os.Exit(1)
		}

		// Terminalden gelen her şey metindir (string). Kullanıcı "1" yazsa bile bu sayı değildir.
		// strconv.Atoi ile bunu metinden sayıya (int) çeviriyoruz.
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			// Eğer kullanıcı sayı yerine harf girerse (mesela todo done abc) çökmesini engelledik.
			fmt.Fprintf(os.Stderr, "Hata: Geçersiz ID: %q bir sayı değil\n", os.Args[2])
			os.Exit(1)
		}

		// Görevi listede bulmak için döngüyle tek tek bakıyoruz
		found := false // Görevi bulup bulmadığımızı takip etmek için bayrak (flag)
		for i, t := range store.Tasks {
			if t.ID == id { // Görevi bulduk!
				found = true
				if t.Done {
					// Zaten bitmiş bir görevse hata verme, sadece uyar
					fmt.Printf("Bilgi: Görev zaten tamamlanmış: %q\n", t.Title)
					os.Exit(0)
				}

				// Görevi tamamlandı yapıyoruz
				store.Tasks[i].Done = true
				now := time.Now()
				store.Tasks[i].CompletedAt = &now // & işareti ile bellekteki adresini gösteriyoruz

				// Değişikliği dosyaya kaydet
				if err := Save(store); err != nil {
					fmt.Fprintf(os.Stderr, "Hata: Kaydedilemedi: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Görev tamamlandı: %q\n", t.Title)
				break // Döngüden çık, diğer görevlere bakmaya gerek kalmadı
			}
		}

		// Döngü bittiğinde found hala false ise demek ki o ID listede yok
		if !found {
			fmt.Fprintf(os.Stderr, "Hata: %d numaralı görev bulunamadı\n", id)
			os.Exit(1)
		}

	case "delete":
		// delete komutunun yanına ID girilmiş mi?
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Hata: Görev ID'si gerekli. Örnek: todo delete 1")
			os.Exit(1)
		}

		// Metni yine sayıya çevir
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hata: Geçersiz ID: %q bir sayı değil\n", os.Args[2])
			os.Exit(1)
		}

		// Silinecek elemanın sırasını (index) bulmamız lazım
		foundIndex := -1
		var deletedTitle string
		for i, t := range store.Tasks {
			if t.ID == id {
				foundIndex = i // Elemanın listedeki kaçıncı sırada olduğunu not aldık
				deletedTitle = t.Title
				break
			}
		}

		// Eğer index hala -1 ise o ID'de görev yoktur
		if foundIndex == -1 {
			fmt.Fprintf(os.Stderr, "Hata: %d numaralı görev bulunamadı\n", id)
			os.Exit(1)
		}

		// BURASI ÖNEMLİ: Pythonda bir listeden eleman silmek için .remove() vb. kullanıyordum
		// ama Go dilinde slice (dilim) mantığı biraz farklıymış.
		// Taktik şu: Silinecek elemanın sol tarafını al, sağ tarafını al, ikisini birleştir.
		// Sondaki ... işareti sağdaki dilimin içindeki tüm elemanları tek tek ekle demekmiş.
		store.Tasks = append(store.Tasks[:foundIndex], store.Tasks[foundIndex+1:]...)

		// Yeni listeyi kaydet
		if err := Save(store); err != nil {
			fmt.Fprintf(os.Stderr, "Hata: Silinme kaydedilemedi: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Görev silindi: %q\n", deletedTitle)

	case "help":
		printHelp()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Hata: Bilinmeyen komut: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`todo - basit bir görev yöneticisi

Kullanım:
  todo add <metin>     Yeni görev ekler
  todo list [--all]    Görevleri listeler
  todo done <id>       Görevi tamamlandı işaretler
  todo delete <id>     Görevi siler
  todo help            Bu mesajı gösterir`)
}
