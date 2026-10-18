# gRPC Polyglot Lab

Praktikum **Komunikasi Inter-service dengan gRPC (Polyglot)**.
Satu kontrak `user.proto` dipakai oleh satu gRPC Server (Go) dan dua gRPC Client (Python dan Node.js).

Perintah di bawah dijalankan pada **Git Bash** (Windows). Port server: `50051`.

## Daftar Isi

1. [Struktur Project](#struktur-project)
2. [Prasyarat](#prasyarat)
3. [Langkah 1 - Setup Folder dan Git](#langkah-1---setup-folder-dan-git)
4. [Langkah 2 - Menulis Kontrak user.proto](#langkah-2---menulis-kontrak-userproto)
5. [Langkah 3 - gRPC Server Go](#langkah-3---grpc-server-go)
6. [Langkah 4 - gRPC Client Python](#langkah-4---grpc-client-python)
7. [Langkah 5 - gRPC Client Node.js](#langkah-5---grpc-client-nodejs)
8. [Langkah 6 - Pengujian dan Push ke GitHub](#langkah-6---pengujian-dan-push-ke-github)
9. [Langkah 7 - Tugas Evaluasi Mandiri](#langkah-7---tugas-evaluasi-mandiri)
10. [Analisis Pendekatan Kompilasi](#analisis-pendekatan-kompilasi)

## Struktur Project

```
grpc-polyglot-lab/
├── .gitignore
├── README.md
├── user.proto          # kontrak layanan (Single Source of Truth)
├── go-server/          # gRPC Server (Go)
│   ├── go.mod
│   ├── main.go
│   └── pb/             # kode Go hasil kompilasi user.proto
├── python-client/      # gRPC Client (Python)
│   ├── .venv/
│   ├── user_pb2.py             # hasil kompilasi
│   ├── user_pb2_grpc.py        # hasil kompilasi
│   └── client.py
└── node-client/        # gRPC Client (Node.js)
    ├── package.json
    └── client.js
```

## Prasyarat

Pastikan Git, Python 3.10+, Node.js (LTS), Go, dan `protoc` sudah terpasang. Verifikasi:

```bash
git --version
python --version
node -v
npm -v
go version
protoc --version
```

## Langkah 1 - Setup Folder dan Git

```bash
cd ~/Documents
mkdir grpc-polyglot-lab
cd grpc-polyglot-lab
git init

cat > .gitignore <<'EOF'
node_modules/
__pycache__/
.venv/
*.exe
EOF

git add .gitignore
git commit -m "Initial commit: Add gitignore rules"
```

## Langkah 2 - Menulis Kontrak user.proto

Buat `user.proto` di folder utama proyek:

```proto
syntax = "proto3";

package user;

option go_package = "./pb";

message UserRequest {
  int32 user_id = 1;
}

message UserResponse {
  string name = 1;
  string email = 2;
}

service UserService {
  rpc GetUserProfile(UserRequest) returns (UserResponse);
}
```

## Langkah 3 - gRPC Server Go

### 3.1 Inisialisasi modul dan dependensi

```bash
mkdir go-server
cd go-server
go mod init go-server

go get google.golang.org/grpc
go get google.golang.org/protobuf

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"

protoc-gen-go --version
protoc-gen-go-grpc --version
```

### 3.2 Kompilasi user.proto menjadi kode Go

```bash
mkdir -p pb
protoc --proto_path=.. \
  --go_out=pb --go_opt=paths=source_relative \
  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
  ../user.proto
ls -la pb
```

Hasil: `pb/user.pb.go` dan `pb/user_grpc.pb.go`.

### 3.3 Menulis main.go

Buat `go-server/main.go`:

```go
package main

import (
	"context"
	"log"
	"net"

	pb "go-server/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUserProfile(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log.Printf("Request diterima untuk user_id=%d", req.GetUserId())

	if req.GetUserId() == 1 {
		return &pb.UserResponse{
			Name:  "Budi Santoso",
			Email: "budi@universitas.ac.id",
		}, nil
	}

	return nil, status.Error(codes.NotFound, "user tidak ditemukan")
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Gagal membuka port: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{})

	log.Println("gRPC server berjalan di port 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Server gagal: %v", err)
	}
}
```

> Ubah `Name` dan `Email` sesuai data yang ingin ditampilkan.

### 3.4 Menjalankan server

```bash
go mod tidy
go run main.go
```

Server mencetak `gRPC server berjalan di port 50051`. Biarkan terminal ini tetap terbuka.

## Langkah 4 - gRPC Client Python

### 4.1 Virtual environment dan library

```bash
cd ~/Documents/grpc-polyglot-lab
mkdir python-client
cd python-client

python -m venv .venv
source .venv/Scripts/activate

python -m pip install grpcio grpcio-tools
python -m pip show grpcio grpcio-tools
```

### 4.2 Kompilasi user.proto menjadi modul Python

```bash
python -m grpc_tools.protoc --proto_path=../ --python_out=. --grpc_python_out=. ../user.proto
```

Hasil: `user_pb2.py` dan `user_pb2_grpc.py`.

### 4.3 Menulis client.py

```python
import grpc
import user_pb2
import user_pb2_grpc


def main():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = user_pb2_grpc.UserServiceStub(channel)

        request = user_pb2.UserRequest(user_id=1)
        response = stub.GetUserProfile(request)

        print("Nama:", response.name)
        print("Email:", response.email)


if __name__ == "__main__":
    main()
```

## Langkah 5 - gRPC Client Node.js

### 5.1 Inisialisasi npm dan library

```bash
cd ~/Documents/grpc-polyglot-lab
mkdir node-client
cd node-client

npm init -y
npm install @grpc/grpc-js @grpc/proto-loader
```

### 5.2 Menulis client.js

Client Node.js membaca `../user.proto` langsung saat program berjalan (tanpa `protoc`).

```javascript
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");

const PROTO_PATH = "../user.proto";

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const userPackage = protoDescriptor.user;

const client = new userPackage.UserService(
    "localhost:50051",
    grpc.credentials.createInsecure()
);

client.GetUserProfile({ user_id: 1 }, (error, response) => {
    if (error) {
        console.error("Error:", error.message);
        client.close();
        return;
    }

    console.log("Nama:", response.name);
    console.log("Email:", response.email);

    client.close();
});
```

## Langkah 6 - Pengujian dan Push ke GitHub

### 6.1 Pengujian

Buka tiga terminal.

**Terminal 1 - Server Go**

```bash
cd go-server
go run main.go
```

**Terminal 2 - Client Python**

```bash
cd python-client
source .venv/Scripts/activate
python client.py
```

**Terminal 3 - Client Node.js**

```bash
cd node-client
node client.js
```

Kriteria berhasil:

1. Tidak ada error crash (misalnya `Connection refused` atau `Cannot find module`) di ketiga terminal.
2. Data profil tercetak di terminal Python dan Node.js, padahal data hanya ditulis (hardcoded) di server Go.
3. Terminal Go Server menampilkan log `Request diterima untuk user_id=1` setiap client dijalankan.

### 6.2 Push ke GitHub

Buat repositori kosong baru di akun GitHub, lalu jalankan dari folder utama proyek:

```bash
cd ~/Documents/grpc-polyglot-lab

git add .
git commit -m "Feat: Add full polyglot implementation of Go server, Python client, and Node client"

git branch -M main
git remote add origin https://github.com/nisatridasari/Praktikum-gRPC-polyglot.git
git push -u origin main
```


## Langkah 7 - Tugas Evaluasi Mandiri

### 7.1 Menambah field `phone`

1. Tambahkan `string phone = 3;` ke `message UserResponse` di `user.proto`:

   ```proto
   message UserResponse {
     string name = 1;
     string email = 2;
     string phone = 3;
   }
   ```

2. Kompilasi ulang untuk Go dan Python:

   ```bash
   # dari folder go-server
   protoc --proto_path=.. \
     --go_out=pb --go_opt=paths=source_relative \
     --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
     ../user.proto

   # dari folder python-client (venv aktif)
   python -m grpc_tools.protoc --proto_path=../ --python_out=. --grpc_python_out=. ../user.proto
   ```

3. Server Go: isi field `Phone` pada `UserResponse` di `main.go`.
4. Client Python: tambahkan `print("Telepon:", response.phone)` di `client.py`.
5. Client Node.js: tambahkan `console.log("Telepon:", response.phone);` di `client.js`. Node.js tidak perlu kompilasi ulang karena membaca `user.proto` saat program berjalan.
6. Jalankan ulang server dan kedua client, lalu pastikan nomor telepon tercetak di terminal.

### 7.2 Analisis dan push README

Analisis ada pada bagian di bawah ini. Simpan ke Git:

```bash
git add README.md
git commit -m "Docs: Add README with compilation approach analysis"
git push
```

## Analisis Pendekatan Kompilasi

Pada pendekatan **static code generation** (Go dan Python), file `user.proto` diproses lebih dahulu
oleh `protoc` menjadi file kode sebelum program dijalankan. Pada Go, hasilnya adalah `user.pb.go` dan
`user_grpc.pb.go` di folder `pb`, yang kemudian diimpor oleh `main.go`. Pada Python, hasilnya adalah
modul `user_pb2` dan `user_pb2_grpc` yang diimpor oleh `client.py`. Dengan alur ini, setiap perubahan
pada `user.proto` harus diikuti kompilasi ulang agar kode ikut berubah.

Pada pendekatan **dynamic runtime generation** (Node.js), tidak ada langkah `protoc` terpisah.
`client.js` memuat `../user.proto` ketika program dijalankan melalui `@grpc/proto-loader`
(`protoLoader.loadSync`), lalu membentuk definisi layanan dengan `grpc.loadPackageDefinition`.
Dengan demikian, pemrosesan skema terjadi saat program berjalan, bukan pada langkah kompilasi
terpisah. Kedua alur menghasilkan keluaran yang sama pada praktikum ini. Kecepatan dan penggunaan
sumber daya tidak diukur, sehingga analisis ini hanya membandingkan perbedaan alur kerja.
