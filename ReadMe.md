> Ứng dụng GoHugo là một ứng dụng web nhỏ được viết bằng Golang để tự động hoá việc phát hành web site Hugo.


## Hugo là gì?

Hugo là phần mềm viết bằng Golang sẽ dịch hệ thống file Markdown sang web site tĩnh HTML, CSS.
Hugo sử dụng tập themes để thay đổi giao diện của web site tĩnh.

Hugo có 2 chế độ:
1. Live render Markdown và phục vụ như là một web server
2. Render Markdown ra HTML hoặc minified HTML, cần dùng web server khác như Nginx, Go Caddy để phục vụ

## Work flow của ứng dụng GoHugo
Tác giả web site sẽ soạn các file Markdown trong một thư mục Hugo rồi push lên Github.

Cứ mỗi lần push git, thì Github tạo ra một event hook thực chất là một HTTP request gọi đến một địa chỉ cấu hình. Nếu địa chỉ được GoHugo hứng thì GoHugo sẽ đọc nội dung của HTTP request do Github thực hiện, rồi thực hiện:

- git pull <repo chứa Markdown>
- 


## File cấu hình config.yml
Trong file [config.yml](config.yml) sẽ liệt kê link git repo, thư mục mã nguồn markdown Hugo, thư mục minified HTML. Khi github repo được push code, github sẽ có event hook gọi vào GoHugo, GoHugo sẽ chạy lệnh hugo để minify.