Vấn đề:
    * Trong mỗi môi trường hay các hệ điều hành sẽ có những cách chuyển đường dẫn khác nhau
        + Ta có thể sử dụng filepath:
            - join: đường dẫn đến điểm đích dùng cho đa hệ điều hành đa nền tảng
        
witches create example --db=mysql
// Có Folder chứa file .env file cấu hình DB đơn giản 
    example 
        example/.env

cd example 
// Sửa .env

witches database docker-up

witches init

// Tạo Templates

witches install
// Chạy migrate đầu
// Cài các dependency


