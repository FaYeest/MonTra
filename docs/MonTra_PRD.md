# MonTra — Product Requirements Document (PRD)

**Versi:** 1.1  
**Status:** Draft untuk implementasi MVP  
**Tanggal:** 18 Oktober 2026  
**Produk:** MonTra (Money Tracker)  
**Kategori:** Personal Finance, AI/ML, FinTech  
**Target pasar awal:** Indonesia  
**Platform MVP:** Android  
**Bahasa antarmuka awal:** Bahasa Indonesia  
**Mata uang default:** Rupiah Indonesia (IDR)

---

## Daftar Isi

1. [Ringkasan Produk](#1-ringkasan-produk)
2. [Latar Belakang dan Masalah](#2-latar-belakang-dan-masalah)
3. [Visi, Tujuan, dan Non-Goals](#3-visi-tujuan-dan-non-goals)
4. [Target Pengguna dan Jobs to Be Done](#4-target-pengguna-dan-jobs-to-be-done)
5. [Prinsip Produk](#5-prinsip-produk)
6. [Ruang Lingkup dan Prioritas](#6-ruang-lingkup-dan-prioritas)
7. [Persyaratan Fungsional](#7-persyaratan-fungsional)
8. [Persyaratan AI/ML](#8-persyaratan-aiml)
9. [Alur Pengguna](#9-alur-pengguna)
10. [Persyaratan UI/UX](#10-persyaratan-uiux)
11. [Arsitektur Sistem dan Teknologi](#11-arsitektur-sistem-dan-teknologi)
12. [Desain Data](#12-desain-data)
13. [API Specification Awal](#13-api-specification-awal)
14. [Keamanan dan Privasi](#14-keamanan-dan-privasi)
15. [Persyaratan Nonfungsional](#15-persyaratan-nonfungsional)
16. [Metrik Keberhasilan](#16-metrik-keberhasilan)
17. [Strategi Pengujian dan Acceptance Criteria](#17-strategi-pengujian-dan-acceptance-criteria)
18. [Roadmap Implementasi](#18-roadmap-implementasi)
19. [Risiko dan Mitigasi](#19-risiko-dan-mitigasi)
20. [Kriteria Rilis MVP](#20-kriteria-rilis-mvp)
21. [Asumsi dan Keputusan yang Masih Terbuka](#21-asumsi-dan-keputusan-yang-masih-terbuka)
22. [Glosarium](#22-glosarium)

---

## 1. Ringkasan Produk

### 1.1 Deskripsi

**MonTra (Money Tracker)** adalah aplikasi pengelolaan keuangan pribadi yang membantu pengguna mencatat transaksi, mengelola akun keuangan, merencanakan anggaran, memantau arus kas, dan memahami kebiasaan finansial melalui analisis data dan fitur AI/Machine Learning.

MonTra bukan hanya buku kas digital. Aplikasi ini mengubah catatan transaksi menjadi informasi yang bisa digunakan: kategori pengeluaran otomatis, deteksi perubahan pola belanja, proyeksi pengeluaran bulanan, serta insight yang dapat ditelusuri kembali ke data sumber.

### 1.2 Ringkasan Solusi

MonTra memiliki empat kemampuan utama:

1. **Track:** mencatat pemasukan, pengeluaran, transfer, akun, dan anggaran.
2. **Understand:** menyajikan laporan dan pola pengeluaran.
3. **Predict:** memperkirakan total pengeluaran sampai akhir periode ketika data memadai.
4. **Improve:** menyarankan tindakan berdasarkan data tanpa mengambil keputusan atas nama pengguna.

### 1.3 Diferensiasi Produk

- Kategorisasi transaksi berbantuan ML.
- Insight yang disertai bukti numerik dan periode pembanding.
- Prediksi pengeluaran dengan informasi keterbatasan dan ketidakpastian.
- Arsitektur yang memisahkan pencatatan finansial deterministik dari model AI.
- Pengguna tetap mengendalikan transaksi, anggaran, dan keputusan finansial.

---

## 2. Latar Belakang dan Masalah

Pengguna sering menghadapi kendala dalam mengelola keuangan pribadi karena:

- Tidak tahu kategori mana yang paling banyak menghabiskan uang.
- Tidak mencatat transaksi secara konsisten.
- Tidak menyadari adanya kenaikan pengeluaran dari waktu ke waktu.
- Menetapkan anggaran tanpa melihat riwayat pengeluaran.
- Sulit memperkirakan apakah uang akan cukup hingga akhir bulan.
- Harus membaca laporan sendiri untuk memperoleh insight.

Aplikasi pencatatan keuangan dasar dapat menyimpan transaksi dan menampilkan grafik, tetapi insight yang relevan sering kali masih harus ditafsirkan secara manual.

MonTra bertujuan mengurangi beban tersebut dengan menggabungkan pengelolaan transaksi yang konsisten dan analitik yang dapat dijelaskan.

---

## 3. Visi, Tujuan, dan Non-Goals

### 3.1 Product Vision

Menjadi asisten pengelolaan keuangan pribadi yang memudahkan pengguna memahami ke mana uangnya pergi, memperkirakan kebutuhan finansial, dan membentuk kebiasaan pengelolaan uang yang lebih terencana.

### 3.2 Product Goals

- Mempermudah pencatatan transaksi sehari-hari.
- Menyediakan saldo dan ringkasan finansial yang konsisten dengan data transaksi.
- Mengurangi pekerjaan kategorisasi manual.
- Menampilkan tren dan perubahan pengeluaran yang relevan.
- Menyediakan prediksi pengeluaran dengan evaluasi yang terukur.
- Meningkatkan transparansi data dan kendali pengguna.
- Menjaga privasi catatan keuangan.

### 3.3 Non-Goals MVP

MonTra pada MVP **tidak** akan:

- Memproses pembayaran atau transfer uang secara nyata.
- Menyimpan dana pengguna.
- Mengakses rekening bank secara langsung tanpa integrasi resmi dan persetujuan pengguna.
- Menyediakan pinjaman, asuransi, atau layanan perbankan.
- Menjalankan perdagangan saham, kripto, atau aset investasi.
- Menjamin hasil finansial pengguna.
- Memberikan nasihat investasi profesional.
- Membiarkan LLM menghitung saldo atau mengubah transaksi keuangan.

---

## 4. Target Pengguna dan Jobs to Be Done

### 4.1 Persona A — Mahasiswa

**Karakteristik**
- Mengandalkan uang saku, beasiswa, atau pendapatan tambahan.
- Memiliki pengeluaran rutin untuk makanan, transportasi, pendidikan, dan hiburan.
- Ingin mengetahui berapa uang yang masih tersedia hingga akhir bulan.

**Kebutuhan**
- Input transaksi singkat.
- Anggaran per kategori.
- Peringatan ketika anggaran hampir habis.
- Insight sederhana tanpa istilah finansial yang rumit.

### 4.2 Persona B — Pekerja Awal Karier

**Karakteristik**
- Menerima penghasilan secara berkala.
- Memiliki tagihan, langganan, atau cicilan.
- Ingin mengendalikan pengeluaran dan meningkatkan tabungan.

**Kebutuhan**
- Pelacakan pemasukan dan pengeluaran.
- Transaksi berulang dan pengingat.
- Analisis tren bulanan.
- Proyeksi pengeluaran akhir bulan.

### 4.3 Persona C — Pengguna yang Ingin Memperbaiki Kebiasaan Finansial

**Karakteristik**
- Belum memiliki kebiasaan mencatat keuangan.
- Sering lupa mencatat transaksi kecil.
- Tidak mengetahui penyebab utama pembengkakan pengeluaran.

**Kebutuhan**
- Saran kategori otomatis.
- Deteksi perubahan pengeluaran.
- Rekomendasi anggaran berdasarkan riwayat.
- Penjelasan yang berbasis pada transaksi aktual.

### 4.4 Jobs to Be Done

- Ketika selesai melakukan transaksi, pengguna ingin mencatatnya dengan cepat agar catatan keuangan tetap akurat.
- Ketika mendekati akhir bulan, pengguna ingin memperkirakan kecukupan uangnya.
- Ketika pengeluaran meningkat, pengguna ingin mengetahui kategori penyebabnya.
- Ketika membuat target tabungan, pengguna ingin menilai apakah target tersebut realistis berdasarkan pola transaksi.
- Ketika memiliki banyak transaksi, pengguna ingin mendapatkan insight tanpa menganalisis semua data secara manual.

---

## 5. Prinsip Produk

1. **Financial correctness first.** Semua saldo, total pemasukan, total pengeluaran, dan transfer harus dihitung melalui business logic deterministik.
2. **AI is assistive.** Model AI/ML membantu pengguna, bukan mengambil alih kendali.
3. **Explainability.** Insight harus menjelaskan periode, angka, dan perbandingan yang menjadi dasar.
4. **Graceful degradation.** Gangguan service AI tidak boleh menghentikan pencatatan transaksi inti.
5. **Privacy by design.** Data yang dikirim ke komponen AI dibatasi pada informasi yang dibutuhkan.
6. **Progressive intelligence.** Model menyesuaikan tingkat kecanggihan dengan ketersediaan data, bukan memaksakan prediksi pada pengguna baru.
7. **User correction.** Pengguna dapat membetulkan kategori atau mengabaikan peringatan.
8. **No fabricated precision.** Sistem tidak boleh menampilkan angka prediksi seolah-olah merupakan hasil yang pasti.

---

## 6. Ruang Lingkup dan Prioritas

### 6.1 Prioritas MVP

| ID | Fitur | Prioritas | Keterangan |
|---|---|---|---|
| MVP-01 | Registrasi dan login | P0 | Akses aman ke akun pengguna |
| MVP-02 | Profil dan preferensi | P0 | Mata uang, zona waktu, pengaturan dasar |
| MVP-03 | CRUD transaksi | P0 | Pemasukan, pengeluaran, transfer |
| MVP-04 | Akun keuangan virtual | P0 | Tunai, rekening, e-wallet sebagai catatan |
| MVP-05 | Kategori transaksi | P0 | Kategori sistem dan kategori pribadi |
| MVP-06 | Dashboard | P0 | Ringkasan keuangan dan transaksi terbaru |
| MVP-07 | Anggaran | P0 | Anggaran bulanan dan per kategori |
| MVP-08 | Laporan | P0 | Ringkasan, tren, kategori, perbandingan periode |
| MVP-09 | ML categorization | P0 | Saran kategori dengan fallback manual |
| MVP-10 | Insight statistik | P0 | Temuan berdasarkan data yang dapat diverifikasi |
| MVP-11 | Forecast pengeluaran | P1 | Aktif jika data memadai; baseline sederhana dahulu |
| MVP-12 | Deteksi anomali | P1 | Tidak boleh diklaim sebagai deteksi penipuan |
| MVP-13 | Transaksi berulang | P1 | Pembuatan transaksi berdasarkan jadwal |
| MVP-14 | Ekspor CSV | P1 | Ekspor data milik pengguna |
| MVP-15 | Notifikasi anggaran | P1 | Ambang yang dapat dikonfigurasi |

### 6.2 Post-MVP

- AI financial assistant berbasis bahasa natural.
- OCR struk dan ekstraksi item belanja.
- Impor mutasi bank dari CSV dengan pemetaan kolom.
- Integrasi rekening melalui penyedia resmi yang sesuai ketentuan.
- Target tabungan dan estimasi waktu pencapaian.
- Pengelolaan utang dan rencana pelunasan.
- Sinkronisasi lintas perangkat.
- Web application.
- Analisis keuangan rumah tangga atau kelompok.

Fitur Post-MVP tidak boleh menunda penyelesaian pencatatan, saldo, anggaran, dan laporan yang akurat.

---

## 7. Persyaratan Fungsional

### FR-01 — Authentication dan User Profile

**Deskripsi:** Pengguna dapat membuat dan mengakses akun MonTra secara aman.

**Persyaratan:**
- Registrasi dan login menggunakan email dan kata sandi.
- Proses reset kata sandi.
- Verifikasi email jika kebijakan implementasi memerlukannya.
- Pengaturan nama tampilan, mata uang, zona waktu, dan preferensi aplikasi.
- Logout dan pencabutan sesi sesuai implementasi autentikasi.
- Permintaan penghapusan akun dan data.

**Acceptance criteria:**
- Pengguna yang belum terautentikasi tidak dapat membaca data pribadi.
- Setiap endpoint memeriksa otorisasi resource, bukan hanya validitas token.
- Pengguna A tidak dapat membaca atau mengubah data pengguna B dengan mengganti ID resource.
- Kegagalan autentikasi tidak mengungkap apakah kata sandi atau identitas tertentu benar.

### FR-02 — Dashboard

**Deskripsi:** Menyediakan ringkasan kondisi finansial pada periode terpilih.

**Komponen:**
- Total pemasukan.
- Total pengeluaran.
- Arus kas bersih pada periode tersebut.
- Saldo tiap akun dan total saldo akun yang dikelola.
- Visualisasi pengeluaran berdasarkan kategori.
- Perbandingan dengan periode sebelumnya.
- Kemajuan anggaran.
- Insight utama jika tersedia.
- Transaksi terbaru.

**Acceptance criteria:**
- Pengguna dapat mengganti periode laporan.
- Data dashboard konsisten dengan transaksi sumber.
- Transfer internal tidak dihitung sebagai pemasukan atau pengeluaran eksternal.
- Keadaan kosong ditangani dengan instruksi untuk menambahkan transaksi, bukan grafik kosong yang membingungkan.

### FR-03 — Transaction Management

**Deskripsi:** Pengguna dapat membuat, membaca, mengubah, menghapus, mencari, dan memfilter transaksi.

**Atribut minimum:**
- ID transaksi.
- Pengguna pemilik.
- Akun terkait.
- Jenis: `income`, `expense`, atau `transfer` sesuai model domain.
- Nominal dan kode mata uang.
- Deskripsi.
- Kategori jika relevan.
- Tanggal transaksi.
- Catatan opsional.
- Sumber: manual, import, atau OCR pada masa depan.
- Status/sumber kategori: manual atau prediksi model.
- Timestamp pembuatan dan perubahan.

**Persyaratan:**
- Validasi nominal dan tanggal.
- Filter menurut tanggal, tipe, akun, dan kategori.
- Pengurutan dan pagination pada daftar transaksi.
- Koreksi kategori setelah transaksi disimpan.
- Perubahan atau penghapusan transaksi memperbarui laporan turunan.
- Sediakan mekanisme idempotensi untuk operasi create yang dapat terkirim ulang.

**Acceptance criteria:**
- Transaksi berhasil disimpan sebelum sistem menampilkan konfirmasi sukses.
- Perubahan nominal, tanggal, akun, jenis, atau kategori tercermin pada laporan yang relevan.
- Penghapusan transaksi tidak meninggalkan agregasi yang salah.
- Proses gagal tidak boleh menyebabkan data setengah tersimpan.

### FR-04 — Internal Transfers

**Deskripsi:** Pengguna dapat mencatat pemindahan dana di antara akun virtual miliknya.

**Persyaratan:**
- Memiliki akun sumber dan akun tujuan yang berbeda.
- Mendukung nominal transfer dan biaya transfer opsional.
- Memvalidasi bahwa kedua akun dimiliki pengguna yang sama.
- Mengurangi saldo sumber dan menambah saldo tujuan secara atomik.
- Biaya transfer, jika ada, dicatat sebagai beban eksternal.

**Acceptance criteria:**
- Transfer bukan pemasukan atau pengeluaran eksternal.
- Nilai bersih seluruh akun tidak berubah karena pokok transfer.
- Kegagalan saat menyimpan transfer tidak boleh memperbarui hanya salah satu akun.

### FR-05 — Accounts and Balance

**Deskripsi:** Pengguna mengatur akun virtual yang digunakan untuk melacak lokasi uang.

**Jenis akun contoh:**
- Tunai.
- Rekening bank.
- E-wallet.
- Akun khusus lainnya.

**Persyaratan:**
- Pengguna dapat membuat, mengubah nama, mengarsipkan, atau menghapus akun jika tidak melanggar integritas data.
- Saldo dimulai dari saldo awal dan dipengaruhi transaksi yang relevan.
- Akun dan transaksi menyimpan konteks mata uang.
- Jika akun sudah digunakan, penghapusan harus dicegah atau dilakukan melalui prosedur yang mempertahankan integritas histori.

**Acceptance criteria:**
- Saldo dapat direkonsiliasi terhadap saldo awal dan pergerakan transaksi.
- Sistem menjelaskan bahwa saldo merupakan catatan pengguna, bukan saldo bank yang terverifikasi, kecuali tersedia integrasi resmi.

### FR-06 — Categories

**Deskripsi:** Mengelola kategori pemasukan dan pengeluaran.

**Contoh pengeluaran:**
- Makanan dan minuman.
- Transportasi.
- Belanja.
- Tagihan.
- Pendidikan.
- Kesehatan.
- Hiburan.
- Langganan.
- Cicilan dan utang.
- Lainnya.

**Contoh pemasukan:**
- Gaji.
- Uang saku.
- Freelance.
- Bisnis.
- Hadiah.
- Pemasukan lainnya.

**Persyaratan:**
- Menyediakan kategori sistem.
- Memungkinkan pengguna membuat kategori pribadi.
- Membedakan jenis transaksi yang dapat memakai kategori.
- Menangani kategori yang sudah digunakan melalui arsip atau pemetaan ulang, bukan menghapus histori secara diam-diam.

### FR-07 — Budget Management

**Deskripsi:** Pengguna mengatur batas pengeluaran per periode.

**Persyaratan:**
- Mendukung anggaran keseluruhan dan per kategori.
- Mendukung periode mingguan atau bulanan; MVP dapat memulai dengan bulanan.
- Menampilkan nominal terpakai, sisa anggaran, dan persentase pemakaian.
- Mengizinkan ambang peringatan.
- Menampilkan riwayat penggunaan anggaran.
- Memberi saran nominal anggaran berdasarkan data historis jika tersedia.
- Tidak mengubah anggaran secara otomatis tanpa tindakan pengguna.

**Acceptance criteria:**
- Penggunaan anggaran dihitung dari transaksi pengeluaran yang relevan dalam periode yang benar.
- Transfer internal tidak mengurangi anggaran.
- Perubahan periode, zona waktu, atau tanggal transaksi ditangani konsisten.
- Anggaran yang terlampaui ditandai tanpa menyembunyikan nominal aktual.

### FR-08 — Reports and Analytics

**Deskripsi:** Menampilkan laporan keuangan yang berasal dari transaksi sumber.

**Laporan minimum:**
- Ringkasan pemasukan dan pengeluaran.
- Arus kas bersih.
- Pengeluaran berdasarkan kategori.
- Tren pengeluaran per hari atau bulan.
- Perbandingan antarperiode.
- Ringkasan per akun.
- Ekspor CSV sesuai filter.

**Acceptance criteria:**
- Total laporan dapat direkonsiliasi ke transaksi sumber.
- Periode laporan ditampilkan secara eksplisit.
- Transfer internal tidak dihitung dua kali sebagai aktivitas eksternal.
- Ekspor hanya memuat data pengguna yang terautentikasi.

### FR-09 — Recurring Transactions

**Deskripsi:** Pengguna dapat membuat template transaksi berulang.

**Contoh:** gaji bulanan, tagihan internet, atau langganan.

**Persyaratan:**
- Mendukung jadwal berulang yang ditetapkan dengan jelas.
- Pengguna dapat menjeda, mengubah, atau menonaktifkan jadwal.
- Sistem harus menangani tanggal yang tidak ada pada bulan tertentu.
- Pekerjaan terjadwal harus idempoten agar satu kejadian tidak dibuat dua kali.
- Pengguna diberi tahu apakah transaksi dibuat otomatis atau hanya sebagai pengingat; pilihan ini harus eksplisit.

### FR-10 — Notifications

**Deskripsi:** Memberi peringatan atas anggaran atau insight yang memerlukan perhatian.

**Persyaratan:**
- Notifikasi dapat dinonaktifkan atau dikonfigurasi.
- Pengguna dapat memilih ambang anggaran.
- Deduplication mencegah notifikasi identik terlalu sering.
- Jangan mengungkap isi transaksi secara sensitif pada layar kunci secara default.

### FR-11 — Data Export and Account Deletion

**Deskripsi:** Pengguna dapat mengekspor data tertentu dan meminta penghapusan data.

**Persyaratan:**
- Ekspor transaksi milik pengguna dalam CSV.
- Validasi dan otorisasi permintaan ekspor.
- Permintaan penghapusan akun memiliki konfirmasi yang jelas.
- Hapus atau anonimisasi data turunan dan artefak terkait sesuai kebijakan retensi yang dinyatakan.
- Jelaskan keterbatasan penghapusan dari backup berdasarkan kebijakan retensi.

---

## 8. Persyaratan AI/ML

AI/ML harus menjadi layanan pendukung yang dapat dimatikan secara terpisah dari alur transaksi inti.

### AI-01 — Automatic Transaction Categorization

**Tujuan:** Mengusulkan kategori transaksi dari deskripsi dan fitur yang diizinkan.

**Contoh:**
- Deskripsi: `KOPI KENANGAN`
- Jenis: pengeluaran
- Nominal: Rp25.000
- Prediksi: Makanan dan minuman

**Pipeline:**
1. Go menyimpan transaksi sebagai sumber kebenaran.
2. Go mengirim fitur minimum yang dibutuhkan ke Python.
3. Python menormalisasi input dan menjalankan model.
4. Python mengembalikan prediksi, confidence jika tersedia, dan versi model.
5. Go memvalidasi respons dan menerapkan kebijakan confidence.
6. Prediksi berkeyakinan tinggi dapat menjadi kategori yang disarankan atau diterapkan sesuai aturan produk.
7. Prediksi berkeyakinan rendah meminta konfirmasi pengguna.
8. Koreksi pengguna dapat dikumpulkan sebagai label kandidat setelah validasi dan sesuai persetujuan data.

**Baseline yang disarankan:** TF-IDF + Logistic Regression dari scikit-learn.

**Model lanjutan:** Transformer untuk klasifikasi teks hanya jika evaluasi membuktikan peningkatan yang relevan terhadap biaya dan kompleksitas.

**Fitur input yang dipertimbangkan:**
- Deskripsi/merchant yang dinormalisasi.
- Tipe transaksi.
- Mata uang.
- Fitur tambahan yang tidak sensitif dan terbukti meningkatkan kinerja.

Nominal tidak boleh diasumsikan sebagai penentu kategori tunggal. Hindari fitur sensitif yang tidak diperlukan.

**Evaluasi:**
- Macro-F1 sebagai metrik utama.
- Precision, recall, dan confusion matrix per kategori.
- Kalibrasi confidence jika confidence digunakan untuk keputusan otomatis.
- Evaluasi terpisah untuk kategori langka.

**Target awal:** Macro-F1 minimal 0,80 pada data uji representatif. Ini merupakan target pengembangan, bukan klaim performa yang sudah dicapai.

**Acceptance criteria:**
- Training, validation, dan test dipisahkan dengan benar.
- Pengguna bisa memperbaiki kategori setiap saat.
- Kategori berkeyakinan rendah tidak diterapkan otomatis tanpa aturan persetujuan yang eksplisit.
- Ketika service/model gagal, transaksi masih dapat disimpan secara manual.

### AI-02 — Spending Anomaly Detection

**Tujuan:** Mengidentifikasi perubahan pola pengeluaran yang patut diperiksa.

**Jenis temuan yang mungkin:**
- Pengeluaran kategori meningkat dibanding periode pembanding.
- Nominal transaksi berbeda dari pola pengguna.
- Frekuensi belanja meningkat.
- Ada transaksi yang mungkin duplikat.

**Pendekatan awal:** aturan statistik dan perbandingan historis. Uji Robust Z-score atau Isolation Forest bila data dan pola sesuai. Model lebih kompleks tidak otomatis lebih baik.

**Persyaratan:**
- Setiap peringatan menyertakan alasan terukur.
- Model membedakan transaksi besar dari perubahan pola kategori.
- Pengguna dapat menandai peringatan relevan/tidak relevan.
- Sistem mengendalikan frekuensi peringatan serupa.
- Tidak menyebut transaksi sebagai penipuan atau fraud hanya berdasarkan anomali.

**Metrik evaluasi:** precision per peringatan, false-positive rate, dan umpan balik pengguna.

### AI-03 — Monthly Expense Forecasting

**Tujuan:** Memperkirakan total pengeluaran hingga akhir periode.

**Pendekatan bertahap:**
1. Baseline rata-rata historis atau proyeksi sederhana berbasis hari.
2. Moving average atau exponential smoothing.
3. Model time-series yang lebih kompleks hanya jika mengungguli baseline pada evaluasi yang realistis.

**Persyaratan:**
- Model harus mengetahui periode tanggal dan cakupan riwayat yang tersedia.
- Sistem tidak memprediksi terlalu percaya diri ketika riwayat minim.
- Jika data belum memadai, tampilkan status bahwa prediksi belum tersedia.
- Prediksi menyertakan estimasi ketidakpastian ketika pendekatan mendukungnya.
- Prediksi dibandingkan dengan nilai aktual setelah periode selesai.
- Nilai aktual pengeluaran tetap dihitung dari transaksi, bukan hasil prediksi.

**Metrik:** MAE dan WAPE atau metrik relatif lain yang cocok dengan distribusi data.

**Target awal:** Mengurangi error setidaknya 10% dibanding baseline pada data evaluasi yang representatif. Jika tidak tercapai, gunakan baseline yang lebih andal dan dokumentasikan hasilnya.

### AI-04 — Personalized Financial Insights

**Tujuan:** Mengubah hasil analitik menjadi insight yang dapat ditindaklanjuti.

**Contoh:**
- Pengeluaran makanan naik 18% dibanding bulan sebelumnya.
- Tiga kategori tertentu menyumbang sebagian besar pengeluaran periode ini.
- Berdasarkan tren saat ini, anggaran transportasi berpotensi terlampaui.
- Total pengeluaran langganan pada periode ini mencapai Rp180.000.

**Pipeline:**
1. Go atau modul analitik menghitung agregasi angka yang terverifikasi.
2. Rules engine/analitik memilih temuan yang relevan.
3. Sistem membuat payload fakta yang terstruktur.
4. Jika digunakan, LLM hanya memformulasikan fakta tersebut menjadi bahasa natural.
5. Validasi memeriksa angka, periode, dan jenis insight sebelum ditampilkan.

**Persyaratan:**
- Insight harus bisa ditelusuri ke agregasi yang menjadi dasarnya.
- Insight menyatakan periode dan pembanding secara eksplisit.
- Tidak menghasilkan insight saat data tidak cukup.
- Rekomendasi tidak memaksa dan tidak mengubah anggaran otomatis.
- Pengguna dapat menyembunyikan insight atau memberi feedback.

### AI-05 — Natural Language Financial Assistant (Post-MVP)

**Contoh pertanyaan:**
- "Berapa pengeluaran makan saya bulan ini?"
- "Kategori mana yang paling meningkat?"
- "Apakah saya berpotensi melewati anggaran bulan ini?"
- "Berapa yang perlu saya sisihkan setiap minggu untuk mencapai target tabungan?"

**Arsitektur:**
1. Pengguna mengirim pertanyaan ke backend Go.
2. Go mengautentikasi pengguna dan memvalidasi permintaan.
3. Intent dipetakan ke fungsi backend terdaftar dengan skema parameter yang ketat.
4. Backend menghitung atau mengambil agregasi yang diotorisasi.
5. LLM menyusun jawaban hanya dari hasil tersebut.
6. Jawaban melewati pemeriksaan format dan konteks sebelum dikirim.

**Pengamanan:**
- LLM tidak mendapat akses SQL bebas ke database.
- Fungsi yang tersedia memiliki parameter dan otorisasi ketat.
- Jangan pernah mengirim password, PIN, OTP, token autentikasi, atau kredensial bank.
- Jika data tidak tersedia, asisten harus menyatakan keterbatasan.
- Jawaban tidak boleh mengklaim hasil pasti untuk prediksi.

### AI-06 — Receipt OCR (Post-MVP)

Pengguna dapat mengambil atau mengunggah foto struk untuk mengekstrak merchant, tanggal, total, dan—jika didukung—item pembelian.

Semua hasil OCR harus masuk ke formulir konfirmasi. OCR tidak boleh langsung membuat transaksi final tanpa kesempatan pengguna memeriksa hasilnya.

### 8.1 Strategi Data ML

Sumber data awal:
- Transaksi yang dimasukkan pengguna secara manual.
- Koreksi kategori pengguna yang telah diproses sesuai kebijakan privasi.
- Data sintetis untuk pengujian awal.
- Dataset publik yang lisensinya sesuai jika diperlukan.

Data sintetis tidak boleh dianggap sebagai representasi sempurna perilaku transaksi nyata.

### 8.2 Siklus Hidup Model

1. Menentukan tujuan dan baseline.
2. Menetapkan skema label.
3. Membersihkan dan menormalkan data.
4. Memisahkan train/validation/test dengan mencegah leakage.
5. Melatih dan mengevaluasi model.
6. Menyimpan versi model dan hasil evaluasi.
7. Menguji inference dan compatibility API.
8. Melakukan rollout dengan pemantauan.
9. Menonaktifkan atau mengembalikan versi model jika terjadi regresi.
10. Melatih ulang hanya setelah pemeriksaan kualitas, privasi, dan performa.

### 8.3 Cold Start

Saat pengguna baru belum memiliki histori yang cukup:
- Gunakan model global atau aturan sederhana untuk kategorisasi.
- Tampilkan laporan historis tanpa prediksi yang berlebihan.
- Jelaskan bahwa insight personal membutuhkan data lebih banyak.
- Aktifkan fitur prediksi hanya setelah ambang data dan kualitas yang ditentukan tercapai.

---

## 9. Alur Pengguna

### 9.1 First-Time User Flow

1. Pengguna membuka MonTra.
2. Pengguna registrasi/login.
3. Pengguna memilih preferensi dasar dan zona waktu.
4. Pengguna membuat akun virtual dan memasukkan saldo awal jika diperlukan.
5. Pengguna memilih atau melewati pengaturan anggaran awal.
6. Pengguna mencatat transaksi pertama.
7. Dashboard menampilkan ringkasan yang tersedia.
8. Aplikasi menjelaskan bahwa insight personal akan berkembang seiring bertambahnya data.

### 9.2 Transaction Flow

1. Pengguna menekan tombol tambah transaksi.
2. Pengguna memilih pemasukan, pengeluaran, atau transfer.
3. Pengguna mengisi nominal, deskripsi, akun, dan tanggal.
4. Untuk pemasukan/pengeluaran, ML dapat menyarankan kategori.
5. Pengguna menerima atau mengoreksi saran.
6. Transaksi disimpan melalui Go backend.
7. Laporan dan saldo diperbarui berdasarkan transaksi berhasil.
8. Kegagalan AI tidak menggagalkan pencatatan manual.

### 9.3 Insight Flow

1. Pengguna membuka halaman Insights.
2. Aplikasi memuat insight dan rentang waktunya.
3. Aplikasi menampilkan angka, periode pembanding, dan konteks.
4. Pengguna dapat membuka transaksi atau laporan pendukung.
5. Pengguna memberi feedback atau menyembunyikan insight.

### 9.4 Forecast Flow

1. Pengguna membuka detail prediksi bulanan.
2. Aplikasi menampilkan data aktual hingga saat ini.
3. Aplikasi menampilkan estimasi akhir bulan dan ketidakpastian jika tersedia.
4. Aplikasi menyebutkan jumlah/riwayat data yang digunakan bila bermanfaat.
5. Pengguna dapat membandingkan prediksi dengan anggaran.

---

## 10. Persyaratan UI/UX

### 10.1 Design Direction

- Gaya modern, minimalis, dan berorientasi pada keterbacaan data.
- Mendukung light mode dan dark mode.
- Kontras teks dan komponen yang memadai.
- Warna tidak menjadi satu-satunya penanda status.
- Bahasa utama antarmuka adalah Bahasa Indonesia.
- Nominal rupiah diformat konsisten.
- Grafik sederhana, berlabel jelas, dan tidak menyesatkan.

### 10.2 Navigasi Utama

1. **Home:** ringkasan kondisi keuangan.
2. **Transactions:** riwayat dan pengelolaan transaksi.
3. **Budgets:** anggaran dan progres.
4. **Insights:** analisis dan prediksi AI/ML.
5. **Profile:** pengaturan akun, preferensi, dan privasi.

Aksi menambahkan transaksi harus mudah ditemukan dari dashboard.

### 10.3 Daftar Layar MVP

1. Splash/loading.
2. Login.
3. Registrasi.
4. Reset kata sandi.
5. Onboarding.
6. Home dashboard.
7. Form tambah/edit transaksi.
8. Transaction list.
9. Transaction detail.
10. Category management.
11. Account management.
12. Budget overview.
13. Budget detail.
14. Financial reports.
15. AI insights.
16. Forecast detail bila fitur tersedia.
17. Notification center.
18. Profile dan settings.
19. Privacy/data management.
20. Empty, loading, error, dan offline states.

### 10.4 Panduan Penyajian AI

- Tandai prediksi sebagai estimasi, bukan fakta aktual.
- Tampilkan periode analisis.
- Sediakan dasar atau alasan untuk insight.
- Hindari kepastian palsu dan istilah yang mengintimidasi.
- Sediakan koreksi manual dan mekanisme feedback.
- Jangan menyajikan peringatan anomali sebagai bukti penipuan.

---

## 11. Arsitektur Sistem dan Teknologi

### 11.1 Technology Stack

| Lapisan | Teknologi | Tanggung jawab |
|---|---|---|
| Mobile | Flutter + Dart | UI, navigasi, form, state, komunikasi API |
| Backend utama | Golang | REST API, autentikasi, otorisasi, business logic finansial |
| Database | PostgreSQL | Sumber kebenaran data keuangan |
| AI/ML service | Python + FastAPI | Inference, evaluasi, prediksi, deteksi anomali |
| ML libraries | scikit-learn, pandas, NumPy | Preprocessing dan model baseline |
| Validasi ML API | Pydantic | Validasi request/response Python |
| Go database driver | pgx | Akses PostgreSQL |
| Go router | Chi atau Gin | Routing HTTP; pilih satu |
| Migrations | golang-migrate atau setara | Versioning schema database |
| Model serialization | Joblib untuk model yang kompatibel | Penyimpanan artefak model |
| Deployment MVP | Docker Compose | Menjalankan service secara konsisten |
| Reverse proxy | Caddy atau Nginx | TLS dan routing eksternal |
| Logging | `log/slog` atau logger terstruktur | Observability tanpa membocorkan data sensitif |

### 11.2 Komponen Arsitektur

```text
Flutter Mobile App
        |
        | HTTPS / REST API
        v
+----------------------------------+
| Go Backend                       |
|                                  |
| Authentication & Authorization   |
| Transaction Management           |
| Accounts & Transfers             |
| Categories & Budgets             |
| Reports & Financial Calculations |
| Insight Orchestration            |
| Python ML Client                 |
+----------------------------------+
       |                     |
       | SQL                 | Internal HTTP/JSON
       v                     v
+------------------+  +--------------------------+
| PostgreSQL       |  | Python FastAPI Service   |
|                  |  |                          |
| Users            |  | Category Classification |
| Accounts         |  | Anomaly Detection        |
| Transactions     |  | Expense Forecasting      |
| Budgets          |  | Structured Insights      |
| Predictions      |  | Model Evaluation         |
+------------------+  +--------------------------+
                              |
                              v
                      +------------------+
                      | Versioned Models |
                      +------------------+
```

### 11.3 Tanggung Jawab Go Backend

Go merupakan backend utama dan sumber eksekusi business logic finansial.

Modul yang disarankan:
- `auth`: autentikasi dan sesi.
- `users`: profil dan preferensi.
- `transactions`: transaksi pemasukan dan pengeluaran.
- `accounts`: akun virtual dan saldo.
- `transfers`: transfer antar-akun.
- `categories`: kategori sistem dan pengguna.
- `budgets`: anggaran dan ambang peringatan.
- `reports`: ringkasan dan agregasi keuangan.
- `insights`: pengambilan insight terverifikasi.
- `mlclient`: komunikasi internal ke Python.
- `notifications`: notifikasi dan deduplikasi.

Go harus memvalidasi input, mengotorisasi pemilik resource, melakukan transaksi database atomik, dan menerapkan aturan fallback bila service AI tidak tersedia.

### 11.4 Tanggung Jawab Python Service

Python menangani:
- Preprocessing dan normalisasi fitur ML.
- Inference klasifikasi kategori.
- Deteksi anomali dan prediksi pengeluaran.
- Validasi output model.
- Versioning dan metadata model.
- Pipeline training/evaluation yang terpisah dari API produksi.

Python tidak boleh mengubah transaksi atau saldo otoritatif secara langsung. Untuk MVP, Go memasok fitur dan agregasi minimum yang diperlukan; Python tidak perlu memiliki akses langsung ke database transaksi.

### 11.5 Komunikasi Service

Komunikasi awal menggunakan HTTP/JSON pada jaringan internal Docker.

Persyaratan:
- Python tidak diekspos ke internet publik.
- Go menerapkan timeout request dan response.
- Semua request/response divalidasi dengan schema.
- Terapkan autentikasi antarlayanan atau kontrol jaringan yang sesuai.
- Jangan mengirim kata sandi, token pengguna, PIN, OTP, atau kredensial eksternal ke Python.
- Batasi data transaksi yang dikirim ke model.
- Kegagalan ML tidak menggagalkan penyimpanan transaksi.
- Retry hanya diterapkan pada operasi yang aman untuk diulang.
- Gunakan versi API internal yang eksplisit.

### 11.6 Struktur Repository yang Disarankan

```text
montra/
├── AGENTS.md
├── README.md
├── PRD.md
├── .env.example
├── .gitignore
├── docker-compose.yml
├── apps/
│   └── mobile/
│       ├── lib/
│       ├── test/
│       └── pubspec.yaml
├── services/
│   ├── api/
│   │   ├── cmd/
│   │   │   └── api/
│   │   ├── internal/
│   │   │   ├── auth/
│   │   │   ├── users/
│   │   │   ├── transactions/
│   │   │   ├── accounts/
│   │   │   ├── transfers/
│   │   │   ├── categories/
│   │   │   ├── budgets/
│   │   │   ├── reports/
│   │   │   ├── insights/
│   │   │   ├── notifications/
│   │   │   └── mlclient/
│   │   ├── migrations/
│   │   ├── tests/
│   │   ├── go.mod
│   │   └── Dockerfile
│   └── ml/
│       ├── app/
│       │   ├── main.py
│       │   ├── api/
│       │   ├── schemas/
│       │   ├── services/
│       │   └── models/
│       ├── training/
│       ├── evaluation/
│       ├── artifacts/
│       ├── tests/
│       ├── pyproject.toml
│       └── Dockerfile
├── packages/
│   └── api-spec/
│       └── openapi.yaml
└── infrastructure/
    ├── docker/
    └── proxy/
```

Model artefak besar atau sensitif tidak otomatis di-commit ke Git. Tetapkan penyimpanan artefak dan kontrol akses ketika ukuran dan kebutuhan model meningkat.

### 11.7 Prinsip Arsitektur

- PostgreSQL menyimpan sumber kebenaran finansial.
- Go mengelola perhitungan bisnis dan operasi data transaksi.
- Python menghasilkan output ML yang merupakan data turunan.
- Semua hasil prediksi harus memiliki metadata model dan waktu pembuatan.
- Prediksi dan insight dapat dihitung ulang.
- Kegagalan AI tidak boleh merusak catatan transaksi.
- Mulai dengan modular monolith untuk backend Go dan satu service Python terpisah; jangan menambah infrastruktur terdistribusi yang belum dibutuhkan.

---

## 12. Desain Data

Desain berikut adalah baseline konseptual. Tipe kolom, indeks, constraint, dan kebijakan retensi harus ditetapkan pada migrasi database.

### 12.1 Entitas Inti

#### `users`
- `id`
- `email`
- `password_hash` atau referensi identity provider
- `display_name`
- `preferred_currency`
- `timezone`
- `created_at`
- `updated_at`
- `deleted_at` jika soft-delete dibutuhkan

#### `accounts`
- `id`
- `user_id`
- `name`
- `account_type`
- `opening_balance`
- `currency`
- `created_at`
- `updated_at`
- `archived_at` nullable

#### `categories`
- `id`
- `user_id` nullable untuk kategori sistem
- `name`
- `transaction_type`
- `parent_category_id` nullable
- `is_system`
- `is_archived`
- `created_at`
- `updated_at`

#### `transactions`
- `id`
- `user_id`
- `account_id`
- `category_id` nullable
- `type` (`income` atau `expense`; transfer dapat dikelola melalui entitas khusus)
- `amount`
- `currency`
- `description`
- `transaction_date`
- `notes` nullable
- `source` (`manual`, `import`, `ocr`)
- `category_source` (`manual`, `model`, `system`)
- `category_model_version` nullable
- `created_at`
- `updated_at`

#### `transfers`
- `id`
- `user_id`
- `source_account_id`
- `destination_account_id`
- `amount`
- `currency`
- `fee_amount` default 0
- `transaction_date`
- `description` nullable
- `created_at`
- `updated_at`

#### `budgets`
- `id`
- `user_id`
- `category_id` nullable untuk anggaran keseluruhan
- `amount`
- `currency`
- `period_type`
- `start_date`
- `end_date`
- `alert_threshold_percent`
- `created_at`
- `updated_at`

#### `recurring_transactions`
- `id`
- `user_id`
- `account_id`
- `category_id` nullable
- `type`
- `amount`
- `currency`
- `recurrence_rule`
- `next_run_at`
- `is_active`
- `created_at`
- `updated_at`

#### `ml_predictions`
- `id`
- `user_id`
- `transaction_id` nullable
- `prediction_type`
- `predicted_value` dalam format terstruktur yang sesuai
- `confidence` nullable
- `model_version`
- `created_at`

#### `insights`
- `id`
- `user_id`
- `insight_type`
- `period_start`
- `period_end`
- `evidence` sebagai data terstruktur
- `generated_text`
- `model_version` nullable
- `created_at`
- `expires_at` nullable

#### `user_feedback`
- `id`
- `user_id`
- `entity_type`
- `entity_id`
- `feedback_type`
- `created_at`

#### `notification_events`
- `id`
- `user_id`
- `event_type`
- `deduplication_key`
- `payload` terstruktur dan minimum
- `created_at`
- `read_at` nullable

### 12.2 Aturan Integritas Data

- Nominal uang menggunakan tipe numerik yang sesuai dan tidak memakai floating-point untuk nilai finansial otoritatif.
- Simpan kode mata uang secara eksplisit.
- Tetapkan konvensi nominal secara konsisten, misalnya nominal transaksi selalu non-negatif dan tipe transaksi menentukan arah arus.
- Transfer harus dimodelkan konsisten dan tidak dihitung sebagai pemasukan/pengeluaran eksternal.
- Terapkan foreign key dan constraint kepemilikan melalui backend serta kontrol database yang sesuai.
- Operasi yang memengaruhi lebih dari satu baris terkait harus atomik.
- Gunakan indeks untuk kueri umum seperti `user_id`, `transaction_date`, `account_id`, dan `category_id` setelah pola akses dikonfirmasi.
- Definisikan kebijakan saat kategori atau akun yang telah digunakan diarsipkan.

### 12.3 Prinsip Kepemilikan Data

Setiap query terhadap data privat harus dibatasi pada `user_id` dari konteks autentikasi. Jangan mempercayai `user_id` yang dikirim bebas oleh client sebagai bukti kepemilikan.

---

## 13. API Specification Awal

Seluruh endpoint publik menggunakan prefix `/api/v1`. Nama dan payload berikut adalah rancangan awal yang perlu difinalisasi melalui OpenAPI sebelum implementasi penuh.

### 13.1 Authentication

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `POST /api/v1/auth/forgot-password`
- `POST /api/v1/auth/reset-password`

### 13.2 User Profile

- `GET /api/v1/me`
- `PATCH /api/v1/me`
- `DELETE /api/v1/me`

### 13.3 Transactions

- `GET /api/v1/transactions`
- `POST /api/v1/transactions`
- `GET /api/v1/transactions/{id}`
- `PATCH /api/v1/transactions/{id}`
- `DELETE /api/v1/transactions/{id}`

Query list dapat mendukung filter tanggal, kategori, tipe, akun, search, sort, dan pagination.

### 13.4 Accounts and Transfers

- `GET /api/v1/accounts`
- `POST /api/v1/accounts`
- `GET /api/v1/accounts/{id}`
- `PATCH /api/v1/accounts/{id}`
- `POST /api/v1/transfers`
- `GET /api/v1/transfers`
- `GET /api/v1/transfers/{id}`
- `PATCH /api/v1/transfers/{id}`
- `DELETE /api/v1/transfers/{id}`

### 13.5 Categories

- `GET /api/v1/categories`
- `POST /api/v1/categories`
- `PATCH /api/v1/categories/{id}`
- `DELETE /api/v1/categories/{id}` atau mekanisme arsip yang aman

### 13.6 Budgets

- `GET /api/v1/budgets`
- `POST /api/v1/budgets`
- `GET /api/v1/budgets/{id}`
- `PATCH /api/v1/budgets/{id}`
- `DELETE /api/v1/budgets/{id}`

### 13.7 Reports and Insights

- `GET /api/v1/reports/summary`
- `GET /api/v1/reports/categories`
- `GET /api/v1/reports/trends`
- `GET /api/v1/reports/accounts`
- `GET /api/v1/reports/export.csv`
- `GET /api/v1/insights`
- `GET /api/v1/forecasts/monthly`
- `POST /api/v1/insights/{id}/feedback`

### 13.8 Python Internal API

Endpoint ini hanya untuk komunikasi internal antarlayanan dan bukan bagian dari API publik.

- `GET /internal/v1/health`
- `POST /internal/v1/predict/category`
- `POST /internal/v1/detect/anomalies`
- `POST /internal/v1/predict/expense-forecast`
- `POST /internal/v1/insights/compose` jika lapisan generasi bahasa natural diaktifkan

### 13.9 Contoh Request Kategorisasi

`POST /internal/v1/predict/category`

```json
{
  "description": "Kopi Kenangan",
  "transaction_type": "expense",
  "currency": "IDR",
  "amount": 25000
}
```

Contoh respons:

```json
{
  "predicted_category": "food_and_drink",
  "confidence": 0.91,
  "model_version": "category-v1"
}
```

Angka confidence pada contoh hanya ilustrasi, bukan klaim kinerja model. Sistem produksi harus memakai output model yang dievaluasi dan confidence yang dikalibrasi bila digunakan untuk pengambilan keputusan.

### 13.10 API Contract Requirements

- Definisikan schema sukses dan error secara konsisten.
- Gunakan kode HTTP sesuai hasil operasi.
- Validasi tipe, rentang nominal, tanggal, dan enum.
- Terapkan pagination untuk daftar besar.
- Terapkan idempotency key untuk operasi tertentu yang dapat diulang.
- Jangan menampilkan stack trace atau informasi internal dalam error publik.
- Dokumentasikan autentikasi, otorisasi, batas request, dan perilaku retry di OpenAPI.

---

## 14. Keamanan dan Privasi

### SEC-01 — Authentication dan Authorization

- Jika autentikasi dikelola sendiri, gunakan password hashing yang kuat dan algoritma yang sesuai praktik terkini.
- Gunakan sesi atau token yang dikelola secara aman.
- Terapkan otorisasi pada setiap endpoint dan resource.
- Terapkan rate limiting untuk autentikasi dan operasi yang berisiko.
- Jangan memasukkan secret ke dalam repository.

### SEC-02 — Data Protection

- Gunakan HTTPS untuk semua komunikasi eksternal.
- Lindungi database, backup, dan artefak model dengan kontrol akses.
- Minimalkan data sensitif yang disimpan.
- Jangan menulis deskripsi transaksi, nominal sensitif, token, atau secret ke log secara default.
- Validasi unggahan file jika OCR ditambahkan.
- Terapkan kebijakan backup dan prosedur pemulihan yang diuji.

### SEC-03 — AI Privacy

- Jangan mengirim data finansial ke penyedia LLM eksternal tanpa dasar pemrosesan dan persetujuan yang sesuai.
- Kirim agregasi atau fitur minimum yang dibutuhkan, bukan seluruh riwayat transaksi.
- Jangan menggunakan data pribadi pengguna untuk melatih model global secara otomatis.
- Sediakan dokumentasi penggunaan, retensi, dan penghapusan data.
- Pastikan data yang digunakan untuk training memiliki dasar yang sesuai dan sudah melalui proses minimisasi/de-identifikasi jika diperlukan.

### SEC-04 — Internal Service Security

- Service Python tidak boleh diekspos langsung ke internet publik.
- Terapkan autentikasi antarlayanan atau kontrol jaringan yang memadai.
- Validasi payload dan ukuran request.
- Batasi concurrency dan request duration.
- Hindari log payload transaksi mentah.
- Pisahkan kredensial database dari service yang tidak membutuhkannya.

### SEC-05 — Audit dan Monitoring

- Catat kejadian keamanan yang relevan.
- Monitor kegagalan autentikasi dan akses tidak sah.
- Catat versi model dan hasil operasi ML tanpa menyimpan data pribadi yang tidak diperlukan.
- Sediakan mekanisme mematikan fitur ML melalui konfigurasi atau feature flag.

Kewajiban hukum dan regulasi perlu ditinjau sesuai bentuk peluncuran, lokasi pemrosesan, penyedia layanan, serta apakah nantinya MonTra terhubung dengan lembaga keuangan.

---

## 15. Persyaratan Nonfungsional

### 15.1 Performance

- Target awal p95 response time API baca sederhana: maksimal 500 ms pada lingkungan uji yang ditetapkan, tidak termasuk jaringan client.
- Operasi penulisan transaksi harus memberikan konfirmasi setelah commit database berhasil.
- Pemanggilan AI tidak boleh memblokir jalur simpan transaksi secara wajib.
- Daftar transaksi menggunakan pagination dan query yang terindeks sesuai pola akses.

### 15.2 Reliability

- Kegagalan ML tidak boleh merusak transaksi.
- Operasi multi-record finansial harus atomik.
- Semua pekerjaan terjadwal harus idempoten atau dilindungi dari pemrosesan ganda.
- Backup database dan prosedur pemulihan harus diuji sebelum produksi.

### 15.3 Scalability

- Mendukung pertumbuhan pengguna dan transaksi tanpa mengubah aturan domain utama.
- Service Python dapat diskalakan terpisah jika beban inference meningkat.
- Gunakan background jobs untuk pekerjaan yang tidak perlu berada di request path.
- Jangan menambah message broker atau orchestration kompleks sebelum ada kebutuhan yang terukur.

### 15.4 Maintainability

- Modul domain dipisahkan dengan jelas.
- Schema database dikelola melalui migration.
- API didokumentasikan menggunakan OpenAPI.
- Konfigurasi diatur melalui environment atau secret manager.
- Model memiliki versi, metrik evaluasi, dan prosedur rollback.
- Perubahan aturan finansial memerlukan test dan review.

### 15.5 Accessibility and Usability

- Nominal mudah dibaca dan memakai format lokal konsisten.
- Status tidak hanya ditandai menggunakan warna.
- Form transaksi menghindari input yang tidak perlu.
- Empty, loading, error, dan unavailable states harus dirancang.
- Pengguna dapat mengoperasikan fitur inti meskipun AI tidak tersedia.

---

## 16. Metrik Keberhasilan

### 16.1 Product Metrics

- **Activation:** persentase pengguna baru yang mencatat transaksi pertama.
- **Engagement:** pengguna aktif mingguan dan frekuensi pencatatan transaksi.
- **Retention:** day-7 dan day-30 retention.
- **Budget adoption:** persentase pengguna aktif yang membuat setidaknya satu anggaran.
- **Insight engagement:** persentase pengguna yang membuka detail insight atau mengikuti tautan ke laporan sumber.
- **Data quality:** tingkat duplikasi yang terdeteksi dan transaksi yang dikoreksi pengguna.

### 16.2 AI/ML Metrics

- Kategorisasi: macro-F1, precision, recall, confusion matrix.
- Confidence: kalibrasi dan ketepatan keputusan otomatis.
- Anomali: precision alert, false-positive rate, feedback relevansi.
- Forecast: MAE/WAPE terhadap nilai aktual dan perbandingan dengan baseline.
- Insight: kecocokan angka/periode dengan sumber, keterlacakan, dan feedback pengguna.
- AI assistant bila kelak diterapkan: grounded answer rate, error rate, dan pelanggaran otorisasi.

### 16.3 Target Eksperimen Awal

- 70% pengguna baru menyelesaikan transaksi pertama pada sesi awal.
- 60% pengguna yang dianggap aktif mencatat transaksi secara berulang sesuai definisi analitik yang ditetapkan.
- Macro-F1 kategorisasi minimal 0,80 pada data uji yang representatif.
- Forecast mengurangi error minimal 10% dibanding baseline ketika jumlah dan kualitas data mendukung.

Semua target merupakan hipotesis awal. Target harus diperbarui setelah pilot dan tidak boleh dipromosikan sebagai hasil yang telah dicapai sebelum pengukuran.

---

## 17. Strategi Pengujian dan Acceptance Criteria

### 17.1 Unit Testing

**Go:**
- Perhitungan saldo.
- Pemasukan, pengeluaran, dan transfer.
- Aturan anggaran dan periode.
- Validasi nominal, kategori, tanggal, dan akun.
- Otorisasi resource.
- Parsing dan pemetaan respons ML.

**Python:**
- Preprocessing dan normalisasi teks.
- Validasi request/response.
- Output kategori dan metadata model.
- Logika prediksi, fallback, dan error handling.
- Pipeline evaluasi model.

**Flutter:**
- Validasi form transaksi.
- Format nominal dan tanggal.
- UI state untuk loading, empty, success, dan error.
- Widget dan state management tests.

### 17.2 Integration Testing

- Go dengan PostgreSQL.
- Go dengan Python service.
- Kategorisasi transaksi end-to-end.
- Pembuatan dan pengeditan transaksi memengaruhi laporan secara benar.
- Transfer yang gagal tidak meninggalkan pembaruan sebagian.
- Otorisasi mencegah akses silang pengguna.

### 17.3 AI/ML Evaluation Tests

- Dataset uji terpisah dari training.
- Uji class imbalance dan kategori langka.
- Uji teks kosong, singkatan, ejaan salah, merchant tidak dikenal, dan deskripsi bising.
- Uji data historis minim untuk forecasting.
- Perbandingan model dengan baseline.
- Regresi metrik sebelum mengganti versi model.
- Validasi bahwa insight numerik sesuai evidence yang diberikan.

### 17.4 Security Testing

- Broken access control/IDOR.
- Input validation dan SQL injection.
- Authentication abuse dan rate limiting.
- Kebocoran data melalui log atau error.
- Akses eksternal langsung ke Python service.
- Secret scanning pada repository.
- Pemulihan backup dan penghapusan data sesuai kebijakan.

### 17.5 Acceptance Criteria Global

MVP dianggap siap ketika:
1. Alur autentikasi dan otorisasi lulus uji.
2. CRUD transaksi berjalan dan tervalidasi.
3. Saldo, transfer, anggaran, dan laporan dapat direkonsiliasi ke data sumber.
4. Semua operasi finansial multi-record menggunakan transaksi database yang benar.
5. ML categorization telah dievaluasi dan memiliki fallback.
6. Insight memiliki periode dan evidence yang valid.
7. Kegagalan Python tidak memblokir pencatatan manual.
8. Tidak terdapat masalah keamanan kritis yang belum ditangani.
9. Backup dan restore telah diuji.
10. Dokumentasi API dan deployment mencerminkan implementasi aktual.

---

## 18. Roadmap Implementasi

### Phase 1 — Foundation

**Fokus:**
- Setup repository dan CI dasar.
- Docker Compose untuk lingkungan lokal.
- PostgreSQL dan migrations.
- Authentication/authorization.
- CRUD account, category, transaction.
- Transfer internal dan perhitungan saldo.

**Definition of Done:**
- Alur transaksi dasar berfungsi.
- Perhitungan saldo/transfer punya unit dan integration tests.
- Akses antar-pengguna telah diuji.

### Phase 2 — Core Personal Finance

**Fokus:**
- Dashboard.
- Budget management.
- Reports dan visualisasi.
- Filter transaksi.
- Ekspor CSV.
- Recurring transactions dan notification jika waktu memungkinkan.

**Definition of Done:**
- Angka dashboard dapat direkonsiliasi ke data sumber.
- Periode dan filter laporan konsisten.
- Aplikasi menangani empty/error states dengan baik.

### Phase 3 — ML Foundation

**Fokus:**
- Dataset dan skema kategori.
- Pipeline preprocessing.
- Baseline classifier.
- Evaluasi model.
- FastAPI internal service.
- Go `mlclient` dan confidence handling.

**Definition of Done:**
- Ada hasil evaluasi pada test set.
- Model memiliki versi dan metadata.
- Pengguna dapat mengoreksi kategori.
- Service gagal tidak memblokir transaksi.

### Phase 4 — Intelligent Insights

**Fokus:**
- Spending analytics.
- Anomaly detection.
- Expense forecast.
- Insight cards dengan evidence.

**Definition of Done:**
- Setiap insight punya dasar numerik yang dapat ditelusuri.
- Forecast diuji terhadap baseline.
- Peringatan anomali tidak dinyatakan sebagai kepastian penipuan.
- Fitur bisa dimatikan jika kualitas turun.

### Phase 5 — Controlled Release

**Fokus:**
- Security review.
- Performance testing.
- Observability.
- Backup/restore drill.
- Pilot terbatas.
- Perbaikan bug dan validasi UX.

**Definition of Done:**
- Tidak ada masalah keamanan kritis terbuka.
- Alur inti stabil.
- Restore backup telah teruji.
- Metrik produk dan model mulai dikumpulkan secara sah dan terukur.

---

## 19. Risiko dan Mitigasi

| Risiko | Dampak | Mitigasi |
|---|---|---|
| Histori pengguna terlalu sedikit | Prediksi dan insight tidak akurat | Cold-start fallback, threshold data, tampilkan keterbatasan |
| Kategorisasi salah | Laporan menyesatkan | Confidence policy, koreksi manual, evaluasi per kategori |
| Forecast meleset | Pengguna kehilangan kepercayaan | Bandingkan dengan baseline, tampilkan ketidakpastian, ukur error setelah periode selesai |
| Insight mengandung angka salah | Keputusan pengguna dapat keliru | Hitung angka secara deterministik dan validasi terhadap evidence |
| Data finansial bocor | Kerugian privasi dan kepercayaan | Otorisasi ketat, minimisasi data, enkripsi, pembatasan log, pengujian keamanan |
| AI service tidak tersedia | Fitur cerdas tidak berfungsi | Timeout, fallback, feature flag, isolasi dari jalur transaksi inti |
| Scope MVP melebar | Rilis tertunda | P0 harus selesai sebelum fitur Post-MVP |
| Data training tidak representatif | Kinerja buruk pada pengguna nyata | Evaluasi terpisah, pemantauan drift, pengumpulan data sesuai izin |
| Peringatan berlebihan | Pengguna mengabaikan notifikasi | Deduplikasi, batas frekuensi, feedback relevansi |
| Transfer salah dihitung | Saldo dan laporan tidak konsisten | Model transfer terpisah, transaksi database atomik, test rekonsiliasi |

---

## 20. Kriteria Rilis MVP

Rilis MVP dapat dipertimbangkan jika seluruh kriteria berikut terpenuhi:

### Product
- Pengguna dapat registrasi, login, dan mengelola profil.
- Pengguna dapat membuat dan mengelola akun virtual.
- Pengguna dapat mencatat pemasukan, pengeluaran, dan transfer.
- Pengguna dapat membuat anggaran dan melihat laporan.
- Pengguna dapat mengoreksi kategori transaksi.

### Data Integrity
- Saldo dan laporan konsisten dengan data sumber.
- Transfer internal tidak dihitung ganda.
- Constraint database dan otorisasi diuji.
- Backup serta pemulihan berhasil diuji.

### AI/ML
- Kategorisasi sudah dievaluasi pada data uji yang terpisah.
- Model memiliki versioning dan fallback.
- Insight numerik dapat ditelusuri.
- Forecast/anomaly hanya diaktifkan setelah memenuhi gerbang kualitas yang ditetapkan.

### Security and Operations
- Tidak ada secret yang tersimpan di repository.
- Python service tidak dapat diakses langsung dari internet.
- Logging tidak membocorkan data transaksi secara tidak perlu.
- Health checks, error handling, dan timeout tersedia.
- Tidak ada kerentanan kritis yang diketahui dan belum ditangani.

---

## 21. Asumsi dan Keputusan yang Masih Terbuka

PRD ini memakai asumsi berikut untuk menjaga pekerjaan dapat dimulai:

1. Flutter digunakan sebagai frontend mobile MVP.
2. Golang menjadi satu-satunya backend publik untuk API produk.
3. PostgreSQL menjadi sumber kebenaran transaksi dan data keuangan.
4. Python/FastAPI berjalan sebagai internal AI/ML service.
5. MVP dimulai dengan pencatatan manual; integrasi bank dan OCR bukan prasyarat rilis.
6. Mata uang awal IDR, namun desain data tetap menyimpan kode mata uang.
7. Insight awal menggunakan analitik deterministik; LLM untuk penjelasan bahasa natural bersifat opsional dan bukan kebutuhan awal.
8. Target metrik adalah sasaran pengembangan dan perlu divalidasi lewat pilot.

Keputusan yang wajib dikunci sebelum produksi:
- Mekanisme autentikasi dan manajemen sesi/token.
- Detail skema transaksi dan strategi penghitungan saldo.
- Strategi backup, recovery, dan retensi data.
- Hosting, domain, serta pengelolaan secrets.
- Kebijakan penggunaan data koreksi untuk pelatihan model.
- Batas minimal histori sebelum forecast diaktifkan.
- Threshold confidence dan kebijakan auto-apply kategorisasi.
- Kebijakan notifikasi dan retention.

---

## 22. Glosarium

- **AI:** Artificial Intelligence, teknologi untuk menghasilkan prediksi atau keluaran cerdas.
- **ML:** Machine Learning, metode untuk mempelajari pola dari data.
- **Inference:** Proses menjalankan model terlatih untuk menghasilkan prediksi.
- **Training:** Proses melatih model menggunakan data.
- **Confidence score:** Skor keyakinan keluaran model; tidak otomatis merupakan probabilitas terkalibrasi.
- **Macro-F1:** Rata-rata F1-score antar-kelas dengan bobot sama untuk tiap kelas.
- **MAE:** Mean Absolute Error, rata-rata nilai absolut selisih prediksi dan aktual.
- **WAPE:** Weighted Absolute Percentage Error, metrik error relatif berdasarkan total nilai aktual.
- **Data leakage:** Kebocoran informasi antara data training dan evaluasi yang membuat metrik tampak lebih baik dari kemampuan sebenarnya.
- **Cold start:** Kondisi ketika pengguna atau sistem belum memiliki data historis yang cukup.
- **Source of truth:** Sumber data otoritatif untuk menentukan nilai yang benar.
- **Idempotency:** Sifat operasi yang dapat diulang tanpa menghasilkan efek duplikat yang tidak diinginkan.
- **Internal transfer:** Pemindahan dana tercatat di antara akun milik pengguna sendiri.

---

## Definisi Akhir Produk

MonTra adalah aplikasi pengelolaan keuangan pribadi yang menggabungkan pencatatan transaksi, anggaran, laporan, dan analitik AI/ML. Golang mengelola API, keamanan, dan seluruh business logic finansial. PostgreSQL menjadi sumber kebenaran data. Python/FastAPI menjalankan model kategorisasi, deteksi anomali, dan prediksi secara terpisah.

Prinsip yang tidak boleh dikompromikan adalah akurasi perhitungan finansial, keterlacakan insight, privasi data, serta kendali pengguna. Kesuksesan MonTra dinilai dari kemampuan menghasilkan informasi yang akurat, relevan, dapat dievaluasi, dan membantu pengguna memahami kondisi keuangannya.
