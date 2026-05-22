[Cobra] https://cobra.dev/docs/tutorials/first-cli/

Init
Install
Chinh DB_TYPE trong witches.env gồm [mysql, mssql, postgres]
witches database up
Chinh DB_URL == trong witches.env URL (Trả lại)Return từ database up
migrate up and down 
Migration dùng để quản lý schema/database SQL.
Khuyến nghị chạy migration ở terminal riêng, tách biệt với terminal chạy ứng dụng.
Run

Development Notes
Database lifecycle được quản lý qua Docker Compose.
Migration nên được chạy độc lập với runtime application.
witches.env chứa toàn bộ cấu hình hệ thống.