## DỰ ÁN: XỔ SỐ KIẾN THIẾT CON GÀ TRỐNG
![alt text](https://i.pinimg.com/originals/5f/ea/b8/5feab8ab8c379fb5a1e59cb2bdb83f6d.jpg)

### CHỨC NĂNG:
1. [X] Đăng ký tài khoản người dùng
2. [X] Đăng nhập tài khoản
3. [X] Đặt cược
4. [X] Hiển thị danh sách đặt cược

### CÀI ĐẶT:
1. Clone dự án về máy và giải nén file LotteryClient.rar
2. Tại folder chứa source code, mở terminal hoặc cmd và chạy lệnh: docker-compose up -d

![image](https://github.com/lphngdo/when_was_the_golden_era_of_video_games/assets/145152948/48dad4cf-7fef-461e-ae1d-9580a9a1c3d9)

3. Mở Docker lên nếu thành công sẽ có trạng thái như hình bên dưới

![image](https://github.com/lphngdo/when_was_the_golden_era_of_video_games/assets/145152948/b63831dc-ee54-4428-b2e3-a08e99077994)

4. Mở app như Tableplus để làm visualize cho MYSQL từ DOCKER để tiện theo dõi DATABASE
5. Sau khi hoàn thành thì giao diện Tableplus sẽ có kết quả như hình bên dưới

![image](https://github.com/lphngdo/when_was_the_golden_era_of_video_games/assets/145152948/7f369727-7fc2-4ea4-9db7-a5cc7210c1b8)

---
## CÂU HỎI: TRIỂN KHAI MONITORING
![image](https://github.com/lphngdo/when_was_the_golden_era_of_video_games/assets/145152948/16e7772c-223d-4841-90dd-d9be65c57d3d)

### CÂU TRẢ LỜI:

Nếu phải monitor performance của ứng dụng xổ số, tôi sẽ track những metrics sau:
1. Tổng số lượng các yêu cầu (requests):
   - http_requests_total: Đếm tổng số lượng các yêu cầu HTTP đến ứng dụng xổ số.
2. Thời gian xử lý yêu cầu (Request duration):
   - http_request_duration_seconds: Đo thời gian mà hệ thống mất để xử lý mỗi yêu cầu. Điều này giúp đánh giá hiệu suất của ứng dụng.
3. Lượng bộ nhớ sử dụng (Memory usage):
   - process_resident_memory_bytes: Đo lượng bộ nhớ thực sự mà tiến trình của ứng dụng đang sử dụng.
4. CPU sử dụng:
   - process_cpu_seconds_total: Đo tổng thời gian CPU mà tiến trình đã sử dụng.
5. Số lượng lỗi (Error rates):
   - http_requests_errors_total: Đếm số lượng yêu cầu gặp lỗi.
6. Thời gian phản hồi (Response time):
    - http_request_duration_seconds: Đo thời gian mà mỗi yêu cầu mất từ khi được nhận đến khi trả về kết quả.
7. Số lượng người dùng đồng thời (Concurrent users):
    - http_requests_in_progress: Đếm số lượng yêu cầu đang được xử lý đồng thời.
8. Số lượng kết nối đến cơ sở dữ liệu (Database connections):
    - database_connections: Đếm số lượng kết nối đến cơ sở dữ liệu, nếu ứng dụng sử dụng cơ sở dữ liệu.
9. Sự hoạt động của các thành phần hệ thống khác (External services):
    - Đo lường thời gian gọi tới các dịch vụ bên ngoài, đảm bảo rằng chúng hoạt động đúng cách.
10. Lưu lượng mạng (Network traffic):
    - http_request_bytes: Đo lượng dữ liệu đang được truyền qua mạng.
      
---
>[!NOTE]
>**_Công nghệ sử dụng:_**
>
>Client: C# (Winform)
>
>Server: Golang (Gin Framework)
>
>Database: MySQL
