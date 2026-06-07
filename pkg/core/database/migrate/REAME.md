1. Migration là gì?
- Một tập hợp các thay đổi cần thực hiện đối với cấu trúc các đối tượng trong cơ sở dữ liệu quan hệ.
- Đây là một phương pháp để quản lý và thực hiện các thay đổi tăng dần đối với cấu trúc dữ liệu một cách có kiểm soát và lập trình. Những thay đổi này thường có thể đảo ngược, nghĩa là chúng có thể được hoàn tác hoặc khôi phục lại nếu cần.
- Giúp thay đổi lược đồ cơ sở dữ liệu(database schema) từ trạng thái hiện tại sang trạng thái mới mong muốn, cho dù đó là thêm bảng và cột, xóa các phần tử, tách trường hoặc thay đổi kiểu dữ liệu và ràng buộc.
- Bằng cách quản lý những thay đổi này theo phương pháp lập trình, việc duy trì tính nhất quán và độ chính xác trong cơ sở dữ liệu trở nên dễ dàng hơn, cũng như việc theo dõi lịch sử các sửa đổi đã thực hiện.
[Nguồn]https://www.freecodecamp.org/news/database-migration-golang-migrate/

2. Cách tạo một Migrate mới
Tạo các tập tin di chuyển bằng lệnh sau:
$ migrate create -ext sql -dir database/migration/ -seq init_mg
- Bạn sử dụng -seq để tạo phiên bản tuần tự và init_mg là tên của quá trình di chuyển.

Một quá trình di chuyển thường bao gồm hai tệp riêng biệt, một tệp để chuyển cơ sở dữ liệu sang trạng thái mới (được gọi là "up") và một tệp khác để hoàn tác các thay đổi đã thực hiện về trạng thái trước đó (được gọi là "down").

Tệp "up" được sử dụng để thực hiện các thay đổi mong muốn đối với cơ sở dữ liệu, trong khi tệp "down" được sử dụng để hoàn tác các thay đổi đó và đưa cơ sở dữ liệu trở về trạng thái trước đó.

Định dạng của các tệp đó đối với SQL là:
{version}_{title}.down.sql
{version}_{title}.up.sql

=> Khi bạn tạo các tệp di chuyển, chúng sẽ trống theo mặc định. Để thực hiện các thay đổi bạn muốn, bạn cần điền vào đó các truy vấn SQL thích hợp.

3. Cách chạy quá trình Migration

- Để thực thi các câu lệnh SQL trong các tệp di chuyển, công cụ migrate yêu cầu kết nối hợp lệ đến cơ sở dữ liệu.
- Để thực hiện điều này, bạn cần cung cấp chuỗi kết nối ở định dạng phù hợp.

$ migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable" -verbose up

Giờ đây, trong trình shell Postgres của bạn, bạn có thể kiểm tra các bảng mới được tạo bằng cách sử dụng các lệnh sau:

\d+

\d+ table_name DESCRIBE TABLE

4. Cách khắc phục lỗi di chuyển dữ liệu

Nếu một quá trình di chuyển chứa lỗi và vẫn được thực thi, lệnh `migrate` sẽ ngăn không cho bất kỳ quá trình di chuyển nào khác được chạy trên cùng cơ sở dữ liệu đó.

Ngay cả sau khi lỗi trong quá trình di chuyển đã được khắc phục, một thông báo lỗi như "Dirty database version 1. Fix and force version" sẽ vẫn hiển thị. Điều này cho thấy cơ sở dữ liệu "bị lỗi" và cần được điều tra.

Cần phải xác định xem quá trình di chuyển đã được áp dụng một phần hay chưa. Sau khi xác định được điều này, cần phải buộc phiên bản cơ sở dữ liệu phản ánh trạng thái thực của nó bằng lệnh force.

$ migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable" force <VERSION>

5. Cách thêm lệnh vào Makefile

migration_up: migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable" -verbose up

migration_down: migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable" -verbose down

migration_fix: migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable" force VERSION


[Kết Luận]:
    Các hệ thống di chuyển dữ liệu thường tạo ra các tập tin có thể được chia sẻ giữa các nhà phát triển và nhiều nhóm. Chúng cũng có thể được áp dụng cho nhiều cơ sở dữ liệu và được quản lý bằng hệ thống kiểm soát phiên bản.

    Việc ghi chép lại các thay đổi đối với cơ sở dữ liệu giúp theo dõi lịch sử các sửa đổi đã được thực hiện. Bằng cách này, lược đồ cơ sở dữ liệu và sự hiểu biết của ứng dụng về cấu trúc đó có thể cùng phát triển.

    Như vậy là chúng ta đã kết thúc phần thảo luận về việc di chuyển cơ sở dữ liệu. Tôi hy vọng bạn đã thấy những thông tin này hữu ích và bổ ích.

    Nếu bạn thấy bài viết này hữu ích, hãy chia sẻ với đồng nghiệp và bạn bè trên mạng xã hội. Ngoài ra, hãy theo dõi tôi trên Twitter để cập nhật thêm thông tin về công nghệ và lập trình. Cảm ơn bạn đã đọc!



docker run -v {{ migration dir }}:/migrations --network host migrate/migrate
    -path=/migrations/ -database postgres://localhost:5432/database up 2

    docker run --rm   -v "$(pwd)/migrator/migrations:/migrations"   migrate/migrate   -path=/migrations   -database "postgres://docker:docker@host.docker.internal:7557/test_db?sslmode=disable" up 1
