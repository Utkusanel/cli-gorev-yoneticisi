package main

import (
	"errors"
	"time"
)

// Görev yapısı. CompletedAt için pointer (*time.Time) kullandım çünkü
// görev bitmemişse nil (boş) bırakmam gerekiyor. Go'da normal time.Time nil olamıyor.
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Done        bool       `json:"done"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

// Yeni görev ekleme
// store değişkenini *Store (pointer) olarak aldık. Eğer * koymazsam Go kopyasını alıyor
// ve asıl listemiz güncellenmiyor. Bunu kavramam biraz zaman aldı.
func AddTask(store *Store, title string) Task {
	newTask := Task{
		ID:          store.NextID,
		Title:       title,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	// Slice'a yeni eleman ekle ve bir sonraki ID'yi hazırla
	store.Tasks = append(store.Tasks, newTask)
	store.NextID++

	return newTask
}

// Görevi bitirme
func CompleteTask(store *Store, id int) (Task, error) {
	for i, t := range store.Tasks {
		if t.ID == id {
			if t.Done {
				return t, errors.New("zaten tamamlanmış") // Zaten bitmişse uyarmak için hata dön
			}

			// Görevi bitir ve o anki saati ata.
			// Saatin adresini (&now) veriyoruz çünkü struct'ta pointer ile tanımladık.
			store.Tasks[i].Done = true
			now := time.Now()
			store.Tasks[i].CompletedAt = &now
			return store.Tasks[i], nil
		}
	}
	return Task{}, errors.New("bulunamadı")
}

// Görevi silme
func DeleteTask(store *Store, id int) (Task, error) {
	for i, t := range store.Tasks {
		if t.ID == id {
			// Pythondaki remove gibi kolay değilmiş :)
			// Silinecek elemanın solundaki listeyle sağındaki listeyi birleştiriyoruz.
			// ... işareti sağdaki dilimin elemanlarını tek tek ekliyor.
			store.Tasks = append(store.Tasks[:i], store.Tasks[i+1:]...)
			return t, nil
		}
	}
	return Task{}, errors.New("bulunamadı")
}