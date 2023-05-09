<a name="readme-top"></a>
<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![Apache License][license-shield]][license-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">

[//]: # (  <a href="https://github.com/ChatDan/ChatDanBackend">)

[//]: # (    <img src="images/logo.png" alt="Logo" width="80" height="80">)

[//]: # (  </a>)

<h3 align="center">ChatDan Backend</h3>

  <p align="center">
    a message box and 'biaobai' platform for Fudaners
    <br />
    <a href="https://github.com/ChatDan/ChatDanBackend"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/ChatDan/ChatDanBackend">View Demo</a>
    ·
    <a href="https://github.com/ChatDan/ChatDanBackend/issues">Report Bug</a>
    ·
    <a href="https://github.com/ChatDan/ChatDanBackend/issues">Request Feature</a>
  </p>
</div>



## About The Project

[//]: # ([![Product Name Screen Shot][product-screenshot]]&#40;https://example.com&#41;)
a message box and 'biaobai' platform for Fudaners

### Built With

[![Go][go.dev]][go-url]
[![Swagger][swagger.io]][swagger-url]

## Getting Started

### Prerequisites

Install mysql and redis. Quick start using docker.

```shell
docker run -d --name mysql \
  -e MYSQL_PASSWORD={MYSQL_PASSWORD} \
  mysql:latest
  
docker run -d --name redis redis:latest
```

then Install Apisix for API Gateway and jwt-auth. Please follow [Apisix Documentation](https://apisix.apache.org/zh/docs/apisix/getting-started/README/)


### Installation

```shell
docker run -d --name chatdan_backend \
  -e DB_URL={DB_URL} \
  -e REDIS_URL={REDIS_URL} \
  -e APISIX_URL={APISIX_URL} \
  -e APISIX_ADMIN_KEY={APISIX_ADMIN_KEY} \
  jingyijun3104/chatdan_backend:latest
```

## Usage

_For more examples, please refer to the [Documentation](https://chatdan-test.jingyijun.xyz:8443/docs)_

## Roadmap

- [x] user management
- [ ] message box
- [ ] chat
- [ ] hole and floor
- [ ] wall

See the [open issues](https://github.com/ChatDan/ChatDanBackend/issues) for a full list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the Apache 2.0 License. See `LICENSE.txt` for more information.

## Contact

JingYiJun - JingYiJun3104@outlook.com

Project Link: [https://github.com/ChatDan/ChatDanBackend](https://github.com/ChatDan/ChatDanBackend)

[//]: # (https://www.markdownguide.org/basic-syntax/#reference-style-links)
[contributors-shield]: https://img.shields.io/github/contributors/ChatDan/ChatDanBackend.svg?style=for-the-badge
[contributors-url]: https://github.com/ChatDan/ChatDanBackend/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/ChatDan/ChatDanBackend.svg?style=for-the-badge
[forks-url]: https://github.com/ChatDan/ChatDanBackend/network/members
[stars-shield]: https://img.shields.io/github/stars/ChatDan/ChatDanBackend.svg?style=for-the-badge
[stars-url]: https://github.com/ChatDan/ChatDanBackend/stargazers
[issues-shield]: https://img.shields.io/github/issues/ChatDan/ChatDanBackend.svg?style=for-the-badge
[issues-url]: https://github.com/ChatDan/ChatDanBackend/issues
[license-shield]: https://img.shields.io/github/license/ChatDan/ChatDanBackend.svg?style=for-the-badge
[license-url]: https://github.com/ChatDan/ChatDanBackend/blob/main/LICENSE
[go.dev]: https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white
[go-url]: https://go.dev
[swagger.io]: https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white
[swagger-url]: https://swagger.io