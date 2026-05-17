Cấu trúc logger:
        +----------------+
        |   Your App     |
        | (logger.Info)  |
        +-------+--------+
                |
                v
        +----------------+
        |     Logger     |
        |  (zap.Logger)  |
        +-------+--------+
                |
                v
        +----------------+
        |    zapcore     |
        | (engine log)   |
        +-------+--------+
                |
     +----------+----------+
     |          |          |
     v          v          v
+---------+ +---------+ +-----------+
| Console | |  File   | |   ELK /   |
| stdout  | | app.log | | Loki etc. |
+---------+ +---------+ +-----------+

Nguồn bạn có thể tham khảo: 
[Zap]: https://pkg.go.dev/go.uber.org/zap#section-documentation

