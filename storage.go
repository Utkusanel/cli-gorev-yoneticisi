package main

import (
	"encoding/json"
	"errors"
	"os"
)

// Veriyi tutacağımız dosya. Sabit yaptık ki her yerde aynı ismi kullanalım.
const dataFile = "tasks.json"

// JSON'ın en dışındaki ana yapı. Düz liste yerine struct yaptık çünkü
// hoca ileride yeni özellikler eklemek istersek düz listede patlarız demişti.
type Store struct {
	Tasks  []Task `json:"tasks"`
	NextID int    `json:"nextID"`
}

// Load: Program açılır açılmaz verileri dosyadan çekiyor
func Load() (Store, error) {
	var store Store

	data, err := os.ReadFile(dataFile)
	if err != nil {
		// os.ErrNotExist: Yani böyle bir dosya henüz diskte yok.
		// Program ilk kez açılıyorsa çökmek yerine boş bir store veriyoruz.
		if errors.Is(err, os.ErrNotExist) {
			store.NextID = 1
			return store, nil
		}
		// İzin hatası falan varsa gerçekten patlasın
		return store, err
	}

	// Dosya var ama içi tamamen boşsa (E11 kuralı için)
	if len(data) == 0 {
		store.NextID = 1
		return store, nil
	}

	// JSON'ı bizim Store struct'ına çevir.
	// DİKKAT: &store yazmayı unutunca çalışmıyor. Pointer ile adresini vermemiz lazım ki
	// asıl değişkenin içini doldurabilsin, kopyasıyla uğraşmasın.
	err = json.Unmarshal(data, &store)
	if err != nil {
		// İçine saçma sapan metin yazılırsa (E10 kuralı)
		return store, errors.New("veri dosyası okunamadı: geçersiz format")
	}

	return store, nil
}

// Save: Bellekteki son hali alıp diske json olarak eziyor
func Save(store Store) error {
	// MarshalIndent JSON dosyasının içini alt alta okunabilir yapıyor, yoksa tek satır çorba oluyor.
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	// 0644 klasik dosya okuma/yazma izni
	return os.WriteFile(dataFile, data, 0644)
}
