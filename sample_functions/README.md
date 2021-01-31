# Sample Functions

Consists of 36% web server, 27% data processing, 20% third-party integration and 17% for internal tooling.
This ratio is determined based on the serverless [community survey](https://serverless.com/blog/2018-serverless-community-survey-huge-growth-usage/) result introduced in [A Berkeley View on Serverless Computing](https://arxiv.org/abs/1902.03383).

Application request frequency follows the exponential distribution and applications with short execution times are set to executed more frequently.

| Name | Type                 | Language        |  Req. Frequency | Execution Time | Library  | Size |
|----|------------------------|-----------------|-----------------|----------------|----------|------|
| W1 | WebServer              | Java8 (OpenJDK) | 1.635/s         | ~1s            | spring   | 9.7MB |
| W2 | WebServer              | Java8 (OpenJDK) | 1.273/s         | ~1s            | spring   | 9.7MB |
| W3 | WebServer              | Python3.7       | 0.992/s         | ~1s            | requests | 54MB  |
| W4 | WebServer              | NodeJS 12       | 0.772/s         | ~1s            | express  | 2MB   |
| W5 | WebServer              | Python3.7       | 0.601/s         | ~1s            | numpy    | 26MB  |
| W6 | WebServer              | Java8 (OpenJDK) | 0.468/s         | ~1s            | spring   | 9.7MB |
| W7 | WebServer              | NodeJS 12       | 0.364/s         | ~1s            | jimp     | 17MB  |
| W8 | WebServer              | Python3.7       | 0.284/s         | ~1s            | requests | 54MB  |
| W9 | WebServer              | NodeJS 12       | 0.221/s         | ~1s            | express  | 2MB   |
| W10 | WebServer             | Java8 (OpenJDK) | 0.172/s         | ~1s            | spring   | 9.7MB |
| W11 | WebServer             | Java8 (OpenJDK) | 0.134/s         | ~1s            | spring   | 9.7MB |
| T1 | 3rd Party Integration  | Python3.7       | 0.104/s         | ~3s            | boto3    | 98MB  |
| T2 | 3rd Party Integration  | Python3.7       | 0.0814/s        | ~3s            | boto3    | 98MB  |
| T3 | 3rd Party Integration  | NodeJS 12       | 0.0634/s        | ~3s            | axios    | 496KB |
| T4 | 3rd Party Integration  | Python3.7       | 0.0494/s        | ~3s            | requests | 54MB  |
| T5 | 3rd Party Integration  | NodeJS 12       | 0.0384/s        | ~3s            | axios    | 496KB |
| T6 | 3rd Party Integration  | NodeJS 12       | 0.0300/s        | ~3s            | axios    | 496KB |
| I1 | Internal Tooling       | Python3.7       | 0.0233/s        | 3s~            | requests | 54MB  |
| I2 | Internal Tooling       | Python3.7       | 0.0182/s        | 3s~            | requests | 54MB  |
| I3 | Internal Tooling       | NodeJS 12       | 0.0141/s        | 3s~            | -        | 14KB  |
| I4 | Internal Tooling       | NodeJS 12       | 0.0110/s        | 3s~            | -        | 14KB  |
| I5 | Internal Tooling       | Python3.7       | 0.0086/s        | 3s~            | -        | 15KB  |
| D1 | Data Processing        | Python3.7       | 0.0066/s        | 20s~           | numpy    | 131MB |
| D2 | Data Processing        | Python3.7       | 0.0052/s        | 20s~           | numpy    | 131MB |
| D3 | Data Processing        | Python3.7       | 0.0041/s        | 20s~           | numpy    | 131MB |
| D4 | Data Processing        | NodeJS 12       | 0.0031/s        | 20s~           | -        | 14KB  |
| D5 | Data Processing        | NodeJS 12       | 0.0024/s        | 20s~           | -        | 14KB  |
| D6 | Data Processing        | Python3.7       | 0.0019/s        | 20s~           | -        | 14KB  |
| D7 | Data Processing        | Python3.7       | 0.0015/s        | 20s~           | -        | 14KB  |
| D8 | Data Processing        | NodeJS 12       | 0.0012/s        | 20s~           | -        | 14KB  |
