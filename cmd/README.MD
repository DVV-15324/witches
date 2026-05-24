Vấn đề:
    * Trong mỗi môi trường hay các hệ điều hành sẽ có những cách chuyển đường dẫn khác nhau
        + Ta có thể sử dụng filepath:
            - join: đường dẫn đến điểm đích dùng cho đa hệ điều hành đa nền tảng
        
witches init example --db=mysql

cd example

witches install

witches database docker-up

witches database docker-down