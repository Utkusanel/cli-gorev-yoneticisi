package main

import (
	"encoding/json"
	"errors"
	"os"
)

// Verileri kaydedeceğimiz dosyanın adını sabit (const) bir değişkene atadık.
const dataFile = "tasks.json"

// JSON dosyasının en dışındaki ana yapı.
// Neden düz dizi yapmadık? İleride başka ayarlar (mesela NextID) eklersek
// buraya koymak daha kolay diye struct yaptım.
type Store struct {
	Tasks  []Task `json:"tasks"`
	NextID int    `json:"nextID"`
}

// Load: Dosyadan verileri okuyup bize dolu bir Store döndürür.
func Load() (Store, error) {
	var store Store

	// 1. Dosyayı okumayı dene
	data, err := os.ReadFile(dataFile)
	if err != nil {
		// os.ErrNotExist hatası "böyle bir dosya henüz yok" demekmiş.
		// Eğer dosya yoksa program çökmesin. İlk defa çalışıyordur.
		if errors.Is(err, os.ErrNotExist) {
			store.NextID = 1  // İlk görev ekleneceği zaman 1 numara olsun
			return store, nil // Hata yokmuş gibi boş store dönüyoruz
		}
		// Dosya var ama yetki yoksa falan gerçek hatadır, bunu döndür.
		return store, err
	}

	// 2. Dosya var ama içi tamamen boşsa (0 byte durumu - E11 kuralı).
	if len(data) == 0 {
		store.NextID = 1
		return store, nil
	}

	// 3. Dosyadaki JSON metnini, bizim Go'daki Store yapısına çeviriyoruz (Unmarshal).
	// DİKKAT: &store diyerek pointer (adres) verdik. & işareti koymazsak Go
	// değişkenin kopyasını oluşturuyor ve asıl store'un içini dolduramıyor.
	err = json.Unmarshal(data, &store)
	if err != nil {
		// Dosyanın içine elle saçma sapan bir şey yazılmışsa program patlamasın (E10 kuralı).
		return store, errors.New("veri dosyası okunamadı: geçersiz format")
	}

	return store, nil
}

// Save: Bellekteki verileri (Store) alır ve diske kaydeder.
func Save(store Store) error {
	// JSON dosyası tek satır olmasın, alt alta düzgün (okunabilir) görünsün diye MarshalIndent kullanıyoruz.
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err // Dönüştürme başarısız olursa hatayı geri fırlat
	}

	// Dosyayı diske yaz.
	// 0644 dosya izinleridir (Sahibi okur ve yazar, diğerleri sadece okur).
	err = os.WriteFile(dataFile, data, 0644)
	if err != nil {
		return err
	}

	// Her şey sorunsuz bittiyse nil (boş) dön, yani hata yok.
	return nil
}
