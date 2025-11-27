# Sistem Pelaporan Prestasi Mahasiswa – Backend API

Sistem ini merupakan backend untuk aplikasi pelaporan prestasi mahasiswa yang memungkinkan mahasiswa melaporkan prestasi, dosen wali memverifikasi, dan admin mengelola pengguna. Sistem dibangun untuk mendukung proses akademik secara digital dengan arsitektur modern dan aman.

---

## Penjelasan Sistem
Sistem Pelaporan Prestasi Mahasiswa menyediakan API untuk:

- Pelaporan prestasi oleh mahasiswa  
- Verifikasi prestasi oleh dosen wali  
- Pengelolaan user, role, dan permissions oleh admin  
- Penyimpanan prestasi dengan data dinamis berdasarkan jenis prestasi  
- Workflow status prestasi: draft → submitted → verified/rejected  

Dua database digunakan untuk mendukung fleksibilitas data:
- **PostgreSQL** – data relasional: users, roles, mahasiswa, dosen  
- **MongoDB** – data prestasi dinamis  

---

## Tools
- **Golang (Go)** – backend utama  
- **PostgreSQL** – relational database  
- **MongoDB** – dynamic schema storage  
- **JWT Authentication** – autentikasi  
- **RBAC (Role-Based Access Control)** – hak akses  
- **RESTful API** – arsitektur komunikasi  

---

## Fitur Sistem
- Login & autentikasi JWT  
- Role-Based Access Control  
- CRUD prestasi mahasiswa  
- Field prestasi dinamis (MongoDB)  
- Upload lampiran prestasi  
- Submit – Verify – Reject workflow  
- Manajemen user & role oleh admin  
- Statistik & laporan prestasi  

---

## Pembuat
**Aisha Laily Purwanto**  
**NIM:** 43231052  
**Kelas:** TI-C8  
UAS Pemrograman Backend Lanjut  
DIV Teknik Informatika – Universitas Airlangga  


