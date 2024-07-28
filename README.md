<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>


<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://astanait.edu.kz/">
    <img src="https://static.tildacdn.pro/tild3764-6633-4663-b138-303730646233/aitu-logo__2.png" alt="Logo" height="80">
  </a>

<h3 align="center">AITU UCMS comments service</h3>

</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li><a href="#protofiles">Protofiles</a></li>
    <li><a href="#websockets-api">Websockets API</a></li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

[//]: # ([![Product Name Screen Shot][product-screenshot]]&#40;https://example.com&#41;)

This is a simple comments service for [AITU UCMS][aitu-ucms-url] project. It is a part of the project, which is a web application for managing the university's clubs. The service allows users to leave comments on the content of the [AITU UCMS][aitu-ucms-url].

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* [![Go][go-shield]][go-url]
* [![MongoDB][mongodb-shield]][mongodb-url]
* [![Docker][docker-shield]][docker-url]
* [![Docker Compose][docker-compose-shield]][docker-compose-url]
* [![RabbitMQ][rabbitmq-shield]][rabbitmq-url]
* [![Centrifuge][centrifuge-shield]][centrifuge-url]
* [![GRPC][grpc-shield]][grpc-url]


<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- PROTOFILES -->
## Protofiles

* [Protofiles Repository][protofiles-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## [Websockets API](docs/websocket.md) 



<!-- GETTING STARTED -->
## Getting Started
### Prerequisites

* Go version 1.22.3
* Docker 26.1.4
* Docker Compose 2.27.1
* [Taskfile 3](https://taskfile.dev/installation/) 

  ```sh
    go version
    docker --version
    docker-compose --version
  ```

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/ARUMANDESU/uniclubs-comments-service.git
   ```
2. Change directory
   ```sh
   cd uniclubs-comments-service
   ```
3. Write the environment variables in the `.env` file
   ```dotenv
   ENV=dev

   START_TIMEOUT=
   SHUTDOWN_TIMEOUT=
   
   HTTP_ADDRESS=
   HTTP_TIMEOUT=
   HTTP_IDLE_TIMEOUT=
   
   GRPC_PORT=
   GRPC_TIMEOUT=
    
   MONGODB_URI=mongodb://<user>:<password>@<host>:<port>
   MONGODB_PING_TIMEOUT=10s
   MONGODB_DATABASE_NAME=<your_database_name>
    
   RABBITMQ_USER=<user>
   RABBITMQ_PASSWORD=<password>
   RABBITMQ_HOST=<host>
   RABBITMQ_PORT=<port>
   
   USER_SERVICE_ADDRESS=<host>:<port>
   USER_SERVICE_TIMEOUT=10s
   USER_SERVICE_RETRIES_COUNT=2

   JWT_SECRET=
   ```
4. Run the service
   ```sh
   task r:e
   ```


<p align="right">(<a href="#readme-top">back to top</a>)</p>




<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[aitu-url]: https://astanait.edu.kz/
[aitu-ucms-url]: https://www.ucms.space/
[protofiles-url]: https://github.com/ARUMANDESU/uniclubs-protos

[go-url]: https://golang.org/
[mongodb-url]: https://www.mongodb.com/
[docker-url]: https://www.docker.com/
[docker-compose-url]: https://docs.docker.com/compose/
[rabbitmq-url]: https://www.rabbitmq.com/
[websockets-url]: https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API
[centrifuge-url]: https://github.com/centrifugal/centrifuge
[grpc-url]: https://grpc.io/

[go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[mongodb-shield]: https://img.shields.io/badge/MongoDB-47A248?style=for-the-badge&logo=mongodb&logoColor=white
[docker-shield]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[docker-compose-shield]: https://img.shields.io/badge/Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white
[rabbitmq-shield]: https://img.shields.io/badge/RabbitMQ-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white
[websockets-shield]: https://img.shields.io/badge/Websockets-777777?style=for-the-badge&logo=websocket&logoColor=white
[centrifuge-shield]: https://img.shields.io/badge/Centrifuge-FF6600?style=for-the-badge&logo=centrifuge&logoColor=white
[grpc-shield]: https://img.shields.io/badge/GRPC-00ADD8?style=for-the-badge&logo=grpc&logoColor=white