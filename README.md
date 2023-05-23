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

[//]: # (  <a href="https://github.com/ChatDan/chatdan_backend">)

[//]: # (    <img src="images/logo.png" alt="Logo" width="80" height="80">)

[//]: # (  </a>)

<h3 align="center">ChatDan Backend</h3>

  <p align="center">
    a message box and 'biaobai' platform for Fudaners
    <br />
    <a href="https://github.com/ChatDan/chatdan_backend"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/ChatDan/chatdan_backend">View Demo</a>
    ·
    <a href="https://github.com/ChatDan/chatdan_backend/issues">Report Bug</a>
    ·
    <a href="https://github.com/ChatDan/chatdan_backend/issues">Request Feature</a>
  </p>
</div>

## About The Project

[//]: # ([![Product Name Screen Shot][product-screenshot]]&#40;https://example.com&#41;)
a message box and 'biaobai' platform for Fudaners

### Built With

[![Go][go.dev]][go-url]
[![Swagger][swagger.io]][swagger-url]

## Getting Started

### Build and Run locally

#### run

```shell
git clone https://github.com/ChatDan/chatdan_backend.git
cd chatdan_backend
# install swag and generate docs
go install github.com/swaggo/swag/cmd/swag@latest
swag init --parseInternal --parseDepth 1 # to generate the latest docs, this should be run before compiling
# build for debug
go build -o chatdan.exe
# build for release
go build -ldflags "-s -w" -o chatdan.exe
# run
export STANDALONE=true
./treehole.exe
```

For documentation, please open http://localhost:8000/docs after running app

#### test

```shell
export MODE=test
go test -v ./tests/...
```

#### benchmark

```shell
export MODE=bench
go test -v -benchmem -cpuprofile=cpu.out -benchtime=1s ./benchmarks/... -bench .
```

### Production Deploy

Install mysql and redis. Quick start using docker.

```shell
docker run -d --name mysql \
  -e MYSQL_PASSWORD={MYSQL_PASSWORD} \
  mysql:latest
  
docker run -d --name redis redis:latest
```

then Install Apisix for API Gateway and jwt-auth. Please
follow [Apisix Documentation](https://apisix.apache.org/zh/docs/apisix/getting-started/README/)

then start container

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

See the [open issues](https://github.com/ChatDan/chatdan_backend/issues) for a full list of proposed features (and known
issues).

## Contributors

This project exists thanks to all the people who contribute.

<a href="https://github.com/ChatDan/chatdan_backend/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ChatDan/chatdan_backend"  alt="contributors"/>
</a>

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any
contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also
simply open an issue with the tag "enhancement".
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

Project Link: [https://github.com/ChatDan/chatdan_backend](https://github.com/ChatDan/chatdan_backend)

[//]: # (https://www.markdownguide.org/basic-syntax/#reference-style-links)

[contributors-shield]: https://img.shields.io/github/contributors/ChatDan/chatdan_backend.svg?style=for-the-badge

[contributors-url]: https://github.com/ChatDan/chatdan_backend/graphs/contributors

[forks-shield]: https://img.shields.io/github/forks/ChatDan/chatdan_backend.svg?style=for-the-badge

[forks-url]: https://github.com/ChatDan/chatdan_backend/network/members

[stars-shield]: https://img.shields.io/github/stars/ChatDan/chatdan_backend.svg?style=for-the-badge

[stars-url]: https://github.com/ChatDan/chatdan_backend/stargazers

[issues-shield]: https://img.shields.io/github/issues/ChatDan/chatdan_backend.svg?style=for-the-badge

[issues-url]: https://github.com/ChatDan/chatdan_backend/issues

[license-shield]: https://img.shields.io/github/license/ChatDan/chatdan_backend.svg?style=for-the-badge

[license-url]: https://github.com/ChatDan/chatdan_backend/blob/main/LICENSE

[go.dev]: https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white

[go-url]: https://go.dev

[swagger.io]: https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white

[swagger-url]: https://swagger.io