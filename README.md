<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

[//]: # ([![Contributors][contributors-shield]][contributors-url])

[//]: # ([![Forks][forks-shield]][forks-url])

[//]: # ([![Stargazers][stars-shield]][stars-url])

[//]: # ([![Issues][issues-shield]][issues-url])

[//]: # ([![MIT License][license-shield]][license-url])

[//]: # ([![LinkedIn][linkedin-shield]][linkedin-url])



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
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
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


<p align="right">(<a href="#readme-top">back to top</a>)</p>



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
   GRPC_PORT=44046
   GRPC_TIMEOUT=10s
    
   MONGODB_URI=mongodb://<user>:<password>@localhost:<port>
   MONGODB_PING_TIMEOUT=10s
   MONGODB_DATABASE_NAME=your_database_name
    
   RABBITMQ_USER=user
   RABBITMQ_PASSWORD=password
   RABBITMQ_HOST=localhost
   RABBITMQ_PORT=5672
   ```
4. Run the service
   ```sh
   task r:e
   ```


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage




<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/github_username/repo_name.svg?style=for-the-badge
[contributors-url]: https://github.com/github_username/repo_name/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/github_username/repo_name.svg?style=for-the-badge
[forks-url]: https://github.com/github_username/repo_name/network/members
[stars-shield]: https://img.shields.io/github/stars/github_username/repo_name.svg?style=for-the-badge
[stars-url]: https://github.com/github_username/repo_name/stargazers
[issues-shield]: https://img.shields.io/github/issues/github_username/repo_name.svg?style=for-the-badge
[issues-url]: https://github.com/github_username/repo_name/issues
[license-shield]: https://img.shields.io/github/license/github_username/repo_name.svg?style=for-the-badge
[license-url]: https://github.com/github_username/repo_name/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/linkedin_username
[product-screenshot]: images/screenshot.png
[aitu-url]: https://astanait.edu.kz/
[aitu-ucms-url]: https://www.ucms.space/
[go-url]: https://golang.org/
[mongodb-url]: https://www.mongodb.com/
[docker-url]: https://www.docker.com/
[docker-compose-url]: https://docs.docker.com/compose/
[rabbitmq-url]: https://www.rabbitmq.com/
[websockets-url]: https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API
[centrifuge-url]: https://github.com/centrifugal/centrifuge

[go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[mongodb-shield]: https://img.shields.io/badge/MongoDB-47A248?style=for-the-badge&logo=mongodb&logoColor=white
[docker-shield]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[docker-compose-shield]: https://img.shields.io/badge/Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white
[rabbitmq-shield]: https://img.shields.io/badge/RabbitMQ-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white
[websockets-shield]: https://img.shields.io/badge/Websockets-777777?style=for-the-badge&logo=websocket&logoColor=white
[centrifuge-shield]: https://img.shields.io/badge/Centrifuge-FF6600?style=for-the-badge&logo=centrifuge&logoColor=white