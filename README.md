# CLI Görev Yöneticisi (gotodo)

## 1. Proje Nedir
Bu proje, Go ile sıfırdan geliştirilmiş, terminal üzerinden çalışan basit bir komut satırı görev yöneticisidir (Todo App). Harici hiçbir kütüphane kullanılmadan yazılmış olup, verileri JSON formatında kalıcı olarak saklar.


## 1. Proje Komutları
go run . //
            add        / Yeni Görev Ekleme
            done       / Görev Bitirme 
            delete     / Görev Silme
            list       / Görev Listeleme
            list all   / Detaylı Görev Listeleme 
7. Neler Öğrendim
   Python'dan farklı olarak, Go'da listelerden (slice) bir eleman silmek için özel bir remove fonksiyonu olmadığını; bunun yerine silinecek elemanın solundaki ve sağındaki dilimlerin birleştirildiğini öğrendim.

Go'da hataların try-catch blokları yerine her adımda err != nil ile yakalanıp manuel olarak işlenmesi mantığını kavradım.

Bir yapının (struct) güncellenebilmesi için fonksiyonlara o yapının kopyası yerine bellek adresinin (pointer &) gönderilmesi gerektiğini tecrübe ettim.

Tamamlanmamış görevlerin tarihlerini sıfır değeri (0001-01-01) yerine null bırakabilmek için pointer (*time.Time) kullanımının önemini anladım.

text/tabwriter paketini kullanarak terminal çıktılarını sütunlar halinde nasıl düzgün hizalayabileceğimi öğrendim.
