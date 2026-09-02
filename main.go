package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// Görevlerin özelliklerini tutan yapı (struct).
// CompletedAt için pointer (*time.Time) kullandık. Çünkü görev henüz bitmediyse
// varsayılan bir tarih (0001-01-01) atamak yerine 'nil' yani boş bırakmak istiyoruz.
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Done        bool       `json:"done"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

func main() {
	// os.Args[0] her zaman programın kendi adı ("gotodo") olur.
	// Kullanıcı sadece "todo" yazıp enter'a bastıysa uzunluk 1'dir, yardım menüsünü göster.
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	// Kullanıcının girdiği ilk ana komutu alıyoruz (add, list, vb.)
	command := os.Args[1]

	// Hangi komut olursa olsun, program başlar başlamaz önce kayıtlı verileri okuyalım.
	store, err := Load()
	if err != nil {
		// Hoca dokümanda hataları hep os.Stderr'e basmamızı istemişti (Kural K3).
		fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
		os.Exit(1)
	}

	// Gelen komuta göre yönlendirme yapıyoruz
	switch command {
	case "add":
		// add komutu "todo add metin" şeklinde yazılır, yani en az 3 kelime (argüman) lazım.
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Hata: Görev metni gerekli. Örnek: todo add \"Görev adı\"")
			os.Exit(1)
		}

		// Kullanıcı yanlışlıkla sadece boşluk (space) girdiyse, TrimSpace bunu temizler (E5 kuralı).
		title := strings.TrimSpace(os.Args[2])

		// Metin tamamen boşsa kaydetmeyi engelle (E4 kuralı).
		if title == "" {
			fmt.Fprintln(os.Stderr, "Hata: Görev metni boş olamaz.")
			os.Exit(1)
		}

		// Yeni görevimizi hazırlıyoruz
		newTask := Task{
			ID:          store.NextID,
			Title:       title,
			Done:        false,
			CreatedAt:   time.Now(),  // O anki saati al
			CompletedAt: nil,         // Henüz tamamlanmadı
		}

		// Python'daki liste "append" mantığıyla birebir aynı.
		// Go'da slice'a (dinamik diziye) yeni eleman ekliyoruz.
		store.Tasks = append(store.Tasks, newTask)

		// Bir sonraki eklenecek görev için ID'yi 1 artırıp hazırlayalım.
		store.NextID++

		// Yaptığımız değişikliği dosyaya (tasks.json) kaydediyoruz.
		err = Save(store)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hata: Görev kaydedilemedi: %v\n", err)
			os.Exit(1) // Kayıt başarısızsa program hata kodu (1) ile çıksın.
		}

		// İşlem başarılı!
		fmt.Printf("Görev eklendi (ID: %d)\n", newTask.ID)

	case "list":
		// Hiç görev eklenmemişse, boş tablo çizmek yerine kullanıcıya mesaj verelim.
		if len(store.Tasks) == 0 {
			fmt.Println("Henüz görev yok. Eklemek için:\ntodo add \"görev metni\"")
			os.Exit(0)
		}

		// Tablonun sütunları düzgün ve hizalı çıksın diye tabwriter kullanıyoruz.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tDURUM\tGÖREV\tOLUŞTURULMA\t")

		// Slice (dizi) içindeki tüm görevleri tek tek dönüyoruz.
		// '_' olan index numarası (kullanmayacağımız için alt tire yaptık).
		// 't' olan ise o anki görevimiz.
		for _, t := range store.Tasks {
			status := "[ ]" // Varsayılan durum bitmedi
			if t.Done {
				status = "[x]" // Bittiyse içini doldur
			}

			// Tarihi okunabilir formata çeviriyoruz. (Go'da referans tarih hep 2006-01-02'dir)
			timeStr := t.CreatedAt.Format("2006-01-02 15:04")

			// Her bir görevi alt alta tabloya yazıyoruz.
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n", t.ID, status, t.Title, timeStr)
		}
		// Tabloyu ekrana dök.
		w.Flush()

	case "done":
		fmt.Println("done komutu Aşama 4'te yapılacak")
	case "delete":
		fmt.Println("delete komutu Aşama 4'te yapılacak")
	case "help":
		printHelp()
		os.Exit(0)
	default:
		// "todo falanfilan" gibi saçma bir komut girilirse.
		fmt.Fprintf(os.Stderr, "Hata: Bilinmeyen komut: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

// Yardım menüsünü ekrana basan fonksiyon.
func printHelp() {
	fmt.Println(`todo - basit bir görev yöneticisi

Kullanım:
  todo add <metin>     Yeni görev ekler
  todo list [--all]    Görevleri listeler
  todo done <id>       Görevi tamamlandı işaretler
  todo delete <id>     Görevi siler
  todo help            Bu mesajı gösterir`)
}