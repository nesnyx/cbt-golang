
# CBT Golang

Proyek ini merupakan implementasi sistem Computer-Based Test (CBT) menggunakan bahasa pemrograman Go (Golang).
Dirancang untuk menyediakan platform ujian daring yang efisien, aman, dan mudah digunakan.

## Fitur

- Manajemen pengguna dengan peran seperti admin, peserta, dan institusi.
- Autentikasi dan otorisasi pengguna.
- Pembuatan dan pengelolaan soal ujian.
- Penyimpanan dan evaluasi jawaban peserta.
- Pengelolaan sesi ujian untuk institusi.
- Penyajian hasil ujian dan analisis performa.

## Struktur Direktori

- `admin/` - Modul untuk manajemen admin.
- `auth/` - Modul untuk autentikasi dan otorisasi.
- `handler/` - Handler untuk routing HTTP.
- `helper/` - Fungsi-fungsi pembantu.
- `institution/` - Modul untuk manajemen institusi.
- `login/` - Modul untuk proses login pengguna.
- `peserta/` - Modul untuk manajemen peserta ujian.
- `result/` - Modul untuk pengolahan hasil ujian.
- `role/` - Modul untuk manajemen peran pengguna.
- `schema/` - Definisi skema database.
- `sessionsInstitution/` - Modul untuk sesi ujian institusi.
- `soal/` - Modul untuk manajemen soal ujian.
- `tmpAnswer/` - Modul untuk penyimpanan sementara jawaban peserta.
- `utils/` - Utilitas umum yang digunakan di seluruh aplikasi.

## Prasyarat

- Go versi 1.16 atau lebih baru.
- Database PostgreSQL atau MySQL.
- Git untuk kontrol versi.

## Instalasi

1. Klon repositori ini:
   ```bash
   git clone https://github.com/nesnyx/cbt-golang.git
   cd cbt-golang
   ```

2. Inisialisasi dan unduh dependensi:
   ```bash
   go mod tidy
   ```

3. Konfigurasikan koneksi database di file konfigurasi.

4. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

## Penggunaan

Setelah aplikasi berjalan, Anda dapat mengakses antarmuka pengguna melalui browser di `http://localhost:8080`.
Gunakan kredensial admin untuk masuk dan mulai mengelola ujian.

## Kontribusi

Kontribusi sangat diterima!
Silakan buat *issue* untuk melaporkan bug atau mengusulkan fitur baru.
Pull request juga sangat dihargai.

## Lisensi

Proyek ini dilisensikan di bawah MIT License.
Lihat file [LICENSE](LICENSE) untuk informasi lebih lanjut.
