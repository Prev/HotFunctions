# Sample Functions

Consists of 36% web server, 27% data processing, 20% third-party integration and 17% for internal tooling.
This ratio is determined based on the serverless [community survey](https://serverless.com/blog/2018-serverless-community-survey-huge-growth-usage/) result introduced in [A Berkeley View on Serverless Computing](https://arxiv.org/abs/1902.03383).

Application request frequency follows the exponential distribution and applications with short execution times are set to executed more frequently.

| Name | Type                 | Language        |  Req. Frequency | Execution Time | Library |
|----|------------------------|-----------------|-------------|-------------|------|
| W1 | WebServer              | Java8 (OpenJDK) | 1.635/s    | ~1s | spring |
| W2 | WebServer              | Java8 (OpenJDK) | 1.273/s  | ~1s | spring |
| W3 | WebServer              | Python3.7       | 0.992/s  | ~1s | requests |
| W4 | WebServer              | NodeJS 12       | 0.772/s  | ~1s | express |
| W5 | WebServer              | Python3.7       | 0.601/s  | ~1s | numpy |
| W6 | WebServer              | Java8 (OpenJDK) | 0.468/s  | ~1s | spring |
| W7 | WebServer              | NodeJS 12       | 0.364/s  | ~1s | jimp |
| W8 | WebServer              | Python3.7       | 0.284/s  | ~1s | requests |
| W9 | WebServer              | NodeJS 12       | 0.221/s  | ~1s | express |
| W10 | WebServer             | Java8 (OpenJDK) | 0.172/s  | ~1s | spring |
| W11 | WebServer             | Java8 (OpenJDK) | 0.134/s   | ~1s | spring |
| T1 | 3rd Party Intergration | Python3.7       | 0.104/s  | ~3s | boto3 |
| T2 | 3rd Party Intergration | Python3.7       | 0.0814/s  | ~3s | boto3 |
| T3 | 3rd Party Intergration | NodeJS 12       | 0.0634/s  | ~3s | axios |
| T4 | 3rd Party Intergration | Python3.7       | 0.0494/s  | ~3s | requests |
| T5 | 3rd Party Intergration | NodeJS 12       | 0.0384/s  | ~3s | axios |
| T6 | 3rd Party Intergration | NodeJS 12       | 0.0300/s  | ~3s | axios |
| I1 | Internal Tooling       | Python3.7       | 0.0233/s  | 3s~ | requests |
| I2 | Internal Tooling       | Python3.7       | 0.0182/s  | 3s~ | requests |
| I3 | Internal Tooling       | NodeJS 12       | 0.0141/s  | 3s~ |
| I4 | Internal Tooling       | NodeJS 12       | 0.0110/s  | 3s~ |
| I5 | Internal Tooling       | Python3.7       | 0.0086/s  | 3s~ |
| D1 | Data Processing        | Python3.7       | 0.0066/s  | 20s~ | numpy |
| D2 | Data Processing        | Python3.7       | 0.0052/s  | 20s~ | numpy |
| D3 | Data Processing        | Python3.7       | 0.0041/s  | 20s~ | numpy |
| D4 | Data Processing        | NodeJS 12       | 0.0031/s  | 20s~ |
| D5 | Data Processing        | NodeJS 12       | 0.0024/s  | 20s~ |
| D6 | Data Processing        | Python3.7       | 0.0019/s  | 20s~ |
| D7 | Data Processing        | Python3.7       | 0.0015/s  | 20s~ |
| D8 | Data Processing        | NodeJS 12       | 0.0012/s  | 20s~ |
