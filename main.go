package main

import (
	"fmt"
	"os"
)

func main() {
	// 1. os.Args'ı ekrana basıp inceleme aşaması
	// Terminalde ne geldiğini görmek için alttaki satırı açıp test edebilirsin:
	// fmt.Println("Gelen argümanlar:", os.Args)

	// Hiç argüman girilmeden sadece "todo" yazılırsa (E1 kuralı) help göster ve çık
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	// os.Args[0] çalıştırılan programın adıdır.
	// Kullanıcının girdiği ilk komut os.Args[1] içindedir.
	command := os.Args[1]

	// 2. Argümana göre dallanan switch yapısı
	switch command {
	case "add":
		fmt.Println("add komutu çağrıldı")
	case "list":
		fmt.Println("list komutu çağrıldı")
	case "done":
		fmt.Println("done komutu çağrıldı")
	case "delete":
		fmt.Println("delete komutu çağrıldı")
	case "help":
		printHelp()
		os.Exit(0)
	default:
		// 3. Bilinmeyen bir komut girilirse (örneğin: todo foo) stderr'e yaz ve exit code 1 ile çık
		fmt.Fprintf(os.Stderr, "Hata: Bilinmeyen komut: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

// Yardım metnini ekrana basan fonksiyon
func printHelp() {
	fmt.Println(`todo basit bir görev yöneticisi

Kullanım:
  todo add <metin>     Yeni görev ekler
  todo list [--all]    Görevleri listeler
  todo done <id>       Görevi tamamlandı işaretler
  todo delete <id>     Görevi siler
  todo help            Bu mesajı gösterir`)
}