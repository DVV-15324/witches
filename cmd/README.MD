

Vấn đề:
    * Trong mỗi môi trường hay các hệ điều hành sẽ có những cách chuyển đường dẫn khác nhau
        + Ta có thể sử dụng filepath:
            - join: đường dẫn đến điểm đích dùng cho đa hệ điều hành đa nền tảng
        


witches init example --db=mysql

witches install

cd example

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